# GitHub Actions Workflows for Local Runner

This directory contains GitHub Actions workflows optimized for **self-hosted runners**.

## Workflows

### ci-cd.yml
Main CI/CD pipeline that runs on your local runner:
- **Lint**: Code formatting and static analysis checks
- **Build**: Compiles the Go application
- **Test**: Runs tests (if any)
- **Run**: Smoke test to verify the application starts

**Runs on:** `self-hosted` runner

### build-and-deploy.yml
Build and deployment workflow:
- **Build**: Compiles the application
- **Deploy**: Deploys the binary to your local system

**Usage:**
- Manual trigger via GitHub Actions UI
- Option to build-only (skip deployment)

## Setting Up Self-Hosted Runner

1. **Install GitHub Actions Runner:**
   ```bash
   # Create a folder
   mkdir actions-runner && cd actions-runner
   
   # Download the latest runner package
   curl -o actions-runner-linux-x64-2.311.0.tar.gz -L https://github.com/actions/runner/releases/download/v2.311.0/actions-runner-linux-x64-2.311.0.tar.gz
   
   # Extract the installer
   tar xzf ./actions-runner-linux-x64-2.311.0.tar.gz
   ```

2. **Configure the runner:**
   ```bash
   ./config.sh --url https://github.com/YOUR_USERNAME/YOUR_REPO --token YOUR_TOKEN
   ```

3. **Install as a service (optional):**
   ```bash
   sudo ./svc.sh install
   sudo ./svc.sh start
   ```

4. **For macOS:**
   - Download: `actions-runner-osx-x64-2.311.0.tar.gz`
   - Follow similar steps

## Requirements

- **Go 1.21+** installed on your runner machine
- **Git** installed
- **Network access** to GitHub

## Workflow Execution

### Automatic Triggers
- Push to `main`, `master`, or `develop` branches
- Pull requests to `main` or `master`

### Manual Triggers
- Go to Actions tab → Select workflow → Run workflow

## Notes

- Docker is **not required** - workflows build and run Go directly
- All builds happen on your local machine
- Binaries are saved in the `bin/` directory
- Artifacts are uploaded to GitHub (optional)

## Troubleshooting

**Runner not picking up jobs:**
- Check runner status: `./run.sh` or service status
- Verify runner is online in GitHub: Settings → Actions → Runners

**Go not found:**
- Install Go 1.21 or higher
- Verify: `go version`
- Add Go to PATH if needed

**Build failures:**
- Check Go version compatibility
- Verify all dependencies: `go mod download`
- Check for linting errors
