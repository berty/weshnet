package queue

type MetricsTracer[T any] interface {
	ItemQueued(name string, item T)
	ItemPop(name string, item T)
}

type noopTracer[T any] struct{}

func (*noopTracer[T]) ItemQueued(name string, item T) {}
func (*noopTracer[T]) ItemPop(name string, item T)    {}
