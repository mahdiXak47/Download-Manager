# CI/CD Documentation

This document describes the CI/CD pipeline setup for the Download Manager application.

## Overview

The CI/CD pipeline automates:
- Code quality checks (linting, formatting)
- Building for multiple platforms
- Testing across different environments
- Docker image creation and publishing
- Automated deployments
- Release management

## GitHub Actions Workflows

### 1. Main CI/CD Pipeline (`.github/workflows/ci-cd.yml`)

**Triggers:**
- Push to `main`, `master`, or `develop` branches
- Pull requests to `main` or `master`
- Release creation

**Jobs:**

#### Lint Job
- Runs `gofmt` to check code formatting
- Runs `go vet` for static analysis
- Runs `golangci-lint` for comprehensive linting

#### Build Job
Builds the application for multiple platforms:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

Artifacts are uploaded and retained for 30 days.

#### Release Job
- Creates GitHub releases when a tag is created
- Packages all platform binaries
- Generates checksums for verification
- Auto-generates release notes

#### Docker Job
- Builds multi-platform Docker images (linux/amd64, linux/arm64)
- Pushes to Docker Hub (on main/master branches)
- Uses build cache for faster builds

### 2. Test Suite (`.github/workflows/test.yml`)

**Triggers:**
- Push to main branches
- Pull requests

**Features:**
- Tests across Ubuntu, macOS, and Windows
- Tests with Go 1.21 and 1.22
- Code coverage reporting via Codecov
- Caching for faster builds

### 3. Deployment Workflow (`.github/workflows/deploy.yml`)

**Triggers:**
- Manual workflow dispatch
- Push to main/master branches

**Features:**
- Deploys to staging or production environments
- SSH-based deployment
- Health checks after deployment
- Slack notifications (optional)

## Local Development

### Using Make

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Format code
make fmt

# Run linters
make lint

# Run all checks
make check

# Clean build artifacts
make clean
```

### Using Docker

```bash
# Build image
make docker-build

# Run container
make docker-run

# View logs
make docker-logs

# Stop container
make docker-stop
```

## Deployment

### Docker Compose Deployment

1. **Configure environment:**
   ```bash
   # Edit docker-compose.yml if needed
   # Set environment variables
   export VERSION=1.0.0
   export COMMIT=$(git rev-parse --short HEAD)
   ```

2. **Deploy:**
   ```bash
   docker-compose up -d
   ```

3. **Verify:**
   ```bash
   docker-compose ps
   docker-compose logs -f download-manager
   ```

### Manual Server Deployment

1. **Prepare server:**
   ```bash
   # On server
   sudo mkdir -p /opt/download-manager/{downloads,config,logs}
   ```

2. **Deploy:**
   ```bash
   # From local machine or CI/CD
   ./deploy.sh staging
   # or
   ./deploy.sh production latest
   ```

3. **Verify:**
   ```bash
   # Check service status
   systemctl status download-manager
   # or
   docker-compose ps
   ```

## Configuration

### GitHub Secrets

Configure these in: Settings → Secrets and variables → Actions

| Secret | Description | Required |
|--------|-------------|----------|
| `DOCKER_USERNAME` | Docker Hub username | Yes (for Docker builds) |
| `DOCKER_PASSWORD` | Docker Hub password/token | Yes (for Docker builds) |
| `DEPLOY_HOST` | Deployment server hostname | Yes (for deployments) |
| `DEPLOY_USER` | SSH username | Yes (for deployments) |
| `DEPLOY_SSH_KEY` | SSH private key | Yes (for deployments) |
| `DEPLOY_PORT` | SSH port (default: 22) | No |
| `REPO_URL` | Git repository URL | Yes (for deployments) |
| `SLACK_WEBHOOK_URL` | Slack webhook URL | No (for notifications) |

### Environment Variables

For Docker deployment:

```bash
DOWNLOAD_PATH=/app/downloads
CONFIG_PATH=/home/downloader/.config/download-manager
LOG_PATH=/app/logs
```

## Release Process

1. **Create a release tag:**
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

2. **GitHub Actions will:**
   - Build for all platforms
   - Create a GitHub release
   - Upload binaries and checksums
   - Generate release notes

3. **Download releases:**
   - Go to GitHub Releases page
   - Download the binary for your platform
   - Verify checksums

## Monitoring

### Health Checks

The Docker container includes a health check:
```bash
docker inspect --format='{{.State.Health.Status}}' download-manager
```

### Logs

**Docker:**
```bash
docker-compose logs -f download-manager
```

**Systemd:**
```bash
journalctl -u download-manager -f
```

## Troubleshooting

### Build Failures

1. Check Go version compatibility
2. Verify all dependencies are available
3. Check for linting errors

### Deployment Failures

1. Verify SSH access to deployment server
2. Check server disk space
3. Verify Docker is installed (if using Docker)
4. Check service logs

### Docker Issues

1. Verify Docker credentials
2. Check Docker Hub rate limits
3. Ensure multi-platform build support is enabled

## Best Practices

1. **Always test locally before pushing:**
   ```bash
   make check
   ```

2. **Use semantic versioning for releases:**
   - Major.Minor.Patch (e.g., 1.0.0)

3. **Keep secrets secure:**
   - Never commit secrets to repository
   - Rotate secrets regularly
   - Use least privilege principle

4. **Monitor deployments:**
   - Set up alerts for failed deployments
   - Monitor application health
   - Review logs regularly

5. **Version control:**
   - Tag releases properly
   - Write clear commit messages
   - Use feature branches for development

## Support

For issues or questions about the CI/CD pipeline:
- Open an issue on GitHub
- Check workflow logs in Actions tab
- Review this documentation

