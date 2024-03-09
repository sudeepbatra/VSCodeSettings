package smartapi

import (
	"fmt"

	SmartApi "github.com/angel-one/smartapigo"
)

func smartApiSessionToken() {
	ABClient := SmartApi.New("ClientCode", "Password", "API Key")

	fmt.Println("Client :- ", ABClient)

}
