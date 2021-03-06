package xrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// This method requires a random number of bits.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate key pair: %v", err)
	}

	// The public key is part of the PrivateKey struct
	return privateKey, &privateKey.PublicKey, nil
}

// Export public key as a string in PEM format
func ExportPubKeyAsPEMStr(pubkey *rsa.PublicKey) string {
	pubKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pubkey),
		},
	))
	return pubKeyPem
}

// Export private key as a string in PEM format
func ExportPrivKeyAsPEMStr(privkey *rsa.PrivateKey) string {
	privKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privkey),
		},
	))
	return privKeyPem

}

// Save string to a file
func SaveKeyToFile(keyPem, filename string) {
	pemBytes := []byte(keyPem)
	ioutil.WriteFile(filename, pemBytes, 0400)
}

// Decode private key struct from PEM string
func ExportPEMStrToPrivKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	return key
}

// Decode public key struct from PEM string
func ExportPEMStrToPubKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	key, _ := x509.ParsePKCS1PublicKey(block.Bytes)
	return key
}
