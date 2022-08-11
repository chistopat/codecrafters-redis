package main

import (
	"log"

	net "redis/app/network"
	"redis/app/redis"
)

const (
	PORT = 6379
	HOST = "0.0.0.0"
)

func main() {
	log.Println("My simple redis started!")
	cache := redis.NewMyRedis()
	server := net.NewNetworkServer(HOST, PORT, cache.OnConnect)
	if err := server.Run(); err != nil {
		log.Fatalf("can't start server: %v\n", err)
	}
}
