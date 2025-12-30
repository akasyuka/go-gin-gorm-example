package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type Config struct {
	App        AppConfig        `yaml:"app"`
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Kafka      KafkaConfig      `yaml:"kafka"`
	Monitoring MonitoringConfig `yaml:"metrics"`
	Auth       AuthConfig       `yaml:"auth"`
}

type AppConfig struct {
	Name string `yaml:"name"`
	Env  string `yaml:"env"`
}

// ===== Server config =====
type ServerConfig struct {
	HTTP HTTPServerConfig `yaml:"http"`
	GRPC GRPCServerConfig `yaml:"grpc"` // новый блок для gRPC
}

type HTTPServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type GRPCServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// ===== Database config =====
type DatabaseConfig struct {
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	Host       string             `yaml:"host"`
	Port       int                `yaml:"port"`
	Name       string             `yaml:"name"`
	User       string             `yaml:"user"`
	Password   string             `yaml:"password"`
	SSLMode    string             `yaml:"ssl_mode"`
	Pool       DBPool             `yaml:"pool"`
	Migrations PostgresMigrations `yaml:"migrations"`
}

type PostgresMigrations struct {
	Path string `yaml:"path"`
}

type DBPool struct {
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// ===== Kafka config =====
type KafkaConfig struct {
	BootstrapServers []string      `yaml:"bootstrap_servers"`
	ClientID         string        `yaml:"client_id"`
	Timeouts         KafkaTimeouts `yaml:"timeouts"`
	Producer         KafkaProducer `yaml:"producer"`
	Consumer         KafkaConsumer `yaml:"consumer"`
	Topics           KafkaTopics   `yaml:"topics"`
	Security         KafkaSecurity `yaml:"security"`
}

type KafkaTimeouts struct {
	Dial  time.Duration `yaml:"dial_ms"`
	Read  time.Duration `yaml:"read_ms"`
	Write time.Duration `yaml:"write_ms"`
}

type KafkaProducer struct {
	Acks              string `yaml:"acks"`
	Retries           int    `yaml:"retries"`
	RetryBackoffMs    int    `yaml:"retry_backoff_ms"`
	LingerMs          int    `yaml:"linger_ms"`
	BatchSize         int    `yaml:"batch_size"`
	MaxRequestSize    int    `yaml:"max_request_size"`
	CompressionType   string `yaml:"compression_type"`
	EnableIdempotence bool   `yaml:"enable_idempotence"`
}

type KafkaConsumer struct {
	GroupID              string `yaml:"group_id"`
	AutoOffsetReset      string `yaml:"auto_offset_reset"`
	EnableAutoCommit     bool   `yaml:"enable_auto_commit"`
	AutoCommitIntervalMs int    `yaml:"auto_commit_interval_ms"`
	MaxPollRecords       int    `yaml:"max_poll_records"`
	SessionTimeoutMs     int    `yaml:"session_timeout_ms"`
	HeartbeatIntervalMs  int    `yaml:"heartbeat_interval_ms"`
	MaxPollIntervalMs    int    `yaml:"max_poll_interval_ms"`
	IsolationLevel       string `yaml:"isolation_level"`
}

type KafkaTopics struct {
	Orders TopicConfig `yaml:"orders"`
}

type TopicConfig struct {
	Name              string `yaml:"name"`
	Partitions        int    `yaml:"partitions"`
	ReplicationFactor int    `yaml:"replication_factor"`
}

type KafkaSecurity struct {
	Protocol string    `yaml:"protocol"`
	SASL     KafkaSASL `yaml:"sasl"`
}

type KafkaSASL struct {
	Mechanism string `yaml:"mechanism"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

// ===== Monitoring / Prometheus =====
type MonitoringConfig struct {
	Prometheus PrometheusConfig `yaml:"prometheus"`
}

type PrometheusConfig struct {
	Enabled        bool          `yaml:"enabled"`
	MetricsPath    string        `yaml:"metrics_path"`
	Port           int           `yaml:"port"`
	ScrapeInterval time.Duration `yaml:"scrape_interval"` // теперь time.Duration
	JobName        string        `yaml:"job_name"`
}

// ===== Auth / Keycloak =====
type AuthConfig struct {
	Keycloak KeycloakConfig `yaml:"keycloak"`
}

type KeycloakConfig struct {
	JWKSURL string `yaml:"jwks_url"`
}
