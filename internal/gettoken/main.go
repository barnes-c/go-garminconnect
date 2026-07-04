// gettoken logs into Garmin Connect and prints the access token and display name
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/barnes-c/go-garminconnect/internal/login"
)

func main() {
	tokenFile := flag.String("token-file", ".garmin_token.json", "token cache file")
	flag.Parse()

	c, err := login.Client(*tokenFile)
	if err != nil {
		log.Fatalf("login: %v", err)
	}
	fmt.Printf("%s\n%s\n", c.Token(), c.DisplayName())
}
