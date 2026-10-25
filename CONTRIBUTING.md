# Contributing

## Ground rules
- **Synthetic data only.** Never commit or paste real child, parent or specialist data (PRV-08).
- **No legal-compliance claims.** Write "assumption / open question" (see `docs/kazakhstan-privacy-legal-checklist.md`).
- **MVP scope is fixed** (SF-01): no video calls, payments, school/clinic integrations or clinical advice.
- Sources of truth: `docs/openapi.yaml` (API), `docs/srs.md` (requirements; wins over `requirements.md`), `docs/rbac-and-threat-model.md` (access matrix). Conflicts are resolved by an ADR in `docs/decisions/`.

## Workflow
1. Create a ticket `SSP-XX`, then a branch `SSP-XX-short-name` from `main`.
2. Small commits: `SSP-XX: imperative summary`.
3. Run `make check` before every push. It runs lint + typecheck + test + build for backend and frontend and fails on any error — nothing is skipped.
4. Open a PR using the template. **Every PR needs at least one review from a code owner** (`.github/CODEOWNERS`); CI must be green.
5. Squash or rebase merge; delete the branch.

## Local setup
Requirements: Go 1.24+, Node 22 (see `frontend/.nvmrc`, minimum 20.19), npm.

```sh
make tools      # golangci-lint, gitleaks
make install    # npm ci
make check
```

### Pre-commit hooks
```sh
pipx install pre-commit   # or: pip install pre-commit
pre-commit install
pre-commit run --all-files
```
Hooks: gitleaks, gofmt, goimports, go vet, eslint, large-file guard (1 MB), whitespace fixes.

## Binary files
Presentations, documents and diagram exports live in `docs/assets/` and are tracked with Git LFS for new files (`.gitattributes`). Editable diagram sources (Mermaid in `docs/architecture.md`, `docs/diagrams/*.mmd`, `*.puml`) are the source of truth.
