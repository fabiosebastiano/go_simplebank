package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/fabiosebastiano/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var request createTransferRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validateAccount(ctx, request.FromAccountID, request.Currency) {
		return
	}
	if !server.validateAccount(ctx, request.ToAccountID, request.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: request.FromAccountID,
		ToAccountID:   request.ToAccountID,
		Amount:        request.Amount,
	}

	account, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) validateAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, "Account non trovato a DB")
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] ha currency diverse: %s Vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
