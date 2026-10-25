# Binary assets

Presentations, Word documents, PDFs and rendered diagram exports. These are **not** sources of truth.

- Diagram sources: Mermaid blocks in `docs/architecture.md`, `docs/diagrams/dfd.mmd`, `docs/diagrams/usecase.puml`.
- `diagrams/` holds historic exports; CI (`.github/workflows/diagrams.yml`) renders fresh SVGs from `*.mmd` as a build artifact.
- New `.pptx/.docx/.pdf` files are stored via Git LFS (`.gitattributes`). Files committed before SSP-50 stay in normal history (no history rewrite by decision).
