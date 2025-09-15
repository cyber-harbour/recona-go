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

	var IPv4 = "1.1.1.1"
	details, err := client.Host.GetDetails(context.Background(), IPv4)
	if err != nil {
		log.Fatal(err.Error())
	}
	if details != nil {
		printIpInfo(details)
	} else {
		log.Printf("IP %s not found", IPv4)
	}

}

func printIpInfo(details *models.Host) {
	examplesToPrint := 3
	fmt.Println("Details about " + details.IP)
	fmt.Printf("Updated at: %s\n", details.UpdatedAt)
	fmt.Printf("Country: %s\n", details.Geo.Country)
	fmt.Printf("CIDR: %s\n", details.Isp.Network)
	fmt.Printf("Open ports (%d)\n", len(details.Ports))
	if len(details.Ports) >= examplesToPrint {
		for i := 0; i < examplesToPrint; i++ {
			fmt.Printf("  - %d\n", details.Ports[i].Port)
		}
		fmt.Printf("    --- and %d more ---\n", len(details.Ports)-examplesToPrint)
	}
	fmt.Printf("Affected by %d CVEs\n", len(details.CVEList))
	fmt.Printf("PTR record: %s\n", details.PtrRecord.Value)
	fmt.Printf("Technologies detected (%d)\n", len(details.Technologies))
	if len(details.Technologies) >= examplesToPrint {
		for i := range details.Technologies {
			fmt.Printf("  - %s, version %s on port %d\n",
				details.Technologies[i].Name,
				details.Technologies[i].Version,
				details.Technologies[i].Port,
			)
			if i >= examplesToPrint {
				fmt.Printf("      --- and %d more ---\n", len(details.Technologies)-examplesToPrint)
				break
			}
		}
	}
	fmt.Println()
}
