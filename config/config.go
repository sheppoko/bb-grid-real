package config

const (
	Debug = 0
)

const (
	CoinName                     = "xrp"
	OrderNumInOnetime            = 2
	UnSoldBuyPositionLogFileName = "unsold_buy_position.json"
)

var (
	BuyRange               = 0.0001
	TakeProfitRange        = 0.0001
	MaxPositionCount       = 2000
	PositionMaxDownPercent = 30.0
)
