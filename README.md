### 1. Estructura de carpetas recomendada


```text
Tem0r_Ransomware/

 README.md                 
 go.mod                    <-- (Archivo de configuración de Go)
 go.sum                    <-- (Archivo de dependencias)

 0_createdummy.go          <-- Generador de entorno
 1_generatekeys.go         <-- Generador de claves
 2_encrypt.go              <-- Cifrador local (Parte 1)
 3_decrypt.go              <-- Descifrador (Parte 1)
 server.go                 <-- Servidor del atacante (Parte 2)
 4_ransomware_client.go    <-- Cliente exfiltrador (Parte 2)

 informe/
     Informe_Forense.pdf   

```

---

### 2. Archivo `README.md`

Crea un archivo llamado `README.md` y pega este contenido. Explica qué es cada archivo para que el profesor lo vea claro.

```markdown
# Práctica de Análisis de Ransomware (PoC "Tem0r")

Este repositorio contiene el código fuente y la documentación de la práctica **"Análisis de un Ransomware"**. El objetivo es puramente académico y educativo, simulando el comportamiento de un ransomware en un entorno controlado Linux para comprender sus mecanismos de ataque y defensa.

##  Disclaimer / Aviso Legal
**USO ESTRICTAMENTE EDUCACIONAL.** Este software es una Prueba de Concepto (PoC) diseñada para entornos de laboratorio aislados. No debe utilizarse en sistemas de producción ni con fines malintencionados.

##  Contenido del Proyecto

El proyecto se divide en dos fases y consta de los siguientes scripts desarrollados en **Go**:

### Fase 1: Cifrado y Descifrado Local
* **`0_createdummy.go`**: Script de preparación. Genera 200 archivos aleatorios en `/tmp/dummy` para simular los datos de la víctima.
* **`1_generatekeys.go`**: Genera el par de claves RSA-2048 (`public.key` y `private.key`) necesarias para el cifrado híbrido.
* **`2_encrypt.go`**: Realiza el ataque local. Cifra los archivos con AES, protege la clave AES con RSA y elimina los originales, dejando una nota de rescate.
* **`3_decrypt.go`**: Herramienta de recuperación. Utiliza la clave privada para revertir el proceso y restaurar los archivos.

### Fase 2: Doble Extorsión (Exfiltración)
* **`server.go`**: Simula el servidor C2 (Command & Control) del atacante. Escucha en el puerto 8080 y recibe las claves y archivos robados.
* **`4_ransomware_client.go`**: Versión avanzada del malware. Antes de cifrar, conecta vía WebSocket con el servidor y exfiltra la clave privada y los datos de la víctima.

##  Requisitos y Ejecución

* **Sistema Operativo:** Linux (Ubuntu recomendado).
* **Lenguaje:** Go (Golang).

### Instalación de dependencias
```bash
go mod init temor_project
go get [github.com/gorilla/websocket](https://github.com/gorilla/websocket)
go get [github.com/golang-jwt/jwt](https://github.com/golang-jwt/jwt)

```

### Ejemplo de uso (Fase de Exfiltración)

1. Iniciar el servidor atacante:
```bash
go run server.go

```


2. Ejecutar el cliente en la máquina víctima:
```bash
go run 4_ransomware_client.go

```



##  Informe

El análisis forense completo con las evidencias de la ejecución se encuentra en la carpeta `/informe`.

---

**Alumno:** Ruben Ferrer Marquez
**Asignatura:** Puesta en Producción Segura

```

---






