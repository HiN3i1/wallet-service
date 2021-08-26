package v1

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/HiN3i1/wallet-service/cybavo"
	"github.com/HiN3i1/wallet-service/db"
	"github.com/HiN3i1/wallet-service/types/api"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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

	apiCodeObj, err := db.GetWalletById(cb.WalletID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.NewServerError(err.Error()))
		return
	}

	checksum := c.Request.Header.Get("X-CHECKSUM")
	payload := string(postBody) + apiCodeObj.ApiSecret
	sha, _ := calcSHA256([]byte(payload))
	checksumVerf := base64.URLEncoding.EncodeToString(sha)

	if checksum != checksumVerf {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.NewServerError("bad checksum"))
		return
	}

	logrus.Debug("Callback => %s", postBody)

	cbType := cybavo.CallbackType(cb.Type)
	if cbType == cybavo.DepositCallback {
		//
		// deposit unique ID
		// uniqueID := fmt.Sprintf("%s_%d", cb.TXID, cb.VOutIndex)
		//
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
	c.JSON(http.StatusOK, api.RequestOK)
}
