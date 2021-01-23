package mongodb

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"

	"github.com/eveisesi/athena"
	"github.com/newrelic/go-agent/_integrations/nrmongo"
	"github.com/pkg/errors"
	"github.com/volatiletech/null"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, uri *url.URL) (*mongo.Client, error) {

	monitor := nrmongo.NewCommandMonitor(nil)

	opts := options.Client().ApplyURI(uri.String()).SetMonitor(monitor)
	opts.SetRegistry(customCodecRegistery().Build())
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to mongo db")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping mongo db")
	}

	return client, err

}

// Mongo Operators
const (
	equal            string = "$eq"
	greaterthan      string = "$gt"
	greaterthanequal string = "$gte"
	in               string = "$in"
	lessthan         string = "$lt"
	lessthanequal    string = "$lte"
	notequal         string = "$ne"
	and              string = "$and"
	or               string = "$or"
	exists           string = "$exists"
)

func BuildFilters(operators ...*athena.Operator) primitive.D {

	var ops = make(primitive.D, 0)
	for _, a := range operators {
		switch a.Operation {
		case athena.EqualOp:
			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: equal, Value: a.Value}}})
		case athena.NotEqualOp:
			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: notequal, Value: a.Value}}})
		case athena.GreaterThanOp:
			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: greaterthan, Value: a.Value}}})
		case athena.GreaterThanEqualToOp:
			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: greaterthanequal, Value: a.Value}}})
		case athena.LessThanOp:
			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: lessthan, Value: a.Value}}})
		case athena.LessThanEqualToOp:
			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: lessthanequal, Value: a.Value}}})
		case athena.ExistsOp:
			ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: exists, Value: a.Value.(bool)}}})
		case athena.OrOp:
			switch o := a.Value.(type) {
			case []*athena.Operator:
				arr := make(primitive.A, 0)

				for _, op := range o {
					arr = append(arr, BuildFilters(op))
				}

				ops = append(ops, primitive.E{Key: or, Value: arr})
			default:
				panic(fmt.Sprintf("invalid type %#T supplied, expected one of []*athena.Operator", o))
			}

		case athena.AndOp:
			switch o := a.Value.(type) {
			case []*athena.Operator:
				arr := make(primitive.A, 0)
				for _, op := range o {
					arr = append(arr, BuildFilters(op))
				}

				ops = append(ops, primitive.E{Key: and, Value: arr})
			default:
				panic(fmt.Sprintf("invalid type %#T supplied, expected one of []*athena.Operator", o))
			}

		case athena.InOp:
			v := reflect.ValueOf(a.Value)
			switch v.Kind() {
			case reflect.Slice, reflect.Array:
				arr := make(primitive.A, v.Len())
				for i := 0; i < v.Len(); i++ {
					if !v.Index(i).IsValid() {
						continue
					}
					arr = append(arr, v.Index(i).Interface())
				}

				ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: in, Value: arr}}})
			default:
				panic(fmt.Sprintf("invalid type %#T supplied, expected one of []*athena.OpValue", a.Value))
			}

			// case athena.NotInOp:
			// 	v := reflect.ValueOf(a.Value)
			// 	switch v.Kind() {
			// 	case reflect.Slice, reflect.Array:
			// 		arr := make(primitive.A, v.Len())
			// 		for i := 0; i < v.Len(); i++ {
			// 			if !v.Index(i).IsValid() {
			// 				continue
			// 			}
			// 			arr = append(arr, v.Index(i).Interface())
			// 		}

			// 		ops = append(ops, primitive.E{Key: a.Column, Value: primitive.D{primitive.E{Key: notin, Value: arr}}})
			// 	default:
			// 		panic(fmt.Sprintf("invalid type %#T supplied, expected one of []*athena.OpValue", a.Value))
			// 	}
		}
	}

	return ops

}

func BuildFindOptions(ops ...*athena.Operator) *options.FindOptions {
	var opts = options.Find()
	for _, a := range ops {
		switch a.Operation {
		case athena.LimitOp:
			opts.SetLimit(a.Value.(int64))
		case athena.SkipOp:
			opts.SetSkip(a.Value.(int64))
		case athena.OrderOp:
			opts.SetSort(primitive.D{primitive.E{Key: a.Column, Value: a.Value}})
		}
	}

	return opts
}

func newBool(b bool) *bool {
	return &b
}
func newString(s string) *string {
	return &s
}

var (
	typeNullString  = reflect.TypeOf(null.String{})
	typeNullTime    = reflect.TypeOf(null.Time{})
	typeNullFloat64 = reflect.TypeOf(null.Float64{})
	typeNullInt     = reflect.TypeOf(null.Int{})
	typeNullUint    = reflect.TypeOf(null.Uint{})
	typeNullInt64   = reflect.TypeOf(null.Int64{})
	typeNullUint64  = reflect.TypeOf(null.Uint64{})
)

