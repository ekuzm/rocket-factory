package closer

import (
	"context"
	"errors"
	"fmt"
)

func New() *Closer {
	return &Closer{}
}

type Closer struct {
	funcs []func() error
}

func (c *Closer) Add(fn func() error) {
	c.funcs = append(c.funcs, fn)
}

func (c *Closer) Close(ctx context.Context) error {
	var out error

	for len(c.funcs) > 0 {
		select {
		case <-ctx.Done():
			return errors.Join(ctx.Err(), out)
		default:
		}

		curr := c.funcs[len(c.funcs)-1]
		c.funcs = c.funcs[:len(c.funcs)-1]

		out = errors.Join(out, safety(curr))
	}

	return out
}

func safety(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during close: %v", r)
		}
	}()

	return fn()
}

func (c *Closer) IsEmpty() bool {
	return len(c.funcs) == 0
}
