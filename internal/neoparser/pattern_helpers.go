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

func SepBy[T, S any](
	item Pattern[T],
	separator Pattern[S],
) Pattern[[]T] {
	return PatternFunc[[]T](func(ctx *Context) Match[[]T] {
		start := ctx.Pos
		values := []T{}

		first := item.Match(ctx)

		if first.Kind == NoMatch {
			ctx.Reset(start)

			return Match[[]T]{
				Kind:  Matched,
				Value: values,
				Start: start,
				End:   start,
			}
		}

		if first.Kind == Failed {
			return Match[[]T]{
				Kind:  Failed,
				Start: start,
				End:   first.End,
				Err:   first.Err,
			}
		}

		values = append(values, first.Value)

		for {
			beforeSeparator := ctx.Pos

			sep := separator.Match(ctx)

			if sep.Kind == NoMatch {
				ctx.Reset(beforeSeparator)
				break
			}

			if sep.Kind == Failed {
				return Match[[]T]{
					Kind:  Failed,
					Start: start,
					End:   sep.End,
					Err:   sep.Err,
				}
			}

			next := item.Match(ctx)

			if next.Kind != Matched {
				return Match[[]T]{
					Kind:  Failed,
					Start: start,
					End:   ctx.Pos,
					Err: &ParseError{
						Position: ctx.Pos,
						Message:  "expected item after separator",
					},
				}
			}

			values = append(values, next.Value)
		}

		return Match[[]T]{
			Kind:  Matched,
			Value: values,
			Start: start,
			End:   ctx.Pos,
		}
	})
}

func SepBy1[T, S any](
	item Pattern[T],
	separator Pattern[S],
) Pattern[[]T] {
	return PatternFunc[[]T](func(ctx *Context) Match[[]T] {
		start := ctx.Pos

		result := SepBy(item, separator).Match(ctx)

		if result.Kind != Matched {
			return result
		}

		if len(result.Value) == 0 {
			ctx.Reset(start)

			return Match[[]T]{
				Kind:  NoMatch,
				Start: start,
				End:   start,
			}
		}

		return result
	})
}

func Between[L, T, R any](
	left Pattern[L],
	content Pattern[T],
	right Pattern[R],
) Pattern[T] {
	return PatternFunc[T](func(ctx *Context) Match[T] {
		start := ctx.Pos

		l := left.Match(ctx)
		if l.Kind != Matched {
			ctx.Reset(start)

			return Match[T]{
				Kind:  l.Kind,
				Start: start,
				End:   l.End,
				Err:   l.Err,
			}
		}

		value := content.Match(ctx)
		if value.Kind != Matched {
			return Match[T]{
				Kind:  Failed,
				Start: start,
				End:   value.End,
				Err:   value.Err,
			}
		}

		r := right.Match(ctx)
		if r.Kind != Matched {
			return Match[T]{
				Kind:  Failed,
				Start: start,
				End:   r.End,
				Err:   r.Err,
			}
		}

		return Match[T]{
			Kind:  Matched,
			Value: value.Value,
			Start: start,
			End:   ctx.Pos,
		}
	})
}

func OneOf[T any](patterns ...Pattern[T]) Pattern[T] {
	return PatternFunc[T](func(ctx *Context) Match[T] {
		start := ctx.Pos

		var best *Match[T]

		for _, pattern := range patterns {
			ctx.Reset(start)

			result := pattern.Match(ctx)

			if result.Kind == Matched {
				return result
			}

			if result.Kind == Failed {
				if best == nil || result.End > best.End {
					copy := result
					best = &copy
				}
			}
		}

		ctx.Reset(start)

		if best != nil {
			return *best
		}

		return Match[T]{
			Kind:  NoMatch,
			Start: start,
			End:   start,
		}
	})
}