var allTypes = []reflect.Type{typeNullString, typeNullTime, typeNullFloat64, typeNullUint, typeNullInt, typeNullInt64, typeNullUint64}

func customCodecRegistery() *bsoncodec.RegistryBuilder {

	var primitiveCodecs bson.PrimitiveCodecs
	rb := bsoncodec.NewRegistryBuilder()
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)

	for _, t := range allTypes {
		rb.RegisterTypeEncoder(t, bsoncodec.ValueEncoderFunc(EncodeNullValue))
		rb.RegisterTypeDecoder(t, bsoncodec.ValueDecoderFunc(DecodeNullValue))
	}

	primitiveCodecs.RegisterPrimitiveCodecs(rb)
	return rb

}

func EncodeNullValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {

	matched := false
	for _, rt := range allTypes {
		if val.Type() == rt {
			matched = true
			break
		}
	}

	if !val.IsValid() || !matched {
		return bsoncodec.ValueEncoderError{Name: "EncodeNullValue", Types: allTypes, Received: val}
	}

	switch v := val.Interface().(type) {
	case null.String:
		if v.Valid {
			return vw.WriteString(v.String)
		}

		return vw.WriteNull()
	case null.Time:
		if v.Valid {
			return vw.WriteString(v.Time.Format(time.RFC3339))
		}

		return vw.WriteNull()
	case null.Uint:
		if v.Valid {
			return vw.WriteInt32(int32(v.Uint))
		}

		return vw.WriteInt32(0)
	case null.Float64:
		if v.Valid {
			return vw.WriteDouble(v.Float64)
		}

		return vw.WriteDouble(0.00)
	case null.Int:
		if v.Valid {
			return vw.WriteInt32(int32(v.Int))
		}

		return vw.WriteInt32(0)
	case null.Int64:
		var val int64
		if v.Valid {
			val = v.Int64
		}

		return vw.WriteInt64(val)
	case null.Uint64:
		var val int64
		if v.Valid {
			val = int64(v.Uint64)
		}

		return vw.WriteInt64(val)
	default:
		panic(fmt.Sprintf("EncodeNullValue: unaccounted for type in switch %v ", v))
	}
}

func DecodeNullValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {

	if !val.CanSet() || val.Kind() != reflect.Struct {
		return bsoncodec.ValueDecoderError{Name: "DecodeNullStringValue", Kinds: []reflect.Kind{reflect.Struct}, Received: val}
	}

	switch vr.Type() {
	case bsontype.String:
		str, err := vr.ReadString()
		if err != nil {
			return err
		}

		switch val.Interface().(type) {
		case null.String:
			val.Set(reflect.ValueOf(null.NewString(str, true)))
		case null.Time:

			t, err := time.Parse(time.RFC3339, str)
			if err != nil {
				return err
			}
			val.Set(reflect.ValueOf(null.NewTime(t, true)))

		}
	case bsontype.Double:
		d, err := vr.ReadDouble()
		if err != nil {
			return err
		}

		val.Set(reflect.ValueOf(null.NewFloat64(d, true)))
	case bsontype.Int32:
		d, err := vr.ReadInt32()
		if err != nil {
			return err
		}

		switch val.Interface().(type) {
		case null.Uint:
			val.Set(reflect.ValueOf(null.NewUint(uint(d), true)))
		default:
			return fmt.Errorf("[DecodeNullValue]: unhandled integer conversion %s to %T", vr.Type(), val.Interface())
		}
	case bsontype.DateTime:
		d, err := vr.ReadDateTime()
		if err != nil {
			return err
		}

		val.Set(reflect.ValueOf(null.NewTime(time.Unix(0, d*int64(time.Millisecond)), true)))

	case bsontype.Null:
		err := vr.ReadNull()
		if err != nil {
			return err
		}

		switch val.Interface().(type) {
		case null.String:
			val.Set(reflect.ValueOf(null.StringFromPtr(nil)))
		case null.Time:
			val.Set(reflect.ValueOf(null.TimeFromPtr(nil)))
		case null.Uint:
			val.Set(reflect.ValueOf(null.UintFromPtr(nil)))
		case null.Uint64:
			val.Set(reflect.ValueOf(null.Uint64FromPtr(nil)))
		case null.Float64:
			val.Set(reflect.ValueOf(null.Float64FromPtr(nil)))
		}

	default:
		return fmt.Errorf("[DecodeNullValue]: don't know how to decode %s", vr.Type())

	}

	return nil

}
