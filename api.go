package binanceapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	binance "github.com/adshao/go-binance/v2"
)

var (
	apiKey    = os.Getenv("BINANCE_API_KEY")
	secretKey = os.Getenv("BINANCE_SECRET_KEY")
)

type ApiInfo struct {
	Client  *binance.Client
	lotSize SLotSize
}

type SLotSize struct {
	FilterType string
	MinQty     string
	MaxQty     string
	StepSize   string
}

type SExchangeInfo struct {
	Timezone   string `json:"timezone"`
	Servertime int64  `json:"serverTime"`
	Ratelimits []struct {
		Ratelimittype string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		Intervalnum   int    `json:"intervalNum"`
		Limit         int    `json:"limit"`
	} `json:"rateLimits"`
	Exchangefilters []interface{} `json:"exchangeFilters"`
	Symbols         []struct {
		Symbol                     string   `json:"symbol"`
		Status                     string   `json:"status"`
		Baseasset                  string   `json:"baseAsset"`
		Baseassetprecision         int      `json:"baseAssetPrecision"`
		Quoteasset                 string   `json:"quoteAsset"`
		Quoteprecision             int      `json:"quotePrecision"`
		Quoteassetprecision        int      `json:"quoteAssetPrecision"`
		Basecommissionprecision    int      `json:"baseCommissionPrecision"`
		Quotecommissionprecision   int      `json:"quoteCommissionPrecision"`
		Ordertypes                 []string `json:"orderTypes"`
		Icebergallowed             bool     `json:"icebergAllowed"`
		Ocoallowed                 bool     `json:"ocoAllowed"`
		Quoteorderqtymarketallowed bool     `json:"quoteOrderQtyMarketAllowed"`
		Isspottradingallowed       bool     `json:"isSpotTradingAllowed"`
		Ismargintradingallowed     bool     `json:"isMarginTradingAllowed"`
		Filters                    []struct {
			Filtertype       string `json:"filterType"`
			Minprice         string `json:"minPrice,omitempty"`
			Maxprice         string `json:"maxPrice,omitempty"`
			Ticksize         string `json:"tickSize,omitempty"`
			Multiplierup     string `json:"multiplierUp,omitempty"`
			Multiplierdown   string `json:"multiplierDown,omitempty"`
			Avgpricemins     int    `json:"avgPriceMins,omitempty"`
			Minqty           string `json:"minQty,omitempty"`
			Maxqty           string `json:"maxQty,omitempty"`
			Stepsize         string `json:"stepSize,omitempty"`
			Minnotional      string `json:"minNotional,omitempty"`
			Applytomarket    bool   `json:"applyToMarket,omitempty"`
			Limit            int    `json:"limit,omitempty"`
			Maxnumorders     int    `json:"maxNumOrders,omitempty"`
			Maxnumalgoorders int    `json:"maxNumAlgoOrders,omitempty"`
		} `json:"filters"`
		Permissions []string `json:"permissions"`
	} `json:"symbols"`
}

func New() ApiInfo {
	apiInfo := ApiInfo{Client: binance.NewClient(apiKey, secretKey)}

	//Sync time with binance server time
	apiInfo.Client.NewSetServerTimeService().Do(context.Background())

	return apiInfo
}

// getDepth get the blockchain history
// symbol = "BTCUSDT"
// interval = refresh every X minutes, hours, ...
// startTime = from a timestamp (Optional)
// limit = number of iteration (Optional)
func (a ApiInfo) GetDepth(symbol, interval string, startTime int64, limit int) ([]*binance.Kline, error) {

	klines := a.Client.NewKlinesService()
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
func (a ApiInfo) GetCurrentPrice(pair string) (*binance.SymbolPrice, error) {

	prices, err := a.Client.NewListPricesService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	for _, p := range prices {
		if p.Symbol == pair {
			return p, err
		}
	}

	log.Println("symbol not found:", pair)
	return nil, err
}

// PlaceOrderLimit will buy when it will reach a limit
// Example : resp, err := client.PlaceOrderLimit(binance.SideTypeSell, pair, "0.000204", "51700", "real")
func (a ApiInfo) PlaceOrderLimit(action binance.SideType, pair, quantity, price, mode string) (*binance.CreateOrderResponse, error) {

	var order *binance.CreateOrderResponse
	var err error

	req := a.Client.NewCreateOrderService().Symbol(pair).
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
// Quantity is in coin quantity.
// Example : buy AUDIOBTC, must init quantity with BTC amount
func (a ApiInfo) PlaceOrderMarket(action binance.SideType, pair, quantity, mode string) (*binance.CreateOrderResponse, error) {

	var order *binance.CreateOrderResponse
	var err error

	req := a.Client.NewCreateOrderService().Symbol(pair).
		Side(action).Type(binance.OrderTypeMarket).QuoteOrderQty(quantity)

	if mode == "real" {
		order, err = req.Do(context.Background())
	} else if mode == "test" {
		err = req.Test(context.Background())
	}

	return order, err

}

func (a ApiInfo) PlaceOrderMarketQuantity(action binance.SideType, pair, quantity, mode string) (*binance.CreateOrderResponse, error) {

	var order *binance.CreateOrderResponse
	var err error

	req := a.Client.NewCreateOrderService().Symbol(pair).
		Side(action).Type(binance.OrderTypeMarket).QuoteOrderQty(quantity)

	if mode == "real" {
		order, err = req.Do(context.Background())
	} else if mode == "test" {
		err = req.Test(context.Background())
	}

	return order, err

}

func (a ApiInfo) PlaceOrderMarketAmount(action binance.SideType, pair, quantity, mode string) (*binance.CreateOrderResponse, error) {

	var order *binance.CreateOrderResponse
	var err error

	req := a.Client.NewCreateOrderService().Symbol(pair).
		Side(action).Type(binance.OrderTypeMarket).Quantity(quantity)

	if mode == "real" {
		order, err = req.Do(context.Background())
	} else if mode == "test" {
		err = req.Test(context.Background())
	}

	return order, err

}

func (a ApiInfo) GetOrders() ([]*binance.Order, error) {
	return a.Client.NewListOrdersService().Do(context.Background())
}

func (a ApiInfo) GetOpenOrders() ([]*binance.Order, error) {
	return a.Client.NewListOpenOrdersService().Do(context.Background())
}

//func ListOpenOrder(pair string) {
//}

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

// GetBalances return all more than 0 balance)
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
		locked, err := strconv.ParseFloat(elem.Locked, 64)
		if err != nil {
			return balances, err
		}

		if free > 0 || locked > 0 {
			balances = append(balances, elem)
		}
	}

	return balances, err
}

