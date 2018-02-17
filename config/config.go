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
	BuyRange               = 0.0001
	TakeProfitRange        = 0.005
	MaxPositionCount       = 2232
	PositionMaxDownPercent = 20.0
)
