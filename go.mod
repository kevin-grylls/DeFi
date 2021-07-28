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
