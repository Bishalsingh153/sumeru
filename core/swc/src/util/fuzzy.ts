/** Subsequence / prefix fuzzy match (shared by app launcher and IM user search). */

export function fuzzyScore(query: string, text: string): number {
  const q = query.trim().toLowerCase();
  const t = text.trim().toLowerCase();
  if (!q) return 1;
  if (!t) return 0;
  if (t === q) return 100;
  if (t.startsWith(q)) return 80;
  if (t.includes(q)) return 60;

  let qi = 0;
  for (let i = 0; i < t.length && qi < q.length; i++) {
    if (t[i] === q[qi]) qi++;
  }
  return qi === q.length ? 40 : 0;
}

export interface FuzzyRanked<T> {
  item: T;
  score: number;
}

export function rankByFuzzy<T>(query: string, items: T[], fields: (item: T) => string[]): FuzzyRanked<T>[] {
  const q = query.trim();
  if (!q) return items.map((item) => ({ item, score: 1 }));
  return items
    .map((item) => {
      const vals = fields(item);
      const score = Math.max(0, ...vals.map((v) => fuzzyScore(q, v)));
      return { item, score };
    })
    .filter(({ score }) => score > 0)
    .sort((a, b) => b.score - a.score);
}
