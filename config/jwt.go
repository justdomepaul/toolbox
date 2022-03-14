package config

// JWT type
type JWT struct {
	EcdsaPrivateKeyPath string `split_words:"true" default:"./es256_private.pem"`
	EcdsaPrivateKey     string `split_words:"true" default:""`
	RsaPrivateKeyPath   string `split_words:"true" default:"./rs256-private.pem"`
	RsaPrivateKey       string `split_words:"true" default:""`
	HmacSecretKeyPath   string `split_words:"true" default:"./hs-secret.pem"`
	HmacSecretKey       string `split_words:"true" default:""`
}
