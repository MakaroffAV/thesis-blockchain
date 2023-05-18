//  closed: true
//  author: makarov aleksei
//  target: this is part of the package
//          core that stores code which
//          describes the blockchain structure and interaction operations
// todo:    logging

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

package core

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"log"

	"github.com/boltdb/bolt"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

const (

	// dbFilename is filename of database
	// for storing blocks of blockchain and ensuring consistency
	dbFilename = "blockchain.db"

	// blocksBucket is name of table in
	// database for storing full blockchain blocks data
	tableBlocks = "blocks"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// blockchain is data structure
// for describing the blocks sequence
type blockchain struct {
	Db  *bolt.DB
	Tip []byte
}

// ------------------------------------------------------------------------ //

// getIterator is method of blockchain structure
// it is just access point for creating new blockchainIterator structure
func (bc *blockchain) getIterator() *blockchainIterator {
	return &blockchainIterator{
		db:          bc.Db,
		currentHash: bc.Tip,
	}
}

// ------------------------------------------------------------------------ //

// addBlock is method of blockchain structure
// it generate new block in chain, save it in database
// and replace info about last block to new generated one
func (bc *blockchain) addBlock(transactions []*transaction) *block {

	var lastBlockHash []byte

	for _, currTransaction := range transactions {
		if bc.checkTransaction(currTransaction) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	// Get last saved Block Hash
	// todo: handler error
	_ = bc.Db.View(func(tx *bolt.Tx) error {
		lastBlockHash = tx.Bucket([]byte(tableBlocks)).Get([]byte("l"))
		return nil
	})

	// Generate new Block structure
	newBlock := newBlock(transactions, lastBlockHash)

	// Replace previous last Block Hash in database to new one
	_ = bc.Db.Update(func(tx *bolt.Tx) error {
		blocksTable := tx.Bucket([]byte(tableBlocks))
		// todo: handler error

		newBlockFmt, _ := newBlock.marshal()

		if err := blocksTable.Put(newBlock.Hash, newBlockFmt); err != nil {
			return err
		}

		if err := blocksTable.Put([]byte("l"), newBlock.Hash); err != nil {
			return err
		}

		bc.Tip = newBlock.Hash

		return nil
	})

	return newBlock

}

// ------------------------------------------------------------------------ //

// findTransactionWithId is method of
// blockchain structure, it iterates over
// all blockchain and finds transaction with passed id
func (bc *blockchain) findTransactionWithId(id []byte) (transaction, error) {

	var (
		// define blockchain iterator
		blockchainIterator = bc.getIterator()
	)

	// iterate over
	// all blocks in blockchain
	for {

		// fetch following
		// block in blockchain
		block := blockchainIterator.next()

		// iterate over all
		// transactions in block
		for _, currTransaction := range block.Transactions {

			// compare current
			// transaction id with passed transaction id
			if bytes.Compare(currTransaction.Id, id) == 0 {
				return *currTransaction, nil
			}

		}

		// if we have reached
		// genesis Block, go to next Block in blockchain
		if len(block.PrevBlockHash) == 0 {
			break
		}

	}

	// for loop have not been broken,
	// transaction with passed id is not found
	return transaction{}, errors.New("transaction is not found")

}

// ------------------------------------------------------------------------ //

// checkTransaction is method of blockchain structure
// it get toCheckTransaction and iterate over all referenced
// transaction outputs and check existence of toCheckTransaction
func (bc *blockchain) checkTransaction(toCheckTransaction *transaction) bool {

	if toCheckTransaction.isCoinBase() {
		return true
	}

	var (
		// prevTransactions is map for
		// storing key value pair, transaction
		// id and  transaction data  structure,
		// which are linked to toCheckTransaction inputs
		prevTransactions = make(map[string]transaction)
	)

	// iterate over all
	// transactions, which are
	// linked to toCheckTransaction inputs
	for _, input := range toCheckTransaction.Input {

		// check for existing
		// toCheckTransaction
		// with referenced transaction input ids
		prevTransaction, err := bc.findTransactionWithId(input.Id)
		if err != nil {
			log.Panic(err)
		}

		prevTransactions[hex.EncodeToString(prevTransaction.Id)] = prevTransaction

	}

	// todo: comment
	return toCheckTransaction.check(prevTransactions)

}

// ------------------------------------------------------------------------ //

// findUnspentTransactions is method of blockchain structure
// it iterate above all blocks, Transactions in blockchain network
// and finds the unspent Transactions and generates the map of key
// value pairs, transaction id and array of spent transaction Output ids
func (bc *blockchain) findUnspentTransactions() map[string]transactionOutputs {

	var (

		// blockchain iterator
		bcIterator = bc.getIterator()

		// map of spent transaction,
		// outputs it contains key Value pair
		// as transaction id and array of spent transaction outputs
		spentTransactionsOutputs = make(map[string][]int)

		// array of unspent Transactions
		unspentTransactionOutputs = make(map[string]transactionOutputs)
	)

	for {

		// iterate over all
		// blocks in blockchain
		currBlock := bcIterator.next()

		// iterate over all
		// Transactions within the Block
		for _, currTransaction := range currBlock.Transactions {

			// convert transaction
			// id from bytes array to string
			currTransactionId := hex.EncodeToString(currTransaction.Id)

		Output:

			// iterate over all
			// transaction outputs
			for outputId, output := range currTransaction.Output {

				// check if transaction
				// already defined as spent
				if spentTransactionsOutputs[currTransactionId] != nil {

					// iterate over all
					// Transactions defined as spent
					for _, spentOutput := range spentTransactionsOutputs[currTransactionId] {

						// if transaction Output is spent
						// goto another transaction Output
						if spentOutput == outputId {
							continue Output
						}

					}

				}

				outs := unspentTransactionOutputs[currTransactionId]
				outs.Outputs = append(outs.Outputs, output)
				unspentTransactionOutputs[currTransactionId] = outs

			}

			// check if transaction
			// is not coinbase ("genesis transaction")
			if currTransaction.isCoinBase() == false {

				// iterate over all
				// transaction inputs
				for _, input := range currTransaction.Input {

					// convert transaction
					// id from bytes array to string
					inputTransactionId := hex.EncodeToString(input.Id)

					// add to spent Transactions map key Value pair,
					// transaction Id, transaction's Input Output Id
					spentTransactionsOutputs[inputTransactionId] = append(
						spentTransactionsOutputs[inputTransactionId], input.Output)

				}

			}

		}

		// if we have reached
		// genesis Block, go to next Block in blockchain
		if len(currBlock.PrevBlockHash) == 0 {
			break
		}

	}

	return unspentTransactionOutputs

}

// ------------------------------------------------------------------------ //

/*
// findUnspentTransactionsOutputs is method
// of blockchain structure it iterated over
// all user's unspent Transactions and select all
// outputs, which could be unlocked by user's private key
func (bc *blockchain) findUnspentTransactionsOutputs(pubKeyHash []byte) []transactionOutput {

    var (

        // unspentTransactions is
        // returned Value from findUnspentTransactions func
        unspentTransactions = bc.findUnspentTransactions(pubKeyHash)

        // unspentTransactionOutputs is array of transaction
        // outputs, which do not connect with another transaction inputs
        unspentTransactionOutputs []transactionOutput
    )

    // iterate over all
    // unspent Transactions in blockchain
    for _, currTransaction := range unspentTransactions {

        // iterate over all
        // unspent transaction outputs
        for _, output := range currTransaction.Output {

            // check if user's private
            // key could unlock transaction Output
            if output.isLockedWithKey(pubKeyHash) {

                // if transaction Output could be
                // unlocked by user's private key,
                // user can  connect  it with next
                // Input and create new transaction
                unspentTransactionOutputs = append(unspentTransactionOutputs, output)

            }

        }

    }

    return unspentTransactionOutputs

}

*/

// ------------------------------------------------------------------------ //

// signTransaction is method of blockchain structure
// it got transaction and user private key and signing transaction with it
func (bc *blockchain) signTransaction(toSignTransaction *transaction, prKey ecdsa.PrivateKey) {

	var (
		// prevTransactions is map for
		// storing key value pair, transaction
		// id and  transaction data  structure,
		// which are linked to toSignTransaction inputs
		prevTransactions = make(map[string]transaction)
	)

	// iterate over all
	// transactions, which are
	// linked to toSignTransaction inputs
	for _, input := range toSignTransaction.Input {

		// check for existing
		// transaction with referenced transaction input id
		prevTransaction, err := bc.findTransactionWithId(input.Id)
		if err != nil {
			log.Panic(err)
		}

		prevTransactions[hex.EncodeToString(prevTransaction.Id)] = prevTransaction

	}

	// sign passed transaction with user
	// private key and array of transaction
	// inputs which were referenced to toSignTransaction
	toSignTransaction.sign(prKey, prevTransactions)

}

// ------------------------------------------------------------------------ //

/*
// findSpendableTransactionOutputs is method of
// blockchain structure, it iterates over all unspent
// user's transaction outputs and counts the number of coins and spendable transaction outputs
func (bc *blockchain) findSpendableTransactionOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {

    var (

        // coinSum is start Value of user
        // coins number on spendable transaction outputs
        coinSum = 0

        // unspentTransactions is
        // returned Value from findUnspentTransactions func
        unspentTransactions = bc.findUnspentTransactions(pubKeyHash)

        // unspentTransactionOutputs is map of key
        // Value pairs of transaction ids and spendable transaction outputs
        unspentTransactionOutputs = make(map[string][]int)
    )

Work:

    // iterate over all
    // unspent user's Transactions
    for _, currTransaction := range unspentTransactions {

        // convert transaction id
        // from bytes array to string
        currTransactionId := hex.EncodeToString(currTransaction.Id)

        // iterate over all
        // transaction outputs
        for outputId, output := range currTransaction.Output {

            // check if transaction Output
            // could be unlocked by user private key
            if output.isLockedWithKey(pubKeyHash) && coinSum < amount {

                coinSum += output.Value
                unspentTransactionOutputs[currTransactionId] = append(unspentTransactionOutputs[currTransactionId], outputId)

                // check whether the
                // amount of coins found is sufficient
                if coinSum >= amount {
                    break Work
                }

            }

        }

    }

    return coinSum, unspentTransactionOutputs

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

*/
