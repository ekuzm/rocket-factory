package closer

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
	"github.com/sirupsen/logrus"
)

type Func func(context.Context) error

var closer = &Closer{}

type Closer struct {
	mtx   sync.Mutex
	once  sync.Once
	funcs []Func
}

func Add(fn Func) {
	closer.mtx.Lock()
	defer closer.mtx.Unlock()

	closer.funcs = append(closer.funcs, fn)
}

func CloseAll(ctx context.Context) error {
	var out error

	closer.once.Do(func() {
		closer.mtx.Lock()
		funcs := append([]Func{}, closer.funcs...)
		closer.funcs = nil
		closer.mtx.Unlock()

		if len(funcs) == 0 {
			logger.Debug("No closer functions registered")

			return
		}

		var errs []string
		done := make(chan struct{})

		logger.WithFields(logrus.Fields{
			"closerCount": len(funcs),
		}).Debug("Starting closer")

		go func() {
			defer func() {
				close(done)
			}()

			for i := len(funcs) - 1; i >= 0; i-- {
				func() {
					defer func() {
						if r := recover(); r != nil {
							errs = append(errs, fmt.Sprintf("panic: %v", r))
						}
					}()
					if err := funcs[i](ctx); err != nil {
						errs = append(errs, err.Error())
					}
				}()
			}
		}()

		select {
		case <-ctx.Done():
			out = ctx.Err()
		case <-done:
			logger.WithFields(logrus.Fields{
				"closerCount": len(funcs),
				"errorCount":  len(errs),
			}).Debug("Finished closer")

			if len(errs) > 0 {
				out = fmt.Errorf("closed funcs: %v", strings.Join(errs, " | "))
			}
		}
	})

	return out
}
