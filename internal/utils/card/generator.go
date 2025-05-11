package card

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateCardNumber() string {
	prefix := "4"

	partial := prefix
	for i := 0; i < 14; i++ {
		partial += fmt.Sprintf("%d", rand.Intn(10))
	}

	// Вычисление контрольной цифры
	checkDigit := luhnCheckDigit(partial)
	return partial + fmt.Sprintf("%d", checkDigit)
}

func GenerateExpiryDate() string {
	year := (time.Now().Year() + 3) % 100
	month := time.Now().Month()

	return fmt.Sprintf("%02d/%02d", month, year)
}

func GenerateCVV() string {
	return fmt.Sprintf("%03d", rand.Intn(1000))
}

func luhnCheckDigit(number string) int {
	sum := 0
	alternate := false

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if alternate {
			digit *= 2
			if digit > 9 {
				digit = (digit % 10) + 1
			}
		}
		sum += digit
		alternate = !alternate
	}

	return (10 - (sum % 10)) % 10
}
