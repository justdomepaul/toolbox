package authenticate

import (
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	pb "github.com/grpc-ecosystem/go-grpc-middleware/testing/testproto"
	"github.com/justdomepaul/toolbox/definition"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/justdomepaul/toolbox/jwt"
	"github.com/justdomepaul/toolbox/key"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
	"time"
)

type testService struct {
	pb.UnimplementedTestServiceServer
}

func (t *testService) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println(md)
	return &pb.PingResponse{}, nil
}

func (t *testService) PingList(req *pb.PingRequest, srv pb.TestService_PingListServer) error {
	log.Println(req)
	return srv.Send(&pb.PingResponse{})
}

func newHubCredentials(jwtToken string) (cred metadataCredentials) {
	cred.metadata = map[string]string{definition.AuthorizationKey: definition.AuthorizationType + " " + jwtToken}
	return
}

type metadataCredentials struct {
	metadata map[string]string
}

func (p metadataCredentials) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return p.metadata, nil
}

func (p metadataCredentials) RequireTransportSecurity() bool {
	return false
}

type InterceptorSuite struct {
	suite.Suite
	hs384key   string
	hs384Token string
}

func (suite *InterceptorSuite) SetupSuite() {
	suite.hs384key = "HwwCUU7B00LU000-NINrJPtuglzBF54fWmvAc4K8wlma4hRz4lB2HBxEFpPld8bA"
	hs384, err := jwt.NewHS384JWT(key.ToBinaryRunes(suite.hs384key))
	suite.NoError(err)
	suite.T().Log(hs384)
}

func (suite *InterceptorSuite) TestInterceptorAllowed() {
	ctx := context.Background()
	clientID := []byte("clientID")
	fullPath := "/mwitkow.testproto.TestService/Ping"

	authorization := &testIAuthorization{}
	authorization.On("GetID").Return(clientID)
	authorization.On("GetClaim").Return(mock.AnythingOfType("interface{}"))
	authenticate := &testIAuthenticate{}
	authenticate.On("Authenticate", mock.AnythingOfType("func() (string, error)"), fullPath).
		Return(authorization, errorhandler.ErrInWhitelist)

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(authenticate)),
		grpc.StreamInterceptor(StreamServerInterceptor(authenticate)),
	)
	pb.RegisterTestServiceServer(srv, &testService{})

	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	// it is here to properly stop the server
	defer func() { time.Sleep(10 * time.Millisecond) }()
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	defer srv.Stop()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(getBufDialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(newHubCredentials(suite.hs384Token)),
	)
	suite.NoError(err)
	defer conn.Close()
	client := pb.NewTestServiceClient(conn)
	resp, err := client.Ping(ctx, &pb.PingRequest{})
	suite.NoError(err)
	suite.T().Log(resp)
}

func (suite *InterceptorSuite) TestInterceptor() {
	ctx := context.Background()
	clientID := []byte("clientID")
	fullPath := "/mwitkow.testproto.TestService/Ping"

	authorization := &testIAuthorization{}
	authorization.On("GetID").Return(clientID)
	authorization.On("GetClaim").Return(mock.AnythingOfType("interface{}"))
	authenticate := &testIAuthenticate{}
	authenticate.On("Authenticate", mock.AnythingOfType("func() (string, error)"), fullPath).
		Return(authorization, nil)

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(authenticate)),
		grpc.StreamInterceptor(StreamServerInterceptor(authenticate)),
	)
	pb.RegisterTestServiceServer(srv, &testService{})

	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	// it is here to properly stop the server
	defer func() { time.Sleep(10 * time.Millisecond) }()
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	defer srv.Stop()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(getBufDialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(newHubCredentials(suite.hs384Token)),
	)
	suite.NoError(err)
	defer conn.Close()
	client := pb.NewTestServiceClient(conn)
	resp, err := client.Ping(ctx, &pb.PingRequest{})
	suite.NoError(err)
	suite.T().Log(resp)
}

