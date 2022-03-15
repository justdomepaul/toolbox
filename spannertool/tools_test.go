package spannertool

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/gocql/gocql"
	"github.com/justdomepaul/toolbox/errorhandler"
	itestutil "github.com/justdomepaul/toolbox/mock/go/testutil"
	spannerutil "github.com/justdomepaul/toolbox/mock/spanner/testutil"
	"github.com/justdomepaul/toolbox/timestamp"
	"github.com/stretchr/testify/suite"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"testing"
)

func setupMockedTestServer(t *testing.T) (server *spannerutil.MockedSpannerInMemTestServer, client *spanner.Client, teardown func()) {
	return setupMockedTestServerWithConfig(t, spanner.ClientConfig{})
}

func setupMockedTestServerWithConfig(t *testing.T, config spanner.ClientConfig) (server *spannerutil.MockedSpannerInMemTestServer, client *spanner.Client, teardown func()) {
	return setupMockedTestServerWithConfigAndClientOptions(t, config, []option.ClientOption{})
}

func setupMockedTestServerWithConfigAndClientOptions(t *testing.T, config spanner.ClientConfig, clientOptions []option.ClientOption) (server *spannerutil.MockedSpannerInMemTestServer, client *spanner.Client, teardown func()) {
	grpcHeaderChecker := &itestutil.HeadersEnforcer{
		OnFailure: t.Fatalf,
		Checkers: []*itestutil.HeaderChecker{
			{
				Key: "x-goog-api-client",
				ValuesValidator: func(token ...string) error {
					if len(token) != 1 {
						return status.Errorf(codes.Internal, "unexpected number of api client token headers: %v", len(token))
					}
					if !strings.HasPrefix(token[0], "gl-go/") {
						return status.Errorf(codes.Internal, "unexpected api client token: %v", token[0])
					}
					if !strings.Contains(token[0], "gccl/") {
						return status.Errorf(codes.Internal, "unexpected api client token: %v", token[0])
					}
					return nil
				},
			},
		},
	}
	clientOptions = append(clientOptions, grpcHeaderChecker.CallOptions()...)
	server, opts, serverTeardown := spannerutil.NewMockedSpannerInMemTestServer(t)
	opts = append(opts, clientOptions...)
	ctx := context.Background()
	formattedDatabase := fmt.Sprintf("projects/%s/instances/%s/databases/%s", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	client, err := spanner.NewClientWithConfig(ctx, formattedDatabase, config, opts...)
	if err != nil {
		t.Fatal(err)
	}
	return server, client, func() {
		client.Close()
		serverTeardown()
	}
}

type mockBatchMutate struct {
	BAR int `spanner:"BAR" json:"BAR,omitempty"`
}

type BatchMutateSuite struct {
	suite.Suite
}

func (suite *BatchMutateSuite) TestBatchMutate() {
	ctx := context.Background()
	_, client, teardown := setupMockedTestServer(suite.T())
	defer teardown()

	var inputs []mockBatchMutate
	for i := 1; i < 20000; i++ {
		inputs = append(inputs, mockBatchMutate{
			BAR: i,
		})
	}

	suite.NoError(BatchMutate(ctx, client, inputs, "FOO", []string{"BAR"}))
}

func TestBatchMutateSuite(t *testing.T) {
	suite.Run(t, new(BatchMutateSuite))
}

type templateSpannerMock struct {
	ID           gocql.UUID `spanner:"id" json:"id"`
	PartyToken   string     `spanner:"party_token" json:"party_token"`
	HashPhone    string     `spanner:"hash_phone" json:"hash_phone"`
	RawData      []byte     `spanner:"raw_data" json:"raw_data"` // json will to base64 string
	ReceivedTime int64      `spanner:"received_time" json:"received_time"`
	CreatedTime  int64      `spanner:"created_time" json:"created_time"`
}

type FetchSpannerTagValueSuite struct {
	suite.Suite
}

func (suite *FetchSpannerTagValueSuite) TestGetColumns() {
	id := gocql.TimeUUID()
	columns, placeholder, params := FetchSpannerTagValue(templateSpannerMock{
		ID:           id,
		PartyToken:   "testPlatformToken",
		HashPhone:    "testHashPhone",
		RawData:      []byte("ttt"),
		ReceivedTime: 0,
	}, true, "created_time")
	suite.Equal(Columns{"id", "party_token", "hash_phone", "raw_data"}, columns)
	suite.Equal(Placeholders{"@id", "@party_token", "@hash_phone", "@raw_data"}, placeholder)
	suite.Equal(Parameters{
		"id":          id,
		"party_token": "testPlatformToken",
		"hash_phone":  "testHashPhone",
		"raw_data":    []byte("ttt"),
	}, params)
	suite.T().Log(columns)
	suite.T().Log(placeholder)
	suite.T().Log(params)
}

func (suite *FetchSpannerTagValueSuite) TestGetColumns1() {
	columns, placeholder, params := FetchSpannerTagValue(templateSpannerMock{
		ID:         gocql.TimeUUID(),
		PartyToken: "testPlatformToken",
	}, true, "created_time")
	suite.Equal(Columns{"id", "party_token"}, columns)
	suite.Equal(Placeholders{"@id", "@party_token"}, placeholder)
	suite.T().Log(params)
}

func (suite *FetchSpannerTagValueSuite) TestGetColumns2() {
	columns, placeholder, params := FetchSpannerTagValue(templateSpannerMock{
		ID:           gocql.TimeUUID(),
		PartyToken:   "testPlatformToken",
		ReceivedTime: timestamp.GetNowTimestamp(),
	}, true, "created_time")
	suite.Equal(Columns{"id", "party_token", "received_time"}, columns)
	suite.Equal(Placeholders{"@id", "@party_token", "@received_time"}, placeholder)
	suite.T().Log(params)
}

func (suite *FetchSpannerTagValueSuite) TestGetColumns3() {
	columns, placeholder, params := FetchSpannerTagValue(templateSpannerMock{
		ID:           gocql.TimeUUID(),
		PartyToken:   "testPlatformToken",
		ReceivedTime: timestamp.GetNowTimestamp(),
		CreatedTime:  timestamp.GetNowTimestamp(),
	}, true, "created_time")
	suite.Equal(Columns{"id", "party_token", "received_time", "created_time"}, columns)
	suite.Equal(Placeholders{"@id", "@party_token", "@received_time", "CURRENT_TIMESTAMP()"}, placeholder)
	suite.T().Log(params)
}

func TestFetchSpannerTagValueSuite(t *testing.T) {
	suite.Run(t, new(FetchSpannerTagValueSuite))
}

type MockSpannerIterator struct {
	err error
}

func (m *MockSpannerIterator) Next() (*spanner.Row, error) {
	return &spanner.Row{}, m.err
}

type MockEntity struct {
}

type SpannerIteratorRowSuite struct {
	suite.Suite
}

func (suite *SpannerIteratorRowSuite) TestGetIteratorFirstRow() {
	entity := MockEntity{}
	suite.NoError(GetIteratorFirstRow(&MockSpannerIterator{}, &entity))
}

func (suite *SpannerIteratorRowSuite) TestGetIteratorFirstRowIteratorDone() {
	entity := MockEntity{}
	suite.Error(GetIteratorFirstRow(&MockSpannerIterator{
		err: iterator.Done,
	}, &entity), iterator.Done)
}

func (suite *SpannerIteratorRowSuite) TestGetIteratorFirstRowHaveError() {
	entity := MockEntity{}
	suite.Error(GetIteratorFirstRow(&MockSpannerIterator{
		err: errors.New("got error"),
	}, &entity), errors.New("got error"))
}

func TestSpannerIteratorRowSuite(t *testing.T) {
	suite.Run(t, new(SpannerIteratorRowSuite))
}

type ExecuteAndCheckEffectRowOverZeroSuite struct {
	suite.Suite
}

func (suite *ExecuteAndCheckEffectRowOverZeroSuite) TestExecuteAndCheckEffectRowOverZero() {
	ctx := context.Background()
	_, client, teardown := setupMockedTestServer(suite.T())
	defer teardown()
	_, err := client.ReadWriteTransaction(ctx, func(c context.Context, txn *spanner.ReadWriteTransaction) error {
		return ExecuteAndCheckEffectRowOverZero(txn, c, spanner.Statement{SQL: `UPDATE FOO SET BAR=1 WHERE BAZ=2`})
	})
	suite.NoError(err)
}

func (suite *ExecuteAndCheckEffectRowOverZeroSuite) TestExecuteAndCheckEffectRowOverZeroNoEffect() {
	ctx := context.Background()
	_, client, teardown := setupMockedTestServer(suite.T())
	defer teardown()
	_, err := client.ReadWriteTransaction(ctx, func(c context.Context, txn *spanner.ReadWriteTransaction) error {
		return ExecuteAndCheckEffectRowOverZero(txn, c, spanner.Statement{SQL: `UPDATE FOO SET BAR=1 WHERE BAZ=1`})
	})
	suite.Error(err, errorhandler.ErrUpdateNoEffect)
}

func (suite *ExecuteAndCheckEffectRowOverZeroSuite) TestExecuteAndCheckEffectRowOverZeroError() {
	ctx := context.Background()
	_, client, teardown := setupMockedTestServer(suite.T())
	defer teardown()
	_, err := client.ReadWriteTransaction(ctx, func(c context.Context, txn *spanner.ReadWriteTransaction) error {
		return ExecuteAndCheckEffectRowOverZero(txn, c, spanner.Statement{})
	})
	suite.Error(err)
}

func TestExecuteAndCheckEffectRowOverZeroSuite(t *testing.T) {
	suite.Run(t, new(ExecuteAndCheckEffectRowOverZeroSuite))
}
