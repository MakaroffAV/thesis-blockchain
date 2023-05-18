//  closed: true
//  author: makarov aleksei
//  target: this is part of the package
//          core that stores code which
//          describes the Block data structure and its methods
// todo:    logging

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

package core

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

import (
	"bytes"

	"github.com/MakaroffAV/thesis-blockchain/pkg/hashlib"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// block is data structure for
// describing the smallest part of blockchain network
// block stores information about users and their transactions in network
type block struct {
	Hash          []byte
	Nonce         int
	Timestamp     int64
	Transactions  []*transaction
	PrevBlockHash []byte
}

// ------------------------------------------------------------------------ //

// hashTransactions is method of block structure
// fetch all Transactions from block, get transaction ids,
// join them in bytes array with no separator, hash it with sha256 algorithm
func (b block) hashTransactions() []byte {

	var transactionIds [][]byte

	for _, transaction := range b.Transactions {
		transactionIds = append(transactionIds, transaction.Id)
	}

	// mTree := NewMerkleTree(transactions)
	// return mTree.RootNode.Data

	return hashlib.Sha256(bytes.Join(transactionIds, []byte{}))

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
