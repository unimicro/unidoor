package main

import (
	"crypto/tls"
	"github.com/stianeikeland/go-rpio"
	"golang.org/x/crypto/acme/autocert"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var doorRemote = rpio.Pin(2)

var tokens = make(map[string]string)

const (
	DOOR_LOG_PATH     = "access.log"
	CERTIFICATE_CACHE = "certs"
	DOMAIN            = "unidoor.space"
	TOKENS_FILE_PATH  = "tokens"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatal("You need to supply port as an argument")
		os.Exit(1)
	}
	port := os.Args[1]

	openDoorRemoteGPIO()
	defer closeDoorRemoteGPIO()

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(DOMAIN),
		Cache:      autocert.DirCache(CERTIFICATE_CACHE),
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/token/", tokenHandler)
	http.HandleFunc("/token", tokenHandler)

	server := &http.Server{
		Addr: ":" + port,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	log.Print("Starting server on port ", port)
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
		os.Exit(3)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-type", "text/html")
		indexFile, err := ioutil.ReadFile("index.html")
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		} else {
			w.Write(indexFile)
		}
	default:
		w.WriteHeader(404)
	}
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	tokens := parseTokenFile(readFile(TOKENS_FILE_PATH))
	var token string
	switch r.Method {
	case "GET":
		parts := strings.SplitN(r.URL.String(), "/token/", 2)
		if len(parts) == 2 {
			token = parts[1]
		}
	case "POST":
		token = r.Header.Get("token")
	default:
		w.WriteHeader(404)
		return
	}

	if username, ok := tokens[token]; ok {
		go openDoor()
		w.Write([]byte("OPEN"))
		appendFile(
			DOOR_LOG_PATH,
			time.Now().Format("2006-01-02T15:04:05")+" Open for "+username+"\n",
		)
	} else {
		w.WriteHeader(401)
	}
}

func openDoorRemoteGPIO() {
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
	doorRemote.Output()
	doorRemote.High()
}

func closeDoorRemoteGPIO() {
	rpio.Close()
}

func readFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		os.Exit(10)
	}
	return data
}

func appendFile(path string, text string) {
	if f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660); err == nil {
		defer f.Close()
		f.Write([]byte(text))
	} else {
		log.Fatal(err)
	}
}

func parseTokenFile(tokenFile []byte) map[string]string {
	tokens := make(map[string]string)
	t := string(tokenFile)
	lines := strings.Split(t, "\n")
	for _, line := range lines {
		tokenAndName := strings.SplitN(line, " ", 2)
		if len(tokenAndName) == 2 {
			tokens[tokenAndName[0]] = tokenAndName[1]
		}
	}
	return tokens
}

func openDoor() {
	doorRemote.Low()
	time.Sleep(time.Second * 2)
	doorRemote.High()
}
