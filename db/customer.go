package db

import (
	"github.com/go-pg/pg/orm"
	log "github.com/sirupsen/logrus"
)

// Customer is using saving user's data
type Customer struct {
	Id         int64
	CustomerID string       `pg:",pk" json:"customer_id"`
	SubWallets []*SubWallet `pg:"rel:has-many,fk:customer_id"`
}

func CreateCustomer(customerID string) (customer *Customer, err error) {
	db := GetDBClient()
	customer = new(Customer)
	customer.CustomerID = customerID
	_, err = db.Model(customer).Insert()
	if err != nil {
		log.Error("Failed to create customer ", err)
		return nil, err
	}
	return customer, nil
}

func CreateCustomerWallet(customerID string, address string, wallet *Wallet) (err error) {
	db := GetDBClient()
	customer := new(Customer)
	err = db.Model(customer).Where("customer_id = ?", customerID).Select()
	if err != nil {
		customer, err = CreateCustomer(customerID)
		if err != nil {
			log.Error("Failed to create customer ", err)
			return err
		}
	}

	subWallet := &SubWallet{
		CustomerID: customer.Id,
		Address:    address,
		Memo:       address,
		WalletID:   wallet.Id,
	}

	_, err = db.Model(subWallet).Insert()
	if err != nil {
		log.Error("Failed to create subWallet ", err)
		return err
	}
	return nil
}

func GetSubWalletByCustomerID(customerID string, coin string) (subWallet *SubWallet, err error) {
	db := GetDBClient()

	customer := new(Customer)
	err = db.Model(customer).Where("customer_id = ?", customerID).Select()
	if err != nil {
		log.Error("Failed to fetch customer ", err)
		return nil, err
	}

	subWallet = new(SubWallet)
	coinType := CoinTypes[coin]
	err = db.Model(subWallet).
		Where("customer_id = ?", customer.Id).
		Relation("Wallet", func(q *orm.Query) (*orm.Query, error) {
			return q.Where("coin_type = ?", coinType), nil
		}).First()

	return subWallet, err
}

func GetWalletByCustomerID(customerId string, coin string) (wallet *Wallet, err error) {
	subWallet, err := GetSubWalletByCustomerID(customerId, coin)
	if subWallet != nil {
		wallet = subWallet.Wallet
	}
	return wallet, err
}

func GetAddressByCustomerID(customerId string, coin string) (address string, err error) {
	subWallet, err := GetSubWalletByCustomerID(customerId, coin)
	if subWallet != nil {
		address = subWallet.Address
	}
	return address, err
}
