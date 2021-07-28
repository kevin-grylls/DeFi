# DeFi

### What does it means?

- Decentralized Finance
- 블록체인 네트워크에서 동작하는 탈중앙화된 금융 어플리케이션

### How to use

- 로컬 golang 개발 디렉토리에 소스코드를 클론합니다.
- go.mod 파일의 상대경로를 수정합니다.

```go
// go.mod

module github.com/kevin/My-DeFi

go 1.16

require (
	github.com/ethereum/go-ethereum v1.10.6
	github.com/kevin/My-DeFi/contracts v0.0.0
	github.com/kevin/My-DeFi/controllers v0.0.0
	github.com/labstack/echo/v4 v4.4.0
)

replace (
	github.com/kevin/My-DeFi/contracts v0.0.0 => ./contracts
	github.com/kevin/My-DeFi/controllers v0.0.0 => ./controllers
)

```

- 이후 다음의 커맨드로 의존성을 설치하고 실행합니다.

```sh
// golang 의존성 설치
go mod init
go mod tidy
go run main.go
```

- Token.go 생성시 필요할 시 아래의 커멘드를 이용합니다.

```sh
// token 생성
solc --abi ./contracts/insurance.sol -o ./contracts/dist
cd ./contracts
abigen --abi=./dist/InsuranceToken.abi --pkg=token --out=token.go
```

### About Smart Contract

- [이 주소](https://ropsten.etherscan.io/address/0xeb0422465f7b484187baf716795f74d9f86ab027)에서 스마트 컨트랙트를 확인할 수 있습니다.

### Use Case

- 직장인 A 씨는 DeFi 결합형 만기환급형 암보험에 가입하였습니다.
- 정기 납입금의 일부를 이더리움으로 납입할 수 있는 상품이며,
- 이더 -> 현금 전환 비율은 보험사 기준에 의해 정기적으로 업데이트 됩니다.
- 적립된 현금과 마찬가지로 이더리움도 만기 환급 받을 수 있으며,
- 적립 시기의 이더리움 가격 변동을 평균하여, 현금 + 이더리움이 20년 만기 지급됩니다.
- 만기 이전 보험급 환급시 이더리움은 기준 전환 비율의 0.005% 수수료를 지급합니다.
- 이 상품은 배우자 및 자녀에게 납입금을 상속하거나 증여할 수 있습니다.

### Admin APIs

- 보험 납입액 총액 조회

```json
http://localhost:1323/balance [GET]
{
  "Balance": 33000000 // 잔액
}
```

- 보험 납입액 이전

```json
http://localhost:1323/transfer [PUT]
{
    "PrivateKey":"5872d1dffdb62f86f1565947b7121d5f4e697a9b964ee10d812f0ec040fcc210", // 보내는 사용자의 Private Key
    "From":"0x95306b3CF6c0d087ec4C95469750cD455e6074B1", // 보내는 사용자
    "To":"0x2A26311AcdCEAaE71B18C376e0D2143665e28B0B", // 받는 사용자
    "Amount":4000000 // 납입금, 40만원
}

// response
{
    "Status": true,
    "Message": "적립금 소유주 변경 신청이 완료 되었습니다."
}
```

- 설정된 이더리움 시세 조회

```json
http://localhost:1323/rates/all [PUT]
{
  "PrivateKey": "adf8a28c37d92511d5c67cc42bd775bc5a6b23aa0c03d83a55c3baf78d2e308c"
}

// response
{
    "Exchange": "₩2,600,000",
    "Revenue": "₩2,613,000",
    "Average": "₩2,400,000"
}
```

- 이더리움 기준 시세 변경

```json
http://localhost:1323/rate/exchange [PUT]
{
    "PrivateKey": "adf8a28c37d92511d5c67cc42bd775bc5a6b23aa0c03d83a55c3baf78d2e308c",
    "Rate": 26000000
}

// response
{
    "Status": true,
    "Message": "이더리움 적립 시세를 변경하였습니다."
}

http://localhost:1323/rate/average [PUT]
{
    "PrivateKey": "adf8a28c37d92511d5c67cc42bd775bc5a6b23aa0c03d83a55c3baf78d2e308c",
    "Rate": 26000000
}

// response
{
    "Status": true,
    "Message": "이더리움 만기적립 시세를 조정하였습니다."
}

http://localhost:1323/rate/revenue [PUT]
{
    "PrivateKey": "adf8a28c37d92511d5c67cc42bd775bc5a6b23aa0c03d83a55c3baf78d2e308c",
    "Rate": 26000000
}

// response
{
    "Status": true,
    "Message": "이더리움 조기적립 시세를 조정하였습니다."
}


```

### User APIs

- 개인 고객 이더리움 적립금 확인

```json
http://localhost:1323/balance/:addr [GET]

// response
{
    "Balance": 2600000
}

```

- 이더리움으로 보험료 납부

```json
http://localhost:1323/buy [POST]
{
  "From": "0x95306b3CF6c0d087ec4C95469750cD455e6074B1", // 사용자 지갑주소
  "PrivateKey": "5872d1dffdb62f86f1565947b7121d5f4e697a9b964ee10d812f0ec040fcc210", // 프라이빗 키
  "Amount": 5 // 이더 (5, 0.01, 10...)
}

// response
{
    "Status": true,
    "Message": "이번달 보험료로 5Eth 납부 하였습니다."
}
```

- 이더리움 적립급 조기 지급 신청

```json
http://localhost:1323/refund/early [POST]
{
  "PrivateKey": "5872d1dffdb62f86f1565947b7121d5f4e697a9b964ee10d812f0ec040fcc210",
  "Amount": 3000000 // 적립금액
}

// response
{
    "Status": true,
    "Message": "₩3,000,000원의 적립된 이더리움을 인출하였습니다."
}
```

- 이더리움 만기 지급 신청

```json
http://localhost:1323/refund/final [POST]
{
    "PrivateKey":"5872d1dffdb62f86f1565947b7121d5f4e697a9b964ee10d812f0ec040fcc210"
}

// response
{
    "Status": true,
    "Message": "만기 환급금 신청을 완료하였습니다."
}
```
