package main

import (
	cryptorand "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// Inicializar la semilla aleatoria
	rand.Seed(time.Now().UnixNano())

	// Crear el directorio si no existe
	err := os.MkdirAll("/tmp/dummy", 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// Definir extensiones para archivos de texto y binarios
	textExtensions := []string{".txt", ".log", ".csv"}
	binaryExtensions := []string{".bin", ".dat", ".jpg"}

	// Generar 200 archivos
	for i := 0; i < 200; i++ {
		// Decidir el tipo de archivo (50% probabilidad)
		isText := rand.Intn(2) == 0

		var extension string
		if isText {
			extension = textExtensions[rand.Intn(len(textExtensions))]
		} else {
			extension = binaryExtensions[rand.Intn(len(binaryExtensions))]
		}

		// Crear el nombre del archivo
		filename := filepath.Join("/tmp/dummy", fmt.Sprintf("dummy_%d%s", i, extension))

		// Crear el archivo
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Error creating file:", err)
			continue
		}

		// Escribir contenido basado en el tipo
		if isText {
			content := generateRandomText()
			file.WriteString(content)
		} else {
			// Generar contenido binario aleatorio (tamaño entre 100KB y 1MB)
			size := rand.Intn(900000) + 100000
			io.CopyN(file, cryptorand.Reader, int64(size))
		}

		file.Close()
	}
	fmt.Println("200 dummy files created in /tmp/dummy")
}

// generateRandomText crea una frase aleatoria para archivos de texto
func generateRandomText() string {
	words := []string{"example", "random", "data", "test", "file", "content", "dummy", "text"}
	sentenceLength := rand.Intn(20) + 5
	var sentence []string
	for i := 0; i < sentenceLength; i++ {
		sentence = append(sentence, words[rand.Intn(len(words))])
	}
	return strings.Join(sentence, " ") + "\n"
}
