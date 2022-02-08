package util

import (
	"fmt"
	"regexp"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

func GenerateUniqeId() string {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return id.String()
}

func ParseExpiryDate(expiry string) time.Time {
	r := regexp.MustCompile(`(\d\d)\/(\d\d)`)
	matches := r.FindStringSubmatch(expiry)

	month := matches[1]
	year := matches[2]

	layout := "2006-01-02"
	date := fmt.Sprintf("20%s-%s-15", year, month)

	parsedDate, err := time.Parse(layout, date)
	if err != nil {
		panic(err)
	}

	return parsedDate
}
