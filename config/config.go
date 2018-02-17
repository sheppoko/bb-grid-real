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
	BuyRange               = 0.005
	TakeProfitRange        = 0.005
	MaxPositionCount       = 45
	PositionMaxDownPercent = 20.0
)
