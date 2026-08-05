package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func VerifyHMAC(
	ticketID string,
	expTime int64,
	provided string,
	secret string,
) bool {

	message := fmt.Sprintf(
		"%s.%d",
		ticketID,
		expTime,
	)

	mac := hmac.New(
		sha256.New,
		[]byte(secret),
	)

	mac.Write([]byte(message))

	expected := hex.EncodeToString(
		mac.Sum(nil),
	)

	return hmac.Equal(
		[]byte(expected),
		[]byte(provided),
	)
}
