package config

import "time"

type AuthConfig struct {
	AccessToken  AccessTokenConfig
	RefreshToken RefreshTokenConfig
	OTP          OTPConfig
	Challenge    ChallengeConfig
}

type AccessTokenConfig struct {
	TTL    time.Duration `env:"AUTH_ACCESS_TTL" envDefault:"15m"`
	Secret string        `env:"AUTH_JWT_SECRET,required"`
	Issuer string        `env:"AUTH_ACCESS_ISSUER" envDefault:"megan-messenger"`
}

type RefreshTokenConfig struct {
	TTL           time.Duration `env:"AUTH_REFRESH_TTL" envDefault:"720h"`
	KeyPrefix     string        `env:"AUTH_REFRESH_KEY_PREFIX" envDefault:"refresh:"`
	UserKeyPrefix string        `env:"AUTH_REFRESH_USER_KEY_PREFIX" envDefault:"user:"`
}

type OTPConfig struct {
	TTL            time.Duration `env:"AUTH_OTP_TTL" envDefault:"5m"`
	ResendCooldown time.Duration `env:"AUTH_OTP_RESEND_COOLDOWN" envDefault:"60s"`
	MaxAttempts    int           `env:"AUTH_OTP_MAX_ATTEMPTS" envDefault:"5"`
}

type ChallengeConfig struct {
	TTL       time.Duration `env:"AUTH_CHALLENGE_TTL" envDefault:"5m"`
	KeyPrefix string        `env:"AUTH_CHALLENGE_KEY_PREFIX" envDefault:"auth:challenge:"`
}
