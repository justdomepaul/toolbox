package main

import (
	"fmt"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/database/postgres"
	"time"
)

func main() {
	cfgSet, err := config.NewSet()
	if err != nil {
		panic(err)
	}
	cfgSet.Postgres.PostgresURL = ""
	session, err := postgres.NewSession(cfgSet.Postgres)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	var now time.Time
	if err := session.QueryRowx("SELECT NOW();").Scan(&now); err != nil {
		panic(err)
	}
	fmt.Println(now)
}
