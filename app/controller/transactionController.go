package controller

import (
	"fmt"
	"log"
	"net/http"
	"ppob-backend/app/dto"
	"ppob-backend/app/services/transactionServices"
	"ppob-backend/helper"
	"strconv"

	"github.com/labstack/echo/v4"
)

type (
	transactionController struct {
		txnServices transactionServices.TransactionServices
	}

	TransactionController interface {
		GetBalance(ctx echo.Context) error
		TopupWallet(ctx echo.Context) error
		BuyPulsa(ctx echo.Context) error
		Webhook(ctx echo.Context) error
		TransactionHistory(ctx echo.Context) error
	}
)

func NewTransactionController(txnServices transactionServices.TransactionServices) TransactionController {
	return &transactionController{
		txnServices: txnServices,
	}
}

func (c *transactionController) GetBalance(ctx echo.Context) error {
	var (
		request dto.GetBalanceRequest
	)

	err := ctx.Bind(&request)

	if err != nil {
		ctx.Logger().Warnf("Error on binding: %v", err.Error())
		return echo.ErrBadRequest
	}

	response, err := c.txnServices.GetUserBalanceByUserID(request)

	if err != nil {
		ctx.Logger().Warnf("GetBalance controller: %v", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Something went wrong: %v", err.Error()))
	}

	responseBuild := helper.BuildResponse(http.StatusOK, "success", response)

	return ctx.JSON(http.StatusOK, responseBuild)
}

func (c *transactionController) TopupWallet(ctx echo.Context) error {
	var (
		request dto.TopupWalletRequest
	)

	err := ctx.Bind(&request)

	if err != nil {
		ctx.Logger().Warnf("Error on binding: %v", err.Error())
		return echo.ErrBadRequest
	}

	response, err := c.txnServices.TopupWallet(request)

	if err != nil {
		ctx.Logger().Warnf("TopupWallet controller: %v", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Something went wrong: %v", err.Error()))
	}

	responseBuild := helper.BuildResponse(http.StatusOK, "success", response)

	return ctx.JSON(http.StatusOK, responseBuild)
}

func (c *transactionController) BuyPulsa(ctx echo.Context) error {
	var (
		request dto.BuyPulsaRequest
	)

	err := ctx.Bind(&request)

	if err != nil {
		ctx.Logger().Warnf("Error on binding: %v", err.Error())
		return echo.ErrBadRequest
	}

	user_id, _ := strconv.ParseUint(request.UserID, 10, 64)

	response, err := c.txnServices.DoTransaction(request.ProductCode, request.CustomerNo, uint(user_id))

	if err != nil {
		ctx.Logger().Warnf("TopupWallet controller: %v", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Something went wrong: %v", err.Error()))
	}

	responseBuild := helper.BuildResponse(http.StatusOK, "success", response)

	return ctx.JSON(http.StatusOK, responseBuild)
}

func (c *transactionController) TransactionHistory(ctx echo.Context) error {
	var (
		request dto.TransactionHistoryRequest
	)

	err := ctx.Bind(&request)

	if err != nil {
		ctx.Logger().Warnf("Error on binding: %v", err.Error())
		return echo.ErrBadRequest
	}

	user_id, _ := strconv.ParseUint(request.UserID, 10, 64)

	response, err := c.txnServices.GetTransactionHistory(user_id)

	if err != nil {
		ctx.Logger().Warnf("TransactionHistory controller: %v", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Something went wrong: %v", err.Error()))
	}

	responseBuild := helper.BuildResponse(http.StatusOK, "success", response)

	return ctx.JSON(http.StatusOK, responseBuild)
}

func (c *transactionController) Webhook(ctx echo.Context) error {
	var (
		request dto.WebhookRequest
		headers dto.WebhookHeaders
	)

	err := ctx.Bind(&request)

	if err != nil {
		log.Printf("Error on binding: %v", err.Error())
		ctx.Logger().Warnf("Error on binding: %v", err.Error())
		return echo.ErrBadRequest
	}

	binder := &echo.DefaultBinder{}

	err = binder.BindHeaders(ctx, &headers)

	if err != nil {
		log.Printf("Error on binding: %s", err.Error())
		ctx.Logger().Warnf("Error on binding: %v", err.Error())
		return echo.ErrBadRequest
	}

	body := dto.WebhookRequestBody(request.Data)

	err = c.txnServices.WebhookHandle(headers, body)

	if err != nil {
		log.Printf("Error on handle: %v", err.Error())
		ctx.Logger().Warnf("Error on handle: %v", err)
		return echo.ErrBadRequest
	}

	return ctx.JSON(http.StatusOK, helper.EmptyResponse{})
}
