package utils

import (
	"context"
	"github.com/justdomepaul/toolbox/definition"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/metadata"
	"testing"
)

type AccessSuite struct {
	suite.Suite
}

func (suite *AccessSuite) TestGetAccessTokenMethod() {
	ctx := context.Background()
	authCtx := metadata.NewIncomingContext(ctx, metadata.Pairs(definition.AuthorizationKey, definition.AuthorizationType+"testToken"))
	clientID, err := GetAccessToken(authCtx)
	suite.NoError(err)
	suite.Equal("testToken", clientID)
}

func (suite *AccessSuite) TestGetAccessTokenMethodErrIncomingMetadataExist() {
	ctx := context.Background()
	clientID, err := GetAccessToken(ctx)
	suite.ErrorIs(err, errorhandler.ErrIncomingMetadataExist)
	suite.Empty(clientID)
}

func (suite *AccessSuite) TestGetAccessTokenMethodErrAuthorizationRequired() {
	ctx := context.Background()
	authCtx := metadata.NewIncomingContext(ctx, metadata.Pairs("test", ""))
	clientID, err := GetAccessToken(authCtx)
	suite.ErrorIs(err, errorhandler.ErrAuthorizationRequired)
	suite.Empty(clientID)
}

func (suite *AccessSuite) TestGetAccessTokenMethodErrAuthorizationTypeBearer() {
	ctx := context.Background()
	authCtx := metadata.NewIncomingContext(ctx, metadata.Pairs(definition.AuthorizationKey, "Basic test-token"))
	clientID, err := GetAccessToken(authCtx)
	suite.ErrorIs(err, errorhandler.ErrAuthorizationTypeBearer)
	suite.Empty(clientID)
}

func TestAccessSuite(t *testing.T) {
	suite.Run(t, new(AccessSuite))
}
