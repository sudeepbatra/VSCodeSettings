package smartapi

import (
	"fmt"

	SmartApi "github.com/angel-one/smartapigo"
)

func smartApiSessionToken() {
	ABClient := SmartApi.New("ClientCode", "Password", "zlndzd6z")

	fmt.Println("Client :- ", ABClient)

}
