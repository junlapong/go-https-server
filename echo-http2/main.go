package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/http2"
)

func main() {

	cer, err := tls.LoadX509KeyPair("cert/localhost.pem", "cert/localhost-key.pem")
	if err != nil {
		log.Println(err)
		return
	}

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		Certificates:             []tls.Certificate{cer},
		CipherSuites: []uint16{
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}

	e := echo.New()
	e.HEAD("/", func(c echo.Context) error {

		return c.String(http.StatusOK, "")
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK\n")
	})

	e.TLSServer = &http.Server{
		Addr:         ":8443",
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	// Enable http2
	http2.ConfigureServer(e.TLSServer, nil)
	e.Logger.Fatal(e.StartServer(e.TLSServer))
}
