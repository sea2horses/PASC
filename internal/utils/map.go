package utils

// Source - https://stackoverflow.com/a/71624929
// Posted by jub0bs, modified by community. See post 'Timeline' for change history
// Retrieved 2026-09-24, License - CC BY-SA 4.0

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}