// GetBalance return the balance of the given asset
func (a ApiInfo) GetBalance(asset string) (binance.Balance, error) {

	balances, err := a.GetBalances()
	if err != nil {
		return binance.Balance{}, err
	}
	for _, elem := range balances {
		if elem.Asset == asset {
			return elem, nil
		}
	}

	return binance.Balance{}, nil
}

func (a ApiInfo) GetInfoService() (*binance.ExchangeInfo, error) {
	return a.Client.NewExchangeInfoService().Do(context.Background())
}

func (a ApiInfo) GetAccountService() (*binance.Account, error) {
	return a.Client.NewGetAccountService().Do(context.Background())
}

func (a ApiInfo) GetTickerPrices() ([]*binance.SymbolPrice, error) {
	return a.Client.NewListPricesService().Do(context.Background())
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

// WsGetCoinPrice return a struct with all real time event
// send adress to price (&price)
func WsGetCoinPrice(pair string, evt *binance.WsAggTradeEvent) {
	wsDepthHandler := func(event *binance.WsAggTradeEvent) {
		*evt = *event
	}
	errHandler := func(err error) {
		log.Println(err)
	}
	_, _, err := binance.WsAggTradeServe(pair, wsDepthHandler, errHandler)
	if err != nil {
		return
	}
}

func (a ApiInfo) GetTransactionPriceAtTime(pair string, time int64) (string, error) {
	trades, err := a.Client.NewAggTradesService().
		Symbol(pair).StartTime(time).EndTime(time).
		Do(context.Background())
	if err != nil {
		return "", err
	}

	average := 0.0
	for _, elem := range trades {
		res, err := strconv.ParseFloat(elem.Price, 64)
		if err != nil {
			return "", err
		}
		average += res
	}

	average = average / float64(len(trades))

	return fmt.Sprintf("%.8f", average), nil
}

func (a ApiInfo) GetTransactionByOrderID(orderID int64, pair string) (*binance.Order, error) {
	return a.Client.NewGetOrderService().Symbol(pair).OrderID(orderID).Do(context.Background())
}

func (a ApiInfo) GetLastFilledTransaction(pair string) (*binance.Order, error) {

	orders, err := a.Client.NewListOrdersService().Symbol(pair).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, nil
	}
	var tmp *binance.Order
	tmp = orders[0]

	for _, elem := range orders {
		if elem.Time > tmp.Time && elem.Status == "FILLED" {
			tmp = elem
		}
	}

	return tmp, nil
}

func (a ApiInfo) GetFreeCoinAmount(asset string) (string, error) {
	res, err := a.GetBalance(asset)
	if err != nil {
		return "", err
	}

	return res.Free, nil
}

func getExchangeInfo() (SExchangeInfo, error) {
	resp, err := http.Get("https://www.binance.com/api/v1/exchangeInfo")
	if err != nil {
		return SExchangeInfo{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return SExchangeInfo{}, err
	}

	var exchangeInfo SExchangeInfo

	if err := json.Unmarshal([]byte(body), &exchangeInfo); err != nil {
		return SExchangeInfo{}, err
	}

	return exchangeInfo, nil
}

func (a ApiInfo) GetLotSize(pair string) (SLotSize, error) {

	exchangeInfo, err := getExchangeInfo()
	if err != nil {
		return SLotSize{}, err
	}

	for _, symbol := range exchangeInfo.Symbols {
		if symbol.Symbol == pair {
			for _, filter := range symbol.Filters {
				if filter.Filtertype == "LOT_SIZE" {
					return SLotSize{
						FilterType: filter.Filtertype,
						MinQty:     filter.Minqty,
						MaxQty:     filter.Maxqty,
						StepSize:   filter.Stepsize,
					}, nil
				}
			}
		}
	}

	return SLotSize{}, nil
}

// func (a ApiInfo) GetPriceAtSpecificTime(symbol string, timestamp int64) error {
// 	ts := int64(timestamp)
// 	trades, err := a.Client.NewAggTradesService().
// 		Symbol(symbol).
// 		// StartTime(ts).EndTime(ts + time.Hour.Milliseconds()).
// 		StartTime(ts).EndTime(ts + time.Hour.Milliseconds()).
// 		Do(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	pp.Println(trades)
// 	// price := trades[0].Price
// 	// fmt.Println(price)

// 	return nil
// }
