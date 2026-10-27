package mocks

import "github.com/ezphlow/go_template/internal/models"

// Greeter is a hand-written models.Greeter for tests (XxxFn field pattern).
type Greeter struct {
	GreetFn func(name string) models.Hello
}

// Greet calls GreetFn when set.
func (g Greeter) Greet(name string) models.Hello {
	if g.GreetFn != nil {
		return g.GreetFn(name)
	}
	return models.Hello{}
}
