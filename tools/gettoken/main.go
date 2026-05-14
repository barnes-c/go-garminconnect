// gettoken logs into Garmin Connect and prints the access token and display
// name on separate lines so shell scripts can capture them.
// A token file is used to cache the session; SSO login is only performed when
// the cached token is missing, expired, or cannot be refreshed.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"
)

func main() {
	tokenFile := flag.String("token-file", ".garmin_token.json", "token cache file")
	flag.Parse()

	email := os.Getenv("GARMIN_EMAIL")
	password := os.Getenv("GARMIN_PASSWORD")

	c := gc.NewClient(*tokenFile)
	if err := c.Login(email, password); err != nil {
		log.Fatalf("login: %v", err)
	}
	fmt.Printf("%s\n%s\n", c.Token(), c.DisplayName())
}
