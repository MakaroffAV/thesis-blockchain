//  closed: false
//  author: makarov aleksei
//  target: this is part of the package
//          core that stores code which
//          describes the blockchain transactions initialization

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

package core

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

import (
    "encoding/hex"
    "fmt"
    "log"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// newCoinBaseTransaction is function
// for creating "genesis" transaction in blockchain network
func newCoinBaseTransaction(to, data string) *transaction {

    if data == "" {
        data = fmt.Sprintf("new transaction to %s", to)
    }

    // set up
    // coinbase transaction input
    coinBaseTransactionInput := transactionInput{[]byte{}, -1, nil, []byte(data)}

    // set up
    // coinbase transaction output
    coinBaseTransactionOutput := NewTXOutput(10, to)

    // set up
    // completely new coinbase transaction
    coinBaseTransactionNew := transaction{
        Id:     nil,
        Input:  []transactionInput{coinBaseTransactionInput},
        Output: []transactionOutput{*coinBaseTransactionOutput},
    }

    coinBaseTransactionNew.Id = coinBaseTransactionNew.hash()

    return &coinBaseTransactionNew

}

// ------------------------------------------------------------------------ //

// new transaction is function
// for creating new transaction between two wallets in blockchain
func newTransaction(from, to string, amount int, UTXOSet *UTXOSet) *transaction {

    var (
        transactionInputs  []transactionInput
        transactionOutputs []transactionOutput
    )

    // fetch all user
    userWallets, userWalletsErr := newWallets()
    if userWalletsErr != nil {
        log.Panic(userWalletsErr)
    }

    fromWallet := userWallets.getWalletWithAddress(from)

    pubKeyHash := HashPubKey(fromWallet.PbKey)
    acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)

    if acc < amount {
        log.Panic("ERROR: Not enough funds")
    }

    // Build a list of inputs
    for txid, outs := range validOutputs {
        txID, err := hex.DecodeString(txid)
        if err != nil {
            log.Panic(err)
        }

        for _, out := range outs {
            input := transactionInput{txID, out, nil, fromWallet.PbKey}
            transactionInputs = append(transactionInputs, input)
        }
    }

    // Build a list of outputs
    transactionOutputs = append(transactionOutputs, *NewTXOutput(amount, to))
    if acc > amount {
        transactionOutputs = append(transactionOutputs, *NewTXOutput(acc-amount, from)) // a change
    }

    tx := transaction{nil, transactionInputs, transactionOutputs}
    tx.Id = tx.hash()
    UTXOSet.Blockchain.signTransaction(&tx, fromWallet.PrKey)

    return &tx
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
