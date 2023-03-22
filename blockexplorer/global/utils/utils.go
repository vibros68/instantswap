package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func SleepPrintMinutes(dur int, item string) {
	then := time.Now().Round(time.Second).Add(time.Minute * time.Duration(dur))
	duration := then.Sub(time.Now().Round(time.Minute))
	fmt.Printf("%s in %s at %v", item, duration.String(), then.Format("15:04:05"))
	sendDelay := dur
	for j := 1; j <= sendDelay; j++ {
		fmt.Printf(".")
		time.Sleep(time.Minute)
	}
	fmt.Println("")
}

func LabelMatching(label, rule string) bool {
	if rule == "*" {
		return true
	}
	if label == rule {
		return true
	}
	var flexLeft, flexRight bool
	ruleLabel := strings.TrimLeft(rule, "*")
	if ruleLabel != rule {
		flexLeft = true
		rule = ruleLabel
	}
	ruleLabel = strings.TrimRight(rule, "*")
	if ruleLabel != rule {
		flexRight = true
	}
	index := strings.Index(label, ruleLabel)
	if index == -1 {
		return false
	}
	if flexLeft && flexRight && index > -1 {
		return true
	}
	if flexRight && index == 0 {
		return true
	}
	if flexLeft && label[len(label)-len(ruleLabel):] == ruleLabel {
		return true
	}

	return false
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
