package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/eveisesi/athena"
	"github.com/urfave/cli"
)

func testCommand(c *cli.Context) error {

	t := athena.MemberFittingItem{}
	ref := "item"

	rt := reflect.TypeOf(t)

	tags := []string{}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tag := f.Tag.Get("db")
		if tag == "-" {
			continue
		}
		parts := strings.Split(tag, ",")
		tags = append(tags, fmt.Sprintf("\"%s\"", parts[0]))
		fmt.Printf("%s.%s,\n", ref, f.Name)
	}

	fmt.Println(strings.Join(tags, ","))

	return nil

}
