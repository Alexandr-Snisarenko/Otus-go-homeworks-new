package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

var (
	ErrOsEnvironmentVariableFormat = errors.New(" environment variable string don't contain '=' ")
	ErrNoCommand                   = errors.New("call without command")
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (string, error) {
	var (
		newEnv []string
		err    error
		osCmd  *exec.Cmd
		out    strings.Builder
	)

	if newEnv, err = parseOsEnv(env); err != nil {
		return "", err
	}

	if len(cmd) == 0 {
		return "", ErrNoCommand
	}

	command := cmd[0]
	if len(cmd) == 1 {
		osCmd = exec.Command(command)
	} else {
		osCmd = exec.Command(command, cmd[1:]...)
	}
	osCmd.Env = newEnv
	osCmd.Stdout = &out
	err = osCmd.Run()
	if err != nil {
		return "", err
	}

	// возвращаем stdout выполненной команды
	return out.String(), nil
}

// функция парсировки текущего системного списка переменных окружения
// и формирования нового списка переменных окружения
// с учетом входного словаря переменных "на замену" (env).
// логика следующая:
// 1 - создаем новый пустой список для пересенных окружения.
// 2 - в новый список переносим все переменные окружения из системного списка, которых нет во входном словаре (env).
// 3 - обрабатываем входной словарь (env). добавляем в новый список те записи словаря env,
// у которых NeedRemove == false.
func parseOsEnv(env Environment) ([]string, error) {
	// получаем список текущих переменных окружения
	osEnv := os.Environ()
	// массив для новых переменных окружения (максимально возможной длинны)
	newEnv := make([]string, 0, len(osEnv)+len(env))

	// проходим по текущему системному списку переменных окружения
	for _, osEnv := range osEnv {
		// определяем имя переменной по символу "="
		// если "=" нет - unexpected error (такого быть не может по идее)
		before, _, ok := strings.Cut(osEnv, "=")
		if !ok {
			return nil, ErrOsEnvironmentVariableFormat
		}

		// ищем переменную окружения в нашем словаре
		// если нашли - пропускаем (они обработаются отдельным блоком)
		// если не нашли - переносим в новый массив as is
		if _, ok := env[before]; !ok {
			newEnv = append(newEnv, osEnv)
		}
	}

	// проходим по нашему словарю переменных окружения
	// если признак NeedRemove == false - добавляем переменную в новый список
	for eVal := range env {
		if !env[eVal].NeedRemove {
			newEnv = append(newEnv, eVal+"="+env[eVal].Value)
		}
	}

	return newEnv, nil
}
