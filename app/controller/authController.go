package controller

import (
	"fmt"
	"net/http"
	"ppob-backend/app/dto"
	authService "ppob-backend/app/services/authServices"
	"ppob-backend/config"
	"ppob-backend/helper"

	"github.com/labstack/echo/v4"
)

type (
	authController struct {
		Config      *config.SystemConfig
		authService authService.AuthService
	}

	AuthController interface {
		Login(ctx echo.Context) error
		Register(ctx echo.Context) error
	}
)

func NewAuthController(config *config.SystemConfig, authService authService.AuthService) AuthController {
	return &authController{
		Config:      config,
		authService: authService,
	}
}

func (c *authController) Login(ctx echo.Context) error {
	var (
		request dto.LoginRequest
	)

	err := ctx.Bind(&request)

	if err != nil {
		ctx.Logger().Warnf("Error on binding: %v", err.Error())
		return echo.ErrBadRequest
	}

	response, err := c.authService.Login(request)

	if err != nil {
		ctx.Logger().Warnf("login controller: %v", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Something went wrong: %v", err.Error()))
	}

	responseBuild := helper.BuildResponse(http.StatusOK, "success", response)

	return ctx.JSON(http.StatusOK, responseBuild)
}

func (c *authController) Register(ctx echo.Context) error {
	var (
		request dto.RegisterRequest
	)

	err := ctx.Bind(&request)

	if err != nil {
		ctx.Logger().Warnf("Error on binding: %v", err.Error())
		return echo.ErrBadRequest
	}

	if err = ctx.Validate(request); err != nil {
		return err
	}

	response, err := c.authService.Register(request)

	if err != nil {
		ctx.Logger().Warnf("register controller: %v", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Something went wrong: %v", err.Error()))
	}

	responseBuild := helper.BuildResponse(http.StatusOK, "success", response)

	return ctx.JSON(http.StatusOK, responseBuild)
}
