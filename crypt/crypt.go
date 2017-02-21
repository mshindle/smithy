package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

// Generate creates a public / private key pair and saves them in the specified file
func Generate(pubFile string, privateFile string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}
	privatePem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	ioutil.WriteFile(privateFile, privatePem, 0600)
	log.WithField("file", privateFile).Info("private key written")

	pubKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	publicPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKey,
	})
	ioutil.WriteFile(pubFile, publicPem, 0644)
	log.WithField("file", pubFile).Info("public key written")
	return nil
}

func extractPemDataBlock(file string) (*pem.Block, error) {
	fileContents, err := ioutil.ReadFile(file)
	if err != nil {
		log.WithField("file", file).Error("could not read file")
		return nil, err
	}
	block, _ := pem.Decode(fileContents)
	if block == nil {
		log.WithField("file", file).Error("could not decode file contents - not PEM encoded")
		return nil, errors.New("data not PEM encoded")
	}
	return block, nil
}
