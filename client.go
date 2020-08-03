package main

// Client requests /hello and prints body to stdouot

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// Request /hello over https port 8443
	// r, err := http.Get("https://devlocal:8443/hello")

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile("devlocal-ca.pem")
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Read the key pair to create certificate
	cert, err := tls.LoadX509KeyPair("devlocal-mtls.cer", "devlocal-mtls.key")
	if err != nil {
		log.Fatal(err)
	}
	// Create a HTTPS client and supply the created CA pool and certificate
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	// Request /hello via the created HTTPS client over port 8443 via GET
	r, err := client.Get("https://devlocal:8443/hello")

	if err != nil {
		log.Fatal(err)
	}

	// Read response body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Print response body to stdout
	fmt.Printf("%s\n", body)
}
