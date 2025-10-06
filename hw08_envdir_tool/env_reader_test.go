package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// подготовка временной директории и тестовых файлов
// для тестирования функции ReadDir.
// готовим два файла.
// первый "REMOVE_ENV" - пустой. для удаления переменной окружения "REMOVE_ENV".
// второй "SET_ENV" для устанавки переменной окружения "SET_ENV" в значение "testEnvValue".
func PrepeareTestDir() (string, error) {
	var (
		dir   string
		fName string
		err   error
		file  *os.File
	)

	dir, err = os.MkdirTemp("", "readDirTest")
	if err != nil {
		return "", fmt.Errorf(" can't create temp dir 'readDirTest': %w", err)
	}

	fName = filepath.Join(dir, "REMOVE_ENV")
	file, err = os.Create(fName)
	if err != nil {
		return "", fmt.Errorf("can't create file REMOVE_ENV: %w", err)
	}
	defer file.Close()

	fName = filepath.Join(dir, "SET_ENV")
	file, err = os.Create(fName)
	if err != nil {
		return "", fmt.Errorf("can't create file SET_ENV: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString("testEnvValue")
	if err != nil {
		return "", fmt.Errorf("can't write string in file 'SET_ENV': %w", err)
	}

	return dir, nil
}

func TestReadDir(t *testing.T) {
	// готовим тестовое окружения
	dir, err := PrepeareTestDir()
	if err != nil {
		t.Errorf(" run PrepeareTestDir failed: %v", err)
	}
	defer os.RemoveAll(dir) // при выходе очищаем тестове окружение

	env, err := ReadDir(dir)
	if err != nil {
		t.Errorf("run ReadDir failed: %v", err)
	}

	if envVal, ok := env["REMOVE_ENV"]; !ok {
		t.Errorf("no variable 'REMOVE_ENV' in result environment set")
	} else {
		require.Equal(t, true, envVal.NeedRemove)
	}

	if envVal, ok := env["SET_ENV"]; !ok {
		t.Errorf("no variable 'SET_ENV' in result environment set")
	} else {
		require.Equal(t, false, envVal.NeedRemove)
		require.Equal(t, "testEnvValue", envVal.Value)
	}
}
