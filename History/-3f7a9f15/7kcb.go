package nse

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sudeepbatra/alpha-hft/logger"
)

func FetchAndSaveNSECorporateActions() {
	logger.Log.Info().Msg("Fetching NSE Corporate Actions")

	url := "https://www.nseindia.com/api/corporates-corporateActions?index=equities"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:121.0)")
	req.Header.Set("Cookie", "nsit=2uy6xga424za2DL-4-aGTGat; nseappid=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcGkubnNlIiwiYXVkIjoiYXBpLm5zZSIsImlhdCI6MTcwNTg0MTU1OCwiZXhwIjoxNzA1ODQ4NzU4fQ.QXeGvB3vNfhhDxi1sS3tIw8Q9946RSmABIyGCEmGXHM")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Parse the JSON response
	var corporateActions []CorporateAction
	err = json.NewDecoder(resp.Body).Decode(&corporateActions)
	if err != nil {
		log.Fatal(err)
	}

	// // Connect to the PostgreSQL database
	// connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	// 	dbHost, dbPort, dbUser, dbPassword, dbName)
	// db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// // Insert data into the PostgreSQL database
	// for _, ca := range corporateActions {
	// 	_, err := db.Exec("INSERT INTO corporate_actions (symbol, series, ind, face_val, subject, ex_date, rec_date, comp, isin) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
	// 		ca.Symbol, ca.Series, ca.Ind, ca.FaceVal, ca.Subject, ca.ExDate, ca.RecDate, ca.Comp, ca.Isin)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	// fmt.Println("Data saved to PostgreSQL successfully.")
}
