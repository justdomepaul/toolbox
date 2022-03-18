package utils

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MetadataSuite struct {
	suite.Suite
}

func (suite *MetadataSuite) TestSetClientIDMethodByBytes() {
	ctx := context.Background()
	uid, err := uuid.NewUUID()
	suite.NoError(err)
	newCtx := SetClientID(ctx, "authorization-clientID", uid[:])
	suite.Equal(uid[:], newCtx.Value("authorization-clientID"))
}

func (suite *MetadataSuite) TestSetClientIDMethodByBytesDuplicateKey() {
	ctx := context.Background()
	uid, err := uuid.NewUUID()
	suite.NoError(err)
	uid2, err := uuid.NewUUID()
	suite.NoError(err)
	newCtx := SetClientID(ctx, "authorization-clientID", uid[:])
	newCtx = SetClientID(newCtx, "authorization-clientID", uid2[:])
	suite.NotEqual(uid[:], newCtx.Value("authorization-clientID"))
	suite.Equal(uid2[:], newCtx.Value("authorization-clientID"))
}

func (suite *MetadataSuite) TestGetClientIDMethodByBytes() {
	ctx := context.Background()
	uid, err := uuid.NewUUID()
	suite.NoError(err)
	newCtx := context.WithValue(ctx, "authorization-clientID", uid[:])
	suite.Equal(uid[:], GetClientID[string, []byte](newCtx, "authorization-clientID"))
}

func (suite *MetadataSuite) TestSetClientIDMethodByString() {
	ctx := context.Background()
	uid, err := uuid.NewUUID()
	suite.NoError(err)
	newCtx := SetClientID(ctx, "authorization-clientID", uid.String())
	suite.Equal(uid.String(), newCtx.Value("authorization-clientID"))
}

func (suite *MetadataSuite) TestSetClientIDMethodByStringDuplicateKey() {
	ctx := context.Background()
	uid, err := uuid.NewUUID()
	suite.NoError(err)
	uid2, err := uuid.NewUUID()
	suite.NoError(err)
	newCtx := SetClientID(ctx, "authorization-clientID", uid.String())
	newCtx = SetClientID(newCtx, "authorization-clientID", uid2.String())
	suite.NotEqual(uid.String(), newCtx.Value("authorization-clientID"))
	suite.Equal(uid2.String(), newCtx.Value("authorization-clientID"))
}

func (suite *MetadataSuite) TestGetClientIDMethodByString() {
	ctx := context.Background()
	uid, err := uuid.NewUUID()
	suite.NoError(err)
	newCtx := context.WithValue(ctx, "authorization-clientID", uid.String())
	suite.Equal(uid.String(), GetClientID[string, string](newCtx, "authorization-clientID"))
}

func TestMetadataSuite(t *testing.T) {
	suite.Run(t, new(MetadataSuite))
}
