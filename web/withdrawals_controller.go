package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
)

// WithdrawalsController can send LINK tokens to another address
type WithdrawalsController struct {
	App services.Application
}

var naz = assets.NewLink(1)

// Create sends LINK from the configured oracle contract to the given address
// See models.WithdrawalRequest for expected inputs. ContractAddress field is
// optional.
//
// Example: "<application>/withdrawals"
func (abc *WithdrawalsController) Create(c *gin.Context) {
	store := abc.App.GetStore()
	txm := store.TxManager
	wr := models.WithdrawalRequest{}

	if err := c.ShouldBindJSON(&wr); err != nil {
		publicError(c, http.StatusBadRequest, err)
		return
	}

	addressWasNotSpecifiedInRequest := utils.ZeroAddress

	if wr.Amount.Cmp(naz) < 0 {
		publicError(c, http.StatusBadRequest, fmt.Errorf(
			"Must withdraw at least %v LINK", naz.String()))
	} else if wr.DestinationAddress == addressWasNotSpecifiedInRequest {
		publicError(c, http.StatusBadRequest, errors.New("Invalid withdrawal address"))
	} else if linkBalance, err := txm.ContractLINKBalance(wr); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	} else if linkBalance.Cmp(wr.Amount) < 0 {
		publicError(c, http.StatusBadRequest, fmt.Errorf(
			"Insufficient link balance. Withdrawal Amount: %v "+
				"Link Balance: %v",
			wr.Amount.String(), linkBalance.String()))
		return
	}

	hash, err := txm.WithdrawLINK(wr)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, hash)
	}
}
