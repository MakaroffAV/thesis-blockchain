//  closed: true
//  author: makarov aleksei
//  target: this is part of the package
//          core that stores code which
//          describes the blockchain iterator architecture pattern

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

package core

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

import "github.com/boltdb/bolt"

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// blockchain is data structure
// for describing the Iterator architectural pattern
type blockchainIterator struct {
    db          *bolt.DB
    currentHash []byte
}

// ------------------------------------------------------------------------ //

// next is method of blockchainIterator structure
// it is just an iterator above all chain of blocks
func (bcI *blockchainIterator) next() *block {

    var blockRaw *block

    _ = bcI.db.View(func(tx *bolt.Tx) error {
        blockFmt := tx.Bucket([]byte(tableBlocks)).Get(bcI.currentHash)
        // todo: handler error
        blockRaw, _ = unMarshal(blockFmt)
        return nil
    })

    // todo: handler error

    bcI.currentHash = blockRaw.PrevBlockHash

    return blockRaw

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
