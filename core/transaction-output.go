//  closed: false
//  author: makarov aleksei
//  target: this is part of the package
//          core that stores code which
//          describes the blockchain transaction's outputs

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

package core

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/MakaroffAV/thesis-blockchain/pkg/base58"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// transactionOutput is data structure for
// describing the amount of coins and the transaction recipient's address
type transactionOutput struct {
	// amount of coins
	Value int
	// recipient's address,
	// it is private key in blockchain network
	PublicKeyHash []byte
}

type transactionOutputs struct {
	Outputs []transactionOutput
}

// ------------------------------------------------------------------------ //

// Lock signs the output
func (tOutput *transactionOutput) lock(address []byte) {
	pubKeyHash := base58.Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	tOutput.PublicKeyHash = pubKeyHash
}

// ------------------------------------------------------------------------ //

// IsLockedWithKey checks if the output can be used by the owner of the pubkey
// v
func (tOutput *transactionOutput) isLockedWithKey(publicKeyHash []byte) bool {
	return bytes.Compare(tOutput.PublicKeyHash, publicKeyHash) == 0
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// NewTXOutput create a new TXOutput
func NewTXOutput(value int, address string) *transactionOutput {
	txo := &transactionOutput{value, nil}
	txo.lock([]byte(address))

	return txo
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// Serialize serializes TXOutputs
func (outs transactionOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// DeserializeOutputs deserializes TXOutputs
func DeserializeOutputs(data []byte) transactionOutputs {
	var outputs transactionOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
