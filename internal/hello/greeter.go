package hello

import "github.com/ezphlow/go_template/internal/models"

// Service is the production Greeter.
type Service struct{}

// Greet returns "hello, <name>", defaulting an empty name to "world".
func (Service) Greet(name string) models.Hello {
	if name == "" {
		name = "world"
	}
	return models.Hello{Message: "hello, " + name}
}
