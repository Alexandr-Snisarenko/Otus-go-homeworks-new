package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	sBldr := strings.Builder{}
	rArr := []rune(s)

	if s == "" {
		return "", nil
	}
	// первый символ не должен быть цифрой
	if isNumber(rArr[0]) {
		return "", fmt.Errorf("start from number: %w", ErrInvalidString)
	}

	// алгоритм следующий: анализируем текущий символ, печатаем предыдущий.
	// начинаем со второго (индекс 1)
	// если текущий - цифра и предыдущий - цифра, то - ошибка
	// если текущий - цифра и предыдущий - не цифра. печатаем предыдущий указанное число раз
	// если текущий - не цифра и предыдущий - цифра, то пропукаем
	// если текущий - не цифра и предыдущий - не цифра, печатаем предыдущий
	for i := 1; i < len(rArr); i++ {
		if repeatCnt, err := strconv.Atoi(string(rArr[i])); err == nil {
			if isNumber(rArr[i-1]) {
				return "", fmt.Errorf("two numbers in a row: %w", ErrInvalidString)
			}
			sBldr.WriteString(strings.Repeat(string(rArr[i-1]), repeatCnt))
		} else if !isNumber(rArr[i-1]) {
			sBldr.WriteRune(rArr[i-1])
		}
	}

	// обрабаотываем крайний символ
	if !isNumber(rArr[len(rArr)-1]) {
		sBldr.WriteRune(rArr[len(rArr)-1])
	}

	return sBldr.String(), nil
}

// ф-я проверки руны на цифру.
func isNumber(r rune) bool {
	_, err := strconv.Atoi(string(r))
	return err == nil
}
