package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

const hmacMaxAge = time.Hour

func ValidateFeedbackHMAC(deviceID, timestamp, signature, secret string) error {
	if deviceID == "" || timestamp == "" || signature == "" || secret == "" {
		return fmt.Errorf("invalid signed feedback request")
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp")
	}

	age := time.Since(time.Unix(ts, 0))
	if age < 0 || age > hmacMaxAge {
		return fmt.Errorf("signature expired")
	}

	expected := computeHMAC(deviceID, timestamp, secret)
	if !hmac.Equal([]byte(signature), []byte(expected)) {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

func computeHMAC(deviceID, timestamp, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(deviceID + "|" + timestamp))
	return hex.EncodeToString(mac.Sum(nil))
}