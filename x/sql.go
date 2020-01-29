package x

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"

	"github.com/masteryconnect/pipe/message"
)

// SQL will transform a message to a query.
// It is batch aware
type SQL struct {
	Table    string   // optional table to apply the mutation to
	MaskKeys []string // useful for logging without sensitive values like passwords
}

// ErrSQLTypeConversionError is the error
var ErrSQLTypeConversionError = errors.New("unknown type to convert to SQL")

// T implements the Tfunc interface
func (s SQL) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	for m := range in {
		var ok bool
		var b message.Batch
		if b, ok = m.(message.Batch); !ok {
			b = message.Batch{m}
		}
		q, err := s.SQLInsertFromBatch(b)
		if err != nil {
			errs <- err
		} else {
			out <- q
		}
	}
}

// I implements the InlineTfunc interface
func (s SQL) I(m interface{}) (interface{}, error) {
	var ok bool
	var b message.Batch
	if b, ok = m.(message.Batch); !ok {
		b = message.Batch{m}
	}
	return s.SQLInsertFromBatch(b)
}

// SQLInsertFromBatch converts a batch message to a bulk INSERT query message
func (s SQL) SQLInsertFromBatch(v message.Batch) (*message.Query, error) {
	var keys []string
	table := s.Table
	sql := "INSERT INTO %s (%s) VALUES %s"
	rowPlaceholders := []string{}
	vals := []interface{}{}
	for _, m := range v {
		switch v := m.(type) {
		case *message.Record:
			if keys == nil {
				keys = v.Keys
			}
			qs := make([]string, len(v.Vals))
			for i := range qs {
				qs[i] = "?"
			}
			rowPlaceholders = append(rowPlaceholders, "("+strings.Join(qs, ",")+")")
			vals = append(vals, v.Vals...)
		default:
			return nil, errors.Wrapf(ErrSQLTypeConversionError, "got type %T", m)
		}
	}
	return &message.Query{SQL: fmt.Sprintf(
		sql,
		table,
		strings.Join(keys, ","),
		strings.Join(rowPlaceholders, ","),
	), Args: vals}, nil
}
