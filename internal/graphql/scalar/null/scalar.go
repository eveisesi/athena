package null

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/volatiletech/null"
)

func MarshalString(ns null.String) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if !ns.Valid {
			_, _ = io.WriteString(w, `null`)
			return
		}

		_, _ = io.WriteString(w, ns.String)
	})
}

func UnmarshalString(i interface{}) (null.String, error) {
	switch v := i.(type) {
	case string:
		if v == "null" {
			return null.NewString("", false), nil
		}
		return null.NewString(v, true), nil
	default:
		return null.NewString("", false), fmt.Errorf("%v is not a valid string", v)
	}
}
