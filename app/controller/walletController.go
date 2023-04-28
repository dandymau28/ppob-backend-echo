package controller

import (
	"fmt"
	"net/http"
	"ppob-backend/app/dto"
	"ppob-backend/app/services/walletServices"
	"ppob-backend/helper"

	"github.com/labstack/echo/v4"
)

type (
	walletController struct {
		walletServices walletServices.WalletServices
	}

	WalletController interface {
		GetBalance(ctx echo.Context) error
		TopupWallet(ctx echo.Context) error
	}
)

func NewWalletController(walletServices walletServices.WalletServices) WalletController {
	return &walletController{
		walletServices: walletServices,
	}
}

func (c *walletController) GetBalance(ctx echo.Context) error {
	var (
		request dto.GetBalanceRequest
	)

	err := ctx.Bind(&request)

	if err != nil {
		ctx.Logger().Warnf("Error on binding: %v", err.Error())
		return echo.ErrBadRequest
	}

	response, err := c.walletServices.GetUserBalanceByUserID(request)

	if err != nil {
		ctx.Logger().Warnf("GetBalance controller: %v", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Something went wrong: %v", err.Error()))
	}

	responseBuild := helper.BuildResponse(http.StatusOK, "success", response)

	return ctx.JSON(http.StatusOK, responseBuild)
}

func (c *walletController) TopupWallet(ctx echo.Context) error {
	var (
		request dto.TopupWalletRequest
	)

	err := ctx.Bind(&request)

	if err != nil {
		ctx.Logger().Warnf("Error on binding: %v", err.Error())
		return echo.ErrBadRequest
	}

	response, err := c.walletServices.TopupWallet(request)

	if err != nil {
		ctx.Logger().Warnf("TopupWallet controller: %v", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Something went wrong: %v", err.Error()))
	}

	responseBuild := helper.BuildResponse(http.StatusOK, "success", response)

	return ctx.JSON(http.StatusOK, responseBuild)
}
