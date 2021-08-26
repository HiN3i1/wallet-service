package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HiN3i1/wallet-service/cybavo"
	"github.com/HiN3i1/wallet-service/db"
	"github.com/HiN3i1/wallet-service/types/api"
	"github.com/gin-gonic/gin"
)

type CreateAddressResp struct {
	Addresses []string `json:"addresses"`
}

func getQueryString(params gin.Params) []string {
	var qs []string
	for _, param := range params {
		s := fmt.Sprintf("%s=%s", param.Key, param.Value)
		qs = append(qs, s)
	}
	return qs
}

func SetAPIToken(c *gin.Context) {
	coin := c.Param("coin")

	if coin == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.EmptyParam)
		return
	}

	var request cybavo.SetAPICodeRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.BadRequest)
		return
	}

	apiCodeParams := db.Wallet{
		APICode:    request.APICode,
		ApiSecret:  request.ApiSecret,
		WalletID:   request.WalletID,
		WalletType: db.Deposit,
	}
	err = db.SetWalletAPICode(&apiCodeParams, coin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.NewServerError(err.Error()))
	}

	c.JSON(http.StatusOK, api.RequestOK)
}

// CreateDepositWalletAddresses
func CreateDepositWalletAddresses(c *gin.Context) {
	customerID := c.Param("customer_id")
	coin := c.Param("coin")

	if customerID == "" || coin == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.EmptyParam)
		return
	}

	wallet, err := db.GetWalletByType(coin)

	if wallet == nil || err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.NewServerError(err.Error()))
		return
	}

	resp, err := cybavo.MakeRequest(wallet, "POST", fmt.Sprintf("/v1/sofa/wallets/%d/addresses", wallet.WalletID),
		nil, c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.NewServerError(err.Error()))
		return
	}

	var m CreateAddressResp
	json.Unmarshal(resp, &m)

	if len(m.Addresses) == 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.NewServerError(err.Error()))
		return
	}

	err = db.CreateCustomerWallet(customerID, m.Addresses[0], wallet)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.NewServerError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, m)
}

func GetDepositWalletAddresses(c *gin.Context) {
	customerID := c.Param("customer_id")
	coin := c.Param("coin")
	if customerID == "" || coin == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.EmptyParam)
		return
	}

	address, err := db.GetAddressByCustomerID(customerID, coin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.NewServerError(err.Error()))
		return
	}

	m := CreateAddressResp{
		Addresses: []string{address},
	}
	c.JSON(http.StatusOK, m)
}
