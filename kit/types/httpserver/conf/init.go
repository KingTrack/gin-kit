package conf

const (
	DefaultRequestBodySizeKB  = 4
	DefaultResponseBodySizeKB = 4
)

type Config struct {
	ServiceName             string `toml:"service_name" json:"service_name" yaml:"service_name"`
	Hostname                string `toml:"hostname" json:"hostname" yaml:"hostname"`
	IP                      string `toml:"ip" json:"ip" yaml:"ip"`
	Port                    int    `toml:"port" json:"port" yaml:"port"`
	ReadTimeoutSec          int    `toml:"read_timeout_sec" json:"read_timeout_sec" yaml:"read_timeout_sec"`
	WriteTimeoutSec         int    `toml:"write_timeout_sec" json:"write_timeout_sec" yaml:"write_timeout_sec"`
	IdleTimeoutSec          int    `toml:"idle_timeout_sec" json:"idle_timeout_sec" yaml:"idle_timeout_sec"`
	PrintRequestBodySizeKB  int    `toml:"print_request_body_size_kb" json:"print_request_body_size_kb" yaml:"print_request_body_size_kb"`    // 默认4MB
	PrintResponseBodySizeKB int    `toml:"print_response_body_size_kb" json:"print_response_body_size_kb" yaml:"print_response_body_size_kb"` // 默认4MB
	CloseRequestBody        bool   `toml:"close_request_body" json:"close_request_body" yaml:"close_request_body"`
	CloseResponseBody       bool   `toml:"close_response_body" json:"close_response_body" yaml:"close_response_body"`
}
