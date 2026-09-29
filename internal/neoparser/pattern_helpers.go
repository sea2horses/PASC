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

func Many[T any](p Pattern[T]) Pattern[[]T] {
	return PatternFunc[[]T](func(ctx *Context) Match[[]T] {
		start := ctx.Pos
		values := []T{}

		for {
			before := ctx.Pos
			result := p.Match(ctx)

			if result.Kind == Failed {
				return Match[[]T]{
					Kind:  Failed,
					Start: start,
					End:   result.End,
					Err:   result.Err,
				}
			}

			if result.Kind == NoMatch {
				ctx.Reset(before)
				break
			}

			if ctx.Pos == before {
				panic("Many() received a pattern that matched without consuming input")
			}

			values = append(values, result.Value)
		}

		return Match[[]T]{
			Kind:  Matched,
			Value: values,
			Start: start,
			End:   ctx.Pos,
		}
	})
}

func Many1[T any](p Pattern[T]) Pattern[[]T] {
	return PatternFunc[[]T](func(ctx *Context) Match[[]T] {
		start := ctx.Pos

		first := p.Match(ctx)

		if first.Kind != Matched {
			ctx.Reset(start)

			return Match[[]T]{
				Kind:  first.Kind,
				Start: start,
				End:   first.End,
				Err:   first.Err,
			}
		}

		values := []T{first.Value}

		rest := Many(p).Match(ctx)

		if rest.Kind == Failed {
			return rest
		}

		values = append(values, rest.Value...)

		return Match[[]T]{
			Kind:  Matched,
			Value: values,
			Start: start,
			End:   ctx.Pos,
		}
	})
}
