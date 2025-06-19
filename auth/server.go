package auth

import (
	"crypto/tls"
	"log/slog"
)

// 服务端加载配置
type sAuth struct{}

func (s sAuth) LoadAuthConfig(caFile authFile) error {
	cert, cerPool, err := load(caFile)
	if err != nil {
		return err
	}
	ServerAuthConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    cerPool,
	}
	return nil
}

// 手动加载配置
func AddServerAuthConfig() {
	serverCaFile := authFile{CrtFile: ServerCrtFile, KeyFile: ServerKeyFile, CaCetFile: CaCertFile}
	sAuth := &sAuth{}
	if !fileExists(ServerCrtFile) || !fileExists(ServerKeyFile) || !fileExists(CaCertFile) {
		slog.Info("指定的证书文件不存在,请核对证书文件", slog.String("crtFile", ServerCrtFile), slog.String("keyFile", ServerKeyFile), slog.String("caFile", CaCertFile))
		return
	}
	sAuth.LoadAuthConfig(serverCaFile)
}
