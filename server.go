package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"

	"github.com/zenazn/goji"
)

var Config struct {
	bind     string
	filesDir string
	siteName string
}

func main() {
	flag.StringVar(&Config.bind, "b", "127.0.0.1:8080",
		"host to bind to (default: 127.0.0.1:8080)")
	flag.StringVar(&Config.filesDir, "d", "files/",
		"path to files directory (default: files/)")
	flag.StringVar(&Config.siteName, "n", "linx",
		"name of the site")
	flag.Parse()

	fmt.Printf("About to listen on http://%s\n", Config.bind)

	nameRe := regexp.MustCompile(`^/(?P<name>[a-z0-9-\.]+)$`)
	selifRe := regexp.MustCompile(`^/selif/(?P<name>[a-z0-9-\.]+)$`)

	goji.Get("/", indexHandler)
	goji.Post("/upload", uploadPostHandler)
	goji.Put("/upload", uploadPutHandler)
	goji.Get(nameRe, fileDisplayHandler)
	goji.Handle(selifRe, http.StripPrefix("/selif/", http.FileServer(http.Dir(Config.filesDir))))

	listener, err := net.Listen("tcp", Config.bind)
	if err != nil {
		log.Fatal("Could not bind: ", err)
	}

	goji.ServeListener(listener)
}