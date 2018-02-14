package main

import (
	"bitbank-grid-trade/adapter"
	"fmt"
)

func main() {

	_, err := adapter.LoadUnSoldStatus()
	if err != nil {
		fmt.Println(err)
		return
	}
	adapter.StartStrategy()

}
