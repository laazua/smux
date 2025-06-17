package auth

import "crypto/tls"

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

	if fileExists(ClientCrtFile) && fileExists(ClientKeyFile) && fileExists(CaCertFile) {
		cAuth.LoadAuthConfig(clientCaFile)
	}
}
