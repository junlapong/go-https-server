# Golang HTTPS/TLS Server

## Generation of self-signed key

using [mkcert](https://github.com/FiloSottile/mkcert)

```
$ mkdir cert
$ cd cert

$ mkcert localhost
The certificate is at "./localhost.pem" and the key at "./localhost-key.pem"
```

## Simple Golang HTTPS/TLS Server

```go
package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
)

func main() {

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
```

## Run

```
$ go run main.go
2019/10/08 23:16:01 ** Service Started on Port 8443 **
```

## Test with cURL

```
$ curl -Ik https://localhost:8443
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 08 Oct 2019 16:18:08 GMT
Content-Length: 16
```

## Check TLS Configuration

using [sslyze](https://github.com/nabla-c0d3/sslyze)

```
$ sslyze --regular localhost:8443

 SCAN RESULTS FOR LOCALHOST:8443 - 127.0.0.1
 -------------------------------------------

 * TLSV1_2 Cipher Suites:
       Forward Secrecy                    OK - Supported
       RC4                                OK - Not Supported

     Preferred:
        TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256                            128 bits                
     Accepted:
        TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256                      256 bits
        TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384                            256 bits
        TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256                            128 bits
```

## References
 - https://github.com/denji/golang-tls
 - https://golang.org/pkg/crypto/tls/#Config
