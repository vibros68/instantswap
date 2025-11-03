package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vibros68/instantswap/instantswap"
	_ "github.com/vibros68/instantswap/instantswap/index"
)

// using this code:
// go run ./_example/exolix/ -exchange=exolix -key=4WRyccNKW***yccNKW

var (
	exchange, apiKey, apiSecret string
)

func init() {
	flag.StringVar(&exchange, "exchange", "exolix", "-exchange=<the exchange you want to trade, default is exolix>")
	flag.StringVar(&apiKey, "key", "", "-key=<your api key>")
	flag.Parse()
}

func main() {
	exchange, err := instantswap.NewExchange(exchange, instantswap.ExchangeConfig{
		Debug:  false,
		ApiKey: apiKey,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	currencies, err := exchange.GetCurrencies()
	fmt.Println("currencies length: ", len(currencies), err)
	currencies, err = exchange.GetCurrenciesToPair("USDC")
	fmt.Println("GetCurrenciesToPair length: ", len(currencies), err)
	res, err := exchange.GetExchangeRateInfo(instantswap.ExchangeRateRequest{
		From:        "USDC",
		To:          "USDT",
		FromNetwork: "BSC",
		ToNetwork:   "BSC",
		Amount:      50,
	})
	fmt.Printf("%+v \n %v \n", res, err)
	order, err := exchange.CreateOrder(instantswap.CreateOrder{
		RefundAddress:   "refund address",        // if the trading fail, the exchange will refund here
		Destination:     "received usdt address", // your received usdt address
		FromCurrency:    "USDC",
		FromNetwork:     "BSC", // set from network (required)
		ToNetwork:       "BSC", // set to network (required)
		OrderedAmount:   0,     // use OrderedAmount or InvoicedAmount
		InvoicedAmount:  50,
		ToCurrency:      "USDT",
		ExtraID:         "",
		Signature:       res.Signature,
		UserReferenceID: "",
		RefundExtraID:   "",
	})
	fmt.Println(order, err)
	// // the exchange will return the rate of exchange is: order.ExchangeRate
	// // you will send btc to order.DepositAddress
	// // use OrderInfo to get order status
	orderInfo, err := exchange.OrderInfo(instantswap.TrackingRequest{OrderId: order.UUID})
	fmt.Println("Order Info: ", orderInfo, err)
	fmt.Println("Order Status: ", orderInfo.InternalStatus.String())
	// when ever the trading is done, the exchange will return the transaction id in
	// orderInfo.TxID
}
