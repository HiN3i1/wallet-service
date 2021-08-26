package v1

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/HiN3i1/wallet-service/cybavo"
	"github.com/HiN3i1/wallet-service/db"
	"github.com/HiN3i1/wallet-service/types/api"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func ActivateAPIToken(c *gin.Context) {

	walletID, err := strconv.ParseInt(c.Param("wallet_id"), 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.EmptyParam)
		return
	}

	var url string
	if walletID == 0 {
		url = "/v1/sofa/wallets/readonly/apisecret/activate"
	} else {
		url = fmt.Sprintf("/v1/sofa/wallets/%d/apisecret/activate", walletID)
	}

	wallet, err := db.GetWalletById(walletID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.NewServerError(err.Error()))
		return
	}
	resp, err := cybavo.MakeRequest(wallet.APICode, wallet.APISecret, "POST", url, nil, c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.NewServerError(err.Error()))
		return
	}

	var m map[string]interface{}
	json.Unmarshal(resp, &m)
	c.JSON(http.StatusOK, m)
}

func calcSHA256(data []byte) (calculatedHash []byte, err error) {
	sha := sha256.New()
	_, err = sha.Write(data)
	if err != nil {
		return
	}
	calculatedHash = sha.Sum(nil)
	return
}

func Callback(c *gin.Context) {

	var cb cybavo.CallbackStruct

	postBody, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.NewServerError(err.Error()))
		return
	}

	err = json.Unmarshal(postBody, &cb)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.NewServerError(err.Error()))
		return
	}

	apiCodeObj, err := db.GetWalletByType(cb.Currency)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.NewServerError(err.Error()))
		return
	}

	checksum := c.Request.Header.Get("X-CHECKSUM")
	payload := string(postBody) + apiCodeObj.APISecret
	sha, _ := calcSHA256([]byte(payload))
	checksumVerf := base64.URLEncoding.EncodeToString(sha)

	if checksum != checksumVerf {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.NewServerError("bad checksum"))
		return
	}

	log.Debug("Callback => %s", string(postBody))

	cbType := cybavo.CallbackType(cb.Type)
	if cbType == cybavo.DepositCallback {

		// deposit unique ID
		err = db.SetDepositCallBack(&cb)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, api.NewServerError(err.Error()))
			return
		}

		if cb.ProcessingState == cybavo.ProcessingStateDone {
			// deposit succeeded, use the deposit unique ID to update your business logic
		}
	} else if cbType == cybavo.WithdrawCallback {
		//
		// withdrawal unique ID
		// uniqueID := cb.OrderID
		//
		if cb.State == cybavo.CallbackStateInChain && cb.ProcessingState == cybavo.ProcessingStateDone {
			// withdrawal succeeded, use the withdrawal uniqueID to update your business logic
		} else if cb.State == cybavo.CallbackStateFailed || cb.State == cybavo.CallbackStateInChainFailed {
			// withdrawal failed, use the withdrawal unique ID to update your business logic
		}
	} else if cbType == cybavo.AirdropCallback {
		//
		// airdrop unique ID
		// uniqueID := fmt.Sprintf("%s_%d", cb.TXID, cb.VOutIndex)
		//
		if cb.ProcessingState == cybavo.ProcessingStateDone {
			// airdrop succeeded, use the airdrop unique ID to update your business logic
		}
	}

	// reply 200 OK to confirm the callback has been processed
	c.JSON(http.StatusOK, "OK")
}
