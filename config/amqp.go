package config

type amqpConfig struct {
	Addr         string `toml:"addr"`
	Username     string `toml:"username"`
	Password     string `toml:"password"`
	QueueName    string `toml:"queue_name"`
	QueueDurable bool   `toml:"queue_durable"`
}
