package dao

func init() {
	initPostgresql()

	initMongo()
	initRedisDb()

}
