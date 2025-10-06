package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var user User
	exit := false
	result := make(DomainStat)
	// псоздаем экземпляр bufio.Reader от интерфейса io.Reader
	rder := bufio.NewReader(r)

	for !exit {
		// читаем строку из буфера rder
		line, err := rder.ReadString('\n')
		if err != nil {
			// если ошибка == EOF ставим метку выхода и обрабатываем роследнюю строку
			if errors.Is(err, io.EOF) {
				exit = true
			} else {
				return nil, err
			}
		}

		// если в строке нт подстроки с доменом - переходим к следующей строке
		if !strings.Contains(line, "."+domain) {
			continue
		}

		// заполняем объект user из json содержимого текущей строки
		if err = jsoniter.Unmarshal([]byte(line), &user); err != nil {
			return nil, err
		}

		// выделяем домен из email (предполагаем, что email - валидный)
		_, userDomain, fnd := strings.Cut(strings.ToLower(user.Email), "@")
		// если не нашли собачки (email пустой, например) - переходим к следующей строке
		if !fnd {
			continue
		}

		// проверяем относится ли домен пользователя к указанному домену первого уровня
		_, matched := strings.CutSuffix(userDomain, "."+domain)
		// если да - инкрементируем запись по домену в мапе
		if matched {
			result[userDomain]++
		}
	}

	return result, nil
}
