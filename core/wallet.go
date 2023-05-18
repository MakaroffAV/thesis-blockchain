//  closed: false
//  author: makarov aleksei
//  target: this is part of the package
//          core that stores code which
//          describes the user wallets in blockchain

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

package core

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

import (
	"crypto/ecdsa"

	"github.com/MakaroffAV/thesis-blockchain/pkg/base58"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

const (
	walletVersion            = byte(0x00)
	walletStorageFilename    = "wallet.dat"
	walletAddressCheckSumLen = 4
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// wallet is data structure for
// describing the pair of users private
// and public keys, with them user could create and sign transactions
type wallet struct {
	PrKey ecdsa.PrivateKey
	PbKey []byte
}

// ------------------------------------------------------------------------ //

// getAddress is method of wallet data structure
// it returnes
func (w wallet) getAddress() []byte {

	pubKeyHash := HashPubKey(w.PbKey)

	versionedPayload := append([]byte{walletVersion}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := base58.Encode(fullPayload)

	return address

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
