package core

// NewWallet creates and returns a Wallet
func NewWallet() *wallet {
	private, public := newKeyPair()
	return &wallet{private, public}
}
