package conf

type Config struct {
	Name           string `toml:"name" json:"name" yaml:"name"`
	Addr           string `toml:"addr" json:"addr" yaml:"addr"`
	Password       string `toml:"password" json:"password" yaml:"password"`
	DB             int    `toml:"db" json:"db" yaml:"db"`
	MinIdleConns   int    `toml:"min_idle_conns" json:"min_idle_conns" yaml:"min_idle_conns"`
	MaxIdleConns   int    `toml:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxActiveConns int    `toml:"max_active_conns" json:"max_active_conns" yaml:"max_active_conns"`
	TimeoutMs      int    `toml:"timeout_ms" json:"timeout_ms" yaml:"timeout_ms"`
}
