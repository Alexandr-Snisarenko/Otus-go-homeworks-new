package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

var ErrNotValidatebleType = errors.New("type of field not validateble")

// перечень поддерживаемых типов полей для валидации.
var validatedTypes = []reflect.Kind{
	reflect.Int,
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
	reflect.Uint,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
	reflect.String,
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	sB := new(strings.Builder)
	sB.WriteString("validation errors: ")
	for _, vErr := range v {
		sB.WriteString(fmt.Sprintf("field:'%s'=> %s ", vErr.Field, vErr.Err))
	}

	return sB.String()
}

func (v ValidationErrors) Unwrap() error {
	if len(v) == 0 {
		return nil
	}
	err := errors.New("validation errors: ")
	for _, vErr := range v {
		err = fmt.Errorf("%w field:'%s'=> %w; ", err, vErr.Field, vErr.Err)
	}
	return err
}

func (v *ValidationErrors) AddErr(field string, err error) {
	*v = append(*v, ValidationError{field, err})
}

// функция валидации. на входе ожидаем структуру или указтельна структуру
// завернутый в interface{}.
func Validate(v interface{}) error {
	// если ничего не передали - выходим с ошибкой
	if v == nil {
		return errors.New("input parameter is not defined (is null)")
	}

	rType := reflect.TypeOf(v)
	rVal := reflect.ValueOf(v)

	// если базовый тип объекта - interface
	// получаем объект представленный интерфейсом
	if rType.Kind() == reflect.Interface {
		rType = rType.Elem()
		rVal = rVal.Elem()
	}

	// если базовый тип  объекта - Pointer, то переходим к объекту
	// на который он указывает
	if rType.Kind() == reflect.Pointer {
		rType = rType.Elem()
		rVal = rVal.Elem()
	}

	// если базовый тип входноо парамтера не структура - выходим с ошибкой
	if rType.Kind() != reflect.Struct {
		return errors.New("input parameter is not a struct")
	}

	// обрабатываем поля структуры
	vErr := make(ValidationErrors, 0, rType.NumField())
	for i := 0; i < rType.NumField(); i++ {
		err := ValidateField(rType.Field(i), rVal.Field(i))
		if err != nil {
			vErr = append(vErr, err...)
		}
	}
	// если были ошибки - возвращаем их в интерфейсе error
	if len(vErr) > 0 {
		return vErr
	}

	return nil
}

// валидация поля структуры.
func ValidateField(cField reflect.StructField, cFieldVal reflect.Value) (vErr ValidationErrors) {
	var (
		vStr     string
		vRuleSet ValidateRuleSet
		err      error
		isArrey  bool
	)

	// если поле не имеет тега 'validate' выходим
	if vStr = cField.Tag.Get("validate"); vStr == "" {
		return nil
	}

	// если базовый тип объекта - interface
	// получаем объект представленный интерфейсом
	if cFieldVal.Kind() == reflect.Interface {
		cFieldVal = cFieldVal.Elem()
	}

	// определяем базовый тип объекта
	cFieldKind := cFieldVal.Kind()

	// если объект слайс или массив
	// то определяем базовый тип элемента массива/слайса
	if cFieldVal.Kind() == reflect.Slice || cFieldVal.Kind() == reflect.Array {
		cFieldKind = cFieldVal.Type().Elem().Kind()
		isArrey = true
	}

	// если тип объекта не содержится в перечне поддерживаемых типов - ошибка
	if !slices.Contains(validatedTypes, cFieldKind) {
		vErr.AddErr(cField.Name, ErrNotValidatebleType)
		return vErr
	}

	// формируем список правил по полю
	if vRuleSet, err = NewValidateRuleSet(vStr); err != nil {
		vErr.AddErr(cField.Name, err)
		return vErr
	}

	// если объект не массив
	// проверяем значение объекта по правилам
	// иначе - выполняем обход массива и проверяем каждый его элемент
	if !isArrey {
		err := vRuleSet.CheckData(cFieldVal.Interface())
		if err != nil {
			vErr.AddErr(cField.Name, err)
		}
	} else {
		for i := 0; i < cFieldVal.Len(); i++ {
			// Получаем элемент среза как reflect.Value
			elem := cFieldVal.Index(i)
			// Проверяем элемент
			err := vRuleSet.CheckData(elem.Interface())
			if err != nil {
				vErr.AddErr(cField.Name, err)
			}
		}
	}

	return vErr
}
