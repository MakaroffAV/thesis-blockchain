/*
package core

import (

	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/MakaroffAV/thesis-blockchain/pkg/base58"

)

// no delete

func CreateWallet() string {

		wallets, _ := newWallets()
		address := wallets.createNewWallet()
		wallets.SaveToFile()

		return address
	}

// no delete

// CLI responsible for processing command line arguments
type CLI struct{}

	func (cli *CLI) printUsage() {
		fmt.Println("Usage:")
		fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
		fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
		fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
		fmt.Println("  listaddresses - Lists all addresses from the wallet file")
		fmt.Println("  printchain - Print all the blocks of the blockchain")
		fmt.Println("  reindexutxo - Rebuilds the UTXO set")
		fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
	}

	func (cli *CLI) validateArgs() {
		if len(os.Args) < 2 {
			cli.printUsage()
			os.Exit(1)
		}
	}

	func (cli *CLI) Run() {
		cli.validateArgs()

		getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
		createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
		createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
		listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
		printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
		reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
		sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

		getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
		createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
		sendFrom := sendCmd.String("from", "", "Source wallet address")
		sendTo := sendCmd.String("to", "", "Destination wallet address")
		sendAmount := sendCmd.Int("amount", 0, "Amount to send")

		switch os.Args[1] {
		case "getbalance":
			err := getBalanceCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "createblockchain":
			err := createBlockchainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "createwallet":
			err := createWalletCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "listaddresses":
			err := listAddressesCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "printchain":
			err := printChainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "send":
			err := sendCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "reindexutxo":
			err := reindexUTXOCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		default:
			cli.printUsage()
			os.Exit(1)
		}

		if getBalanceCmd.Parsed() {
			if *getBalanceAddress == "" {
				getBalanceCmd.Usage()
				os.Exit(1)
			}
			cli.getBalance(*getBalanceAddress)
		}

		if createBlockchainCmd.Parsed() {
			if *createBlockchainAddress == "" {
				createBlockchainCmd.Usage()
				os.Exit(1)
			}
			cli.createBlockchain(*createBlockchainAddress)
		}

		if createWalletCmd.Parsed() {
			cli.createWallet()
		}

		if listAddressesCmd.Parsed() {
			cli.listAddresses()
		}

		if printChainCmd.Parsed() {
			cli.printChain()
		}

		if reindexUTXOCmd.Parsed() {
			cli.reindexUTXO()
		}

		if sendCmd.Parsed() {
			if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
				sendCmd.Usage()
				os.Exit(1)
			}

			cli.send(*sendFrom, *sendTo, *sendAmount)
		}
	}

	func (cli *CLI) createBlockchain(address string) {
		if !ValidateAddress(address) {
			log.Panic("ERROR: Address is not valid")
		}
		bc := createBlockchain(address)

		UTXOSet := UTXOSet{bc}
		UTXOSet.Reindex()

		fmt.Println("Done!")
	}

	func (cli *CLI) createWallet() {
		wallets, _ := newWallets()
		address := wallets.createNewWallet()
		wallets.SaveToFile()

		fmt.Printf("Your new address: %s\n", address)
	}

	func (cli *CLI) getBalance(address string) {
		if !ValidateAddress(address) {
			log.Panic("ERROR: Address is not valid")
		}
		bc := newBlockchain(address)
		UTXOSet := UTXOSet{bc}
		defer bc.Db.Close()

		balance := 0
		pubKeyHash := base58.Decode([]byte(address))
		pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
		UTXOs := UTXOSet.FindUTXO(pubKeyHash)

		for _, out := range UTXOs {
			balance += out.Value
		}

		fmt.Printf("Balance of '%s': %d\n", address, balance)
	}

	func (cli *CLI) listAddresses() {
		wallets, err := newWallets()
		if err != nil {
			log.Panic(err)
		}
		addresses := wallets.getFullAddresses()

		for _, address := range addresses {
			fmt.Println(address)
		}
	}

	func (cli *CLI) printChain() {
		bc := newBlockchain("")
		defer bc.Db.Close()

		bci := bc.getIterator()

		for {
			block := bci.next()

			fmt.Printf("============ Block %x ============\n", block.Hash)
			fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
			pow := newProofOfWork(block)
			fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.check()))
			for _, tx := range block.Transactions {
				fmt.Println(tx)
			}
			fmt.Printf("\n\n")

			if len(block.PrevBlockHash) == 0 {
				break
			}
		}
	}

	func (cli *CLI) send(from, to string, amount int) {
		if !ValidateAddress(from) {
			log.Panic("ERROR: Sender address is not valid")
		}
		if !ValidateAddress(to) {
			log.Panic("ERROR: Recipient address is not valid")
		}

		bc := newBlockchain(from)
		UTXOSet := UTXOSet{bc}
		defer bc.Db.Close()

		tx := newTransaction(from, to, amount, &UTXOSet)
		cbTx := newCoinBaseTransaction(from, "")
		txs := []*transaction{cbTx, tx}

		newBlock := bc.addBlock(txs)
		UTXOSet.Update(newBlock)
		fmt.Println("Success!")

}

	func (cli *CLI) reindexUTXO() {
		bc := newBlockchain("")
		UTXOSet := UTXOSet{bc}
		UTXOSet.Reindex()

		count := UTXOSet.CountTransactions()
		fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
	}
*/
package core
