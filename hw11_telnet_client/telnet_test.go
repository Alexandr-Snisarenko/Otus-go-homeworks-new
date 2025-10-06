package main

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			// для тестирования создаем потоковые данные через объекты io.pipe
			inReader, inWriter := io.Pipe()
			outReader, outWriter := io.Pipe()
			errReader, errWriter := io.Pipe()

			ctx := context.Background()

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, inReader, outWriter, errWriter)
			require.NoError(t, client.Connect(ctx))
			defer func() { require.NoError(t, client.Close()) }()

			_, err = inWriter.Write([]byte("hello\n"))
			require.NoError(t, err)

			reader := bufio.NewReader(outReader)
			s, err := reader.ReadString('\n')
			require.NoError(t, err)
			require.Equal(t, "world\n", s)

			errOut := bufio.NewReader(errReader)
			s, err = errOut.ReadString('\n')
			require.NoError(t, err)
			require.Equal(t, "EOF\n", s)

			inWriter.Close()
			outWriter.Close()
			errReader.Close()
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}
