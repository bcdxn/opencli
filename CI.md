# CI & Deployment

## GitHub Pages — Manual Deployment

The web editor is a fully static SPA that can be hosted on GitHub Pages. Since `wasm_exec.js` and `opencli.wasm` are generated artifacts (not committed to git), the build must happen in CI before deploying.

### One-time Setup

1. In the GitHub repository, go to **Settings → Pages**.
2. Under **Source**, select **GitHub Actions**.
3. (No branch-based source is needed — the workflow deploys from the build artifact.)

### Build & Deploy Steps

The following describes what the GitHub Actions workflow must do. Trigger it manually via **Actions → Deploy Web Editor → Run workflow**, or automate it on push to `main`.

```yaml
# .github/workflows/deploy.yml
name: Deploy Web Editor

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  deploy:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pages: write
      id-token: write
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - uses: actions/setup-node@v4
        with:
          node-version: 22

      - name: Build WASM and UI
        run: make build-ui

      - uses: actions/configure-pages@v5

      - uses: actions/upload-pages-artifact@v3
        with:
          path: web/dist

      - id: deployment
        uses: actions/deploy-pages@v4
```

### Deployed URL

```
https://bcdxn.github.io/open-cli-spec/
```

The `base: './'` in `vite.config.ts` ensures all asset references use relative paths, and the WASM fetch in `web/src/wasm/client.ts` uses `import.meta.env.BASE_URL` to locate `opencli.wasm` correctly under the subdirectory URL.

### Local Build Verification

```sh
# Build WASM and the full UI bundle
make build-ui

# Check that web/dist/ contains opencli.wasm
ls web/dist/

# Preview the production build locally
cd web && npm run preview
```
