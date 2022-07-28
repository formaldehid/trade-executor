package trader

const (
	DefaultDBDataSourceName = "./trader.db"
)

type Config struct {
	Symbol           string
	DBDataSourceName string
}

func NewConfig() *Config {
	return &Config{
		DBDataSourceName: DefaultDBDataSourceName,
	}
}
