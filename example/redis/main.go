package main

import (
	"context"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/database/redis"
	"log"
)

func main() {
	cfgSet, err := config.NewSet()
	if err != nil {
		panic(err)
	}
	session, err := redis.NewSession(cfgSet.Redis)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	log.Println(session.Ping(context.Background()))
}
