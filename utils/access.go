package utils

import (
	"context"
	"github.com/justdomepaul/toolbox/definition"
	"github.com/justdomepaul/toolbox/errorhandler"
	"google.golang.org/grpc/metadata"
	"strings"
)

func GetAccessToken(ctx context.Context) (string, error) {
	md, exist := metadata.FromIncomingContext(ctx)
	if !exist {
		return "", errorhandler.ErrIncomingMetadataExist
	}
	mdCopy := md.Copy()
	authorization := mdCopy.Get(definition.AuthorizationKey)
	if len(authorization) == 0 {
		return "", errorhandler.ErrAuthorizationRequired
	}
	token := ""
	if strings.HasPrefix(authorization[0], definition.AuthorizationType) {
		token = strings.TrimPrefix(authorization[0], definition.AuthorizationType)
	} else {
		return "", errorhandler.ErrAuthorizationTypeBearer
	}
	return token, nil
}
