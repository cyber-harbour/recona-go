package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	recona "github.com/cyber-harbour/recona-go"
	"github.com/cyber-harbour/recona-go/services"
)

func main() {
	accessToken := flag.String("access_token", "", "API personal access token")
	flag.Parse()
	client, err := recona.NewClient(*accessToken)
	if err != nil {
		log.Fatal(err.Error())
	}
	asService := services.NewAccountService(client)

	account, err := asService.GetDetails(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Customer account quotas:")
	fmt.Printf("Subscription: %s\n", *account.SubscriptionName)
	fmt.Printf("Subscription period start at: %s\n", account.StartAt)
	fmt.Printf("Subscription period end at: %s\n", account.EndAt)
	fmt.Printf("Requests limit per day: %d\n", account.RequestLimitPerDay)
	fmt.Printf("Already used requests today: %d\n", account.RequestCount)
	fmt.Printf("RPS: %d\n", account.Permissions.RequestRateLimit)

}
