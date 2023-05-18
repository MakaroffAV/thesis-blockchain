//  closed: false
//  author: makarov aleksei
//  target: this is part of the package
//          core that stores code which
//          describes the blockchain transaction's inputs

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

package core

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

import "bytes"

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// transactionInput is data structure for
// describing the links to previous transactionOutput(s)
// that are now used to create new Transactions in blockchain
type transactionInput struct {
	// transaction Id
	Id []byte
	// previous transaction Id
	Output int
	// sender transaction Signature,
	// it is sender private key, this field guarantee that,
	// only transaction sender could use transaction Output
	Signature []byte
	// PublicKey is field
	PublicKey []byte
}

// ------------------------------------------------------------------------ //

// useKey is method of transactionInput structure
// it compares passed unlocking data and transaction Input's Signature
// after this check we define, could we bind transaction Input with another transaction Output
func (tInput transactionInput) useKey(publicKeyHash []byte) bool {
	return bytes.Compare(HashPubKey(tInput.PublicKey), publicKeyHash) == 0
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
