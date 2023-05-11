package main

import (
	"log"
	"os"

	ding "github.com/ding-live/ding-go"
)

func main() {
	client, err := ding.NewClient(ding.Config{
		CustomerUUID:      os.Getenv("DING_CUSTOMER_UUID"),
		APIKey:            os.Getenv("DING_API_KEY"),
		MaxNetworkRetries: ding.Int(4),
		LeveledLogger: &ding.Logger{
			// Disable default logging
			Level: ding.LevelNull,
		},
	})
	if err != nil {
		panic(err)
	}

	auth, err := client.Retry("5071dbf5-78d0-497a-b844-c1231808c3e9")
	if err != nil {
		panic(err)
	}

	log.Printf("auth: %+v", auth)
}
