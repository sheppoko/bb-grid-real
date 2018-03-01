package config

const (
	Debug = 0
)

const (
	CoinName                     = "ltc"
	CoinPairName                 = "btc"
	OrderNumInOnetime            = 5
	UnSoldBuyPositionLogFileName = "unsold_buy_position.json"
	SecJsonFileName              = "sec.json"
)

var (
	BuyRange               = 0.022
	TakeProfitRange        = 0.022
	MaxPositionCount       = 22
	PositionMaxDownPercent = 40.0
)
