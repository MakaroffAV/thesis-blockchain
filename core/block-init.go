//  closed: true
//  author: makarov aleksei
//  target: this is part of the package
//          core that stores code which
//          describes the create and generate block operations

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

package core

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

import "time"

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// newGenesisBlock is function for
// creating a new genesis block in blockchain network
func newGenesisBlock(coinbase *transaction) *block {
	return newBlock([]*transaction{coinbase}, []byte{})
}

// ------------------------------------------------------------------------ //

// newBlock is function for
// creating a new block of the blockchain network
func newBlock(transactions []*transaction, prevBlockHash []byte) *block {

	// define pointer
	// to block data structure
	block := &block{
		Hash:          []byte{},
		Nonce:         0,
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
	}

	// init proof of
	// work for created block
	blockPow := newProofOfWork(block)

	// run proof of
	// work for created block
	block.Hash, block.Nonce = blockPow.do()

	return block

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
