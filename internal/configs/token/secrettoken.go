package token

import "github.com/tirtahakimpambudhi/restful_api/internal/configs"

type SecretKey struct {
	AccessToken         string `env:"SECRET_KEY_ACCESS_TOKEN,required"`
	RefreshToken        string `env:"SECRET_KEY_REFRESH_TOKEN,required"`
	ForgotPasswordToken string `env:"SECRET_KEY_FP_TOKEN,required"`
}

func NewSecretKey() (*SecretKey, error) {
	var secretKey SecretKey
	if err := configs.GetConfig().Load(&secretKey); err != nil {
		return nil, err
	}
	return &secretKey, nil
}
