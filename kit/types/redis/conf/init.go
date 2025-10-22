package conf

type Config struct {
	Name           string `toml:"name" json:"name" yaml:"name"`
	Addr           string `toml:"addr" json:"addr" yaml:"addr"`
	Password       string `toml:"password" json:"password" yaml:"password"`
	DB             int    `toml:"db" json:"db" yaml:"db"`
	MinIdleConns   int    `toml:"min_idle_conns" json:"min_idle_conns" yaml:"min_idle_conns"`
	MaxIdleConns   int    `toml:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxActiveConns int    `toml:"max_active_conns" json:"max_active_conns" yaml:"max_active_conns"`
	ConnTimeoutMs  int    `toml:"conn_timeout_ms" json:"conn_timeout_ms" yaml:"conn_timeout_ms"`
	ReadTimeoutMs  int    `toml:"read_timeout_ms" json:"read_timeout_ms" yaml:"read_timeout_ms"`
	WriteTimeoutMs int    `toml:"write_timeout_ms" json:"write_timeout_ms" yaml:"write_timeout_ms"`
}
