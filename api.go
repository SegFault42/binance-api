package binance-api

import (
	"context"
	"log"
	"os"
	"strconv"

	binance "github.com/adshao/go-binance/v2"
)

var (
	apiKey    = os.Getenv("BINANCE_API_KEY")
	secretKey = os.Getenv("BINANCE_SECRET_KEY")
)

type ApiInfo struct {
	client *binance.Client
}

func New() ApiInfo {
	apiInfo := ApiInfo{client: binance.NewClient(apiKey, secretKey)}
	apiInfo.client.NewSetServerTimeService().Do(context.Background())

	return apiInfo
}

// func GetNewClient() {
// 	Client = binance.NewClient(apiKey, secretKey)
// }

// getDepth get the blockchain history
// symbol = "BTCUSDT"
// interval = refresh every X minutes, hours, ...
// startTime = from a timestamp (Optional)
// limit = number of iteration (Optional)
func (a ApiInfo) GetDepth(symbol, interval string, startTime int64, limit int) ([]*binance.Kline, error) {

	klines := a.client.NewKlinesService()
	klines = klines.Symbol(symbol)
	klines = klines.Interval(interval)

	if startTime > 0 {
		klines = klines.StartTime(startTime)
	}
	if limit > 0 {
		klines = klines.Limit(limit)
	}

	res, err := klines.Do(context.Background())
	if err != nil {
		return res, nil
	}

	return res, err
}

// GetCurrentPrice This function get the last price on binance
func (a ApiInfo) GetCurrentPrice(symbol string) (*binance.SymbolPrice, error) {

	prices, err := a.client.NewListPricesService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	for _, p := range prices {
		if p.Symbol == symbol {
			return p, err
		}
	}

	log.Println("symbol not found:", symbol)
	return nil, err
}

// PlaceOrderLimit will buy when it will reach a limit
// Example : resp, err := client.PlaceOrderLimit(binance.SideTypeSell, pair, "0.000204", "51700", "real")
func (a ApiInfo) PlaceOrderLimit(action binance.SideType, pair, quantity, price, mode string) (*binance.CreateOrderResponse, error) {

	var order *binance.CreateOrderResponse
	var err error

	req := a.client.NewCreateOrderService().Symbol(pair).
		Side(action).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(quantity).
		Price(price)

	if mode == "real" {
		order, err = req.Do(context.Background())
	} else if mode == "test" {
		err = req.Test(context.Background())
	}

	return order, err

}

// PlaceOrderMarket will buy to the marketPrice
// Quantity is in coin quantity
func (a ApiInfo) PlaceOrderMarket(action binance.SideType, pair, quantity, mode string) (*binance.CreateOrderResponse, error) {

	var order *binance.CreateOrderResponse
	var err error

	req := a.client.NewCreateOrderService().Symbol(pair).
		Side(action).Type(binance.OrderTypeMarket).Quantity(quantity)

	if mode == "real" {
		order, err = req.Do(context.Background())
	} else if mode == "test" {
		err = req.Test(context.Background())
	}

	return order, err

}

func (a ApiInfo) GetOrders(pair string) ([]*binance.Order, error) {
	return a.client.NewListOrdersService().Symbol(pair).Do(context.Background())
}

// func ListOrder(pair string) error {

// 	orders, err := Client.NewListOrdersService().Symbol(pair).Do(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	for _, elem := range orders {
// 		pp.Println(elem)
// 	}

// 	return nil
// }

func (a ApiInfo) GetBalances() ([]binance.Balance, error) {

	var balances []binance.Balance

	res, err := a.GetAccountService()
	if err != nil {
		return balances, err
	}

	for _, elem := range res.Balances {

		free, err := strconv.ParseFloat(elem.Free, 64)
		if err != nil {
			return balances, err
		}
		locked, err := strconv.ParseFloat(elem.Free, 64)
		if err != nil {
			return balances, err
		}

		if free > 0 || locked > 0 {
			balances = append(balances, elem)
		}
	}

	return balances, err
}

func (a ApiInfo) GetBalance(asset string, balances []binance.Balance) binance.Balance {

	for _, elem := range balances {
		if elem.Asset == asset {
			return elem
		}
	}

	return binance.Balance{}
}

func (a ApiInfo) GetInfoService() (*binance.ExchangeInfo, error) {
	return a.client.NewExchangeInfoService().Do(context.Background())
}

func (a ApiInfo) GetAccountService() (*binance.Account, error) {
	return a.client.NewGetAccountService().Do(context.Background())
}

func (a ApiInfo) GetTickerPrices() ([]*binance.SymbolPrice, error) {
	return a.client.NewListPricesService().Do(context.Background())
}

func (a ApiInfo) GetTickerPrice(pair string) (string, error) {
	prices, err := a.GetTickerPrices()
	if err != nil {
		return "", err
	}

	for _, p := range prices {
		if pair == p.Symbol {
			return p.Price, nil
		}
	}
	return "", nil
}