func (suite *InterceptorSuite) TestInterceptorGRPCError() {
	ctx := context.Background()
	clientID := []byte("clientID")
	fullPath := "/mwitkow.testproto.TestService/Ping"

	authorization := &testIAuthorization{}
	authorization.On("GetID").Return(clientID)
	authorization.On("GetClaim").Return(mock.AnythingOfType("interface{}"))
	authenticate := &testIAuthenticate{}
	authenticate.On("Authenticate", mock.AnythingOfType("func() (string, error)"), fullPath).
		Return(authorization, status.Error(codes.Aborted, "got error"))

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(authenticate)),
		grpc.StreamInterceptor(StreamServerInterceptor(authenticate)),
	)
	pb.RegisterTestServiceServer(srv, &testService{})

	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	// it is here to properly stop the server
	defer func() { time.Sleep(10 * time.Millisecond) }()
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	defer srv.Stop()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(getBufDialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(newHubCredentials(suite.hs384Token)),
	)
	suite.NoError(err)
	defer conn.Close()
	client := pb.NewTestServiceClient(conn)
	resp, err := client.Ping(ctx, &pb.PingRequest{})
	suite.Error(err)
	suite.T().Log(resp)
}

func (suite *InterceptorSuite) TestInterceptorAuthenticateErrTokenExpired() {
	ctx := context.Background()
	clientID := []byte("clientID")
	fullPath := "/mwitkow.testproto.TestService/Ping"

	authorization := &testIAuthorization{}
	authorization.On("GetID").Return(clientID)
	authorization.On("GetClaim").Return(mock.AnythingOfType("interface{}"))
	authenticate := &testIAuthenticate{}
	authenticate.On("Authenticate", mock.AnythingOfType("func() (string, error)"), fullPath).
		Return(authorization, fmt.Errorf("%w: "+jwt.ErrTokenExpired.Error(), errorhandler.ErrUnauthenticated))

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(authenticate)),
		grpc.StreamInterceptor(StreamServerInterceptor(authenticate)),
	)
	pb.RegisterTestServiceServer(srv, &testService{})

	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	// it is here to properly stop the server
	defer func() { time.Sleep(10 * time.Millisecond) }()
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	defer srv.Stop()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(getBufDialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(newHubCredentials("anyToken")),
	)
	suite.NoError(err)
	defer conn.Close()
	client := pb.NewTestServiceClient(conn)
	resp, err := client.Ping(ctx, &pb.PingRequest{})
	suite.Error(err, errorhandler.ErrUnauthenticated)
	suite.T().Log(resp)
}

func (suite *InterceptorSuite) TestInterceptorAuthenticateErrRefreshTokenNotExist() {
	ctx := context.Background()
	clientID := []byte("clientID")
	fullPath := "/mwitkow.testproto.TestService/Ping"

	authorization := &testIAuthorization{}
	authorization.On("GetID").Return(clientID)
	authorization.On("GetClaim").Return(mock.AnythingOfType("interface{}"))
	authenticate := &testIAuthenticate{}
	authenticate.On("Authenticate", mock.AnythingOfType("func() (string, error)"), fullPath).
		Return(authorization, errorhandler.ErrNoRefreshToken)

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(authenticate)),
		grpc.StreamInterceptor(StreamServerInterceptor(authenticate)),
	)
	pb.RegisterTestServiceServer(srv, &testService{})

	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	// it is here to properly stop the server
	defer func() { time.Sleep(10 * time.Millisecond) }()
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	defer srv.Stop()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(getBufDialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(newHubCredentials("anyToken")),
	)
	suite.NoError(err)
	defer conn.Close()
	client := pb.NewTestServiceClient(conn)
	resp, err := client.Ping(ctx, &pb.PingRequest{})
	suite.Error(err, errorhandler.ErrUnauthenticated)
	suite.T().Log(resp)
}

func (suite *InterceptorSuite) TestInterceptorAuthenticateErrScopeNotExist() {
	ctx := context.Background()
	clientID := []byte("clientID")
	fullPath := "/mwitkow.testproto.TestService/Ping"

	authorization := &testIAuthorization{}
	authorization.On("GetID").Return(clientID)
	authorization.On("GetClaim").Return(mock.AnythingOfType("interface{}"))
	authenticate := &testIAuthenticate{}
	authenticate.On("Authenticate", mock.AnythingOfType("func() (string, error)"), fullPath).
		Return(authorization, errorhandler.ErrScopeNotExist)

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(authenticate)),
		grpc.StreamInterceptor(StreamServerInterceptor(authenticate)),
	)
	pb.RegisterTestServiceServer(srv, &testService{})

	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	// it is here to properly stop the server
	defer func() { time.Sleep(10 * time.Millisecond) }()
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	defer srv.Stop()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(getBufDialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(newHubCredentials("anyToken")),
	)
	suite.NoError(err)
	defer conn.Close()
	client := pb.NewTestServiceClient(conn)
	resp, err := client.Ping(ctx, &pb.PingRequest{})
	suite.Error(err, errorhandler.ErrInvalidArguments)
	suite.T().Log(resp)
}

