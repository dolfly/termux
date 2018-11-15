package termux

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/eternal-flame-AD/go-termux/internal/api"
	_io "github.com/eternal-flame-AD/go-termux/internal/io"
)

var GlobalTimeout = 15 * time.Second

func tryCloseReader(c io.Reader) {
	if r, ok := c.(_io.CloseReader); ok {
		r.CloseRead()
	}
}
func tryCloseWriter(c io.Writer) {
	if r, ok := c.(_io.CloseWriter); ok {
		r.CloseWrite()
	}
}

func pipe(ai io.Reader, ao io.Writer, bi io.Reader, bo io.Writer) {
	wg := sync.WaitGroup{}

	if bi != nil && ao != nil {
		wg.Add(1)
		go func() {
			io.Copy(ao, bi)
			tryCloseReader(bi)
			tryCloseWriter(ao)
			wg.Done()
		}()
	}
	if ai != nil && bo != nil {
		wg.Add(1)
		go func() {
			io.Copy(bo, ai)
			tryCloseReader(ai)
			tryCloseWriter(bo)
			wg.Done()
		}()
	}
	wg.Wait()
}

func execAction(method string, stdin io.Reader, stdout io.Writer, action string) {
	ctx, cancel := context.WithTimeout(context.Background(), GlobalTimeout)
	defer cancel()
	execActionContext(ctx, stdin, stdout, method, action)
}

func execActionContext(ctx context.Context, stdin io.Reader, stdout io.Writer, method string, action string) {
	call := api.Call{
		Method: method,
		Action: action,
	}

	call.Call(ctx)
	pipe(call, call, stdin, stdout)
}

func exec(stdin io.Reader, stdout io.Writer, method string, args map[string]interface{}, data string) {
	ctx, cancel := context.WithTimeout(context.Background(), GlobalTimeout)
	defer cancel()
	execContext(ctx, stdin, stdout, method, args, data)
}

func execContext(ctx context.Context, stdin io.Reader, stdout io.Writer, method string, args map[string]interface{}, data string) {
	call := api.Call{
		Method: method,
		Args:   args,
		Data:   data,
	}

	call.Call(ctx)
	defer call.Close()
	pipe(call, call, stdin, stdout)
}