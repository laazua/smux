#!/bin/bash

set -e

# 设置工作目录
CERT_DIR="./ssl"
mkdir -p "$CERT_DIR"
cd "$CERT_DIR"

echo "[*] 清理旧文件"
rm -f *.crt *.key *.csr *.srl openssl.cnf

echo "[*] 生成 CA 私钥和证书"
openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 -out ca.crt -subj "/C=CN/ST=SC/L=CD/O=laazua/CN=CA"

echo "[*] 创建 OpenSSL 配置文件带 SAN 支持"
cat > openssl.cnf <<EOF
[ req ]
default_bits       = 2048
prompt             = no
default_md         = sha256
req_extensions     = req_ext
distinguished_name = dn

[ dn ]
C  = CN
ST = Sichuan
L  = Chengdu
O  = laazua
OU = laazua

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = localhost
IP.1  = 127.0.0.1
EOF

### ==== 服务端证书 ====
echo "[*] 生成 Server 私钥和 CSR"
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -subj "/CN=localhost" -config openssl.cnf

echo "[*] 签发 Server 证书（含 SAN）"
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
-out server.crt -days 365 -sha256 -extfile openssl.cnf -extensions req_ext

### ==== 客户端证书 ====
echo "[*] 生成 Client 私钥和 CSR"
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -subj "/CN=client" -config openssl.cnf

echo "[*] 签发 Client 证书（含 SAN）"
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
-out client.crt -days 365 -sha256 -extfile openssl.cnf -extensions req_ext

echo "[✔] 所有证书生成完毕，位于 $CERT_DIR 目录下"
