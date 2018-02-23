package config

const (
	Debug = 0
)

const (
	CoinName                     = "btc"
	OrderNumInOnetime            = 3
	UnSoldBuyPositionLogFileName = "unsold_buy_position.json"
	SecJsonFileName              = "sec.json"
)

var (
	BuyRange               = 0.018
	TakeProfitRange        = 0.018
	MaxPositionCount       = 38
	PositionMaxDownPercent = 50.0
)
