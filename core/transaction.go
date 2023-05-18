//  closed: false
//  author: makarov aleksei
//  target: this is part of the package
//          core that stores code which
//          describes transactions in blockchain

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

package core

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"math/big"
	"strings"

	"fmt"
	"log"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// transaction is data structure for
// describing the set of inputs and outputs
type transaction struct {
	Id     []byte
	Input  []transactionInput
	Output []transactionOutput
}

// ------------------------------------------------------------------------ //

// isCoinBase is method of transaction structure
// it compare current transaction with specially defined coinbase transaction
func (t *transaction) isCoinBase() bool {
	return len(t.Input) == 1 && len(t.Input[0].Id) == 0 && t.Input[0].Output == -1
}

// ------------------------------------------------------------------------ //

// marshal is method of transaction structure
// for converting transaction data structure to array of bytes
func (t *transaction) marshal() []byte {

	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(t)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()

}

// ------------------------------------------------------------------------ //

// Hash returns the hash of the Transaction
// hash is method of transaction structure
// it marshal transaction structure to bytes array and hash it
func (t *transaction) hash() []byte {

	transactionCopy := *t
	transactionCopy.Id = []byte{}

	hash := sha256.Sum256(transactionCopy.marshal())
	return hash[:]

}

// ------------------------------------------------------------------------ //

// Sign signs each input of a Transaction
// sign is method of transaction structure
// it iterates over all
func (t *transaction) sign(prKey ecdsa.PrivateKey, prevTransactions map[string]transaction) {

	// check if transaction
	// is coinbase it should not be signed
	if t.isCoinBase() {
		return
	}

	// iterate over all
	// previous inputs in transaction and check each
	for _, input := range t.Input {
		if prevTransactions[hex.EncodeToString(input.Id)].Id == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := t.simplify()

	for inID, vin := range txCopy.Input {
		prevTx := prevTransactions[hex.EncodeToString(vin.Id)]
		txCopy.Input[inID].Signature = nil
		txCopy.Input[inID].PublicKey = prevTx.Output[vin.Output].PublicKeyHash
		txCopy.Id = txCopy.hash()
		txCopy.Input[inID].PublicKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &prKey, txCopy.Id)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		t.Input[inID].Signature = signature
	}
}

// ------------------------------------------------------------------------ //

// String returns a human-readable representation of a transaction
func (t *transaction) convertToString() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", t.Id))

	for i, input := range t.Input {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Id))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Output))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PublicKey))
	}

	for i, output := range t.Output {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PublicKeyHash))
	}

	return strings.Join(lines, "\n")

}

// ------------------------------------------------------------------------ //

// simplify is method of transaction data structure
// it creates simplified version of transaction data structure
func (t *transaction) simplify() transaction {

	var (
		transactionInputs  []transactionInput
		transactionOutputs []transactionOutput
	)

	for _, input := range t.Input {
		transactionInputs = append(transactionInputs,
			transactionInput{input.Id, input.Output, nil, nil})
	}

	for _, output := range t.Output {
		transactionOutputs = append(transactionOutputs,
			transactionOutput{output.Value, output.PublicKeyHash})
	}

	return transaction{t.Id, transactionInputs, transactionOutputs}

}

// ------------------------------------------------------------------------ //

// check is method of transaction data structure
func (t *transaction) check(prevTransactions map[string]transaction) bool {

	// check if transaction
	// is coinbase it should not be signed
	if t.isCoinBase() {
		return true
	}

	for _, input := range t.Input {
		if prevTransactions[hex.EncodeToString(input.Id)].Id == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := t.simplify()
	curve := elliptic.P256()

	for inID, vin := range t.Input {
		prevTx := prevTransactions[hex.EncodeToString(vin.Id)]
		txCopy.Input[inID].Signature = nil
		txCopy.Input[inID].PublicKey = prevTx.Output[vin.Output].PublicKeyHash
		txCopy.Id = txCopy.hash()
		txCopy.Input[inID].PublicKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(keyLen / 2)])
		y.SetBytes(vin.PublicKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.Id, &r, &s) == false {
			return false
		}
	}

	return true

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
