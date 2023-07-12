package main

import (
	"fmt"

	"github.com/dolfly/termux/pkg/termux"
)

func main() {
	if stat, err := termux.BatteryStatus(); err != nil {
		panic(err)
	} else {
		fmt.Printf("The current battery percentage is %d%%.\n", stat.Percentage)
	}
}
