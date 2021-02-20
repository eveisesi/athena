package null

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/volatiletech/null"
)

func MarshalUint(nu null.Uint) graphql.Marshaler {
	if !nu.Valid {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatUint(uint64(nu.Uint), 10))
	})
}

func UnmarshalUint(i interface{}) (null.Uint, error) {

	if i == nil {
		return null.Uint{}, nil
	}

	switch i := i.(type) {
	case uint:
		return null.NewUint(i, i > 0), nil
	default:
		return null.Uint{}, fmt.Errorf("%v is not a valid uint", i)
	}

}
