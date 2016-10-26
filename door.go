package main

import (
	"github.com/stianeikeland/go-rpio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var doorRemote = rpio.Pin(2)

var tokens = make(map[string]string)

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
	tokens := parseTokenFile(readFile("tokens"))
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
		w.Write([]byte("OPEN"))
		log.Print("Open for ", username)
		openDoor()
	} else {
		w.WriteHeader(401)
	}
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("You need to supply port as an argument")
		os.Exit(1)
	}
	port := os.Args[1]

	if err := rpio.Open(); err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
	defer rpio.Close()
	doorRemote.Output()
	doorRemote.High()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/token/", tokenHandler)
	http.HandleFunc("/token", tokenHandler)
	log.Print("Starting server on port ", port)
	if err := http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", nil); err != nil {
		log.Fatal(err)
		os.Exit(3)
	}
}

func readFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		os.Exit(10)
	}
	return data
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
