package engine

type Config struct {
	Model      string `json:"model"`
	Cfg        string `json:"cfg"`
	Feed       string `json:"feed"`
	Classnames string `json:"classnames"`
}

func (c *Config) LoadConfig() {

}
