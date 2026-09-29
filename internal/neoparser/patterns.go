package neoparser

type Pattern[T any] interface {
	Match(*Context) Match[T]
}

type PatternFunc[T any] func(*Context) Match[T]

func (f PatternFunc[T]) Match(c *Context) Match[T] {
	return f(c)
}
