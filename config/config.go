package config

const (
	Debug = 0
)

const (
	CoinName                     = "xrp"
	OrderNumInOnetime            = 40
	UnSoldBuyPositionLogFileName = "unsold_buy_position.json"
	SecJsonFileName              = "sec.json"
)

var (
	BuyRange               = 0.00005
	TakeProfitRange        = 0.01
	MaxPositionCount       = 3000
	PositionMaxDownPercent = 10.0
)
