---
name: write-lili-learning-notes
description: Create, expand, audit, and update Chinese Markdown learning articles for the Lili software product engineer roadmap. Use for beginner, junior, or intermediate notes that must explain concepts deeply, verify current facts with primary web sources, include practical examples and useful visuals, remain maintainable in Obsidian/GitHub/VS Code, and avoid editorial commentary or unsupported claims.
---

# Write Lili Learning Notes

Produce a self-contained learning tutorial, not a summary card or research log.

## Required reading

Read [references/article-standard.md](references/article-standard.md) completely before writing or auditing any note. Read [references/visuals.md](references/visuals.md) when the topic has a process, hierarchy, state change, architecture, rendered interface, browser output, chart, or other relationship that benefits from a visual.

## Workflow

1. Read the target roadmap item, existing article, adjacent notes, and index.
2. Define the article's exact knowledge boundary. Keep one coherent knowledge point per article; split only when the roadmap item contains independently learnable subjects.
3. Browse current primary sources. Prefer living standards, official documentation, specifications, original papers, authoritative public data, and source repositories. Use secondary material only to discover questions or improve organization.
4. Reduce every claim to its prerequisites and mechanism. Check definitions, defaults, units, inheritance, state transitions, failure conditions, compatibility, deprecation, and version-specific behavior.
5. Write or expand the article using the applicable template in `article-standard.md`.
   When expanding an existing short note, restructure it into one continuous teaching path. Remove or merge the old summary sections; do not append a second explanation and case after the old template.
6. Add examples that can be executed, inspected, calculated, or followed. Never invent output.
7. Add the smallest useful visual. For rendered UI or browser behavior, create a runnable example, render it in a real browser, inspect the DOM/state/errors, and save the resulting image in the repository.
8. Validate links, code fences, local assets, examples, article structure, and repository formatting. Run `scripts/validate_note.py` on every changed note.
9. Run `scripts/validate_code_blocks.py` with an explicit Node.js path when the notes contain JSON, JavaScript, or shell examples; compile TypeScript and other languages with their real toolchains when available.
10. When a directory or stage is complete, run `scripts/audit_corpus.py` on the whole set to detect short notes and repeated long prose across articles.
11. Update the direction index and global counts only when adding, deleting, or splitting articles.

## Hard rules

- Explain what the concept is, why it exists, how it works, how to use it, special behavior, failure modes, boundaries, verification, and related knowledge.
- Keep one explanation path per concept. Repeated “brief version + expanded version”, duplicate cases, and appended template sections are incomplete even when the file passes length checks.
- Expand named attributes, parameters, methods, states, metrics, rules, or variants individually when readers need them to use the concept.
- State the actual rule directly. Do not include research-process narration, source-selection commentary, writing references, or remarks meant only for maintainers.
- Do not add analogies, motivational filler, scene-setting prose, fake quotations, or unsupported absolutes.
- Do not treat a framework, heuristic, convention, or vendor behavior as a universal law.
- Separate standard requirements, implementation behavior, compatibility status, team conventions, and recommendations.
- Put deterministic safety, authorization, validation, accounting, and business invariants in code or controlled systems, not in prompts or visual conventions.
- Prefer semantic HTML, accessible interaction, secure defaults, explicit errors, and production-safe commands.
- Never use screenshots as the only explanation. Pair each image with text, source code, or inspectable state.
- Use ordinary relative Markdown links for repository navigation. Do not use Obsidian-only `[[wiki links]]`; standard Markdown works in GitHub, Obsidian, and VS Code.
- Keep citations in `## 来源`. Do not insert citations as filler after every sentence.

## Level control

- Beginner and junior notes must build vocabulary from zero, enumerate the practical surface, show small complete examples, explain common mistakes, and include exercises with completion criteria.
- Intermediate notes must assume the prerequisite notes, deepen mechanisms, include at least two substantial applications, compare alternatives, show debugging or validation, cover operational and production boundaries, and end with an integration exercise or project acceptance criteria.

## Completion gate

Do not mark an article complete until:

- Every heading and example has substantive content.
- Beginner/junior notes normally reach at least 120 lines and 8 KB of real content; intermediate notes normally reach at least 200 lines and 12 KB. Falling below either floor requires expansion, not template padding.
- A worked case includes concrete input/evidence, step-by-step processing, output, verification, and at least one failure branch.
- Every factual statement that can change has been checked against a current primary source.
- Every important term introduced by the roadmap item is defined and operationalized.
- Code and commands are syntactically valid and safe in their stated context.
- Complete HTML examples pass an HTML5 conformance checker; other executable examples pass the relevant parser, compiler, test, or dry-run check when available.
- Visuals render in both GitHub-compatible Markdown and Obsidian, or have a repository image fallback.
- The article contains no editorial/process commentary forbidden by the standard.
- The validator and `git diff --check` pass.
