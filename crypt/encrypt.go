package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"

	log "github.com/Sirupsen/logrus"
)

// EncryptToString encrypts data and encodes it to a
// standard base64 encoding, as defined in RFC 4648.
func EncryptToString(data []byte, label string, file string) (string, error) {
	ev, err := Encrypt(data, label, file)
	if err != nil {
		return string(ev), err
	}
	s := "ENC[" + base64.StdEncoding.EncodeToString(ev) + "]"
	return s, nil
}

// Encrypt will encrypt the data string using the PEM public
// key extracted from the file
func Encrypt(data []byte, label string, file string) ([]byte, error) {
	block, err := extractPemDataBlock(file)
	if err != nil {
		return nil, err
	}

	pubkey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.WithField("file", file).Error("could not parse public key")
		return nil, err
	}

	encryptedValue, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubkey.(*rsa.PublicKey), data, []byte(label))
	if err != nil {
		log.Error("could not encrypt data")
		return nil, err
	}

	return encryptedValue, nil
}