func (suite *InterceptorSuite) TestInterceptorAuthenticateErrNoRows() {
	ctx := context.Background()
	clientID := []byte("clientID")
	fullPath := "/mwitkow.testproto.TestService/Ping"

	authorization := &testIAuthorization{}
	authorization.On("GetID").Return(clientID)
	authorization.On("GetClaim").Return(mock.AnythingOfType("interface{}"))
	authenticate := &testIAuthenticate{}
	authenticate.On("Authenticate", mock.AnythingOfType("func() (string, error)"), fullPath).
		Return(authorization, fmt.Errorf("%w: "+errorhandler.ErrNoRows.Error(), errorhandler.ErrUnauthenticated))

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(authenticate)),
		grpc.StreamInterceptor(StreamServerInterceptor(authenticate)),
	)
	pb.RegisterTestServiceServer(srv, &testService{})

	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	// it is here to properly stop the server
	defer func() { time.Sleep(10 * time.Millisecond) }()
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	defer srv.Stop()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(getBufDialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(newHubCredentials("anyToken")),
	)
	suite.NoError(err)
	defer conn.Close()
	client := pb.NewTestServiceClient(conn)
	resp, err := client.Ping(ctx, &pb.PingRequest{})
	suite.Error(err, errorhandler.ErrUnauthenticated)
	suite.T().Log(resp)
}

func (suite *InterceptorSuite) TestInterceptorAuthenticateErrScopeNotAllowed() {
	ctx := context.Background()
	clientID := []byte("clientID")
	fullPath := "/mwitkow.testproto.TestService/Ping"

	authorization := &testIAuthorization{}
	authorization.On("GetID").Return(clientID)
	authorization.On("GetClaim").Return(mock.AnythingOfType("interface{}"))
	authenticate := &testIAuthenticate{}
	authenticate.On("Authenticate", mock.AnythingOfType("func() (string, error)"), fullPath).
		Return(authorization, errorhandler.ErrOutOfScopes)

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(authenticate)),
		grpc.StreamInterceptor(StreamServerInterceptor(authenticate)),
	)
	pb.RegisterTestServiceServer(srv, &testService{})

	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	// it is here to properly stop the server
	defer func() { time.Sleep(10 * time.Millisecond) }()
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	defer srv.Stop()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(getBufDialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(newHubCredentials("anyToken")),
	)
	suite.NoError(err)
	defer conn.Close()
	client := pb.NewTestServiceClient(conn)
	resp, err := client.Ping(ctx, &pb.PingRequest{})
	suite.Error(err, errorhandler.ErrNoPermission)
	suite.T().Log(resp)
}

func (suite *InterceptorSuite) TestInterceptorAuthenticateError() {
	ctx := context.Background()
	clientID := []byte("clientID")
	fullPath := "/mwitkow.testproto.TestService/Ping"

	authorization := &testIAuthorization{}
	authorization.On("GetID").Return(clientID)
	authorization.On("GetClaim").Return(mock.AnythingOfType("interface{}"))
	authenticate := &testIAuthenticate{}
	authenticate.On("Authenticate", mock.AnythingOfType("func() (string, error)"), fullPath).
		Return(authorization, errors.New("got error"))

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(authenticate)),
		grpc.StreamInterceptor(StreamServerInterceptor(authenticate)),
	)
	pb.RegisterTestServiceServer(srv, &testService{})

	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	// it is here to properly stop the server
	defer func() { time.Sleep(10 * time.Millisecond) }()
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	defer srv.Stop()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(getBufDialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(newHubCredentials("anyToken")),
	)
	suite.NoError(err)
	defer conn.Close()
	client := pb.NewTestServiceClient(conn)
	resp, err := client.Ping(ctx, &pb.PingRequest{})
	suite.Error(err)
	suite.T().Log(err)
	suite.T().Log(resp)
}

