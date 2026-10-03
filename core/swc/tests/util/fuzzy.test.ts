import { describe, expect, it } from "vitest";
import { fuzzyScore, rankByFuzzy } from "../../src/util/fuzzy.js";

describe("fuzzy", () => {
  it("fuzzyScore prefers prefix matches", () => {
    expect(fuzzyScore("cr", "crm")).toBeGreaterThan(fuzzyScore("mx", "crm"));
  });

  it("rankByFuzzy orders users by name/login", () => {
    const users = [
      { id: 1, name: "Jane Doe", login: "jane" },
      { id: 2, name: "John Smith", login: "john" },
    ];
    const ranked = rankByFuzzy("joh", users, (u) => [u.name, u.login]);
    expect(ranked[0]?.item.id).toBe(2);
  });
});
