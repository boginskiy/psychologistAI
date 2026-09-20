package jwtservice

type JWTConf struct {
	TimeLiveToken  int
	SecretKeyToken string
	HostName       string
}

func NewJWTConfig(timeLiveToken int, secretKeyToken, hostName string) *JWTConf {
	return &JWTConf{}
}

func (c *JWTConf) GetTimeLiveToken() int {
	return c.TimeLiveToken
}
func (c *JWTConf) GetSecretKeyForToken() string {
	return c.SecretKeyToken
}

func (c *JWTConf) GetHostName() string {
	return c.HostName
}
