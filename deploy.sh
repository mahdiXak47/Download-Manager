#!/bin/bash

# Deployment script for Download Manager
# Usage: ./deploy.sh [environment] [version]

set -e

ENVIRONMENT=${1:-staging}
VERSION=${2:-latest}
DEPLOY_DIR="/opt/download-manager"
SERVICE_NAME="download-manager"

echo "🚀 Deploying Download Manager to $ENVIRONMENT (version: $VERSION)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if running as root or with sudo
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}Please run as root or with sudo${NC}"
    exit 1
fi

# Create deployment directory
echo "📁 Creating deployment directory..."
mkdir -p $DEPLOY_DIR/{downloads,config,logs}
chmod 755 $DEPLOY_DIR

# Stop existing service
echo "🛑 Stopping existing service..."
if systemctl is-active --quiet $SERVICE_NAME; then
    systemctl stop $SERVICE_NAME
fi

# Backup current version
if [ -f "$DEPLOY_DIR/download-manager" ]; then
    echo "💾 Backing up current version..."
    cp $DEPLOY_DIR/download-manager $DEPLOY_DIR/download-manager.backup.$(date +%Y%m%d_%H%M%S)
fi

# Copy new binary (if deploying from local build)
if [ -f "./bin/download-manager" ]; then
    echo "📦 Copying new binary..."
    cp ./bin/download-manager $DEPLOY_DIR/download-manager
    chmod +x $DEPLOY_DIR/download-manager
fi

# Or pull from Docker if using containerized deployment
if command -v docker &> /dev/null; then
    echo "🐳 Using Docker deployment..."
    cd $DEPLOY_DIR
    docker-compose pull
    docker-compose build
    docker-compose up -d
    
    # Wait for health check
    echo "⏳ Waiting for service to be healthy..."
    sleep 10
    
    # Check if container is running
    if docker-compose ps | grep -q "Up"; then
        echo -e "${GREEN}✅ Deployment successful!${NC}"
        docker-compose ps
    else
        echo -e "${RED}❌ Deployment failed!${NC}"
        docker-compose logs --tail=50
        exit 1
    fi
else
    # Systemd service deployment
    echo "⚙️  Setting up systemd service..."
    
    # Create systemd service file
    cat > /etc/systemd/system/$SERVICE_NAME.service << EOF
[Unit]
Description=Download Manager
After=network.target

[Service]
Type=simple
User=download-manager
Group=download-manager
WorkingDirectory=$DEPLOY_DIR
ExecStart=$DEPLOY_DIR/download-manager
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

    # Create user if doesn't exist
    if ! id "download-manager" &>/dev/null; then
        useradd -r -s /bin/false download-manager
    fi
    
    # Set permissions
    chown -R download-manager:download-manager $DEPLOY_DIR
    
    # Reload systemd and start service
    systemctl daemon-reload
    systemctl enable $SERVICE_NAME
    systemctl start $SERVICE_NAME
    
    # Check status
    sleep 3
    if systemctl is-active --quiet $SERVICE_NAME; then
        echo -e "${GREEN}✅ Deployment successful!${NC}"
        systemctl status $SERVICE_NAME --no-pager
    else
        echo -e "${RED}❌ Deployment failed!${NC}"
        journalctl -u $SERVICE_NAME --no-pager -n 50
        exit 1
    fi
fi

echo -e "${GREEN}🎉 Deployment complete!${NC}"
echo "📊 Service status:"
if command -v docker &> /dev/null; then
    docker-compose ps
else
    systemctl status $SERVICE_NAME --no-pager
fi

