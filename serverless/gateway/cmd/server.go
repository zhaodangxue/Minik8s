package main

import "minik8s/serverless/gateway/internal"

func main() {
	// Start the server
	gateway := internal.New()
	gateway.RUN()
}
