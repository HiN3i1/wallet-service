package db

import (
	log "github.com/sirupsen/logrus"
)

type WalletType int

const (
	Vault    WalletType = 1
	Deposit  WalletType = 2
	Withdraw WalletType = 3
)

type CoinType uint32

const (
	TypeERC20USDT CoinType = 1
	TypeTRC20USDT CoinType = 2
)

var CoinTypes = map[string]CoinType{
	"ERC20USDT": TypeERC20USDT,
	"TRC20USDT": TypeTRC20USDT,
}

// Wallet is using call cybavo api
type Wallet struct {
	Id         int64
	APICode    string       `pg:"," json:"api_code"`
	ApiSecret  string       `pg:"," json:"api_secret"`
	WalletID   int64        `pg:",pk,unique" json:"wallet_id"`
	CoinType   CoinType     `pg:",notnull" json:"coin_type"`
	WalletType WalletType   `pg:",notnull" json:"wallet_type"`
	SubWallets *[]SubWallet `pg:"rel:has-many,fk:wallet_id"`
}

// Subwallet is user's wallet
type SubWallet struct {
	Id         int64
	CustomerID int64     `pg:",pk"`
	Customer   *Customer `pg:"rel:has-one"`
	Address    string    `pg:"," json:"address"`
	Memo       string    `pg:"," json:"memo"`
	WalletID   int64     `pg:",pk"`
	Wallet     *Wallet   `pg:"rel:has-one"`
}

func SetWalletAPICode(walletRequest *Wallet, coin string) (err error) {
	db := GetDBClient()
	wallet := new(Wallet)
	err = db.Model(wallet).Where("wallet_id = ?", walletRequest.WalletID).Select()
	wallet.APICode = walletRequest.APICode
	wallet.ApiSecret = walletRequest.ApiSecret
	wallet.WalletID = walletRequest.WalletID
	wallet.WalletType = walletRequest.WalletType
	wallet.CoinType = CoinTypes[coin]
	if err != nil {
		_, err = db.Model(wallet).Insert()
		if err != nil {
			log.Error("Failed to insert API token ", err)
			return err
		}
	} else {
		err = db.Update(wallet)
		if err != nil {
			log.Error("Failed to update API token ", err)
			return err
		}
	}
	return nil
}

func GetWalletById(wallet_id int64) (wallet *Wallet, err error) {
	db := GetDBClient()
	wallet = new(Wallet)
	err = db.Model(wallet).Where("wallet_id = ?", wallet_id).Select()
	return wallet, err
}

func GetWalletByType(coin string) (wallet *Wallet, err error) {
	db := GetDBClient()
	wallet = new(Wallet)
	coinType := CoinTypes[coin]
	err = db.Model(wallet).Where("coin_type = ?", coinType).First()
	return wallet, err
}
