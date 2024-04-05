package main

import (
	"PlacesApp/config"
	"context"
)

func main() {
	ctx := context.Background()
	ctx = config.ConnectWithElasticSearch(ctx)
	config.ConfigServer(ctx)
}
