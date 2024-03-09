// Example for placing an order
package main

// func main() {

// 	conf := &go5paisa.Config{
// 		AppName:       "YOUR_APP_NAME_HERE",
// 		AppSource:     "YOUR_APP_SOURCE_HERE",
// 		UserID:        "YOUR_USER_ID_HERE",
// 		Password:      "YOUR_PASSWORD_HERE",
// 		UserKey:       "YOUR_USER_KEY_HERE",
// 		EncryptionKey: "YOUR_ENCRYPTION_KEY_HERE",
// 	}
// 	order := go5paisa.Order{
// 		Exchange:        go5paisa.BSE,
// 		ScripCode:       11111,
// 		ExchangeSegment: go5paisa.CASH,
// 		Qty:             1,
// 		OrderType:       go5paisa.BUY,
// 	}
// 	appConfig := go5paisa.Init(conf)
// 	client, err := go5paisa.Login(appConfig, "xyz@gmail.com", "password", "YYYYMMDD")
// 	if err != nil {
// 		panic(err)
// 	}
// 	res, err := client.PlaceOrder(order)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("%+v\n", res)
// }
