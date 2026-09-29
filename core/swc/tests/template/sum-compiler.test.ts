import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { describe, expect, it } from "vitest";
import { compileSumXml } from "../../src/template/sum/codegen.js";
import { html } from "../../src/template/html.js";
import { forEach, mergeScope, when } from "../../src/template/helpers.js";
import type { SwcEnv } from "../../src/runtime/env.js";

const fixtureDir = dirname(fileURLToPath(import.meta.url));

function runCompiledTemplate(code: string) {
  let body = code.replace(/^import .*;\n/gm, "");
  body = body
    .replace(/: Record<string, unknown>/g, "")
    .replace(/: SwcEnv/g, "")
    .replace(/\bexport function template\b/, "function template");
  return new Function(
    "html",
    "mergeScope",
    "forEach",
    "when",
    "mountComponent",
    "inputValueFromEvent",
    body + "\nreturn template;",
  ) as (...args: unknown[]) => (props: Record<string, unknown>, env: SwcEnv) => ReturnType<typeof html>;
}

describe("sum-template compiler", () => {
  it("compiles t-foreach and t-if", () => {
    const xml = `<div t-foreach="item in items" t-key="item.id"><span t-if="item.active">{{ item.name }}</span></div>`;
    const { code, meta } = compileSumXml(xml, "Demo", "demo.sum.xml");
    expect(code).toContain("forEach");
    expect(code).toContain("when");
    expect(meta.file).toBe("demo.sum.xml");
  });

  it("compiles t-elif, t-else, t-out, and t-model from a fixture", () => {
    const source = readFileSync(join(fixtureDir, "fixtures/branch-out.sum.xml"), "utf8");
    const { code } = compileSumXml(source, "BranchOut", "branch-out.sum.xml");
    expect(code).toContain("when(n === 1");
    expect(code).toContain("[n === 2,");
    expect(code).toContain("[true,");
    expect(code).toContain("${markup}");
    expect(code).toContain("inputValueFromEvent");
  });

  it("compiles t-name, t-call, t-set, and t-att-*", () => {
    const source = readFileSync(join(fixtureDir, "fixtures/composition.sum.xml"), "utf8");
    const { code } = compileSumXml(source, "Composition", "composition.sum.xml");
    expect(code).toContain("function __tpl_chip");
    expect(code).toContain("mergeScope");
    expect(code).toContain("const kind = scope");
    expect(code).toContain('"chip "');
    expect(code).toContain("__tpl_chip");
  });

  it("runs composition fixture at runtime", () => {
    const source = readFileSync(join(fixtureDir, "fixtures/composition.sum.xml"), "utf8");
    const { code } = compileSumXml(source, "Composition", "composition.sum.xml");
    const mountComponent = () => html`<div></div>`;
    const inputValueFromEvent = () => "";
    const env = {} as SwcEnv;
    const template = runCompiledTemplate(code)(html, mergeScope, forEach, when, mountComponent, inputValueFromEvent);
    const root = template({ title: "Hello" }, env).render();
    const span = root.querySelector("span.chip");
    expect(span?.className).toBe("chip primary");
    expect(span?.textContent).toBe("Hello");
  });

  it("t-set applies to all following siblings", () => {
    const xml = `<div><t t-set="x" t-value="'ok'"/><span>{{ x }}</span><em>{{ x }}</em></div>`;
    const { code } = compileSumXml(xml, "SetScope", "set.sum.xml");
    const mountComponent = () => html`<div></div>`;
    const inputValueFromEvent = () => "";
    const env = {} as SwcEnv;
    const template = runCompiledTemplate(code)(html, mergeScope, forEach, when, mountComponent, inputValueFromEvent);
    const root = template({}, env).render();
    expect(root.querySelector("span")?.textContent).toBe("ok");
    expect(root.querySelector("em")?.textContent).toBe("ok");
  });

  it("rejects duplicate t-name", () => {
    const xml = `<t t-name="a"></t><t t-name="a"></t><div/>`;
    expect(() => compileSumXml(xml, "Dup", "dup.sum.xml")).toThrow(/duplicate t-name/);
  });

  it("rejects unknown t-call", () => {
    const xml = `<div><t t-call="missing"/></div>`;
    expect(() => compileSumXml(xml, "Bad", "bad.sum.xml")).toThrow(/unknown t-call/);
  });

  it("rejects empty t-value on t-set", () => {
    const xml = `<div><t t-set="x"/></div>`;
    expect(() => compileSumXml(xml, "Bad", "bad.sum.xml")).toThrow(/t-value/);
  });

  it("rejects recursive t-call between named templates", () => {
    const xml = `<t t-name="a"><t t-call="b"/></t><t t-name="b"><t t-call="a"/></t><div></div>`;
    expect(() => compileSumXml(xml, "Cycle", "cycle.sum.xml")).toThrow(/recursive t-call/);
  });
});
