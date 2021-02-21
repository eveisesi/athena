package mysqldb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/athena"
	"github.com/go-sql-driver/mysql"
)

func New(ctx context.Context, host string, port int, user, pass, dbname string) (*sql.DB, error) {

	config := mysql.Config{
		User:         user,
		Passwd:       pass,
		Net:          "tcp",
		Addr:         fmt.Sprintf("%s:%d", host, port),
		DBName:       dbname,
		Loc:          time.UTC,
		Timeout:      time.Second * 5,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		ParseTime:    true,
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("[MySQL Connect] Failed to connect to mysql server: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("[MySQL Connect] Failed to ping mysql server")
	}

	return db, nil

}

func BuildFilters(s sq.SelectBuilder, operators ...*athena.Operator) sq.SelectBuilder {
	for _, a := range operators {
		if !a.Operation.IsValid() {
			continue
		}

		switch a.Operation {
		case athena.EqualOp:
			s = s.Where(sq.Eq{a.Column: a.Value})
		case athena.NotEqualOp:
			s = s.Where(sq.NotEq{a.Column: a.Value})
		case athena.GreaterThanEqualToOp:
			s = s.Where(sq.GtOrEq{a.Column: a.Value})
		case athena.GreaterThanOp:
			s = s.Where(sq.Gt{a.Column: a.Value})
		case athena.LessThanEqualToOp:
			s = s.Where(sq.LtOrEq{a.Column: a.Value})
		case athena.LessThanOp:
			s = s.Where(sq.Lt{a.Column: a.Value})
		case athena.InOp:
			s = s.Where(sq.Eq{a.Column: a.Value.(interface{})})
		case athena.NotInOp:
			s = s.Where(sq.NotEq{a.Column: a.Value.([]interface{})})
		case athena.LikeOp:
			s = s.Where(sq.Like{a.Column: fmt.Sprintf("%%%v%%", a.Value)})
		case athena.OrderOp:
			s = s.OrderBy(fmt.Sprintf("%s %s", a.Column, a.Value))
		case athena.LimitOp:
			s = s.Limit(uint64(a.Value.(int64)))
		case athena.SkipOp:
			s = s.Offset(uint64(a.Value.(int64)))
		}
	}

	return s

}
