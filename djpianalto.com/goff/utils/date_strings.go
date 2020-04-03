package utils

import (
	"fmt"
	"time"
)

func ParseDateString(inTime time.Time) string {
	d := time.Now().Sub(inTime)
	s := int64(d.Seconds())
	days := s / 86400
	s = s - (days * 86400)
	hours := s / 3600
	s = s - (hours * 3600)
	minutes := s / 60
	seconds := s - (minutes * 60)
	dateString := ""
	if days != 0 {
		dateString += fmt.Sprintf("%v days ", days)
	}
	if hours != 0 {
		dateString += fmt.Sprintf("%v hours ", hours)
	}
	if minutes != 0 {
		dateString += fmt.Sprintf("%v minutes ", minutes)
	}
	if seconds != 0 {
		dateString += fmt.Sprintf("%v seconds ", seconds)
	}
	if dateString != "" {
		dateString += " ago."
	} else {
		dateString = "Now"
	}
	stamp := inTime.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%v\n%v", dateString, stamp)
}
