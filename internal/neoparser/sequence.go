package neoparser

type Pair[A, B any] struct {
	First  A
	Second B
}

func Seq2[A, B any](
	a Pattern[A],
	b Pattern[B],
) Pattern[Pair[A, B]] {
	return PatternFunc[Pair[A, B]](func(ctx *Context) Match[Pair[A, B]] {
		start := ctx.Pos

		ra := a.Match(ctx)
		if ra.Kind != Matched {
			ctx.Reset(start)

			return Match[Pair[A, B]]{
				Kind:  ra.Kind,
				Start: start,
				End:   ra.End,
				Err:   ra.Err,
			}
		}

		rb := b.Match(ctx)
		if rb.Kind != Matched {
			ctx.Reset(start)

			return Match[Pair[A, B]]{
				Kind:  rb.Kind,
				Start: start,
				End:   rb.End,
				Err:   rb.Err,
			}
		}

		return Match[Pair[A, B]]{
			Kind: Matched,
			Value: Pair[A, B]{
				First:  ra.Value,
				Second: rb.Value,
			},
			Start: start,
			End:   ctx.Pos,
		}
	})
}
