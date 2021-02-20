package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/eveisesi/athena"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli"
)

func testCommand(c *cli.Context) error {

	t := athena.Alliance{}

	rt := reflect.TypeOf(t)
	fmt.Printf("type %s @goModel(model: \"%s.%s\") {\n", rt.Name(), rt.PkgPath(), rt.Name())

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tag := f.Tag.Get("json")
		if tag == "-" {
			continue
		}
		tagParts := strings.Split(tag, ",")
		if tagParts[0] == "-" {
			continue
		}
		tag = tagParts[0]
		if tag == "updated_at" || tag == "created_at" {
			continue
		}
		tag = strcase.ToLowerCamel(tag)

		if strings.Contains(tag, "Id") {
			tag = strings.ReplaceAll(tag, "Id", "ID")

		}

		t := graphqlTypeFromGoType(f.Type.String())

		fmt.Printf("\t%s: %s\n", tag, t)
	}

	fmt.Printf("}\n")

	return nil

}

func graphqlTypeFromGoType(t string) string {
	p := strings.Split(t, ".")
	nullable := false
	if len(p) == 2 {
		// If this split into to parts successfully
		// then we got something like time.Time or null.String for t
		t = strings.ToLower(p[1])
		if p[0] == "null" {
			// I'm checking for a very specific null package here
			nullable = true
		}
	}

	out := ""
	switch t {
	case "int8", "int16", "int32", "int64", "int":
		out = "Int"
	case "uint8", "uint16", "uint32", "uint64", "uint":
		out = "Uint"
	case "string":
		out = "String"
	case "float32", "float64":
		out = "Float"
	case "time":
		out = "Time"
	case "bool":
		out = "Bool"
	}

	if !nullable {
		out = out + "!"
	}

	return out
}

// func testCommand(c *cli.Context) error {

// 	contact := &athena.MemberContact{}
// 	d, err := json.Marshal(contact)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(string(d))

// 	//	fmt.Println(len(t), cap(t))
// 	return nil
// }
