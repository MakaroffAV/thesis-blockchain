package hashlib

import (
	"bytes"
	"crypto/sha256"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// Pass vars converted to bytes
// array and calculate on their base sha256 hash
func Sha256(toHashData ...[]byte) []byte {

	fullHash := sha256.Sum256(bytes.Join(toHashData, []byte{}))
	return fullHash[:]

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
