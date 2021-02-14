package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/eveisesi/athena"
	"github.com/urfave/cli"
)

// func testCommand(c *cli.Context) error {

// 	t := athena.CorporationAllianceHistory{}
// 	ref := "history"

// 	rt := reflect.TypeOf(t)

// 	tags := []string{}
// 	fmt.Printf("Insert Values\n")
// 	for i := 0; i < rt.NumField(); i++ {
// 		f := rt.Field(i)
// 		tag := f.Tag.Get("db")
// 		if tag == "-" {
// 			continue
// 		}
// 		parts := strings.Split(tag, ",")
// 		quotedTag := fmt.Sprintf("\"%s\"", parts[0])
// 		tags = append(tags, quotedTag)
// 		fmt.Printf("%s.%s,\n", ref, f.Name)
// 	}
// 	fmt.Printf("\nUpdate Sets\n")
// 	for i := 0; i < rt.NumField(); i++ {
// 		f := rt.Field(i)
// 		tag := f.Tag.Get("db")
// 		if tag == "-" {
// 			continue
// 		}
// 		parts := strings.Split(tag, ",")
// 		quotedTag := fmt.Sprintf("\"%s\"", parts[0])
// 		fmt.Printf("Set(%s, %s.%s).\n", quotedTag, ref, f.Name)
// 	}

// 	fmt.Printf("\n%s\n", strings.Join(tags, ","))

// 	return nil

// }

func testCommand(c *cli.Context) error {

	contact := &athena.MemberContact{}
	d, err := json.Marshal(contact)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(d))

	//	fmt.Println(len(t), cap(t))
	return nil
}
