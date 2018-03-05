package config

const (
	Debug = 0
)

const (
	CoinName                     = "mona"
	CoinPairName                 = "btc"
	OrderNumInOnetime            = 5
	UnSoldBuyPositionLogFileName = "unsold_buy_position.json"
	SecJsonFileName              = "sec.json"
)

var (
	BuyRange               = 0.019
	TakeProfitRange        = 0.019
	MaxPositionCount       = 22
	PositionMaxDownPercent = 47.0
)
