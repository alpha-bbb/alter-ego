package grpc

import (
	"fmt"
	"net"

	"github.com/alpha-bbb/alter-ego/backend/adapter/clock"
	"github.com/alpha-bbb/alter-ego/backend/adapter/database/postgresql"
	"github.com/alpha-bbb/alter-ego/backend/adapter/llm/client"
	"github.com/alpha-bbb/alter-ego/backend/adapter/payment"
	"github.com/alpha-bbb/alter-ego/backend/adapter/ulid"
	"github.com/alpha-bbb/alter-ego/backend/config"
	"github.com/alpha-bbb/alter-ego/backend/database"
	"github.com/alpha-bbb/alter-ego/backend/database/repository"
	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
	"github.com/alpha-bbb/alter-ego/backend/usecase"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(logger *zap.Logger) (*grpc.Server, net.Listener, error) {
	grpcPort := config.GrpcPort()

	grpcClientAddress := config.LlmGrpcClientAddress()

	db, err := postgresql.NewPostgreSQL(config.DSN())
	if err != nil {
		panic(err)
	}

	ulid := ulid.NewULID()
	clock := clock.New()

	llmClient := client.NewLLMGRPCClient(grpcClientAddress)
	stripeDriver := payment.NewStripeDriver(
		config.StripeApiKey(),
		config.StripeEndpointSecret(),
		config.FrontendURL(),
		config.PlaceID(),
	)

	transaction := database.NewGormTransaction(db)
	stateCountRepository := repository.NewStateCountRepository(db)
	subscribeStripeRepository := repository.NewSubscribeStripeRepository(db)
	userAccountLineRepository := repository.NewUserAccountLineRepository(db)
	userAggregateRepository := repository.NewUserAggregateRepository(db)
	userDetailRepository := repository.NewUserDetailRepository(db)
	userRepository := repository.NewUserRepository(db)

	talkUseCase := usecase.NewTalkUseCase(
		ulid,
		clock,
		transaction,
		stateCountRepository,
		subscribeStripeRepository,
		userAccountLineRepository,
		userAggregateRepository,
		userDetailRepository,
		userRepository,
		llmClient,
		stripeDriver,
	)
	stripeSubscribeUseCase := usecase.NewStripeSubscribeUseCase(
		ulid,
		clock,
		transaction,
		stateCountRepository,
		subscribeStripeRepository,
		userAccountLineRepository,
		userAggregateRepository,
		userDetailRepository,
		userRepository,
		llmClient,
		stripeDriver,
	)
	stripeUnsubscribeUseCase := usecase.NewStripeUnsubscribeUseCase(
		clock,
		subscribeStripeRepository,
		userRepository,
		stripeDriver,
	)
	stripeCheckSubscribeUseCase := usecase.NewStripeCheckSubscribeUseCase(
		clock,
		subscribeStripeRepository,
		userRepository,
		stripeDriver,
	)

	backendServer := NewBackendServiceServer(
		talkUseCase,
		stripeSubscribeUseCase,
		stripeUnsubscribeUseCase,
		stripeCheckSubscribeUseCase,
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		return nil, nil, err
	}

	grpcServer := grpc.NewServer()
	backendpb.RegisterBackendServiceServer(grpcServer, backendServer)
	reflection.Register(grpcServer)

	return grpcServer, lis, nil
}
