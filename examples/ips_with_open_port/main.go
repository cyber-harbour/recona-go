package main

import (
	"context"
	"flag"
	"log"

	recona "github.com/cyber-harbour/recona-go"
	"github.com/cyber-harbour/recona-go/models"
)

func main() {
	accessToken := flag.String("access_token", "", "API personal access token")
	flag.Parse()

	client, err := recona.NewClient(*accessToken)
	if err != nil {
		log.Fatal(err.Error())
	}

	var searchPort = "9200" // Elasticsearch port

	results, err := client.Host.SearchAll(context.Background(), models.Search{
		Query: "ports.port.eq: " + searchPort,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Detected %d IPs with open port %s", len(results), searchPort)
}
