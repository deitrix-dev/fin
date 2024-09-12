package main

import (
	"fmt"

	"github.com/deitrix/fin"
	"github.com/rickb777/date"
)

func main() {
	r := fin.Repeat{
		Every:      fin.Week,
		Weekday:    fin.Thursday,
		Multiplier: 5,
		Offset:     1,
	}

	for _, date := range r.DatesUntilN(date.Today(), 10) {
		fmt.Println(date.Format("2006-01-02 Mon"))
	}
}
