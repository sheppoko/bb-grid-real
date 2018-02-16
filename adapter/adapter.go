package adapter

import (
	"bitbank-grid-trade/api"
	"bitbank-grid-trade/config"
	"bitbank-grid-trade/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"time"
)

type UnsoldBuyOrder struct {
	OrderID            int     `json:"order_id"`
	BuyPrice           float64 `json:"buy_price"`
	RemainingBuyAmount float64 `json:"remaining_buy_amount"`
}

var unSoldBuyOrders = []*UnsoldBuyOrder{}
var latestPositionNum int
var activeOrdersCache *api.ActiveOrdersResponse
var highestMarketPrice = 0.0

func StartStrategy() {
	counter := 0
	for {
		time.Sleep(1000 * time.Millisecond) // 休む
		initCache()
		if counter%60 == 0 && false {
			counter = 0
			errCandle := SetRangeFromCandle()
			if errCandle != nil {
				fmt.Println("ロウソク取得中にエラーが発生しました")
				fmt.Println(errCandle)
				continue
			}
		}
		return
		counter++

		_, err := SellCoinIfNeedAndUpdateUnsold()
		if err != nil {
			fmt.Println("売り注文必要チェック及び売り注文作成中にエラーが発生しました")
			fmt.Println(err)
			continue
		}
		//ポジション数を取得
		positionNum, err := GetSellPriceKindNum()
		if err != nil {
			fmt.Println("ポジション数取得時にエラーが発生しました")
			fmt.Println(err)
			continue
		}
		//一番安い売り注文から2段階下げた買い注文を起点に5個入れる
		//1つずつ値計算し、それより高いか同値の注文がある場合は飛ばす
		err = OrderIfNeed(positionNum)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func initCache() {
	activeOrdersCache = nil
}

func SetRangeFromCandle() error {

	candle, _ := api.GetCandle(time.Now().AddDate(0, 0, -260))
	for i := 0; i < 180; i++ {
		candlePart, e := api.GetCandle(time.Now().AddDate(0, 0, -261+i))
		if e != nil {
			fmt.Println(e)
		}
		candle.Data.Candlestick[0].Ohlcv = append(candle.Data.Candlestick[0].Ohlcv, candlePart.Data.Candlestick[0].Ohlcv...)
	}

	bestMoney := 0.0
	bestRange := 0.0
	bestMaxPosition := 0.0
	bestTakeProfitCounter := 0
	songiriCounter := 0
	for buyRange := 0.001; buyRange < 0.01; buyRange = buyRange + 0.001 {
		songiriCounter = 0
		maxHigh := 0.0
		money := 1000000.0
		shouldUpdateLow := false
		basePrice, _ := util.StringToFloat(candle.Data.Candlestick[0].Ohlcv[0][2].(string))
		positions := []float64{basePrice}
		maxPosition := MaxPositionFromRange(buyRange)
		takeProfitCounter := 0
		profitPercentPerTime := (1 / maxPosition) * buyRange
		for _, data := range candle.Data.Candlestick[0].Ohlcv {
			//今回の足で利益を取っていたら計測基準を次の足の安値にする
			high, _ := util.StringToFloat(data[1].(string))
			low, _ := util.StringToFloat(data[2].(string))

			if maxHigh < high {
				maxHigh = high
			}
			if shouldUpdateLow {
				basePrice = low
			}
			if basePrice > low {
				basePrice = low
			}
			if maxHigh*(100-config.PositionMaxDownPercent)/100 > low {
				money = money * (100 - (config.PositionMaxDownPercent / 2)) / 100
				positions = []float64{low}
				fmt.Println("損切り", low, "追加")
				maxHigh = low
				songiriCounter++
			}
			hajime, _ := util.StringToFloat(data[0].(string))
			owari, _ := util.StringToFloat(data[3].(string))
			isYosen := hajime < owari
			nPos := positions

			if isYosen {
				start := positions[len(positions)-1] * (1 - buyRange)
				for i := start; i > low; i = i * (1 - buyRange) {
					fmt.Println(low, i, "追加")
					positions = append(positions, i)
				}

				for _, p := range positions {
					if p*(1+buyRange) < high {
						fmt.Println(high, p, "売却")
						money = money * (1 + profitPercentPerTime)
						nPos = remove(positions, p)
						takeProfitCounter++
					}
				}
				positions = nPos
				if len(positions) == 0 {
					fmt.Println("成行", high, "追加")
					positions = append(positions, high)
				}
			} else {
				for _, p := range positions {
					if p*(1+buyRange) < high {
						fmt.Println(high, p, "売却")
						money = money * (1 + profitPercentPerTime)
						nPos = remove(positions, p)
						takeProfitCounter++
					}
				}
				positions = nPos
				if len(positions) == 0 {
					fmt.Println("成行", high, "追加")
					positions = append(positions, high)
				}
				start := positions[len(positions)-1] * (1 - buyRange)
				for i := start; i > low; i = i * (1 - buyRange) {
					fmt.Println(low, i, "追加")
					positions = append(positions, i)
				}
			}

			// if ((high-basePrice)/basePrice)*0.75 > buyRange && isYosen {
			// 	simCounter += float64(int64(((high - basePrice) / basePrice) * 0.75 / buyRange))
			// 	money = money * (1 + profitPercentPerTime)
			// 	shouldUpdateLow = true
			// }
		}

		if money > bestMoney {
			bestRange = buyRange
			bestMoney = money
			bestTakeProfitCounter = takeProfitCounter
		}
	}
	config.BuyRange = bestRange
	config.TakeProfitRange = bestRange
	config.MaxPositionCount = int(bestMaxPosition)
	fmt.Printf("%f,%f,%d,%d\n", bestMoney, bestRange, bestTakeProfitCounter, songiriCounter)
	return nil

}

func remove(numbers []float64, search float64) []float64 {
	result := []float64{}
	for _, num := range numbers {
		if num != search {
			result = append(result, num)
		}
	}
	return result
}

func MaxPositionFromRange(buyRange float64) float64 {
	pricePer := 100.0
	positionNum := 0.0
	for {
		positionNum++
		pricePer *= 1 - buyRange
		if (100.0 - pricePer) > config.PositionMaxDownPercent {
			break
		}
	}
	return positionNum
}

func LoadUnSoldStatus() (bool, error) {
	raw, err := ioutil.ReadFile(config.UnSoldBuyPositionLogFileName)
	if err != nil {
		fmt.Println("状態の復元に失敗しました")
		return false, err
	} else {
		err = json.Unmarshal(raw, &unSoldBuyOrders)
		fmt.Println("状態を復元しました")
	}
	return true, nil
}

//一番安い売り注文か現在最良Askの高い方から2段階下げた買い注文を起点に5個入れます
//1つずつ値計算し、類似価格の注文がある場合は飛ばします
func OrderIfNeed(nowPositionNum int) error {
	lowestSell, err := GetLowestSellOrderPrice()
	if err != nil {
		return err
	}
	board, errB := api.GetBoard()
	if errB != nil {
		return errB
	}
	boardPrice, _ := util.StringToFloat(board.Data.Asks[0][0])
	startPrice := (lowestSell / (1 + config.TakeProfitRange)) * (1 - config.BuyRange)
	if lowestSell < boardPrice {
		startPrice = (boardPrice / (1 + config.TakeProfitRange))
	}
	orderNum := config.MaxPositionCount - nowPositionNum
	buyMax := orderNum
	if buyMax >= config.OrderNumInOnetime {
		buyMax = config.OrderNumInOnetime
	}
	for i := 0; i < buyMax; i++ {
		p := startPrice * math.Pow((1-config.BuyRange), float64(i))
		hasRangeBuyOrder, err := hasRangeBuyOrder(p)
		if err != nil {
			return err
		}
		if !hasRangeBuyOrder {
			freeJpy, errJpy := api.GetFreeJPY()
			if errJpy != nil {
				return errJpy
			}
			useJpy := freeJpy / float64(orderNum)
			amount := useJpy / p
			BuyCoinAndRegistUnsold(amount, p)
		} else {
			continue
		}
	}
	return nil
}

func hasRangeBuyOrder(price float64) (bool, error) {
	activeOrder, err := GetActiveOrdersFromAPIorCache()
	if err != nil {
		return false, err
	}
	for _, order := range activeOrder.Data.Orders {
		if order.Side == "buy" && (math.Abs(order.Price-price) < price*config.BuyRange) {
			return true, nil
		}
	}
	return false, nil
}

func BuyCoinAndRegistUnsold(amount float64, price float64) (bool, error) {
	amount = util.Round(amount, 4)

	res, err := api.BuyCoin(amount, price)
	if err != nil {
		return false, err
	}
	appendOrder := new(UnsoldBuyOrder)
	appendOrder.OrderID = res.Data.OrderID
	appendOrder.BuyPrice = res.Data.Price
	appendOrder.RemainingBuyAmount = res.Data.RemainingAmount
	unSoldBuyOrders = append(unSoldBuyOrders, appendOrder)
	util.SaveJsonToFile(unSoldBuyOrders, config.UnSoldBuyPositionLogFileName)
	return true, nil
}

//売り注文が出されていない買い注文と現在の状況をチェック�����必要であれば売り注文をなげます。
//戻り値��約定した事が判明し������い注文のうちもっとも低い価��の注文です
func SellCoinIfNeedAndUpdateUnsold() (float64, error) {
	max := 1234567890.0
	lowestPrice := max
	orderIds := []int{}
	var res = &api.OrdersInfoResponse{}
	for i, order := range unSoldBuyOrders {
		orderIds = append(orderIds, order.OrderID)
		if i%5 == 0 || i == len(unSoldBuyOrders)-1 {
			resPart, err := api.GetOrdersInfo(orderIds)
			if err != nil {
				return -1.0, err
			}
			res.Data.Orders = append(res.Data.Orders, resPart.Data.Orders...)
			orderIds = []int{}
		}
	}

	for _, order := range res.Data.Orders {
		for _, unSold := range unSoldBuyOrders {
			//当該注��の残量が��化していた場合は対応��た売り注文を投げる
			if unSold.OrderID == order.OrderID && unSold.RemainingBuyAmount != order.RemainingAmount {

				//一部約���の時に注���������通らない
				sellAmount := unSold.RemainingBuyAmount - order.RemainingAmount
				sellPrice := unSold.BuyPrice * (1 + config.TakeProfitRange)
				fmt.Println("買���注文が約定しているため売り注文を作成し���す...")
				if order.Price < lowestPrice {
					lowestPrice = order.Price
				}
				_, err := api.SellCoin(sellAmount, sellPrice)
				if err != nil {
					fmt.Println(unSold.RemainingBuyAmount, order.RemainingAmount)
					return -1, err
				}

				if order.RemainingAmount <= 0 {
					DeleteUnSoldOrder(order.OrderID)
				} else {
					unSold.RemainingBuyAmount -= sellAmount
				}
				util.SaveJsonToFile(unSoldBuyOrders, config.UnSoldBuyPositionLogFileName)
				fmt.Println("作成しま����た")
			}
		}
	}
	if lowestPrice == max {
		lowestPrice = -1
	}

	return lowestPrice, nil
}

func DeleteUnSoldOrder(orderId int) []*UnsoldBuyOrder {
	orders := []*UnsoldBuyOrder{}
	for _, unSoldOrder := range unSoldBuyOrders {
		if unSoldOrder.OrderID != orderId {
			orders = append(orders, unSoldOrder)
		}
	}
	unSoldBuyOrders = orders
	return orders
}

func DeleteAllUnSoldOrder() {
	orders := []*UnsoldBuyOrder{}
	unSoldBuyOrders = orders
}

//売り買い全ての注文をキャンセルします
func CancelAllOrders() (bool, error) {
	targetOrderId := []int{}
	orders, err := GetActiveOrdersFromAPIorCache()
	if err != nil {
		return false, err
	}
	for _, order := range orders.Data.Orders {
		targetOrderId = append(targetOrderId, order.OrderID)
	}
	_, err = api.CancelOrders(targetOrderId)
	if err != nil {
		return false, err
	}
	DeleteAllUnSoldOrder()
	util.SaveJsonToFile(unSoldBuyOrders, config.UnSoldBuyPositionLogFileName)
	fmt.Println("全ての注文をキャンセルしました")
	return true, nil
}

func GetLowestSellOrderPrice() (float64, error) {

	max := 123456789.00
	lowestSell := max
	activeOrder, err := GetActiveOrdersFromAPIorCache()
	if err != nil {
		return -1.0, err
	}
	for _, order := range activeOrder.Data.Orders {
		if order.Side == "sell" {
			if order.Price < lowestSell {
				lowestSell = order.Price
			}
		}
	}
	if lowestSell == max {
		return 0, nil
	}
	return lowestSell, nil
}

//売り注文の��段���種類数��������������取����������す
func GetSellPriceKindNum() (int, error) {
	res, err := GetActiveOrdersFromAPIorCache()
	prices := []float64{}
	if err != nil {
		return 0, err
	}
	for _, order := range res.Data.Orders {
		if order.Side == "sell" {
			shouldAdd := true
			for _, price := range prices {
				if price == order.Price {
					shouldAdd = false
				}
			}
			if shouldAdd {
				prices = append(prices, order.Price)
			}
		}
	}
	return len(prices), nil
}

func GetBuyOrderNum() (int, error) {
	ret := 0
	res, err := GetActiveOrdersFromAPIorCache()
	if err != nil {
		return 0, err
	}
	for _, order := range res.Data.Orders {
		if order.Side == "buy" {
			ret++
		}
	}
	return ret, nil
}

//買い注文を全てキャンセルします
func CancelAllBuyOrders() (bool, error) {
	targetOrderId := []int{}
	orders, err := GetActiveOrdersFromAPIorCache()
	if err != nil {
		return false, err
	}
	counter := 0
	for _, order := range orders.Data.Orders {
		if order.Side == "buy" {
			targetOrderId = append(targetOrderId, order.OrderID)
		}
		counter++
		if counter == 100 {
			break
		}
	}
	_, errCancel := api.CancelOrders(targetOrderId)
	if errCancel != nil {
		return false, errCancel
	}
	DeleteAllUnSoldOrder()
	util.SaveJsonToFile(unSoldBuyOrders, config.UnSoldBuyPositionLogFileName)
	fmt.Println("全ての買い注文をキ�����ンセ�������������した")
	return true, nil
}

func GetActiveOrdersFromAPIorCache() (*api.ActiveOrdersResponse, error) {
	if activeOrdersCache == nil {
		res, err := api.GetActiveOrders()
		if err != nil {
			return nil, err
		}
		activeOrdersCache = res
		return res, nil
	}

	return activeOrdersCache, nil
}
