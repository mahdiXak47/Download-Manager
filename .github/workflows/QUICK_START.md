# Quick Start - Local Runner CI/CD

## What You Need

✅ Go 1.21+ installed  
✅ Git installed  
✅ GitHub repository  
❌ **Docker NOT required**

## 5-Minute Setup

### 1. Install Runner (One-time setup)

```bash
# Create folder
mkdir ~/actions-runner && cd ~/actions-runner

# Download (Linux example - adjust for your OS)
curl -o runner.tar.gz -L https://github.com/actions/runner/releases/download/v2.311.0/actions-runner-linux-x64-2.311.0.tar.gz
tar xzf runner.tar.gz

# Configure (get token from GitHub: Settings → Actions → Runners → New self-hosted runner)
./config.sh --url https://github.com/YOUR_USERNAME/YOUR_REPO --token YOUR_TOKEN

# Install as service
sudo ./svc.sh install
sudo ./svc.sh start
```

### 2. Verify Runner is Online

- Go to: GitHub → Your Repo → Settings → Actions → Runners
- You should see your runner with a green dot ✅

### 3. Push Code

```bash
git add .
git commit -m "Add CI/CD"
git push
```

### 4. Check Workflow

- Go to: GitHub → Your Repo → Actions tab
- You should see workflows running automatically!

## Workflows

### Automatic (on push/PR)
- **CI/CD Pipeline**: Lint → Build → Test → Run

### Manual (trigger from Actions tab)
- **Build and Deploy**: Build → Deploy to local system

## Output

After build, binary is in:
```
bin/download-manager
```

## Troubleshooting

**Runner offline?**
```bash
cd ~/actions-runner
sudo ./svc.sh restart
```

**Go not found?**
```bash
go version  # Should show 1.21+
```

**Workflow not running?**
- Check runner is online (green dot)
- Check Actions tab for errors
- Verify Go is installed

## That's It!

No Docker, no cloud costs, just your local machine building and running your Go app! 🚀

