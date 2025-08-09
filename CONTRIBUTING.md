# Contributing to MetricHub

## Quick Start

1. Fork and clone the repo.
1. Install pre-commit hooks:

```bash
pip install pre-commit || python -m pip install pre-commit
pre-commit install --install-hooks
pre-commit install --hook-type commit-msg
```

1. Create a feature branch.
1. Run CI checks locally before pushing.

## Commit Messages

Use Conventional Commits (e.g. `feat(api): add deployment listing endpoint`). Commitizen hook will validate.

## Local Quality Gates

- `go fmt ./...` and `golangci-lint run` (backend)
- `npm run lint` & `npm run build` (frontend)
- `go test ./... -race` for backend reliability

## Pull Request Checklist

- [ ] Tests added/updated
- [ ] Docs updated if user-facing change
- [ ] No failing pre-commit hooks
- [ ] CI green

Thank you for helping build MetricHub!
