package utils

import "github.com/shopspring/decimal"

const (
	INTERNAL_SERVER_ERROR        = "internal server error"
	NOT_FOUND_ERROR              = "data not found"
	BAD_REQUEST                  = "bad request"
	STATUS_INTERNAL_SERVER_ERROR = 500
)

func AddDeposit(balance decimal.Decimal, amount decimal.Decimal) decimal.Decimal {
	var result decimal.Decimal

	result = balance.Add(amount)

	return result
}
