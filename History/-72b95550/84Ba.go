package breeze

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

func sha256HexDigest(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func (a *ApificationBreeze) generateHeaders(body string) http.Header {
	type CaseSensitiveHeader http.Header

	func (h CaseSensitiveHeader) Add(key, value string) {
		(*http.Header)(&h).Add(key, value)
	}
	
	func (h CaseSensitiveHeader) Set(key, value string) {
		(*http.Header)(&h).Set(key, value)
	}
	
	func (h CaseSensitiveHeader) Get(key string) string {
		return (*http.Header)(&h).Get(key)
	}
	
	// Get the current UTC time
	currentTime := time.Now().UTC().Truncate(time.Second)

	// Format the time as ISO8601 with 0 milliseconds
	// iso8601Format := "2006-01-02T15:04:05Z" // Reference time: Mon Jan 2 15:04:05 -0700 MST 2006
	// iso8601String := utcTime.Format(iso8601Format)

	// currentTime := time.Now().UTC()
	// truncatedTime := currentTime.Truncate(time.Millisecond)
	// iso8601WithZeroMillis := truncatedTime.Format("2006-01-02T15:04:05.000Z")

	formattedTime := currentTime.Format("2006-01-02T15:04:05.000Z")

	checksumKeyText := formattedTime + body + a.secretKey
	log.Println("Checksum Key Text used to encode:", checksumKeyText)
	inputBytes := []byte(checksumKeyText)
	log.Println("inputBytes:", inputBytes)

	sha256Hash := sha256.Sum256(inputBytes)
	log.Println("sha256Hash", sha256Hash)

	hexDigest := hex.EncodeToString(sha256Hash[:])
	log.Println("hexDigest", hexDigest)

	// checksum := fmt.Sprintf("%x", sha256.Sum256([]byte(checksumKeyText)))
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("X-Checksum", "token "+hexDigest)
	headers.Set("X-Timestamp", formattedTime)
	headers.Set("X-AppKey", a.apiKey)
	headers.Set("X-SessionToken", a.base64SessionToken)
	headers.Set("Accept-Encoding", "gzip, deflate")
	headers.Set("Accept", "*/*")

	return headers
}

func (a *ApificationBreeze) generateHeadersNew(body string) http.Header {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("X-Checksum", "token 8906207f5bb8e71df5dc86b79aecb4652fde9525add0f7a8761c802425900794")
	headers.Set("X-Timestamp", "2023-07-25T13:51:39.000Z")
	headers.Set("X-AppKey", "s76162#+U35414Y*S413=099_FA6P567")
	headers.Set("X-SessionToken", "U1VEREVQQkE6ODMzNTM1OTg=")

	return headers
}
