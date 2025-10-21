package config

type Config struct {
	DB    DBConfig
	HTTP  HTTPConfig
	REDIS RedisClient
}

func NewConfig() Config {
	return Config{
		DB:    LoadDBConfig(),
		HTTP:  LoadHTTPConfig(),
		REDIS: *NewRedisClient(),
	}
}

func LoadBatch() int {
	// batchSizeString := os.Getenv("BATCH_SIZE")
	// batchSize, _ := strconv.Atoi(batchSizeString)
	// if batchSize == 0 {
	// 	batchSize = 100
	// }
	batchSize := 100
	return batchSize
}
