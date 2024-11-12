package util

import (
	"fmt"
	"time"

	customError "github.com/federicodosantos/socialize/pkg/custom-error"
)

func ErrRowsAffected(rows int64) error {
	if rows != 1 {
		return fmt.Errorf("error : %w, got %d rows affected", customError.ErrRowsAffected, rows)
	}

	return nil
}

func ConvertTimeToString(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}
