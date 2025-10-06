package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileExists            = errors.New("destination file already exists")
	ErrParams                = errors.New("not valid parameter")
)

func checkCopyParams(fromPath, toPath string, offset, limit int64, rewrite bool) error {
	var (
		inFileInfo, outFileInfo fs.FileInfo
		err                     error
	)

	// базовые проверки значений параметров
	switch {
	case fromPath == "":
		return fmt.Errorf("%w : name of file to read can't be empty", ErrParams)
	case toPath == "":
		return fmt.Errorf("%w : name of file to write can't be empty", ErrParams)
	case limit < 0:
		return fmt.Errorf("%w : limit can't be negative", ErrParams)
	case offset < 0:
		return fmt.Errorf("%w : offset can't be negative", ErrParams)
	}

	/// проверки параметров файлов ///
	inFileInfo, err = os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("can't open file %s : %w", fromPath, err)
	}

	// обрабатываем только регулярные файлы
	if mode := inFileInfo.Mode(); !mode.IsRegular() {
		return ErrUnsupportedFile
	}

	// если свдиг больше длины файла  - ошибка
	if offset > inFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	// если целевой файл существует - проверяем можно ли его перезаписать (флаг rewrite и тип файла)
	outFileInfo, err = os.Stat(toPath)
	if err == nil {
		// проверяем флаг разрешения перезаписи
		if !rewrite {
			return ErrFileExists
		}
		// проверяем тип файла обрабатываем только регулярные файлы
		if mode := outFileInfo.Mode(); !mode.IsRegular() {
			return fmt.Errorf("can't rewrite destination file %s: %w", toPath, ErrUnsupportedFile)
		}
	}

	return nil
}

func Copy(fromPath, toPath string, offset, limit int64, rewrite bool) error {
	var (
		inFile, outFile *os.File
		inFileInfo      fs.FileInfo
		err             error
	)

	// проверяем корректность параметров и файлов
	if err := checkCopyParams(fromPath, toPath, offset, limit, rewrite); err != nil {
		return err
	}

	inFileInfo, err = os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("can't open file %s : %w", fromPath, err)
	}

	// если исходный файл нулевой длины или offset == file size - копировать нечего, создаём пустой выходнойфайл
	if inFileInfo.Size() == 0 || inFileInfo.Size() == offset {
		outFile, err = os.Create(toPath)
		if err != nil {
			return fmt.Errorf("can't create destination file: %w", err)
		}
		err = outFile.Close()
		return err
	}

	// если данные для копирования есть - выполняем копирование
	inFile, err = os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("can't open source file: %w", err)
	}
	defer inFile.Close()

	// пишем во временный файл. потом переименовываем
	outFile, err = os.Create(toPath + ".tmp")
	if err != nil {
		return fmt.Errorf("can't create destination file: %w", err)
	}
	defer outFile.Close()

	fmt.Println("Copying ", fromPath, " to ", toPath)

	// если задан сдвиг - сдвигаем ридер
	if offset > 0 {
		inFile.Seek(offset, io.SeekStart)
	}

	// если лимит не задан, или лимит, с учетом сдвига, больше размера файла,
	// то limit (размер копируемых данных) = размер файла - сдвиг
	if limit == 0 || offset+limit > inFileInfo.Size() {
		limit = inFileInfo.Size() - offset
	}

	// прогресс бар копирования
	bar := progressbar.DefaultBytes(limit)

	buf := make([]byte, 1024) //  буфер обмена
	exit := false             // флаг завершения копирования
	cntAll := int64(0)        // счетчик скопированных байт

	// основной цикл копирования данных. выходим по лимиту или концу файла
	for !exit {
		cnt, err := inFile.Read(buf)
		cntAll += int64(cnt)

		if err != nil {
			if errors.Is(err, io.EOF) {
				return fmt.Errorf("error read from source file: %w", err)
			}
			exit = true
		}

		// проверяем на limit. если лимит првышен, обрезаем крайнюю порцию по размеру превышения и ставим флаг "на выход"
		if cntAll >= limit {
			cnt -= int(cntAll - limit)
			exit = true
		}

		_, err = outFile.Write(buf[:cnt])
		if err != nil {
			return fmt.Errorf("error write to destination file: %w", err)
		}
		// актуализируем прогрессбар
		bar.Add(cnt)
	}

	if err := outFile.Close(); err != nil {
		return fmt.Errorf("cant close temporary destination file: %w", err)
	}

	if err := os.Rename(outFile.Name(), toPath); err != nil {
		return fmt.Errorf("error rename temporary file to destination file: %w", err)
	}

	fmt.Println("Copying completed")
	return nil
}
