package main

import (
    "crypto/tls"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
    "sync"
    "time"

    rpio "github.com/stianeikeland/go-rpio"
    "golang.org/x/crypto/acme/autocert"
)

var (
    doorRemote = rpio.Pin(2)
    lock       sync.Mutex
)

const (
    doorLogPath      = "access.log"
    certificateCache = "certs"
    domain           = "unidoor.space"
    tokensFilePath   = "tokens"
    transmitTime     = 12 * time.Second
    cacheTime        = 24 * time.Hour
)

func main() {
    openGPIODoorRemote()
    defer closeGPIODoorRemote()

    certManager := autocert.Manager{
        Prompt:     autocert.AcceptTOS,
        HostPolicy: autocert.HostWhitelist(domain),
        Cache:      autocert.DirCache(certificateCache),
    }

    http.HandleFunc("/", rootHandler)
    http.HandleFunc("/token/", tokenHandler)
    http.HandleFunc("/token", tokenHandler)

    go http.ListenAndServe(":http", certManager.HTTPHandler(nil))

    server := &http.Server{
        Addr: ":https",
        TLSConfig: &tls.Config{
            GetCertificate: certManager.GetCertificate,
        },
    }

    log.Print("Starting server on port 80/443")
    if err := server.ListenAndServeTLS("", ""); err != nil {
        log.Fatal(err)
        os.Exit(3)
    }
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        w.Header().Set("Content-type", "text/html")
        w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%.0f", cacheTime.Seconds()))
        w.Write(indexFile)
    default:
        w.WriteHeader(404)
    }
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
    tokens := parseTokenFile(readFile(tokensFilePath))
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
        go appendFile(
            doorLogPath,
            time.Now().Format("2006-01-02T15:04:05")+" Open for "+username+"\n",
        )
    } else {
        w.WriteHeader(401)
    }
}

func openGPIODoorRemote() {
    if err := rpio.Open(); err != nil {
        log.Fatal(err)
        os.Exit(2)
    }
    doorRemote.Output()
    doorRemote.High()
}

func closeGPIODoorRemote() {
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
    lock.Lock()
    doorRemote.Low()
    time.Sleep(transmitTime)
    doorRemote.High()
    lock.Unlock()
}
