package auth

import (
	"crypto/tls"
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
	if fileExists(ServerCrtFile) && fileExists(ServerKeyFile) && fileExists(CaCertFile) {
		sAuth.LoadAuthConfig(serverCaFile)
	}
}
