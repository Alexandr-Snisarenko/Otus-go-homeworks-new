package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

var (
	ErrConnectionAlreadyActive = errors.New("connection is already active")
	ErrConnectionNotActive     = errors.New("connection is not active")
)

// интерфейс чуть поменялся. методы send и receive не включены в интерфейс,
// так как в реализации используются потоковые данные.
type TelnetClient interface {
	Connect(context.Context) error
	Close() error
	IsActive() bool
	Done() <-chan struct{}
	// Send() error
	// Receive() error
}

type telnetClient struct {
	mu        sync.Mutex
	active    bool          // признак активности коннекта
	done      chan struct{} // канал завершения коннекта
	address   string
	timeout   time.Duration
	inReader  io.Reader
	outWriter io.Writer
	errWriter io.Writer
	conn      net.Conn
}

func (t *telnetClient) Connect(ctx context.Context) error {
	if t.active {
		return ErrConnectionAlreadyActive
	}

	// пробуем подключиться с указанным таймаутом
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}
	// если всё ок - переводим  active в true и создаем сигнальный канал
	// код не защищаем - выполняется в начале работы и только из этого метода.
	t.conn = conn
	t.active = true
	t.done = make(chan struct{})

	// запускаем основную горутину клиента
	if err := t.startClient(ctx); err != nil {
		return err
	}

	return nil
}

func (t *telnetClient) Close() error {
	// меняем статус коннекта и соответственно закрываем сигнальный канал в мютексе
	// так как закрытие может выполняться из разных горутин
	// множественный вызов метода Close - допустим и не является проблемой
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.active {
		return nil
	}
	t.active = false
	close(t.done)

	return t.conn.Close()
}

func (t *telnetClient) IsActive() bool {
	return t.active
}

func (t *telnetClient) Done() <-chan struct{} {
	return t.done
}

// создаем новый объект клиента телнет. создается только новая структура без подключения.
func NewTelnetClient(address string, timeout time.Duration,
	inReader io.Reader, outWriter io.Writer, errWriter io.Writer,
) TelnetClient {
	return &telnetClient{
		active:    false,
		address:   address,
		timeout:   timeout,
		inReader:  inReader,
		outWriter: outWriter,
		errWriter: errWriter,
	}
}

// основная горутина объекта.
// запускает две дочерние горутины на чтение и на запись в потоки ввода вывода
// и ждет окончания работы: по сигнальному каналу объекта или по сигналу от контекста.
func (t *telnetClient) startClient(ctx context.Context) error {
	if !t.active {
		return ErrConnectionNotActive
	}

	go func() {
		// запускаем горутину для чтения из входного потока телнет клиента (inReader) во входной поток соединения (connection)
		t.inBuf2outBuff(t.inReader, t.conn)
		// запускаем горутину для чтения из выходного потока соединения в выходной поток объекта (outReader)
		t.inBuf2outBuff(t.conn, t.outWriter)

		select {
		case <-t.done:
			return
		case <-ctx.Done():
			t.Close()
			return
		}
	}()
	return nil
}

// горутина переносит данные из одного буферизированного потока в другой.
func (t *telnetClient) inBuf2outBuff(inBuf io.Reader, outBuf io.Writer) {
	go func() {
		reader := bufio.NewReader(inBuf)
		for {
			select {
			case <-t.done:
				return
			default:
				text, err := reader.ReadString('\n')
				if err != nil {
					t.errWriter.Write([]byte(err.Error() + "\n"))
					t.Close()
					return
				}

				_, err = outBuf.Write([]byte(text))
				if err != nil {
					t.errWriter.Write([]byte(err.Error() + "\n"))
					t.Close()
					return
				}
			}
		}
	}()
}
