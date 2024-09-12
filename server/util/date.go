package util

import "time"

func IsToday(update time.Time) bool {
	ty, tm, td := time.Now().Date()
	uy, um, ud := update.Date()
	return ty == uy && tm == um && td == ud
}
