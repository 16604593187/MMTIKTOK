package config


type Config struct {
	KafkaConfig struct {
		Host     string
		Topic    string
		MinBytes int
		MaxBytes int
	}
	DbConfig struct {
		Path         string `json:"path" yaml:"path"`                     
		Port         int    `json:"port" yaml:"port"`                     
		Config       string `json:"config" yaml:"config"`                 
		Dbname       string `json:"db-name" yaml:"db-name"`               
		Username     string `json:"username" yaml:"username"`             
		Password     string `json:"password" yaml:"password"`             
		MaxIdleConns int    `json:"max-idle-conns" yaml:"max-idle-conns"` 
		MaxOpenConns int    `json:"max-open-conns" yaml:"max-open-conns"` 
	}
	RedisConfig struct {
		Host        string
		Port        int
		Auth        bool
		MaxIdle     int
		Active      int
		IdleTimeout int
	}
	WorkerId uint32
}