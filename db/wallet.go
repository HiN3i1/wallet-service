package db

import (
	"fmt"

	"github.com/HiN3i1/wallet-service/cybavo"
	"github.com/jinzhu/copier"
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
	TypeETH       CoinType = 1
	TypeERC20USDT CoinType = 2
	TypeTRX       CoinType = 101
	TypeTRX20USDT CoinType = 102
)

var CoinTypes = map[string]CoinType{
	"ETH":      TypeETH,
	"ETH-USDT": TypeERC20USDT,
	"TRX":      TypeTRX,
	"TRX-USDT": TypeTRX20USDT,
}

func GetNativeCoin(coin string) (native string) {

	if coin == "ETH" || coin == "ETH-USDT" {
		native = "ETH"
	}
	if coin == "TRX" || coin == "TRX-USDT" {
		native = "TRX"
	}
	return native
}

// Wallet is using call cybavo api
type Wallet struct {
	Id         int64
	APICode    string       `pg:"," json:"api_code"`
	APISecret  string       `pg:"," json:"api_secret"`
	WalletID   int64        `pg:",pk,unique" json:"wallet_id"`
	CoinType   CoinType     `pg:",notnull" json:"coin_type"`
	WalletType WalletType   `pg:",notnull" json:"wallet_type"`
	SubWallets []*SubWallet `pg:"rel:has-many,fk:wallet_id"`
}

// Subwallet is user's wallet
type SubWallet struct {
	Id              int64
	CustomerID      string             `pg:","`
	Customer        *Customer          `pg:"rel:has-one,fk:customer_id"`
	Address         string             `pg:"," json:"address"`
	Memo            string             `pg:"," json:"memo"`
	WalletID        int64              `pg:","`
	Wallet          *Wallet            `pg:"rel:has-one,fk:wallet_id"`
	DepositCallBack []*DepositCallBack `pg:"rel:has-many"`
}

type DepositCallBack struct {
	Id              int64
	Type            int                    `json:"type"`
	UniqueID        string                 `pg:",pk,unique" json:"unique_id"`
	Currency        string                 `json:"currency"`
	TXID            string                 `json:"txid"`
	BlockHeight     int64                  `json:"block_height"`
	Amount          string                 `json:"amount"`
	Memo            string                 `json:"memo"`
	ToAddress       string                 `json:"to_address"`
	WalletID        int64                  `json:"wallet_id"`
	State           cybavo.CallbackState   `json:"state"`
	ConfirmBlocks   int64                  `json:"confirm_blocks"`
	ProcessingState cybavo.ProcessingState `json:"processing_state"`
	Decimals        int                    `json:"decimal"`
	SubWalletID     int64                  `json:"sub_wallet_id"`
	SubWallet       *SubWallet             `pg:"rel:has-one"`
}

func SetDepositCallBack(cb *cybavo.CallbackStruct) (err error) {
	db := GetDBClient()
	dcb := new(DepositCallBack)
	// deposit unique ID
	uniqueID := fmt.Sprintf("%s_%d", cb.TXID, cb.VOutIndex)
	copier.Copy(dcb, cb)

	subWallet, err := GetSubWalletByAddress(cb.ToAddress)

	if err != nil {
		log.Error("Failed to fetch subWallet ", err)
		return err
	}

	dcb.UniqueID = uniqueID
	dcb.SubWalletID = subWallet.Id

	count, err := db.Model((*DepositCallBack)(nil)).Where("unique_id = ?", uniqueID).Count()
	if count == 0 {
		_, err = db.Model(dcb).Insert()
		if err != nil {
			log.Error("Failed to insert deposit callback ", err)
			return err
		}
	} else {
		_, err = db.Model(dcb).WherePK().Update()
		if err != nil {
			log.Error("Failed to update deposit callback ", err)
			return err
		}
	}
	return nil
}

func SetWalletAPICode(walletRequest *Wallet, coin string) (err error) {
	db := GetDBClient()
	wallet := new(Wallet)
	err = db.Model(wallet).Where("wallet_id = ?", walletRequest.WalletID).Select()

	copier.Copy(wallet, walletRequest)
	wallet.CoinType = CoinTypes[coin]

	if err != nil {
		_, err = db.Model(wallet).Insert()
		if err != nil {
			log.Error("Failed to insert API token ", err)
			return err
		}
	} else {
		_, err = db.Model(wallet).WherePK().Update()
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
