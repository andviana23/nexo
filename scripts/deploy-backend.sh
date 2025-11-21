#!/usr/bin/env bash

set -euo pipefail

# Required envs: SSH_HOST, SSH_USER
SSH_HOST=${SSH_HOST:?Informe SSH_HOST (ex: api.meudominio.com)}
SSH_USER=${SSH_USER:?Informe SSH_USER (ex: barber)}
SSH_KEY_PATH=${SSH_KEY_PATH:-$HOME/.ssh/id_rsa}

BACKEND_ARTIFACT=${BACKEND_ARTIFACT:-backend/bin/main}
BACKEND_REMOTE_DIR=${BACKEND_REMOTE_DIR:-/opt/barber-api}
BACKEND_SERVICE=${BACKEND_SERVICE:-barber-api}

if [ ! -f "$BACKEND_ARTIFACT" ]; then
  echo "❌ Artifact not found: $BACKEND_ARTIFACT"
  exit 1
fi

if [ ! -f "$SSH_KEY_PATH" ]; then
  echo "❌ SSH key not found: $SSH_KEY_PATH"
  exit 1
fi

echo "🚀 Deploying backend -> ${SSH_USER}@${SSH_HOST}:${BACKEND_REMOTE_DIR}"
scp -i "$SSH_KEY_PATH" "$BACKEND_ARTIFACT" "${SSH_USER}@${SSH_HOST}:/tmp/barber-api-main"

ssh -i "$SSH_KEY_PATH" "${SSH_USER}@${SSH_HOST}" <<EOF
set -euo pipefail
sudo mkdir -p ${BACKEND_REMOTE_DIR}/backups

if [ -f ${BACKEND_REMOTE_DIR}/main ]; then
  sudo cp ${BACKEND_REMOTE_DIR}/main ${BACKEND_REMOTE_DIR}/backups/main.\$(date +%Y%m%d%H%M%S)
fi

sudo mv /tmp/barber-api-main ${BACKEND_REMOTE_DIR}/main
sudo chown barber:barber ${BACKEND_REMOTE_DIR}/main
sudo chmod +x ${BACKEND_REMOTE_DIR}/main

echo "🔄 Restarting ${BACKEND_SERVICE}"
sudo systemctl restart ${BACKEND_SERVICE}
sleep 3
sudo systemctl status ${BACKEND_SERVICE} --no-pager

echo "✅ Healthcheck (localhost:8080/health)"
curl -f http://localhost:8080/health || true
EOF

echo "✅ Backend deploy finished."
