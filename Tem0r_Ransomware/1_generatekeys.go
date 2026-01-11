package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func generateRSAKey() error {
	// Generates a 2048-bit RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Encodes the private key in PEM format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	}
	privPEM := pem.EncodeToMemory(&privBlock)

	// Write the private key to a file
	err = ioutil.WriteFile("private.key", privPEM, 0600)
	if err != nil {
		return err
	}

	// Encodes the public key in PEM format
	pubDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	pubBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubDER,
	}
	pubPEM := pem.EncodeToMemory(&pubBlock)

	// Write the public key to a file
	err = ioutil.WriteFile("public.key", pubPEM, 0600)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := generateRSAKey()
	if err != nil {
		fmt.Println("Error generating keys:", err)
	} else {
		fmt.Println("RSA keys generated successfully.")
	}
}
