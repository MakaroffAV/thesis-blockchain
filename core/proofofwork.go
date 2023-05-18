//  closed: false
//  author: makarov aleksei
//  target: this is part of the package
//          core that stores code which
//          describes the blockchain proof of work algorithm

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

package core

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

import (
	"bytes"
	"math"
	"math/big"

	"github.com/MakaroffAV/thesis-blockchain/pkg/hashlib"
	"github.com/MakaroffAV/thesis-blockchain/pkg/hexlib"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

const (

	// var for describing the
	// Complexity of the proof-of-work algorithm
	powComplexity = 20
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// proofOfWork is data structure for
// describing the necessary headers for proof of work algorithm execution
type proofOfWork struct {
	Block      *block
	Complexity *big.Int
}

// ------------------------------------------------------------------------ //

// check is method of proofOfWork data structure
// it checks if block in chain generated with proof of work algorithm
func (pow proofOfWork) check() bool {

	var hashInt big.Int

	blockDataHashed := hashlib.Sha256(
		bytes.Join(
			[][]byte{
				pow.Block.PrevBlockHash,
				pow.Block.hashTransactions(),
				hexlib.IntToHex(pow.Block.Timestamp),
				hexlib.IntToHex(int64(powComplexity)),
				hexlib.IntToHex(int64(pow.Block.Nonce)),
			},
			[]byte{},
		),
	)

	return hashInt.SetBytes(blockDataHashed[:]).Cmp(pow.Complexity) == -1

}

// ------------------------------------------------------------------------ //

// do is method of proofOfWork data structure
// it executes necessary work for calculating
// Hash and Nonce Value for generating new block in chain
func (pow proofOfWork) do() ([]byte, int) {

	var (
		currentNonce = 0
		hash         []byte
		hashInt      big.Int
	)

	for currentNonce < math.MaxInt64 {

		hash = hashlib.Sha256(
			bytes.Join(
				[][]byte{
					pow.Block.PrevBlockHash,
					pow.Block.hashTransactions(),
					hexlib.IntToHex(pow.Block.Timestamp),
					hexlib.IntToHex(int64(powComplexity)),
					hexlib.IntToHex(int64(currentNonce)),
				},
				[]byte{},
			),
		)

		if hashInt.SetBytes(hash[:]).Cmp(pow.Complexity) != -1 {
			currentNonce++
		} else {
			break
		}

	}

	return hash[:], currentNonce

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
