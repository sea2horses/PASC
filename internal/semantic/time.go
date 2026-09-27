package semantic

type TimeUnit uint8

const (
	SEGUNDOS TimeUnit = iota
	MILISEGUNDOS
)

func (t TimeUnit) String() string {
	switch t {
	case SEGUNDOS:
		return "segundos"
	case MILISEGUNDOS:
		return "milisegundos"
	default:
		return ""
	}
}
