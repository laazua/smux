package auth

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
)

var (
	ServerAuthConfig *tls.Config
	ClientAuthConfig *tls.Config

	CaCertFile    = ""
	ServerCrtFile = ""
	ServerKeyFile = ""
	ClientCrtFile = ""
	ClientKeyFile = ""
)

type authFile struct {
	CrtFile   string
	KeyFile   string
	CaCetFile string
}

func load(authFile authFile) (tls.Certificate, *x509.CertPool, error) {
	cert, err := tls.LoadX509KeyPair(authFile.CrtFile, authFile.KeyFile)
	if err != nil {
		log.Fatal(err)
		return tls.Certificate{}, nil, err
	}

	caCert, err := os.ReadFile(authFile.CaCetFile)
	if err != nil {
		log.Fatal(err)
		return tls.Certificate{}, nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)
	return cert, certPool, nil
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		// 其他错误（如权限问题）
		return false
	}
}
