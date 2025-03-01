package config

import (
	"log"
	"os"
)

var (
	grpcPort             string
	httpPort             string
	llmGrpcClientAddress string
	stripeApiKey         string
	stripeEndpointSecret string
	frontendURL          string
	placeID              string
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
		httpPort = "50080"
	}
	llmGrpcClientAddress = os.Getenv("LLM_GRPC_CLIENT_ADDRESS")
	if llmGrpcClientAddress == "" {
		log.Print("LLM_GRPC_CLIENT_ADDRESS environment variable is empty")
		llmGrpcClientAddress = "localhost:50051"
	}
	stripeApiKey = os.Getenv("STRIPE_API_KEY")
	if stripeApiKey == "" {
		log.Print("STRIPE_API_KEY environment variable is empty")
	}
	stripeEndpointSecret = os.Getenv("STRIPE_ENDPOINT_SECRET")
	if stripeEndpointSecret == "" {
		log.Print("STRIPE_ENDPOINT_SECRET environment variable is empty")
	}
	frontendURL = os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		log.Print("FRONTEND_URL environment variable is empty")
	}
	placeID = os.Getenv("PLACE_ID")
	if placeID == "" {
		log.Print("PLACE_ID environment variable is empty")
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

func StripeApiKey() string {
	return stripeApiKey
}

func StripeEndpointSecret() string {
	return stripeEndpointSecret
}

func FrontendURL() string {
	return frontendURL
}

func PlaceID() string {
	return placeID
}
