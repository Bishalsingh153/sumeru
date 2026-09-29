import type { SumDocument, SumElement, SumNode } from "./parser.js";
import { parseSumDocument } from "./parser.js";
import { assertSafeTemplateExpr, collectExprIdentifiers } from "./expr.js";
import { templateMeta, type TemplateSourceMeta } from "./meta.js";

function esc(s: string): string {
  return s.replace(/\\/g, "\\\\").replace(/`/g, "\\`").replace(/\$/g, "\\$");
}

function tplFnName(tName: string): string {
  const safe = tName.replace(/\W+/g, "_");
  return `__tpl_${safe}`;
}

function stripTemplateAttrs(attrs: Record<string, string>): Record<string, string> {
  const out: Record<string, string> = {};
  for (const [k, v] of Object.entries(attrs)) {
    if (k.startsWith("t-")) continue;
    out[k] = v;
  }
  return out;
}

function dynamicAttrBindings(attrs: Record<string, string>): string {
  const staticAttrs = stripTemplateAttrs(attrs);
  let out = "";
  for (const [k, v] of Object.entries(attrs)) {
    if (!k.startsWith("t-att-")) continue;
    const attrName = k.slice("t-att-".length);
    if (!attrName) continue;
    assertSafeTemplateExpr(v, k);
    const staticVal = staticAttrs[attrName];
    delete staticAttrs[attrName];
    if (staticVal !== undefined && staticVal !== "") {
      if (attrName === "class") {
        out += ` class=\${${JSON.stringify(`${staticVal} `)} + String(${v})}`;
        continue;
      }
    }
    out += ` ${attrName}=\${${v}}`;
  }
  return out;
}

function attrToHtml(attrs: Record<string, string>): string {
  const bindings = dynamicAttrBindings(attrs);
  const staticOnly = { ...stripTemplateAttrs(attrs) };
  for (const k of Object.keys(attrs)) {
    if (k.startsWith("t-att-")) {
      const attrName = k.slice("t-att-".length);
      delete staticOnly[attrName];
    }
  }
  const staticStr = Object.entries(staticOnly)
    .map(([k, v]) => (v ? `${k}="${esc(v)}"` : k))
    .join(" ");
  return staticStr + bindings;
}

function collectTemplateIdentifiers(root: SumElement): Set<string> {
  const ids = new Set<string>();
  const noteExpr = (expr: string): void => {
    for (const id of collectExprIdentifiers(expr)) {
      ids.add(id);
    }
  };

  const walk = (nodes: SumNode[]): void => {
    for (const node of nodes) {
      if (node.type === "text") continue;
      if (node.type === "interpolation") {
        noteExpr(node.expr);
        continue;
      }
      const el = node;
      for (const [key, raw] of Object.entries(el.attrs)) {
        if (key === "t-else") continue;
        if (key === "t-foreach") {
          const m = raw.match(/^\w+\s+in\s+(.+)$/);
          if (m) noteExpr(m[1]);
          continue;
        }
        if (
          key.startsWith("t-att-") ||
          key === "t-if" ||
          key === "t-elif" ||
          key === "t-esc" ||
          key === "t-raw" ||
          key === "t-out" ||
          key === "t-value" ||
          key === "t-context"
        ) {
          noteExpr(raw);
        }
      }
      walk(el.children);
    }
  };

  if (root.children.length > 0) {
    walk(root.children);
  } else {
    walk([root]);
  }
  return ids;
}

function destructurePrelude(ids: Set<string>): string {
  if (ids.size === 0) return "";
  const lines = [...ids].sort().map((id) => `  const ${id} = scope[${JSON.stringify(id)}];`);
  return `${lines.join("\n")}\n  `;
}

function assertNoCallCycles(templates: Map<string, SumElement>): void {
  for (const [name, el] of templates) {
    const stack = [name];
    const walk = (node: SumElement, chain: string[]): void => {
      const call = node.attrs["t-call"]?.trim();
      if (call) {
        if (call === name || chain.includes(call)) {
          throw new Error(`sum-template: recursive t-call: ${[...chain, call].join(" -> ")}`);
        }
        const target = templates.get(call);
        if (target) walk(target, [...chain, call]);
      }
      for (const child of node.children) {
        if (child.type === "element") walk(child, chain);
      }
    };
    walk(el, stack);
  }
}

function elementBlock(el: SumElement, scopeVar: string, templates: Map<string, SumElement>): string {
  const attrs = { ...el.attrs };
  delete attrs["t-if"];
  delete attrs["t-elif"];
  delete attrs["t-else"];
  delete attrs["t-foreach"];
  delete attrs["t-key"];
  return `html\`<${el.tag} ${attrToHtml(attrs)}>${codegenChildren(el.children, scopeVar, templates)}</${el.tag}>\``;
}

/** Avoid nested html` inside html` when t-set wraps a single ${expr} interpolation. */
function wrapSetterBinding(name: string, expr: string, chunk: string): string {
  const trimmed = chunk.trim();
  if (trimmed.startsWith("${") && trimmed.endsWith("}")) {
    return `((${name}) => ${trimmed.slice(2, -1)})(${expr})`;
  }
  return `((${name}) => html\`${chunk}\`)(${expr})`;
}

function codegenChildren(nodes: SumNode[], scopeVar: string, templates: Map<string, SumElement>): string {
  let out = "";
  let i = 0;
  while (i < nodes.length) {
    const setters: { name: string; expr: string }[] = [];
    while (i < nodes.length) {
      const node = nodes[i];
      if (node.type !== "element" || !node.attrs["t-set"]) break;
      const varName = node.attrs["t-set"].trim();
      if (!varName) throw new Error("sum-template: t-set requires a variable name");
      let expr = node.attrs["t-value"]?.trim();
      if (!expr) {
        const textOnly = node.children.length === 1 && node.children[0].type === "text";
        if (textOnly) {
          expr = JSON.stringify((node.children[0] as { value: string }).value);
        } else if (node.children.length === 0) {
          throw new Error(`sum-template: t-set ${varName} requires t-value or text body`);
        } else {
          throw new Error(`sum-template: t-set ${varName} with nested elements is not supported`);
        }
      } else {
        assertSafeTemplateExpr(expr, `t-set ${varName}`);
      }
      setters.push({ name: varName, expr });
      i += 1;
    }

    if (setters.length > 0) {
      let chunk = codegenChildrenWithoutSetPrefix(nodes.slice(i), scopeVar, templates);
      for (let j = setters.length - 1; j >= 0; j--) {
        const { name, expr } = setters[j];
        chunk = wrapSetterBinding(name, expr, chunk);
      }
      out += `\${${chunk}}`;
      break;
    }

    const piece = codegenSingleSibling(nodes, i, scopeVar, templates);
    out += piece.out;
    i = piece.next;
  }
  return out;
}

/** Siblings only — no leading t-set prefix (used inside t-set scope blocks). */
function codegenChildrenWithoutSetPrefix(
  nodes: SumNode[],
  scopeVar: string,
  templates: Map<string, SumElement>,
): string {
  let out = "";
  let i = 0;
  while (i < nodes.length) {
    const piece = codegenSingleSibling(nodes, i, scopeVar, templates);
    out += piece.out;
    i = piece.next;
  }
  return out;
}

function codegenSingleSibling(
  nodes: SumNode[],
  i: number,
  scopeVar: string,
  templates: Map<string, SumElement>,
): { out: string; next: number } {
  if (i >= nodes.length) return { out: "", next: i };
  const node = nodes[i];
  if (node.type === "element" && node.attrs["t-if"]) {
    const branches: SumElement[] = [node];
    let next = i + 1;
    while (next < nodes.length) {
      const candidate = nodes[next];
      if (candidate.type !== "element") break;
      if ("t-elif" in candidate.attrs) {
        branches.push(candidate);
        next += 1;
        continue;
      }
      if ("t-else" in candidate.attrs) {
        branches.push(candidate);
        next += 1;
      }
      break;
    }
    const first = branches[0];
    assertSafeTemplateExpr(first.attrs["t-if"], "t-if");
    const elifs = branches.slice(1).map((branch) => {
      const cond = "t-else" in branch.attrs ? "true" : branch.attrs["t-elif"];
      if (cond !== "true") assertSafeTemplateExpr(cond, "t-elif");
      return `[${cond}, () => ${elementBlock(branch, scopeVar, templates)}]`;
    });
    const out = `\${when(${first.attrs["t-if"]}, () => ${elementBlock(first, scopeVar, templates)}${
      elifs.length ? `, ${elifs.join(", ")}` : ""
    })}`;
    return { out, next };
  }
  const out = codegenNode(node, scopeVar, templates);
  return { out, next: i + 1 };
}

function codegenNode(node: SumNode, scopeVar: string, templates: Map<string, SumElement>): string {
  if (node.type === "text") return esc(node.value);
  if (node.type === "interpolation") {
    assertSafeTemplateExpr(node.expr, "interpolation");
    return `\${${node.expr}}`;
  }

  const el = node;
  const tCall = el.attrs["t-call"]?.trim();
  if (tCall) {
    if (!templates.has(tCall)) {
      throw new Error(`sum-template: unknown t-call ${tCall}`);
    }
    const tContext = el.attrs["t-context"]?.trim();
    let callScope = scopeVar;
    if (tContext) {
      assertSafeTemplateExpr(tContext, "t-context");
      callScope = `mergeScope(${scopeVar}, ${tContext})`;
    }
    return `\${${tplFnName(tCall)}(${callScope}, env)}`;
  }

  if (el.attrs["t-set"]) {
    return codegenChildren([el], scopeVar, templates);
  }

  const tForeach = el.attrs["t-foreach"];
  const tKey = el.attrs["t-key"];
  const tEsc = el.attrs["t-esc"];
  const tRaw = el.attrs["t-raw"] ?? el.attrs["t-out"];
  const tComponent = el.attrs["t-component"];
  const tModel = el.attrs["t-model"];
  const tRef = el.attrs["t-ref"];
  const tPortal = el.attrs["t-portal"];
  const tSlot = el.attrs["t-slot"];

  if (tForeach) {
    const m = tForeach.match(/^(\w+)\s+in\s+(.+)$/);
    if (!m) throw new Error(`Invalid t-foreach: ${tForeach}`);
    const [, item, collection] = m;
    assertSafeTemplateExpr(collection, "t-foreach");
    const keyExpr = tKey ? tKey.replace(item, `${item}`) : `\`${item}\``;
    const inner = codegenChildren(el.children, scopeVar, templates);
    return `\${forEach(${collection}, (${item}) => ${keyExpr}, (${item}) => html\`<${el.tag} ${attrToHtml(el.attrs)}>${inner}</${el.tag}>\`)}`;
  }

  if (el.attrs["t-if"]) {
    return codegenChildren([el], scopeVar, templates);
  }

  if (tComponent) {
    return `\${mountComponent(${tComponent}, { ...${scopeVar} }, env)}`;
  }

  if (tEsc) {
    assertSafeTemplateExpr(tEsc, "t-esc");
    return `\${${tEsc}}`;
  }
  if (tRaw) {
    assertSafeTemplateExpr(tRaw, "t-raw");
    return `\${${tRaw}}`;
  }

  const attrs = { ...el.attrs };
  delete attrs["t-ref"];
  delete attrs["t-model"];
  delete attrs["t-portal"];
  delete attrs["t-slot"];
  delete attrs["t-out"];
  delete attrs["t-raw"];

  let attrStr = attrToHtml(attrs);
  if (tModel) {
    attrStr += ` @input=\${(event) => { ${tModel} = inputValueFromEvent(event); }} value=\${${tModel}}`;
  }
  if (tRef) attrStr += ` data-ref="${esc(tRef)}"`;
  if (tPortal) attrStr += ` data-portal="${esc(tPortal)}"`;
  if (tSlot) attrStr += ` data-slot="${esc(tSlot)}"`;

  const inner = codegenChildren(el.children, scopeVar, templates);
  const voidTag = ["img", "br", "hr", "input", "meta", "link"].includes(el.tag.toLowerCase());

  if (voidTag) return `<${el.tag} ${attrStr} />`;
  return `<${el.tag} ${attrStr}>${inner}</${el.tag}>`;
}

function codegenTemplateBody(
  template: SumElement,
  scopeVar: string,
  templates: Map<string, SumElement>,
): string {
  const body = codegenChildren(template.children.length ? template.children : [template], scopeVar, templates);
  return template.tag === "t" ? body : codegenNode(template, scopeVar, templates);
}

export interface CodegenResult {
  code: string;
  meta: TemplateSourceMeta;
}

export function codegen(doc: SumDocument, componentName: string, file: string): CodegenResult {
  assertNoCallCycles(doc.templates);

  const tplFns: string[] = [];
  for (const [tName, el] of doc.templates) {
    const fn = tplFnName(tName);
    const body = codegenTemplateBody(el, "scope", doc.templates);
    const ids = collectTemplateIdentifiers(el);
    const prelude = destructurePrelude(ids);
    tplFns.push(`function ${fn}(scope: Record<string, unknown>, env: SwcEnv) {\n${prelude}return html\`${body}\`;\n}`);
  }

  const mainBody = codegenTemplateBody(doc.main, "props", doc.templates);
  const wrapped = doc.main.tag === "t" ? mainBody : codegenNode(doc.main, "props", doc.templates);

  const code = `import { html } from "../../template/html.js";
import { forEach, when, mergeScope } from "../../template/helpers.js";
import { mountComponent } from "../../runtime/component-host.js";
import { inputValueFromEvent } from "../../widgets/field-events.js";
import type { SwcEnv } from "../../runtime/env.js";

${tplFns.join("\n\n")}

export function template(props: Record<string, unknown>, env: SwcEnv) {
  return html\`${wrapped}\`;
}
`;
  return {
    code,
    meta: templateMeta(componentName, file, { snippet: wrapped.slice(0, 200) }),
  };
}

export function compileSumXml(source: string, componentName: string, file: string): CodegenResult {
  return codegen(parseSumDocument(source), componentName, file);
}
