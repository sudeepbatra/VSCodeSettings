// Example to fetch all holdings for a user
package main

import (
	"fmt"

	"github.com/5paisa/go5paisa"
)

func main() {

	conf := &go5paisa.Config{
		AppName:       "5P50603710",
		AppSource:     "YOUR_APP_SOURCE_HERE",
		UserID:        "YOUR_USER_ID_HERE",
		Password:      "YOUR_PASSWORD_HERE",
		UserKey:       "YOUR_USER_KEY_HERE",
		EncryptionKey: "YOUR_ENCRYPTION_KEY_HERE",
	}
	appConfig := go5paisa.Init(conf)
	client, err := go5paisa.Login(appConfig, "xyz@gmail.com", "password", "YYYYMMDD")
	if err != nil {
		panic(err)
	}
	holdings, err := client.GetHoldings()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", holdings)
}
