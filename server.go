package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Write hello world to the response body
	io.WriteString(w, "🍉 Hello, world 🍉")
}

func main() {
	// Set up a /hello resource handler
	http.HandleFunc("/hello", helloHandler)

	// Listen to port 8443 and wait
	// log.Fatal(http.ListenAndServeTLS(":8443", "devlocal-mtls.cer", "devlocal-mtls.key", nil))

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile("devlocal-ca.pem")
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create TLS Config with the CA pool and enable
	// client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	// Create a server instance to listen on port 8443 with the TLS
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	// Listen to https connection with server certificate and wait
	log.Fatal(server.ListenAndServeTLS("devlocal-mtls.cer", "devlocal-mtls.key"))

}
