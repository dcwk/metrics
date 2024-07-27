package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func Encrypt(b []byte, publicKeyPath string) ([]byte, error) {
	publicKeyPEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		panic(err)
	}

	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), b)
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}

func Decrypt(b []byte, privateKeyPath string) ([]byte, error) {
	privateKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	decodedText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, b[:])
	if err != nil {
		return nil, err
	}

	return decodedText, nil
}
