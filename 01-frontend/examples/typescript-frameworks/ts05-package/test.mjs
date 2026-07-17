import assert from "node:assert/strict";
import { formatCents } from "./src/format.ts";

assert.equal(formatCents(1990, "CNY"), "¥19.90");
assert.throws(() => formatCents(19.9, "CNY"), TypeError);
console.log("ts05-package: ok");
