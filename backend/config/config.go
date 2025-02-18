package config

import (
	"log"
	"os"
)

var (
	grpcPort             string
	httpPort             string
	llmGrpcClientAddress string
)

func init() {
	grpcPort = os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		log.Print("GRPC_PORT environment variable is empty")
		grpcPort = "50051"
	}
	httpPort = os.Getenv("HTTP_PORT")
	if httpPort == "" {
		log.Print("HTTP_PORT environment variable is empty")
		httpPort = "8080"
	}
	llmGrpcClientAddress = os.Getenv("LLM_GRPC_CLIENT_ADDRESS")
	if llmGrpcClientAddress == "" {
		log.Print("LLM_GRPC_CLIENT_ADDRESS environment variable is empty")
		llmGrpcClientAddress = "localhost:50051"
	}
}

func GrpcPort() string {
	return grpcPort
}

func HttpPort() string {
	return httpPort
}

func LlmGrpcClientAddress() string {
	return llmGrpcClientAddress
}
