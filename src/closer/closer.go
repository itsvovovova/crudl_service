package closer

import (
	"errors"
	"sync"
)

type Closer struct {
	once  sync.Once
	mu    sync.Mutex
	funcs []func() error
}

func (c *Closer) Add(f func() error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f)
}

func (c *Closer) Close() error {
	var err error
	c.once.Do(func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		var errs []error
		for i := len(c.funcs) - 1; i >= 0; i-- {
			if e := c.funcs[i](); e != nil {
				errs = append(errs, e)
			}
		}
		err = errors.Join(errs...)
	})
	return err
}
