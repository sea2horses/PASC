package neoparser

func Map[A, B any](p Pattern[A], fn func(A) B) Pattern[B] {
	return PatternFunc[B](func(ctx *Context) Match[B] {
		result := p.Match(ctx)

		if result.Kind != Matched {
			return Match[B]{
				Kind:  result.Kind,
				Start: result.Start,
				End:   result.End,
				Err:   result.Err,
			}
		}

		return Match[B]{
			Kind:  Matched,
			Value: fn(result.Value),
			Start: result.Start,
			End:   result.End,
		}
	})
}

type OptionalValue[T any] struct {
	Value T
	Some  bool
}

func Optional[T any](p Pattern[T]) Pattern[OptionalValue[T]] {
	return PatternFunc[OptionalValue[T]](func(ctx *Context) Match[OptionalValue[T]] {
		start := ctx.Pos
		result := p.Match(ctx)

		switch result.Kind {
		case Matched:
			return Match[OptionalValue[T]]{
				Kind: Matched,
				Value: OptionalValue[T]{
					Value: result.Value,
					Some:  true,
				},
				Start: result.Start,
				End:   result.End,
			}
		case Failed:
			return Match[OptionalValue[T]]{
				Kind:  Failed,
				Start: result.Start,
				End:   result.End,
				Err:   result.Err,
			}
		default:
			ctx.Reset(start)

			return Match[OptionalValue[T]]{
				Kind: Matched,
				Value: OptionalValue[T]{
					Some: false,
				},
				Start: start,
				End:   start,
			}
		}
	})
}
