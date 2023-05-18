package core

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

import (
	"errors"
	"fmt"

	"github.com/MakaroffAV/thesis-blockchain/pkg/cslog"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

var (
	errWalletFetchExisting = errors.New("errWalletFetchExisting")
	errW1                  = errors.New("errW1")
	errW2                  = errors.New("errW2")
	errCreateWallet        = errors.New("errCreateWallet")
	errCreateBlockchain    = errors.New("errCreateBlockchain")

	errCheckUserAddress = errors.New("errCheckUserAddress")
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// CreateUserWallet is public function for
// creating new user (wallet) in blockchain network
func CreateUserWallet() (string, error) {

	w, _ := newWallets()

	a, aErr := w.createNewWallet()
	if aErr != nil {
		cslog.Fail(errW2)
		return "", errW2
	}

	w.SaveToFile()

	return a, nil

}

// ------------------------------------------------------------------------ //

// GetBalance is public function for
// fetching user (wallet) balance in blockchain network
func GetUserBalance(address string) (int, error) {

	var (
		balance = 0
	)

	// check the user's wallet for
	// the possibility of existence in the blockchain network
	if !checkAddress(address) {
		return 0, errCheckUserAddress
	}

	// create blockchain duplicate
	b, bErr := newBlockchain(address)
	if bErr != nil {
		return 0, errCheckUserAddress
	}

	t := UTXOSet{b}

	p := GetPbHash(address)

	u := t.FindUTXO(p)

	for _, out := range u {
		balance += out.Value
	}

	return balance, nil

}

// ------------------------------------------------------------------------ //

// CreateNewTransaction is function for
// creating new transaction in blockchain
func CreateNewTransaction(sen, rec string, amount int) {

	if !checkAddress(sen) {
		return
	}

	if !checkAddress(rec) {
		return
	}

	fmt.Println("here2.1")

	b, bErr := newBlockchain(sen)
	if bErr != nil {
		return
	}

	fmt.Println("here2.2")

	u := UTXOSet{b}

	t := []*transaction{newTransaction(sen, rec, amount, &u)}

	fmt.Println("here3")

	newBlock := b.addBlock(t)
	fmt.Println("here4")
	u.Update(newBlock)

}

// ------------------------------------------------------------------------ //

func CreateBlockchain(address string) (bool, error) {

	if !checkAddress(address) {
		return false, errCreateBlockchain
	}

	b, bErr := createBlockchain(address)
	if bErr != nil {
		return false, errCreateBlockchain
	}

	u := UTXOSet{b}
	u.Reindex()

	return true, nil

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
