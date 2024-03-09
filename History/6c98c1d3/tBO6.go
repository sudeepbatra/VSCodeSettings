package smartapi

import (
	"fmt"

	SmartApi "github.com/angel-one/smartapigo"
)

func smartApiSessionToken() {
	ABClient := SmartApi.New("S1632585", "DevDas123", "zlndzd6z")

	fmt.Println("Client :- ", ABClient)

	// User Login and Generate User Session
	session, err := ABClient.GenerateSession("totp here")

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
