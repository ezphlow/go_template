package models

// Hello is the JSON body for GET /v1/hello.
type Hello struct {
	Message string `json:"message"`
}

// Greeter produces a Hello. Implementations live in internal/hello;
// test doubles live in internal/mocks.
type Greeter interface {
	Greet(name string) Hello
}
