package main

import (
	"context"
	"flag"
	"fmt"
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

	domainName := "google.com"

	results, err := client.Domain.SearchAll(context.Background(), models.Search{
		Query: "name.ends_with: " + domainName,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	for i := range results {
		printDomainInfo(results[i])
	}
	if len(results) == 0 {
		log.Printf("Subdomains for domain %s not found", domainName)
	}
}

func printDomainInfo(details *models.Domain) {
	fmt.Println("Details about " + details.Name)
	for _, dnsA := range details.DNSRecords.A {
		fmt.Println("DNS A record:", dnsA)
	}
	fmt.Println("Website title:", details.Extract.Title)
	fmt.Println("Certificate subject org:", details.CertificateSummaries.SubjectDn.O)
	fmt.Println("Certificate issuer org:", details.CertificateSummaries.IssuerDn.O)
	fmt.Println("Updated at:", details.UpdatedAt)
}
