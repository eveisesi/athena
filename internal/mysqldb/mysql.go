package mysqldb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Connect(ctx context.Context, dsn string, driver string) (*sqlx.DB, error) {

	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("[MySQL Connect] Failed to parse provided DSN into valid configuration: %w", err)
	}

	db, err := sql.Open(driver, config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("[MySQL Connect] Failed to connect to mysql server: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("[MySQL Connect] Failed to ping mysql server")
	}

	return sqlx.NewDb(db, driver), nil

}

// func BuildFilters(operators ...*athena.Operator) primitive.D {
// 	// Mongo Like Statement https://stackoverflow.com/questions/55398792/find-entries-via-substring-regex-query-in-mongodb-go-driver
// 	var ops = make(primitive.D, 0)
// 	for _, a := range operators {
// 		switch a.Operation {
// 		case athena.EqualOp:
// 			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: equal, Value: a.Value}}})
// 		case athena.NotEqualOp:
// 			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: notequal, Value: a.Value}}})
// 		case athena.GreaterThanOp:
// 			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: greaterthan, Value: a.Value}}})
// 		case athena.GreaterThanEqualToOp:
// 			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: greaterthanequal, Value: a.Value}}})
// 		case athena.LessThanOp:
// 			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: lessthan, Value: a.Value}}})
// 		case athena.LessThanEqualToOp:
// 			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: lessthanequal, Value: a.Value}}})
// 		case athena.ExistsOp:
// 			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: exists, Value: a.Value.(bool)}}})
// 		case athena.OrOp:
// 			switch o := a.Value.(type) {
// 			case []*athena.Operator:
// 				arr := make(primitive.A, 0)

// 				for _, op := range o {
// 					arr = append(arr, BuildFilters(op))
// 				}

// 				ops = append(ops, primitive.E{Key: or, Value: arr})
// 			default:
// 				panic(fmt.Sprintf("invalid type %#T supplied, expected one of []*athena.Operator", o))
// 			}

// 		case athena.AndOp:
// 			switch o := a.Value.(type) {
// 			case []*athena.Operator:
// 				arr := make(primitive.A, 0)
// 				for _, op := range o {
// 					arr = append(arr, BuildFilters(op))
// 				}

// 				ops = append(ops, primitive.E{Key: and, Value: arr})
// 			default:
// 				panic(fmt.Sprintf("invalid type %#T supplied, expected one of []*athena.Operator", o))
// 			}

// 		case athena.InOp:
// 			v := reflect.ValueOf(a.Value)
// 			switch v.Kind() {
// 			case reflect.Slice, reflect.Array:
// 				arr := make(primitive.A, v.Len())
// 				for i := 0; i < v.Len(); i++ {
// 					if !v.Index(i).IsValid() {
// 						continue
// 					}
// 					arr = append(arr, v.Index(i).Interface())
// 				}

// 				ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: in, Value: arr}}})
// 			default:
// 				panic(fmt.Sprintf("invalid type %#T supplied, expected one of []*athena.OpValue", a.Value))
// 			}

// 			// case athena.NotInOp:
// 			// 	v := reflect.ValueOf(a.Value)
// 			// 	switch v.Kind() {
// 			// 	case reflect.Slice, reflect.Array:
// 			// 		arr := make(primitive.A, v.Len())
// 			// 		for i := 0; i < v.Len(); i++ {
// 			// 			if !v.Index(i).IsValid() {
// 			// 				continue
// 			// 			}
// 			// 			arr = append(arr, v.Index(i).Interface())
// 			// 		}

// 			// 		ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: notin, Value: arr}}})
// 			// 	default:
// 			// 		panic(fmt.Sprintf("invalid type %#T supplied, expected one of []*athena.OpValue", a.Value))
// 			// 	}
// 		}
// 	}

// 	return ops

// }

// func BuildFindOptions(ops ...*athena.Operator) *options.FindOptions {
// 	var opts = options.Find()
// 	for _, a := range ops {
// 		switch a.Operation {
// 		case athena.LimitOp:
// 			opts.SetLimit(a.Value.(int64))
// 		case athena.SkipOp:
// 			opts.SetSkip(a.Value.(int64))
// 		case athena.OrderOp:
// 			opts.SetSort(primitive.D{primitive.E{Key: a.Column, Value: a.Value}})
// 		}
// 	}

// 	return opts
// }
