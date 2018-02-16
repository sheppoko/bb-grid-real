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
	BuyRange               = 0.01
	TakeProfitRange        = 0.01
	MaxPositionCount       = 30
	PositionMaxDownPercent = 20.0
)
