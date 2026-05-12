// gettoken logs into Garmin Connect and prints the access token and display
// name on a single tab-separated line so shell scripts can capture them.
package main

import (
	"fmt"
	"log"
	"os"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"
)

func main() {
	email := os.Getenv("GARMIN_EMAIL")
	password := os.Getenv("GARMIN_PASSWORD")
	if email == "" || password == "" {
		log.Fatal("GARMIN_EMAIL and GARMIN_PASSWORD must be set")
	}
	c := gc.NewClient("")
	if err := c.Login(email, password); err != nil {
		log.Fatalf("login: %v", err)
	}
	fmt.Printf("%s\n%s\n", c.Token(), c.DisplayName())
}
