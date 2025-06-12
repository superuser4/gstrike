package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"gstrike/pkg/util"
	"math/big"
	"os"
	"time"
)

type ServerConfig struct {
	Port uint32	`json:"port"`
}

const configPath = "./config/server.json"
const certPath = "./config/ssl/server.crt"
const keyPath = "./config/ssl/server.key"

func LoadConfig() (ServerConfig, error) {
	var conf ServerConfig
	file, err := os.Open(configPath)
	if err != nil {
		return conf,err
	}
	jsonParser := json.NewDecoder(file)
	if err := jsonParser.Decode(&conf); err != nil {
		return conf,err
	}
	return conf,nil
}

func CheckCert() error {
	_, err := os.Stat(keyPath)
	_, err1 := os.Stat(certPath)

	if err != nil || err1 != nil {
		err2 := CreateCerts()
		if err2 != nil {
			return err
		}
	}
	return nil 
}

func CreateCerts() error {
	fmt.Println(util.PrintStatus + "No SSL certificates found in path, generating default cert and key...")	

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1), 
		Subject: pkix.Name{
			Organization: []string{"GStrike"},
			CommonName:   "localhost",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour), 
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth, 
		},
		DNSNames: []string{"localhost"}, 
	}

	certDER, err := x509.CreateCertificate(
		rand.Reader, &template, &template, &privateKey.PublicKey, privateKey,
	)
	if err != nil {
		return nil
	}

	keyFile, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer keyFile.Close()
	pem.Encode(keyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	certFile, err := os.Create(certPath)
	if err != nil {
		return err
	}
	defer certFile.Close()
	pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})
	return nil
}

