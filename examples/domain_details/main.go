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

	details, err := client.Domain.GetDetails(context.Background(), domainName)
	if err != nil {
		log.Fatal(err.Error())
	}
	if details != nil {
		printDomainInfo(details)
	} else {
		log.Printf("Domain %s not found", domainName)
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
