package spannertool

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
	spannerSession "github.com/justdomepaul/toolbox/database/spanner"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/justdomepaul/toolbox/stringtool"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"reflect"
	"strings"
	"sync"
)

const (
	MaxMutations = 20000
	MaxPartition = 32
)

func BatchMutate[T comparable](
	ctx context.Context,
	session spannerSession.ISession,
	inputs []T,
	inputTable string,
	inputColumns []string,
) error {
	var (
		totalInputs = len(inputs)
		wg          = &sync.WaitGroup{}
		mu          = &sync.Mutex{}
		errContents []error
	)
	inputLen_ := math.Ceil(float64(len(inputs)) / MaxPartition)
	inputLen := int(inputLen_)

	for inputStart := 0; inputStart < totalInputs; inputStart += inputLen {
		inputEnd := inputStart + inputLen
		if inputEnd > totalInputs {
			inputEnd = totalInputs
		}
		wg.Add(1)
		go func(m *sync.Mutex, w *sync.WaitGroup, entities []T, table string, columns []string) {
			defer func() {
				w.Done()
				if e := recover(); e != nil {
					if er, ok := e.(error); ok {
						m.Lock()
						errContents = append(errContents, er)
						m.Unlock()
					} else {
						m.Lock()
						errContents = append(errContents, fmt.Errorf("%v", e))
						m.Unlock()
					}
				}
			}()
			_, err := session.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
				perMutations := MaxMutations / len(columns)
				totalEntity := len(entities)
				for start := 0; start < totalEntity; start += perMutations {
					var mut []*spanner.Mutation
					end := start + perMutations
					if end > totalEntity {
						end = totalEntity
					}
					for _, item := range entities[start:end] {
						prepareMut, err := spanner.InsertOrUpdateStruct(table, item)
						if err != nil {
							return err
						}
						mut = append(mut, prepareMut)
					}
					if err := txn.BufferWrite(mut); err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				panic(err)
			}
		}(mu, wg, inputs[inputStart:inputEnd], inputTable, inputColumns)
	}

	wg.Wait()
	if len(errContents) == 0 {
		return nil
	}
	var errContent error
	for _, errItem := range errContents {
		if status.Code(errItem) == codes.AlreadyExists {
			errItem = fmt.Errorf("%w: %s", errorhandler.ErrAlreadyExists, errItem.Error())
		}
		if err := errors.Wrap(errContent, errItem.Error()); err != nil {
			return err
		}
	}

	return errContent
}

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
	if spanner.ErrCode(err) == codes.Internal &&
		strings.Contains(err.Error(), "No result found for statement") {
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

func ValidListArgument(row, page uint64) error {
	validStruct := struct {
		Row  uint64 `validate:"required,min=1"`
		Page uint64 `validate:"required,min=1"`
	}{
		Row:  row,
		Page: page,
	}
	if err := validator.New().Struct(&validStruct); err != nil {
		return fmt.Errorf("%w: %s", errorhandler.ErrInvalidArguments, err.Error())
	}
	return nil
}
