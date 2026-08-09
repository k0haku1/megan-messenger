package config

type MinIOConfig struct {
	Endpoint  string `env:"APP_MINIO_ENDPOINT" envDefault:"localhost:9000"`
	AccessKey string `env:"APP_MINIO_ACCESS_KEY" envDefault:"minioadmin"`
	SecretKey string `env:"APP_MINIO_SECRET_KEY" envDefault:"minioadmin"`
	Bucket    string `env:"APP_MINIO_BUCKET" envDefault:"megan"`
	UseSSL    bool   `env:"APP_MINIO_USE_SSL" envDefault:"false"`
	// PublicBase is used when generating absolute URLs for clients (optional).
	// When empty, PresignGet returns the signed MinIO URL as-is.
	PublicBase string `env:"APP_MINIO_PUBLIC_BASE" envDefault:""`
}