func (suite *InterceptorSuite) TestInterceptorStreamAllowed() {
	ctx := context.Background()
	clientID := []byte("clientID")
	fullPath := "/mwitkow.testproto.TestService/PingList"

	authorization := &testIAuthorization{}
	authorization.On("GetID").Return(clientID)
	authorization.On("GetClaim").Return(mock.AnythingOfType("interface{}"))
	authenticate := &testIAuthenticate{}
	authenticate.On("Authenticate", mock.AnythingOfType("func() (string, error)"), fullPath).
		Return(authorization, errorhandler.ErrInWhitelist)

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(authenticate)),
		grpc.StreamInterceptor(StreamServerInterceptor(authenticate)),
	)
	pb.RegisterTestServiceServer(srv, &testService{})

	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	// it is here to properly stop the server
	defer func() { time.Sleep(10 * time.Millisecond) }()
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	defer srv.Stop()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(getBufDialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(newHubCredentials(suite.hs384Token)),
	)
	suite.NoError(err)
	defer conn.Close()
	client := pb.NewTestServiceClient(conn)
	stream, err := client.PingList(ctx, &pb.PingRequest{})
	suite.NoError(err)
	defer stream.CloseSend()
	resp, err := stream.Recv()
	suite.NoError(err)
	suite.T().Log(resp)
}

func (suite *InterceptorSuite) TestInterceptorStream() {
	ctx := context.Background()
	clientID := []byte("clientID")
	fullPath := "/mwitkow.testproto.TestService/PingList"

	authorization := &testIAuthorization{}
	authorization.On("GetID").Return(clientID)
	authorization.On("GetClaim").Return(mock.AnythingOfType("interface{}"))
	authenticate := &testIAuthenticate{}
	authenticate.On("Authenticate", mock.AnythingOfType("func() (string, error)"), fullPath).
		Return(authorization, nil)

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(authenticate)),
		grpc.StreamInterceptor(StreamServerInterceptor(authenticate)),
	)
	pb.RegisterTestServiceServer(srv, &testService{})

	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	// it is here to properly stop the server
	defer func() { time.Sleep(10 * time.Millisecond) }()
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	defer srv.Stop()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(getBufDialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(newHubCredentials(suite.hs384Token)),
	)
	suite.NoError(err)
	defer conn.Close()
	client := pb.NewTestServiceClient(conn)
	stream, err := client.PingList(ctx, &pb.PingRequest{})
	suite.NoError(err)
	defer stream.CloseSend()
	resp, err := stream.Recv()
	suite.NoError(err)
	suite.T().Log(resp)
}

func (suite *InterceptorSuite) TestInterceptorStreamAuthenticateError() {
	ctx := context.Background()
	clientID := []byte("clientID")
	fullPath := "/mwitkow.testproto.TestService/PingList"

	authorization := &testIAuthorization{}
	authorization.On("GetID").Return(clientID)
	authorization.On("GetClaim").Return(mock.AnythingOfType("interface{}"))
	authenticate := &testIAuthenticate{}
	authenticate.On("Authenticate", mock.AnythingOfType("func() (string, error)"), fullPath).
		Return(authorization, errors.New("got error"))

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(authenticate)),
		grpc.StreamInterceptor(StreamServerInterceptor(authenticate)),
	)
	pb.RegisterTestServiceServer(srv, &testService{})

	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)
	// it is here to properly stop the server
	defer func() { time.Sleep(10 * time.Millisecond) }()
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	defer srv.Stop()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(getBufDialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(newHubCredentials("anyToken")),
	)
	suite.NoError(err)
	defer conn.Close()
	client := pb.NewTestServiceClient(conn)
	stream, err := client.PingList(ctx, &pb.PingRequest{})
	suite.NoError(err)
	defer stream.CloseSend()
	resp, err := stream.Recv()
	suite.Error(err)
	suite.T().Log(resp)
}

func TestInterceptorSuite(t *testing.T) {
	suite.Run(t, new(InterceptorSuite))
}

func getBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}
