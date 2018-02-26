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
	BuyRange               = 0.03
	TakeProfitRange        = 0.03
	MaxPositionCount       = 17
	PositionMaxDownPercent = 50.0
)
