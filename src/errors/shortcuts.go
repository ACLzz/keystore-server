package errors

import "fmt"


func noField(field string) string {
	return fmt.Sprintf("no '%s' was sent", field)
}

func lengthMax(field string, limit int) string {
	return fmt.Sprintf("%s can't be more than %d symbols", field, limit)
}

func lengthMin(field string, limit int) string {
	return fmt.Sprintf("%s must be at least %d symbols", field, limit)
}

func localeError(field string) string {
	return fmt.Sprintf("%s must contain only english or digit symbols", field)
}
