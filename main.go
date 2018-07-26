package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/dodocat/aliyun-docker-client/models"
)

func main() {

	pool := x509.NewCertPool()

	caCertPath := "ca.pem"

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	cliCrt, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		fmt.Println("Loadx509keypair err:", err)
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            pool,
			Certificates:       []tls.Certificate{cliCrt},
			InsecureSkipVerify: true,
		},
	}

	baseURL, err := url.Parse(os.Getenv("ALI_DOCKER_HOST"))

	if err != nil {
		log.Fatal(err)
	}
	client := models.AliClient{BaseURL: baseURL}
	client.HttpClient = &http.Client{Transport: tr}
	projects, err := client.ListProject()
	if err != nil {
		log.Println(err)
		return
	}

	for _, projet := range projects {
		println(projet.Name)
	}
}
