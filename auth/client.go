package auth

import (
	"crypto/tls"
	"log/slog"
)

// 客户端加载配置
type cAuth struct{}

func (c cAuth) LoadAuthConfig(caFile authFile) error {
	cert, cerPool, err := load(caFile)
	if err != nil {
		return err
	}
	ClientAuthConfig = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            cerPool,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
	}
	return nil
}

// 手动加载配置
func AddClientAuthConfig() {
	clientCaFile := authFile{CrtFile: ClientCrtFile, KeyFile: ClientKeyFile, CaCetFile: CaCertFile}
	cAuth := &cAuth{}

	if !fileExists(ClientCrtFile) || !fileExists(ClientKeyFile) || !fileExists(CaCertFile) {
		slog.Info("指定的证书文件不存在,请核对证书文件", slog.String("crtFile", ClientCrtFile), slog.String("keyFile", ClientKeyFile), slog.String("caFile", CaCertFile))
		return
	}
	cAuth.LoadAuthConfig(clientCaFile)
}
