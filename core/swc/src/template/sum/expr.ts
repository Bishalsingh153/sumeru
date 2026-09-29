/** Reject template expressions that could escape the props scope. */
const UNSAFE_TEMPLATE_EXPR = /[`\\[\];]|=>|\bfunction\b|\bclass\b|\bimport\b|\beval\b|\bnew\b/i;

const EXPR_KEYWORDS = new Set([
  "true",
  "false",
  "null",
  "undefined",
  "props",
  "scope",
  "mergeScope",
  "env",
  "when",
  "forEach",
  "String",
  "Number",
  "Boolean",
]);

/** Bare identifiers referenced in a template expression (for sub-template destructuring). */
export function collectExprIdentifiers(expr: string): string[] {
  const ids: string[] = [];
  const stripped = expr.replace(/'[^']*'|"[^"]*"/g, " ");
  const re = /\b([A-Za-z_]\w*)\b/g;
  let m: RegExpExecArray | null;
  while ((m = re.exec(stripped))) {
    if (!EXPR_KEYWORDS.has(m[1])) {
      ids.push(m[1]);
    }
  }
  return ids;
}

export function assertSafeTemplateExpr(expr: string, where: string): void {
  const trimmed = expr.trim();
  if (!trimmed) {
    throw new Error(`sum-template: empty expression in ${where}`);
  }
  if (UNSAFE_TEMPLATE_EXPR.test(trimmed)) {
    throw new Error(`sum-template: unsafe expression in ${where}: ${trimmed}`);
  }
}
