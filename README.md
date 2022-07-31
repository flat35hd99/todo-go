# Sample Todo application

with

- Golang
  - GORM
  - Echo
- React
  - TypeScript
  - Material UI

## Developer guide

Open this repository using VScode or GitHub codespace because this repository has `.devcontainer` to define common develop environment.

### For backend

This project use:

- testing
- assert
- apitest

to test.

Thus, you will develop following process

1. Write test first
2. Write features
3. Run test by `go test`
4. Before commit, run `go fmt`
5. Commit

### For frontend

Start development server by `yarn dev`.

prettier and eslint are used and run `yarn lint` and `yarn fmt` before commit. You cand build them by `yarn build`.
