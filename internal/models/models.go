package models

type Position struct {
	Offset uint32
	Line   uint32
	Column uint32
}

type Span struct {
	Start Position
	End   Position
}

/* Join spans!! */
func JoinSpans(spans ...Span) Span {
	var lowest_position *Position = nil
	var highest_position *Position = nil
	
	for _, span := range spans {
		if lowest_position == nil || span.Start.Offset < lowest_position.Offset {
			lowest_position = &span.Start
		}

		if highest_position == nil || span.End.Offset > highest_position.Offset {
			highest_position = &span.End
		}
	}

	return Span{
		Start: *lowest_position,
		End: *highest_position,
	}
}
