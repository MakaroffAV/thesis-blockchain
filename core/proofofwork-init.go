//  closed: false
//  author: makarov aleksei
//  target: this is part of the package
//          core that stores code which
//          describes the blockchain proof of work algorithm initialization

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

package core

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

import "math/big"

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// newProofOfWork is access point
// for creating new proofOfWork data structure
func newProofOfWork(b *block) *proofOfWork {

	return &proofOfWork{
		Block:      b,
		Complexity: big.NewInt(1).Lsh(big.NewInt(1), uint(256-powComplexity)),
	}

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
