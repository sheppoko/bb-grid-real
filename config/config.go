package config

const (
	Debug = 0
)

const (
	CoinName                     = "xrp"
	OrderNumInOnetime            = 10
	UnSoldBuyPositionLogFileName = "unsold_buy_position.json"
	SecJsonFileName              = "sec.json"
)

var (
	BuyRange               = 0.003
	TakeProfitRange        = 0.003
	MaxPositionCount       = 119
	PositionMaxDownPercent = 30.0
)
