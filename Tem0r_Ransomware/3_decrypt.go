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
	privateKeyFile := "private.key"
	encryptedDir := "/tmp/dummy"

	// 1. Cargar la clave privada
	privateKeyBytes, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		fmt.Println("Error reading private key file:", err)
		return
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		fmt.Println("Failed to decode PEM block containing the key")
		return
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing private key:", err)
		return
	}

	// 2. Recorrer los archivos cifrados
	err = filepath.Walk(encryptedDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Si es un archivo y tiene extensión .crypted
		if !info.IsDir() && filepath.Ext(path) == ".crypted" {
			fmt.Println("Decrypting file:", path)

			// Leer todo el contenido cifrado
			encryptedData, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading encrypted file %s: %v", path, err)
			}

			// Separar la clave AES cifrada de los datos del archivo
			// La clave RSA de 2048 bits cifra bloques de 256 bytes (Size())
			keySize := privateKey.Size()
			if len(encryptedData) < keySize {
				fmt.Println("Skipping file (too small):", path)
				return nil
			}

			encryptedKey := encryptedData[:keySize]
			fileData := encryptedData[keySize:]

			// 3. Descifrar la clave AES usando RSA privada
			aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedKey, nil)
			if err != nil {
				return fmt.Errorf("error decrypting AES key: %v", err)
			}

			// 4. Descifrar los datos del archivo usando AES
			decryptedData, err := decryptAES(fileData, aesKey)
			if err != nil {
				return fmt.Errorf("error decrypting file data: %v", err)
			}

			// 5. Guardar el archivo restaurado (quitando el .crypted del nombre)
			newPath := path[:len(path)-len(".crypted")]
			err = ioutil.WriteFile(newPath, decryptedData, 0644)
			if err != nil {
				return fmt.Errorf("error writing decrypted file %s: %v", newPath, err)
			}
			fmt.Println("File decrypted and saved as:", newPath)

			// 6. Borrar el archivo cifrado
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("error removing encrypted file %s: %v", path, err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error during decryption process:", err)
		return
	}

	fmt.Println("All files decrypted successfully.")
}

// Función auxiliar para descifrar AES (usada en el paso 4)
func decryptAES(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
