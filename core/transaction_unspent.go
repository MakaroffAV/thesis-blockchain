package core

import (
    "encoding/hex"
    "github.com/boltdb/bolt"
    "log"
)

type UTXOSet struct {
    Blockchain *blockchain
}

const utxoBucket = "chainstate"

// Reindex rebuilds the UTXO set
func (u UTXOSet) Reindex() {
    db := u.Blockchain.Db
    bucketName := []byte(utxoBucket)

    err := u.Blockchain.Db.Update(func(tx *bolt.Tx) error {
        err := tx.DeleteBucket(bucketName)
        if err != nil && err != bolt.ErrBucketNotFound {
            log.Panic(err)
        }

        _, err = tx.CreateBucket(bucketName)
        if err != nil {
            log.Panic(err)
        }

        return nil
    })
    if err != nil {
        log.Panic(err)
    }

    UTXO := u.Blockchain.findUnspentTransactions()

    err = db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket(bucketName)

        for txID, outs := range UTXO {
            key, err := hex.DecodeString(txID)
            if err != nil {
                log.Panic(err)
            }

            err = b.Put(key, outs.Serialize())
            if err != nil {
                log.Panic(err)
            }
        }

        return nil
    })
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (u UTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
    unspentOutputs := make(map[string][]int)
    accumulated := 0
    db := u.Blockchain.Db

    err := db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(utxoBucket))
        c := b.Cursor()

        for k, v := c.First(); k != nil; k, v = c.Next() {
            txID := hex.EncodeToString(k)
            outs := DeserializeOutputs(v)

            for outIdx, out := range outs.Outputs {
                if out.isLockedWithKey(pubkeyHash) && accumulated < amount {
                    accumulated += out.Value
                    unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
                }
            }
        }

        return nil
    })
    if err != nil {
        log.Panic(err)
    }

    return accumulated, unspentOutputs
}

// FindUTXO finds UTXO for a public key hash
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []transactionOutput {
    var UTXOs []transactionOutput
    db := u.Blockchain.Db

    err := db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(utxoBucket))
        c := b.Cursor()

        for k, v := c.First(); k != nil; k, v = c.Next() {
            outs := DeserializeOutputs(v)

            for _, out := range outs.Outputs {
                if out.isLockedWithKey(pubKeyHash) {
                    UTXOs = append(UTXOs, out)
                }
            }
        }

        return nil
    })
    if err != nil {
        log.Panic(err)
    }

    return UTXOs
}

// Update updates the UTXO set with transactions from the Block
// The Block is considered to be the tip of a blockchain
func (u UTXOSet) Update(block *block) {
    db := u.Blockchain.Db

    err := db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(utxoBucket))

        for _, tx := range block.Transactions {
            if tx.isCoinBase() == false {
                for _, vin := range tx.Input {
                    updatedOuts := transactionOutputs{}
                    outsBytes := b.Get(vin.Id)
                    outs := DeserializeOutputs(outsBytes)

                    for outIdx, out := range outs.Outputs {
                        if outIdx != vin.Output {
                            updatedOuts.Outputs = append(updatedOuts.Outputs, out)
                        }
                    }

                    if len(updatedOuts.Outputs) == 0 {
                        err := b.Delete(vin.Id)
                        if err != nil {
                            log.Panic(err)
                        }
                    } else {
                        err := b.Put(vin.Id, updatedOuts.Serialize())
                        if err != nil {
                            log.Panic(err)
                        }
                    }

                }
            }

            newOutputs := transactionOutputs{}
            for _, out := range tx.Output {
                newOutputs.Outputs = append(newOutputs.Outputs, out)
            }

            err := b.Put(tx.Id, newOutputs.Serialize())
            if err != nil {
                log.Panic(err)
            }
        }

        return nil
    })
    if err != nil {
        log.Panic(err)
    }
}

// CountTransactions returns the number of transactions in the UTXO set
func (u UTXOSet) CountTransactions() int {
    db := u.Blockchain.Db
    counter := 0

    err := db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(utxoBucket))
        c := b.Cursor()

        for k, _ := c.First(); k != nil; k, _ = c.Next() {
            counter++
        }

        return nil
    })
    if err != nil {
        log.Panic(err)
    }

    return counter
}
