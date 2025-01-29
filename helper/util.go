package helper

import (
	"fmt"
	"time"
)

func GenerateStringSlug(filename string, id int) string {
	now := time.Now()
	timestamp := now.Format("02-01-2006-15-04-05") // Format: DD-MM-YYYY-HH-MM-SS

	return fmt.Sprintf("%s-%d-%s", filename, id, timestamp)
}
