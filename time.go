package athena

import "github.com/volatiletech/null"

type NullTimeZeroer null.Time

func (t NullTimeZeroer) IsZero() bool {
	return !t.Valid
}
