package config

func NewCassandra(set Set) Cassandra { return set.Cassandra }
func NewCloud(set Set) Cloud         { return set.Cloud }
func NewCockroach(set Set) Cockroach { return set.Cockroach }
func NewCore(set Set) Core           { return set.Core }
func NewFirebase(set Set) Firebase   { return set.Firebase }
func NewGRPC(set Set) GRPC           { return set.GRPC }
func NewJWT(set Set) JWT             { return set.JWT }
func NewMongo(set Set) Mongo         { return set.Mongo }
func NewPubSub(set Set) PubSub       { return set.PubSub }
func NewPostgres(set Set) Postgres   { return set.Postgres }
func NewRedis(set Set) Redis         { return set.Redis }
func NewServer(set Set) Server       { return set.Server }
func NewSpanner(set Set) Spanner     { return set.Spanner }
func NewSet() (Set, error) {
	set := Set{}
	err := LoadFromEnv(
		&set.Cassandra,
		&set.Cloud,
		&set.Cockroach,
		&set.Core,
		&set.Firebase,
		&set.GRPC,
		&set.JWT,
		&set.Mongo,
		&set.PubSub,
		&set.Postgres,
		&set.Redis,
		&set.Server,
		&set.Spanner,
	)
	return set, err
}

type Set struct {
	Cassandra Cassandra
	Cloud     Cloud
	Cockroach Cockroach
	Core      Core
	Firebase  Firebase
	GRPC      GRPC
	JWT       JWT
	Mongo     Mongo
	PubSub    PubSub
	Postgres  Postgres
	Redis     Redis
	Server    Server
	Spanner   Spanner
}
