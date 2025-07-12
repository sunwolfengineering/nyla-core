package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// generateSalt generates a salt string based on the given IP and site ID.
func GenerateSalt(ip, siteID string) string {
	currentDate := time.Now().Format("20060102")
	return ip + "_" + siteID + "_" + currentDate
}

// GeneratePrivateIDHash generates a private ID hash based on the given inputs.
func GeneratePrivateIDHash(ip, userAgent, hostname, siteID string) (string, error) {
	salt := GenerateSalt(ip, siteID)
	data := salt + userAgent + hostname + siteID

	hasher := sha256.New()
	_, err := hasher.Write([]byte(data))
	if err != nil {
		return "", err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash, nil
}
