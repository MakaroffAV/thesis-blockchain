//  closed: false
//  author: makarov aleksei
//  target: this is part of the package
//          core that stores code which
//          initializes collection of user wallets in blockchain

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

package core

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //

// newWallets is function which
// creates new collection of wallets or
// fetches wallet info from file if it exists
func newWallets() (*wallets, error) {
	wallets := wallets{}
	wallets.Wallets = make(map[string]*wallet)

	err := wallets.LoadFromFile()

	return &wallets, err

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ //
