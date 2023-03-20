package main

import (
	"flag"
	"fmt"
	"os"

	_ "code.cryptopower.dev/group/instantswap/exchange/changelly"
	_ "code.cryptopower.dev/group/instantswap/exchange/changenow"
	_ "code.cryptopower.dev/group/instantswap/exchange/coinswitch"
	_ "code.cryptopower.dev/group/instantswap/exchange/flypme"
	_ "code.cryptopower.dev/group/instantswap/exchange/godex"
	_ "code.cryptopower.dev/group/instantswap/exchange/simpleswap"
	_ "code.cryptopower.dev/group/instantswap/exchange/swapzone"
	"code.cryptopower.dev/group/instantswap/instantswap"
)

// using this code:
// go run ./_example/flypme/ -exchange=swapzone -key=4WRyccNKW

var (
	exchange, apiKey, apiSecret string
)

func init() {
	flag.StringVar(&exchange, "exchange", "flypme", "-exchange=<the exchange you want to trade, default is flypme>")
	flag.StringVar(&apiKey, "key", "", "-key=<your api key>")
	flag.StringVar(&apiSecret, "secret", "", "-secret=<your api secret>")
	flag.Parse()
}

func main() {
	exchange, err := instantswap.NewExchange(exchange, instantswap.ExchangeConfig{
		Debug:     false,
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	currencies, err := exchange.GetCurrencies()
	fmt.Println("currencies length: ", len(currencies), err)
	currencies, err = exchange.GetCurrenciesToPair("btc")
	fmt.Println("GetCurrenciesToPair length: ", len(currencies), err)
	res, err := exchange.GetExchangeRateInfo(instantswap.ExchangeRateRequest{
		From:   "BTC",
		To:     "DCR",
		Amount: 5,
	})
	fmt.Printf("%+v \n %v \n", res, err)
	order, err := exchange.CreateOrder(instantswap.CreateOrder{
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
