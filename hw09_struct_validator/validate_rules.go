package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrInvalidValidateRuleString = errors.New("validate rule string is invalid")
	ErrInvalidValidateRuleParam  = errors.New("validate param is incorrect")
	ErrUnknownValidateRule       = errors.New("unknown validate rule")
	ErrInvalidRuleDataType       = errors.New("rule from another data type")
	ErrUnknownDataType           = errors.New("unknown data type")
	// Validate Errors.
	ErrStringLength   = errors.New("string length is not valid")
	ErrStringContent  = errors.New("string content is not valid")
	ErrStringNotInSet = errors.New("string is not in string set")
	ErrNumberValue    = errors.New("namber value is not valid")
	ErrNumberNotInSet = errors.New("number is not in number set")
	ErrDataValidate   = errors.New("data is not valid")
)

// правила валидации.
type ValidateRule uint

const (
	vrUnknown ValidateRule = iota
	vrLen
	vrRegexp
	vrIn
	vrMin
	vrMax
)

// мапинг правил на названия в строке.
var ValidateRuleNameMap = map[string]ValidateRule{
	"len":    vrLen,
	"regexp": vrRegexp,
	"in":     vrIn,
	"min":    vrMin,
	"max":    vrMax,
}

// перечень поддерживаемых правил для строк.
var StringValidateRules = []ValidateRule{
	vrLen,    // "len:32" - длина строки должна быть ровно 32 символа;
	vrRegexp, // "regexp:\\d+" - согласно регулярному выражению строка должна состоять из цифр (\\ - экранирование слэша);
	vrIn,     // "in:foo,bar" - строка должна входить в множество строк {"foo", "bar"}
	vrMin,    // "min:10" - строка не может быть меньше 10 символов;
	vrMax,    // "max:20" - строка не может быть больше 20 символов;
}

// перечень поддерживаемых правил для чисел.
var NumberValidateRules = []ValidateRule{
	vrMin, // "min:10" - число не может быть меньше 10;
	vrMax, // "max:20" - число не может быть больше 20;
	vrIn,  // "in:256,1024" - число должно входить в множество чисел {256, 1024};
}

// набор правил с параметрами кнтроля.
type ValidateRuleSet map[ValidateRule]string

// конструктор для набора правил.
func NewValidateRuleSet(validateString string) (ValidateRuleSet, error) {
	vSet := make(ValidateRuleSet)
	// если строка правил не задана - возвращаем пустой RuleSet
	if validateString == "" {
		return vSet, nil
	}

	// если строка не пустая - разбираем строку и формируем список правил валидации
	// делим строку на катерны по символу "|"
	for _, ruleStr := range strings.Split(validateString, "|") {
		// разделяем правило и контрольные значения по ":"
		sRule, sVal, ok := strings.Cut(ruleStr, ":")
		if !ok {
			return nil, ErrInvalidValidateRuleString
		}
		// если правила нет в общем списке правил - ошибка
		rule, ok := ValidateRuleNameMap[sRule]
		if !ok {
			return nil, ErrUnknownValidateRule
		}
		vSet[rule] = sVal
	}

	return vSet, nil
}

// проверка данных по списку правил.
func (v ValidateRuleSet) CheckData(data interface{}) error {
	// получаем объект за интерфейсом
	rVal := reflect.ValueOf(data)

	// выполняем проверку данных в зависимости от типа данных
	switch rVal.Kind() {
	case reflect.String:
		return checkStringData(v, rVal.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return checkNumberData(v, rVal.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return checkNumberData(v, rVal.Uint())
	default:
		return ErrUnknownDataType
	}
}

// методы для провеки значений по типам реализуются в виде отдельных функций
// ввиду ограничения для дженериков в методах.
func checkStringData(rSet ValidateRuleSet, data string) error {
	for rule, val := range rSet {
		switch rule {
		case vrLen:
			chkVal, err := strconv.Atoi(val)
			if err != nil {
				return ErrInvalidValidateRuleParam
			}
			if len(data) != chkVal {
				return ErrStringLength
			}
		case vrMin:
			chkVal, err := strconv.Atoi(val)
			if err != nil {
				return ErrInvalidValidateRuleParam
			}
			if len(data) < chkVal {
				return ErrStringLength
			}
		case vrMax:
			chkVal, err := strconv.Atoi(val)
			if err != nil {
				return ErrInvalidValidateRuleParam
			}
			if len(data) > chkVal {
				return ErrStringLength
			}
		case vrRegexp:
			rexp, err := regexp.Compile(val)
			if err != nil {
				return ErrInvalidValidateRuleParam
			}

			if !rexp.MatchString(data) {
				return ErrStringContent
			}
		case vrIn:
			if !slices.Contains(strings.Split(val, ","), data) {
				return ErrStringNotInSet
			}
		default:
			return ErrUnknownValidateRule
		}
	}

	return nil
}

// // generics ////
// функция парсировки данных из строки в number (int64 или uint64).
func parseAnyInt[T int64 | uint64](val string) (T, error) {
	var (
		parsVal T
		i64     int64
		ui64    uint64
		err     error
	)
	// определяем тип данных от T
	dataType := fmt.Sprintf("%T", parsVal)

	switch dataType {
	case "int64":
		i64, err = strconv.ParseInt(val, 10, 64)
		if err != nil {
			return parsVal, ErrInvalidValidateRuleParam
		}
		parsVal = T(i64)
	case "uint64":
		ui64, err = strconv.ParseUint(val, 10, 64)
		if err != nil {
			return parsVal, ErrInvalidValidateRuleParam
		}
		parsVal = T(ui64)
	default:
		return parsVal, ErrUnknownDataType
	}

	return parsVal, nil
}

// функця проверки параметров типа int64 или uint64.
func checkNumberData[T int64 | uint64](rSet ValidateRuleSet, data T) error {
	for rule, val := range rSet {
		if !slices.Contains(NumberValidateRules, rule) {
			return ErrInvalidRuleDataType
		}

		switch rule {
		case vrMin:
			chkVal, err := parseAnyInt[T](val)
			if err != nil {
				return ErrInvalidValidateRuleParam
			}
			if data < chkVal {
				return fmt.Errorf("%w less then 'min' value", ErrNumberValue)
			}
		case vrMax:
			chkVal, err := parseAnyInt[T](val)
			if err != nil {
				return ErrInvalidValidateRuleParam
			}

			if data > chkVal {
				return fmt.Errorf("%w greate then 'max' value", ErrNumberValue)
			}
		case vrIn:
			// разбираем строку параметров на элементы
			sSlice := strings.Split(val, ",")
			// создаем пустой слайс под параметры типа int64 нужного размера
			chkSlc := make([]T, 0, len(sSlice))
			// парсим слайс строковых параметров и переводим его в слайс int64
			for _, s := range sSlice {
				chkVal, err := parseAnyInt[T](s)
				if err != nil {
					return ErrInvalidValidateRuleParam
				}
				chkSlc = append(chkSlc, chkVal)
			}
			// проверяем правило "in"
			if !slices.Contains(chkSlc, data) {
				return ErrNumberNotInSet
			}
		default:
			return ErrInvalidRuleDataType
		}
	}
	return nil
}
