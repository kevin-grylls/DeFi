package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	token "github.com/kevin/My-DeFi/contracts"
	controller "github.com/kevin/My-DeFi/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// const ROPSTEN_NETWORK string = "http://127.0.0.1:7545"
// const CONTRACT_ADDR string = "0xfE32d78f7Bcaaf647D83096f7f5f9D21fbECDBC7"
const ROPSTEN_NETWORK = "https://ropsten.infura.io/v3/4d51ae5e895d4d73974e078aacff23aa"
const CONTRACT_ADDR string = "0xEb0422465F7B484187BaF716795F74D9F86aB027"

func main() {
	client, err := ethclient.Dial(ROPSTEN_NETWORK)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(client)

	tokenAddress := common.HexToAddress(CONTRACT_ADDR)
	instance, err := token.NewToken(tokenAddress, client)

	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/balance", func(c echo.Context) error {
		return controller.GetBalance(c, instance)
	})

	e.GET("/balance/:addr", func(c echo.Context) error {

		addr := c.Param("addr")

		return controller.GetBalanceByAddr(c, instance, addr)
	})

	type RateRequest struct {
		PrivateKey string
		Rate int64
	}

	e.PUT("/rate/exchange", func(c echo.Context) error{
		body := new(RateRequest)

		if err := c.Bind(body); err != nil {
			return err
		}

		return controller.SetExchangeRate(c, instance, body.PrivateKey, body.Rate)
	})

	e.PUT("/rate/revenue", func(c echo.Context) error{
		body := new(RateRequest)

		if err := c.Bind(body); err != nil {
			return err
		}

		return controller.SetRevenueRate(c, instance, body.PrivateKey, body.Rate)
	})

	e.PUT("/rate/average", func(c echo.Context) error{
		body := new(RateRequest)

		if err := c.Bind(body); err != nil {
			return err
		}

		return controller.SetAverageRate(c, instance, body.PrivateKey, body.Rate)
	})

	type RateAllRequest struct {
		PrivateKey string
	}

	e.PUT("/rates/all", func(c echo.Context) error {
		body := new(RateAllRequest)

		if err := c.Bind(body); err != nil {
			return err
		}

		return controller.GetAllRates(c, instance, body.PrivateKey)
	})

	type TransferRequest struct {
		PrivateKey string
		From string 
		To string
		Amount int64
	}

	e.PUT("/transfer", func(c echo.Context) error {
		body := new(TransferRequest)

		if err := c.Bind(body); err != nil {
			return err
		}

		return controller.TransterToken(c, instance, client, body.PrivateKey, body.From, body.To, body.Amount)
	})

	type BuyRequest struct {
		PrivateKey string
		From string
		Amount int64
	}

	e.POST("/buy", func(c echo.Context) error {
		body := new(BuyRequest)

		if err := c.Bind(body); err != nil {
			return err
		}

		return controller.BuyInsurance(c, instance, client, body.PrivateKey, body.From, body.Amount)
	})

	type ExchangeRequest struct {
		PrivateKey string
		Amount int64
	}

	e.POST("/refund/early", func(c echo.Context) error {
		body := new(ExchangeRequest)

		if err := c.Bind(body); err != nil {
			return err
		}

		return controller.RefundSavings(c, instance, client, body.PrivateKey, body.Amount)
	})

	type RefundRequest struct {
		PrivateKey string
	}

	e.POST("/refund/final", func(c echo.Context) error {
		body := new(RefundRequest)

		if err := c.Bind(body); err != nil {
			return err
		}

		return controller.RefundFinalSavings(c, instance, client, body.PrivateKey)
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}