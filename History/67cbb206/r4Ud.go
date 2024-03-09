// Example to fetch all holdings for a user
package main

import (
	"fmt"

	"github.com/5paisa/go5paisa"
)

func main() {

	conf := &go5paisa.Config{
		AppName:       "5P50603710",
		AppSource:     "10427",
		UserID:        "uxuZEFys5nv",
		Password:      "7elTHyW0EC3",
		UserKey:       "sR12m8nkT8VEPXtfgLFlspj5BQlSqB51",
		EncryptionKey: "jTS6yEtvhXThvDTYNHQNVXmklWFEaeQj",
	}
	appConfig := go5paisa.Init(conf)
	client, err := go5paisa.Login(appConfig, "9325487506", "574812", "19771013")
	if err != nil {
		panic(err)
	}
	holdings, err := client.GetHoldings()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", holdings)
}
