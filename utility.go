package main

import (
	"fmt"
	"strconv"
	"time"
)

func FormatIntToPrice(value int) string {
	return fmt.Sprintf("€%.2f", float64(value)/float64(100))
}

func TimeOfDay() string {
	c, _ := strconv.Atoi(time.Now().Format("15"))
	switch {
	case c < 12:
		return "Good morning"
	case c < 18:
		return "Good afternoon"
	}
	return "Good evening"
}
