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
	BuyRange               = 0.01
	TakeProfitRange        = 0.01
	MaxPositionCount       = 50
	PositionMaxDownPercent = 40.0
)
