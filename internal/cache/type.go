package cache

type Balance interface {
	GetBalance(walletId string) (int64, error)
}
