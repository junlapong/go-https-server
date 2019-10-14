package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

func main() {

	// redirect every http request to https
	go http.ListenAndServe(":8080", http.HandlerFunc(redirect))

	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		//CipherSuites:             nil,
		CipherSuites: []uint16{
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}

	srv := &http.Server{
		Addr:         ":8443",
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	// Enable http2
	http2.ConfigureServer(srv, nil)

	log.Println("** Service Started on Port 8443 **")
	err := srv.ListenAndServeTLS("cert/localhost.pem", "cert/localhost-key.pem")

	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"status":"ok"}`+"\n")
}

func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	//target := "https://" + req.Host + req.URL.Path
	target := "https://localhost:8443" + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
		// see comments below and consider the codes 308, 302, or 301
		http.StatusPermanentRedirect)
}
