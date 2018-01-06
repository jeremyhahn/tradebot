package converter

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
)

func ConvertDaoUserToCommonUser(user *dao.User) *common.User {
	return &common.User{
		Id:       user.Id,
		Username: user.Username}
}

func ConvertDaoWalletToCommonWallet(wallet *dao.UserWallet) *common.CryptoWallet {
	return &common.CryptoWallet{
		Address:  wallet.Address,
		Currency: wallet.Currency}
}
