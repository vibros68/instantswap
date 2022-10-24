package main

import (
	"code.cryptopower.dev/exchange/lightningswap"
	_ "code.cryptopower.dev/exchange/lightningswap/exchange/flypme"
	"fmt"
	"os"
)

func main() {
	exchange, err := lightningswap.NewExchange("flypme", lightningswap.ExchangeConfig{
		Debug:     false,
		ApiKey:    "",
		ApiSecret: "",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	res, err := exchange.GetExchangeRateInfo(lightningswap.ExchangeRateRequest{
		From:   "BTC",
		To:     "DCR",
		Amount: 5,
	})
	fmt.Printf("%+v \n %v \n", res, err)
	order, err := exchange.CreateOrder(lightningswap.CreateOrder{
		RefundAddress:   "your_btc_address", // if the trading fail, the exchange will refund here
		Destination:     "your_dcr_address", // your received dcr address
		FromCurrency:    "BTC",
		OrderedAmount:   0, // use OrderedAmount or InvoicedAmount
		InvoicedAmount:  0.5,
		ToCurrency:      "DCR",
		ExtraID:         "",
		Signature:       res.Signature,
		UserReferenceID: "",
		RefundExtraID:   "",
	})
	fmt.Println(order, err)

	// the exchange will return the rate of exchange is: order.ExchangeRate
	// you will send btc to order.DepositAddress
	// use OrderInfo to get order status
	orderInfo, err := exchange.OrderInfo(order.UUID)
	fmt.Println(orderInfo, err)

	fmt.Println(orderInfo.InternalStatus.String())
	// when ever the trading is done, the exchange will return the transaction id in
	// orderInfo.TxID
}
