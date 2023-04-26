package main

import (
	"net/http"
	"ppob-backend/app/controller"
	"ppob-backend/app/repository/authRepository"
	"ppob-backend/app/repository/transactionRepository"
	authService "ppob-backend/app/services/authServices"
	dfwebclientservices "ppob-backend/app/services/dfWebClientServices"
	"ppob-backend/app/services/transactionServices"
	"ppob-backend/config"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

var (
	configSystem   *config.SystemConfig                        = config.LoadEnv()
	dbConn         *gorm.DB                                    = config.ConnectDB(configSystem)
	dfWebClient    dfwebclientservices.DfWebClient             = dfwebclientservices.NewDfWebClient(configSystem)
	authRepo       authRepository.AuthRepository               = authRepository.NewAuthRepository(dbConn)
	authServ       authService.AuthService                     = authService.NewAuthService(configSystem, authRepo)
	authController controller.AuthController                   = controller.NewAuthController(configSystem, authServ)
	txnRepo        transactionRepository.TransactionRepository = transactionRepository.NewTransactionRepository(dbConn)
	txnServ        transactionServices.TransactionServices     = transactionServices.NewTransactionServices(configSystem, txnRepo, dfWebClient)
	txnController  controller.TransactionController            = controller.NewTransactionController(txnServ)
)

type (
	CustomValidator struct {
		validator *validator.Validate
	}

	jwtCustomClaims struct {
		Username string `json:"username"`
		jwt.StandardClaims
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api")

	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte(configSystem.JwtSecret),
	}

	//webhook handler
	api.POST("/webhook", txnController.Webhook)

	// insert auth route here
	authRoute := api.Group("/auth")
	authRoute.POST("/login", authController.Login)
	authRoute.POST("/register", authController.Register)

	//insert any route here
	v1 := api.Group("/v1")
	v1.Use(middleware.JWTWithConfig(config))

	trxRoute := v1.Group("/transactions")
	trxRoute.GET("/:user_id/wallet", txnController.GetBalance)
	trxRoute.POST("/:user_id/wallet/topup", txnController.TopupWallet)
	trxRoute.POST("/buy-pulsa", txnController.BuyPulsa)
	trxRoute.GET("/:user_id/history", txnController.TransactionHistory)

	e.Logger.Fatal(e.Start(":1323"))
}
