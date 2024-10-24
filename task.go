package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102" // TODO: move to another block

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("пустое правило повторения")
	}

	startDate, err := time.Parse(DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}

	rule := strings.Split(repeat, " ")
	ruleLiteral := rule[0]

	switch ruleLiteral {
	case "d":
		if len(rule) < 2 {
			return "", fmt.Errorf("не указано количество дней")
		}
		daysN, err := strconv.Atoi(rule[1])
		if err != nil {
			return "", fmt.Errorf("неверное число дней: %v", err)
		}
		if daysN > 400 {
			return "", fmt.Errorf("число дней не может превышать 400")
		}
		newDate := startDate.AddDate(0, 0, daysN)
		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, daysN)
		}
		return newDate.Format(DateFormat), nil

	case "y":
		newDate := startDate.AddDate(1, 0, 0)
		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}
		return newDate.Format(DateFormat), nil

	default:
		return "", fmt.Errorf("некорректный литерал правила")
	}

	// TODO: добавить правило для w
	// TODO: добавить правило для m

}
