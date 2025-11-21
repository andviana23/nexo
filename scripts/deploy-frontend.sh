#!/usr/bin/env bash

set -euo pipefail

# Required envs: SSH_HOST, SSH_USER
SSH_HOST=${SSH_HOST:?Informe SSH_HOST (ex: app.meudominio.com)}
SSH_USER=${SSH_USER:?Informe SSH_USER (ex: barber)}
SSH_KEY_PATH=${SSH_KEY_PATH:-$HOME/.ssh/id_rsa}

FRONT_BUILD_DIR=${FRONT_BUILD_DIR:-frontend/.next/standalone}
FRONT_STATIC_DIR=${FRONT_STATIC_DIR:-frontend/.next/static}
FRONT_PUBLIC_DIR=${FRONT_PUBLIC_DIR:-frontend/public}
FRONT_REMOTE_DIR=${FRONT_REMOTE_DIR:-/opt/barber-frontend}
FRONT_SERVICE=${FRONT_SERVICE:-barber-frontend}
PKG_MGR=${PKG_MGR:-pnpm}
INSTALL_CMD=${INSTALL_CMD:-"$PKG_MGR install --frozen-lockfile"}
BUILD_CMD=${BUILD_CMD:-"$PKG_MGR build"}

# Se SKIP_BUILD != 1 o build é executado aqui
if [ "${SKIP_BUILD:-0}" != "1" ]; then
  echo "🏗️  Building frontend (Next.js standalone)"
  (cd frontend && eval "$INSTALL_CMD" && eval "$BUILD_CMD")
fi

if [ ! -d "$FRONT_BUILD_DIR" ]; then
  echo "❌ Build não encontrado em $FRONT_BUILD_DIR - rode npm run build primeiro."
  exit 1
fi

if [ ! -f "$SSH_KEY_PATH" ]; then
  echo "❌ SSH key not found: $SSH_KEY_PATH"
  exit 1
fi

TMP_PACKAGE=$(mktemp -d)
echo "📦 Empacotando artefatos em $TMP_PACKAGE"

cp -R "${FRONT_BUILD_DIR}/." "$TMP_PACKAGE/"
mkdir -p "$TMP_PACKAGE/.next/static"
cp -R "${FRONT_STATIC_DIR}/." "$TMP_PACKAGE/.next/static/"
cp -R "${FRONT_PUBLIC_DIR}" "$TMP_PACKAGE/public"

tar -czf /tmp/barber-frontend.tar.gz -C "$TMP_PACKAGE" .
rm -rf "$TMP_PACKAGE"

echo "🚀 Enviando frontend -> ${SSH_USER}@${SSH_HOST}:${FRONT_REMOTE_DIR}"
scp -i "$SSH_KEY_PATH" /tmp/barber-frontend.tar.gz "${SSH_USER}@${SSH_HOST}:/tmp/barber-frontend.tar.gz"

ssh -i "$SSH_KEY_PATH" "${SSH_USER}@${SSH_HOST}" <<EOF
set -euo pipefail
sudo mkdir -p ${FRONT_REMOTE_DIR}
sudo tar -xzf /tmp/barber-frontend.tar.gz -C ${FRONT_REMOTE_DIR}
sudo rm /tmp/barber-frontend.tar.gz
sudo chown -R barber:barber ${FRONT_REMOTE_DIR}

echo "🔄 Restarting ${FRONT_SERVICE}"
sudo systemctl restart ${FRONT_SERVICE}
sleep 3
sudo systemctl status ${FRONT_SERVICE} --no-pager
EOF

echo "✅ Frontend deploy finished."
