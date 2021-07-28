package controller

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/Rhymond/go-money"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	token "github.com/kevin/My-DeFi/contracts"
	"github.com/labstack/echo/v4"
)

type CommonResponse struct {
	Status bool
	Message string
}

type Response struct {
	Balance *big.Int
}

func formatJson(balance *big.Int) Response {
	response := &Response{
		Balance: balance,
	}

	return *response
}

// 관리자: 적립금 총액 확인
func GetBalance(c echo.Context, instance *token.Token) error {
	total, err := instance.TotalSupply(&bind.CallOpts{})

	if err != nil {
			return err
	}

	return c.JSONPretty(http.StatusOK, formatJson(total), " ")
}

// 사용자: 개인 적립금 총액 확인
func GetBalanceByAddr(c echo.Context, instance *token.Token, address string) error {
	addr := common.HexToAddress(address)
	balance, err := instance.BalanceOf(&bind.CallOpts{},addr)

	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, formatJson(balance), " ")
}

func formatCommonJson(status bool, message string) CommonResponse {
	response := &CommonResponse{
		Status: status,
		Message: message,
	}

	return *response
}

// 관리자: 이더리움 적립 시세 조정
func SetExchangeRate(c echo.Context, instance *token.Token, adminKey string, amount int64) error {
	
	privateKey, err := crypto.HexToECDSA(adminKey)
	if err != nil {
		log.Fatal(err)
	}

	transactionOpt := bind.NewKeyedTransactor(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	value := big.NewInt(amount)

	instance.TokenTransactor.SetExchangeRate(&bind.TransactOpts{
		From: fromAddress,
		Signer: transactionOpt.Signer,
	},value)

	return c.JSONPretty(http.StatusOK, formatCommonJson(true,"이더리움 적립 시세를 변경하였습니다."), " ")

}

// 관리자: 이더리움 조기 지급 시세 조정
func SetRevenueRate(c echo.Context, instance *token.Token, adminKey string, amount int64) error {
	privateKey, err := crypto.HexToECDSA(adminKey)
	if err != nil {
		log.Fatal(err)
	}

	transactionOpt := bind.NewKeyedTransactor(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	value := big.NewInt(amount)

	instance.TokenTransactor.SetRevenueRate(&bind.TransactOpts{
		From: fromAddress,
		Signer: transactionOpt.Signer,
	},value)

	return c.JSONPretty(http.StatusOK, formatCommonJson(true,"이더리움 조기적립 시세를 조정하였습니다."), " ")

}

// 관리자: 이더리움 장기 지급 시세 조정 (상품 개시 이후의 평균 값)
func SetAverageRate(c echo.Context, instance *token.Token, adminKey string, amount int64) error {
	privateKey, err := crypto.HexToECDSA(adminKey)
	if err != nil {
		log.Fatal(err)
	}

	transactionOpt := bind.NewKeyedTransactor(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	value := big.NewInt(amount)

	instance.TokenTransactor.SetAvgRate(&bind.TransactOpts{
		From: fromAddress,
		Signer: transactionOpt.Signer,
	},value)

	return c.JSONPretty(http.StatusOK, formatCommonJson(true,"이더리움 만기적립 시세를 조정하였습니다."), " ")

}

type RateResponse struct {
	Exchange string
	Revenue string
	Average string
}

// 관리자: 이더리움 설정 시세 목록 확인
func GetAllRates(c echo.Context, instance *token.Token, adminKey string) error { 
	privateKey, err := crypto.HexToECDSA(adminKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	res, err := instance.GetAllRates(&bind.CallOpts{
		From: fromAddress,
	}); 

	if(err != nil) {
		return err
	}

	fmt.Println(res)

	exchange := money.New(res[0].Int64(), "KRW")
	revenue := money.New(res[1].Int64(), "KRW")
	average := money.New(res[2].Int64(), "KRW")

	response := &RateResponse{
		Exchange:exchange.Display(),
		Revenue: revenue.Display(),
		Average: average.Display(),
	}

	return c.JSONPretty(http.StatusOK, response, " ")
}

// 관리자: 적립금 소유주 변경
func TransterToken(c echo.Context, instance *token.Token, client *ethclient.Client, fromKey string, from string, to string, amount int64 ) error {
	privateKey, err := crypto.HexToECDSA(fromKey)
	if err != nil {
		log.Fatal(err)
	}

	transactionOpt := bind.NewKeyedTransactor(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	value := big.NewInt(amount)
	gasLimit := uint64(3000000)
    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err != nil {
        log.Fatal(err)
    }
	toAddress := common.HexToAddress(to)

	res, err := instance.TokenTransactor.Transfer(&bind.TransactOpts{
		From: fromAddress,
		Signer: transactionOpt.Signer,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
	},toAddress,value)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
	
	return c.JSONPretty(http.StatusOK, formatCommonJson(true,"적립금 소유주 변경 신청이 완료 되었습니다."), " ")
}

// 사용자: 이더리움으로 환급형 보험료 납부
func BuyInsurance(c echo.Context, instance *token.Token, client *ethclient.Client, fromKey string, from string, amount int64) error {
	const WEI int64 = 1000000000000000000

	privateKey, err := crypto.HexToECDSA(fromKey)
	if err != nil {
		log.Fatal(err)
	}

	transactionOpt := bind.NewKeyedTransactor(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	value := big.NewInt(amount * WEI)
	gasLimit := uint64(3000000)
    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err != nil {
        log.Fatal(err)
    }


	res, err := instance.TokenTransactor.BuyInsurance(&bind.TransactOpts{
		From: fromAddress,
		Value: value,
		Signer: transactionOpt.Signer,
		GasPrice: gasPrice,
		GasLimit: gasLimit,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
	
	return c.JSONPretty(http.StatusOK, formatCommonJson(true,"이번달 보험료로 " + strconv.FormatInt(amount, 10) + "Eth 납부 하였습니다."), " ")
}

// 사용자: 적립금 조회 환급 신청
func RefundSavings(c echo.Context, instance *token.Token, client *ethclient.Client, pKey string, amount int64) error {

	privateKey, err := crypto.HexToECDSA(pKey)
	if err != nil {
		log.Fatal(err)
	}

	transactionOpt := bind.NewKeyedTransactor(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	value := big.NewInt(amount)
	gasLimit := uint64(3000000)
    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err != nil {
        log.Fatal(err)
    }

	res, err := instance.TokenTransactor.RefundSavings(&bind.TransactOpts{
		From: fromAddress,
		Signer: transactionOpt.Signer,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
	},value)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)

	won := money.New(amount, "KRW")

	return c.JSONPretty(http.StatusOK, formatCommonJson(true,won.Display()+"원의 적립된 이더리움을 인출하였습니다."), " ")
}

// 사용자: 만기 환급 신청
func RefundFinalSavings(c echo.Context, instance *token.Token, client *ethclient.Client,  pKey string) error {

	privateKey, err := crypto.HexToECDSA(pKey)
	if err != nil {
		log.Fatal(err)
	}

	transactionOpt := bind.NewKeyedTransactor(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	gasLimit := uint64(3000000)
    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err != nil {
        log.Fatal(err)
    }

	res, err := instance.TokenTransactor.RefundFinalSavings(&bind.TransactOpts{
		From: fromAddress,
		Signer: transactionOpt.Signer,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)

	return c.JSONPretty(http.StatusOK, formatCommonJson(true,"만기 환급금 신청을 완료하였습니다."), " ")
}