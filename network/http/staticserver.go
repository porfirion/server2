package http

import (
	"errors"
	"log"
	"net"
	"net/http"
)

// ServeStatic starts http server to serve static files (index.html and assets)
func ServeStatic(port int) {
	mux := &http.ServeMux{}
	mux.Handle("/", http.FileServer(http.Dir("templates")))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{Port: port});
	if err != nil {
		log.Fatal(err)
	}

	err = http.Serve(listener, mux);
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("HTTP SERVER STOPPED")
		} else {
			log.Fatal(err)
		}
	}
}