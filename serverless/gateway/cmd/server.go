package main

import "minik8s/serverless/gateway/internal"

func main() {
	// Start the server
	gateway := internal.GetServerlessGatewayInstance()
	gateway.RUN()
}
