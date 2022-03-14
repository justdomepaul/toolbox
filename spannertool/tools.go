package spannertool

import (
	"cloud.google.com/go/spanner"
	"context"
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/justdomepaul/toolbox/stringtool"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	"strings"
)

type Columns []string

func (c Columns) String() string { return strings.Join(c, ", ") }

type Placeholders []string

func (p Placeholders) String() string { return strings.Join(p, ", ") }

type Parameters map[string]interface{}

func FetchSpannerTagValue(input interface{}, ignoreZero bool, createdTimeKeyWord string) (Columns, Placeholders, Parameters) {
	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)
	columns := make([]string, 0)
	placeholder := make([]string, 0)
	params := make(map[string]interface{}, 0)
	for i := 0; i < t.NumField(); i++ {
		if v.Field(i).IsValid() && (!ignoreZero || !v.Field(i).IsZero()) {
			c := t.Field(i).Tag.Get("spanner")
			columns = append(columns, c)
			if c != createdTimeKeyWord {
				placeholder = append(placeholder, stringtool.StringJoin("@", t.Field(i).Tag.Get("spanner")))
			} else {
				placeholder = append(placeholder, "CURRENT_TIMESTAMP()")
				continue
			}
			params[c] = v.Field(i).Interface()
		}
	}
	return columns, placeholder, params
}

func ExecuteAndCheckEffectRowOverZero(txn *spanner.ReadWriteTransaction, ctx context.Context, stmt spanner.Statement) error {
	effectRow, err := txn.Update(ctx, stmt)
	if gRPCErrorStatus, ok := status.FromError(err); ok &&
		gRPCErrorStatus.Code() == codes.Internal &&
		strings.Contains(gRPCErrorStatus.Message(), "No result found for statement") {
		return errorhandler.ErrUpdateNoEffect
	}
	if err != nil {
		return err
	}
	if effectRow == 0 {
		return errorhandler.ErrUpdateNoEffect
	}
	return nil
}

type ISpannerRowIterator interface {
	Next() (*spanner.Row, error)
}

// GetIteratorFirstRow func
// entity must use reference
func GetIteratorFirstRow(iter ISpannerRowIterator, entity interface{}) error {
	row, err := iter.Next()
	if errors.Is(err, iterator.Done) {
		return errorhandler.ErrNoRows
	}
	if err != nil {
		return err
	}
	return row.ToStruct(entity)
}
