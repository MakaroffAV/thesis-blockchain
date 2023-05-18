//  closed: true
//  author: makarov aleksei
//  target: this is part of the package
//          core that stores code which
//          describes the blockchain creating and initializing operations

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

package core

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

import (
	"fmt"
	"log"
	"os"

	"github.com/MakaroffAV/thesis-blockchain/pkg/utils"
	"github.com/boltdb/bolt"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// newBlockchain is function for
// creating a copy of the current blockchain
func newBlockchain(address string) (*blockchain, error) {

	var (

		// bcTip is the last Block
		// Hash in created new one blockchain
		bcTip []byte
	)

	// check if blockchain already exist,
	// and if it exists - creating the new
	// blockchain is impossible, because of the lack of connectivity
	if utils.CheckDbExistence(dbFilename) == false {
		os.Exit(1)
	}

	fmt.Println("here.2.2.1")

	// create new database connection
	// todo: create special Db module
	db, dbErr := bolt.Open(dbFilename, 0600, nil)
	if dbErr != nil {
		log.Panic(dbErr)
		return nil, dbErr
	}

	fmt.Println("here.2.2.2")

	err := db.Update(func(tx *bolt.Tx) error {
		bcTip = tx.Bucket([]byte(tableBlocks)).Get([]byte("l"))
		return nil
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("here.2.2.3")

	return &blockchain{
		Db:  db,
		Tip: bcTip,
	}, nil

}

// ------------------------------------------------------------------------ //

// createBlockchain is function for
// creating new blockchain by passed first blockchain user wallet address
func createBlockchain(address string) (*blockchain, error) {

	var (

		// bcTip is the last Block
		// Hash in created new one blockchain
		bcTip []byte
	)

	// check if blockchain already exist,
	// and if it exists - creating the new
	// blockchain is impossible, because of the lack of connectivity
	if utils.CheckDbExistence(dbFilename) {
		os.Exit(1)
	}

	// create new database connection
	// todo: create special Db module
	db, dbErr := bolt.Open(dbFilename, 0600, nil)
	if dbErr != nil {
		log.Panic(dbErr)
	}

	_ = db.Update(func(tx *bolt.Tx) error {

		coinBaseTransaction := newCoinBaseTransaction(address, "")
		genesisBlock := newGenesisBlock(coinBaseTransaction)

		blocksTable, err := tx.CreateBucket([]byte(tableBlocks))
		if err != nil {
			log.Panic(err)
		}

		genesisBlockFmt, genesisBlockFmtErr := genesisBlock.marshal()
		if genesisBlockFmtErr != nil {
			log.Panic(genesisBlockFmtErr)
		}

		err = blocksTable.Put(genesisBlock.Hash, genesisBlockFmt)
		if err != nil {
			log.Panic(err)
		}

		err = blocksTable.Put([]byte("l"), genesisBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bcTip = genesisBlock.Hash
		return nil
	})

	return &blockchain{
		Db:  db,
		Tip: bcTip,
	}, nil

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
