package core

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// wallets - data structure for describing
// the collection of users wallets in blockchain network
type wallets struct {
	Wallets map[string]*wallet
}

// ------------------------------------------------------------------------ //

// createNewWallet is method of wallets structure
// it creates new user wallet and append it to collection of wallets
func (ws *wallets) createNewWallet() (string, error) {

	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.getAddress())

	ws.Wallets[address] = wallet

	return address, nil
}

// ------------------------------------------------------------------------ //

// getFullAddresses is method of wallets structure
// it returns an array of wallet addresses stored in the wallet file
func (ws *wallets) getFullAddresses() []string {

	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// ------------------------------------------------------------------------ //

// getWalletWithAddress is method of wallets structure
// it fetch user wallet from wallets collection with it address
func (ws *wallets) getWalletWithAddress(address string) wallet {
	return *ws.Wallets[address]
}

// ------------------------------------------------------------------------ //

// LoadFromFile loads wallets from the file
func (ws *wallets) LoadFromFile() error {

	if _, err := os.Stat(walletStorageFilename); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(walletStorageFilename)
	if err != nil {
		log.Panic(err)
	}

	var wallets wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	ws.Wallets = wallets.Wallets

	return nil
}

// ------------------------------------------------------------------------ //

// SaveToFile saves wallets to a file
func (ws *wallets) SaveToFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err, "here1")
	}

	err = ioutil.WriteFile(walletStorageFilename, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
