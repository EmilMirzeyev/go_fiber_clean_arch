package util

import (
	"time"
)

// CalculateAge returns the age in years from birthdate to the current date
func CalculateAge(birthdate time.Time) int {
	now := time.Now()

	age := now.Year() - birthdate.Year()

	// Adjust age if birthday hasn't occurred yet this year
	if now.Month() < birthdate.Month() ||
		(now.Month() == birthdate.Month() && now.Day() < birthdate.Day()) {
		age--
	}

	return age
}

// ParseBirthdate parses a string date in DD.MM.YYYY format to time.Time
func ParseBirthdate(birthdate string) (time.Time, error) {
	layout := "02.01.2006"
	return time.Parse(layout, birthdate)
}
