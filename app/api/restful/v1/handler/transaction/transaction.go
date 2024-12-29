package transaction

import (
	"errors"
	"fmt"
	"net/http"

	v1 "banking/app/api/restful/v1"
	transactionRepo "banking/app/repo/mysql/transaction"
	 "banking/domain"
	"banking/global"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.elastic.co/apm/module/apmzap/v2"
	"go.elastic.co/apm/v2"
)

type TransactionHandler struct {
	transactionService domain.ITransactionService
}

func NewTransactionHandler(TransactionService domain.ITransactionService) domain.ITransactionHandler {
	return &TransactionHandler{
		transactionService: TransactionService,
	}
}

func (h *TransactionHandler) Transfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "TransactionHandler.Transfer", "handler")
		defer span.End()

		var input struct {
			FromUserID uint    `json:"fromUserId" binding:"required,min=1,number"`
			ToUserID   uint    `json:"toUserId" binding:"required,min=1,number"`
			Amount     float64 `json:"amount" binding:"required,gt=0,number"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		if input.FromUserID == input.ToUserID {
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: "fromUserId and toUserId should not be the same",
			})
			return
		}

		user, err := h.transactionService.Transfer(ctx, input.FromUserID, input.ToUserID, decimal.NewFromFloat(input.Amount))
		if err != nil {
			if errors.Is(err, transactionRepo.ErrInsufficientBalance) {
				apm.CaptureError(ctx, err).Send()
				c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
					Msg: transactionRepo.ErrInsufficientBalance.Error(),
				})
				return
			}

			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		global.Logger.With(apmzap.TraceContext(ctx)).Info(fmt.Sprintf("[Transfer]: fromUserId: %d, toUserId: %d, amount: %f", input.FromUserID, input.ToUserID, input.Amount))

		c.JSON(http.StatusOK, &TransferResp{
			Data: user,
		})
	}
}

func (h *TransactionHandler) Deposit() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "TransactionHandler.Deposit", "handler")
		defer span.End()

		var input struct {
			UserID uint    `json:"userId" binding:"required,min=1,number"`
			Amount float64 `json:"amount" binding:"required,gt=0,number"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		user, err := h.transactionService.Deposit(ctx, input.UserID, decimal.NewFromFloat(input.Amount))
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		global.Logger.With(apmzap.TraceContext(ctx)).Info(fmt.Sprintf("[Deposit]: userId: %d, amount: %f", input.UserID, input.Amount))

		c.JSON(http.StatusOK, &DepositResp{
			Data: user,
		})
	}
}

func (h *TransactionHandler) Withdraw() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "TransactionHandler.Withdraw", "handler")
		defer span.End()

		var input struct {
			UserID uint    `json:"userId" binding:"required,min=1,number"`
			Amount float64 `json:"amount" binding:"required,gt=0,number"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		user, err := h.transactionService.Withdraw(ctx, input.UserID, decimal.NewFromFloat(input.Amount))
		if err != nil {
			if errors.Is(err, transactionRepo.ErrInsufficientBalance) {
				apm.CaptureError(ctx, err).Send()
				c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
					Msg: transactionRepo.ErrInsufficientBalance.Error(),
				})
				return
			}

			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		global.Logger.With(apmzap.TraceContext(ctx)).Info(fmt.Sprintf("[Withdraw]: userId: %d, amount: %f", input.UserID, input.Amount))

		c.JSON(http.StatusOK, &WithdrawResp{
			Data: user,
		})
	}
}

func (h *TransactionHandler) GetTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "TransactionHandler.GetTransactions", "handler")
		defer span.End()

		transactions, err := h.transactionService.GetTransactions(ctx)
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		global.Logger.With(apmzap.TraceContext(ctx)).Info(fmt.Sprintf("[GetTransactions]: transactions: %v", transactions))

		c.JSON(http.StatusOK, &GetTransactionsResp{
			Data: transactions,
		})
	}
}
