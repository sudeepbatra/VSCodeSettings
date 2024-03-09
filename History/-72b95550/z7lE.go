package breeze

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/sudeepbatra/alpha-hft/logger"
)

var log = logger.GetLogger()

func sha256HexDigest(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func (a *ApificationBreeze) generateHeaders(body string) http.Header {
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
	log.Debug("Checksum Key Text used to encode:", checksumKeyText)
	inputBytes := []byte(checksumKeyText)
	log.Debug("inputBytes:", inputBytes)

	sha256Hash := sha256.Sum256(inputBytes)
	log.Debug("sha256Hash", sha256Hash)

	hexDigest := hex.EncodeToString(sha256Hash[:])
	log.Debug("hexDigest", hexDigest)

	// checksum := fmt.Sprintf("%x", sha256.Sum256([]byte(checksumKeyText)))
	// headers := http.Header{}
	// headers := make(http.Header)
	// headers["Content-Type"] = []string{"application/json"}
	// headers["X-Checksum"] = []string{"token " + hexDigest}
	// headers["X-Timestamp"] = []string{formattedTime}
	// headers["X-AppKey"] = []string{a.apiKey}
	// headers["X-SessionToken"] = []string{a.base64SessionToken}
	// headers["Accept-Encoding"] = []string{"gzip", "deflate"}
	// headers["Accept"] = []string{"*/*"}

	headers := make(http.Header)
	headers["Content-Type"] = []string{"application/json"}
	headers["X-Checksum"] = []string{"token " + hexDigest}
	headers["X-Timestamp"] = []string{formattedTime}
	headers["X-AppKey"] = []string{a.apiKey}
	headers["X-SessionToken"] = []string{a.base64SessionToken}
	headers["Accept-Encoding"] = []string{"gzip", "deflate"}
	headers["Accept"] = []string{"*/*"}

	return headers
}
