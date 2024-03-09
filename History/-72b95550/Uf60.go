package breeze

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"
)

func sha256HexDigest(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func (a *ApificationBreeze) generateHeaders(body string) http.Header {
	currentTime := time.Now().UTC().Truncate(time.Second)
	formattedTime := currentTime.Format("2006-01-02T15:04:05.000Z")

	checksumKeyText := formattedTime + body + a.secretKey
	sha256HexDigest(checksumKeyText)
	log.Debug("Key Text used to encode:", checksumKeyText)
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
