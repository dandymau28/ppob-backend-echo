package controller

import (
	"fmt"
	"log"
	"net/http"
	"ppob-backend/app/dto"
	"ppob-backend/app/services/transactionServices"
	"ppob-backend/helper"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type (
	transactionController struct {
		txnServices transactionServices.TransactionServices
	}

	TransactionController interface {
		BuyPulsa(ctx echo.Context) error
		PrePurchase(ctx echo.Context) error
		Webhook(ctx echo.Context) error
		TransactionHistory(ctx echo.Context) error
	}
)

func NewTransactionController(txnServices transactionServices.TransactionServices) TransactionController {
	return &transactionController{
		txnServices: txnServices,
	}
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

	user_id := request.UserID

	response, err := c.txnServices.DoTransaction(request.ProductCode, request.CustomerNo, user_id)

	if err != nil {
		ctx.Logger().Warnf("TopupWallet controller: %v", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Something went wrong: %v", err.Error()))
	}

	responseBuild := helper.BuildResponse(http.StatusOK, "success", response)

	return ctx.JSON(http.StatusOK, responseBuild)
}

func (c *transactionController) PrePurchase(ctx echo.Context) error {
	var (
		request dto.PrePurchaseRequest
	)

	err := ctx.Bind(&request)

	if err != nil {
		ctx.Logger().Warnf("Error on binding: %v", err.Error())
		return echo.ErrBadRequest
	}

	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*dto.JwtCustomClaims)
	user_id := claims.UserID

	response, err := c.txnServices.PrePurchase(request.ProductCode, request.CustomerNo, user_id)

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

	user_id := request.UserID

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
