package tst

import (
	"fmt"
	"os"

	"github.com/MakaroffAV/thesis-blockchain/core"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

func OneTest() {

	err1 := os.Remove("blockchain.db")
	if err1 != nil {
		fmt.Println(err1)
		return
	}

	err2 := os.Remove("wallet.dat")
	if err2 != nil {
		fmt.Println(err2)
		return
	}

	w1, w1Err := core.CreateUserWallet()
	if w1Err != nil {
		fmt.Println(w1Err)
		return
	}

	w2, w2Err := core.CreateUserWallet()
	if w1Err != nil {
		fmt.Println(w2Err)
		return
	}
	fmt.Println("here 1.1")

	w3, w3Err := core.CreateUserWallet()
	if w1Err != nil {
		fmt.Println(w3Err)
		return
	}

	fmt.Println("here", w1)

	_, bcErr := core.CreateBlockchain(w1)
	if bcErr != nil {
		fmt.Println(bcErr)
	}

	fmt.Println("here1", w1, w2, w3)

	// w1 10 -4 6
	// w2 0 +4 4
	// w3 0 0 0
	core.CreateNewTransaction(w1, w2, 4)

	fmt.Println("here2", w1)

	// w1 0 0 6
	// w2 4 -3 1
	// w3 0 +3 3
	core.CreateNewTransaction(w2, w3, 3)

	// w1 6 +1 7
	// w2 0 0 1
	// w3 3 -1 2
	core.CreateNewTransaction(w3, w1, 1)

	// w1 7 0 7
	// w2 1 +1 2
	// w3 2 -1 1
	core.CreateNewTransaction(w3, w2, 1)

	b1, b1Err := core.GetUserBalance(w1)
	if b1Err != nil {
		fmt.Println(b1Err)
	}
	if b1 != 7 {
		fmt.Println("incorrect balance w1")
	}

	b2, b2Err := core.GetUserBalance(w2)
	if b2Err != nil {
		fmt.Println(b2Err)
	}
	if b2 != 2 {
		fmt.Println("incorrect balance w2")
	}

	b3, b3Err := core.GetUserBalance(w3)
	if b3Err != nil {
		fmt.Println(b3Err)
	}
	if b3 != 1 {
		fmt.Println("incorrect balance w3")
	}

	fmt.Println(b1, b2, b3)

	fmt.Println("here")

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
