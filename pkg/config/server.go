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
	Port int `json:"port"`
}

const configPath = "./config/server.json"
const CertPath = "./config/ssl/server.crt"
const KeyPath = "./config/ssl/server.key"
const DefaultPort = 443


func createConfig() error {
	
	conf :=  ServerConfig{Port: DefaultPort}
	confJson, err := json.MarshalIndent(conf,"","    ")
	if err != nil {
		return err
	}
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	_, err1 := file.Write(confJson)
	if err1 != nil {
		return err
	}
	return nil
}


func LoadConfig() (*ServerConfig, error) {
	var conf ServerConfig 
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println(util.PrintStatus + "No server config file found, creating default...")
		err1 := createConfig()
		if err1 != nil {
			return &conf, err1
		}
	}
	jsonParser := json.NewDecoder(file)
	if err := jsonParser.Decode(&conf); err != nil {
		return &conf,err
	}
	return nil,nil
}

func CheckCert() error {
	_, err := os.Stat(KeyPath)
	_, err1 := os.Stat(CertPath)

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

	keyFile, err := os.Create(KeyPath)
	if err != nil {
		return err
	}
	defer keyFile.Close()
	pem.Encode(keyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	certFile, err := os.Create(CertPath)
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

