package main

import (
	"log"
	"os"

	ding "github.com/ding-live/ding-go"
)

type Logger struct{}

func (l *Logger) Debugf(msg string, keysAndValues ...interface{}) {
	log.Printf("DEBUG: %s %v\n", msg, keysAndValues)
}

func (l *Logger) Errorf(msg string, keysAndValues ...interface{}) {
	log.Printf("ERROR: %s\n", msg)
}

func (l *Logger) Infof(msg string, keysAndValues ...interface{}) {
	log.Printf("INFO: %s\n", msg)
}

func (l *Logger) Warnf(msg string, keysAndValues ...interface{}) {
	log.Printf("WARN: %s\n", msg)
}

func main() {
	client, err := ding.NewClient(ding.Config{
		CustomerUUID:      os.Getenv("DING_CUSTOMER_UUID"),
		APIKey:            os.Getenv("DING_API_KEY"),
		MaxNetworkRetries: ding.Int(4),
		LeveledLogger:     &Logger{},
	})
	if err != nil {
		panic(err)
	}

	auth, err := client.Authenticate(ding.AuthenticateOptions{
		PhoneNumber: "+33xxxxxxxx",
		IP:          ding.String("192.168.0.1"),
		DeviceType:  &ding.DeviceTypeIOS,
		AppVersion:  ding.String("1.2.0"),
		CallbackURL: ding.String("https://example.com/callback"),
	})
	if err != nil {
		log.Fatalf("unable to authenticate: %s", err)
	}

	log.Printf("auth: %+v", auth)
}
