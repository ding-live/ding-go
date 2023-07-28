package main

import (
	"log"
	"os"

	"github.com/ding-live/ding-go"
)

func main() {
	client, err := ding.NewClient(ding.Config{
		CustomerUUID:      os.Getenv("DING_CUSTOMER_UUID"),
		APIKey:            os.Getenv("DING_API_KEY"),
		MaxNetworkRetries: ding.Int(4),
	})
	if err != nil {
		panic(err)
	}

	auth, err := client.Authenticate(ding.AuthenticateOptions{
		PhoneNumber:     "+xxxxxxxxxxx",
		IP:              ding.String("192.168.0.1"),
		DeviceType:      &ding.DeviceTypeIOS,
		AppVersion:      ding.String("1.2.0"),
		CallbackURL:     ding.String("https://example.com/callback"),
		IsReturningUser: true,
	})
	if err != nil {
		panic(err)
	}

	log.Printf("auth: %+v", auth)
}
