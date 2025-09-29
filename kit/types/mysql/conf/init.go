package conf

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type DBConfig struct {
	Name            string
	DSN             string
	User            string
	Password        string
	Host            string
	Port            string
	DBName          string
	Params          map[string]string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	Timeout         time.Duration
}

type Config struct {
	Name   string   `toml:"name" json:"name" yaml:"name"`
	Master string   `toml:"master" json:"master" yaml:"master"`
	Slaves []string `toml:"slaves" json:"slaves" yaml:"slaves"`
}

// user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s&maxOpenConns=50&maxIdleConns=10&connMaxLifetime=1h

func (c *Config) Parse() (master *DBConfig, slaves []*DBConfig, err error) {
	master, err = parseDSN(c.Name, c.Master)
	if err != nil {
		return nil, nil, err
	}

	for _, dsn := range c.Slaves {
		s, err := parseDSN(c.Name, dsn)
		if err != nil {
			return nil, nil, err
		}
		slaves = append(slaves, s)
	}

	return
}

func parseDSN(name, dsn string) (*DBConfig, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn is empty for %s", name)
	}

	parts := strings.SplitN(dsn, "@tcp(", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid DSN format: %s", dsn)
	}

	userPass := parts[0]
	hostPart := parts[1]

	user := ""
	pass := ""
	if strings.Contains(userPass, ":") {
		up := strings.SplitN(userPass, ":", 2)
		user = up[0]
		pass = up[1]
	} else {
		user = userPass
	}

	idx := strings.Index(hostPart, ")/")
	if idx < 0 {
		return nil, fmt.Errorf("invalid DSN format: %s", dsn)
	}

	hostPort := hostPart[:idx]
	dbAndParams := hostPart[idx+2:]
	dbName := ""
	params := map[string]string{}

	if strings.Contains(dbAndParams, "?") {
		tmp := strings.SplitN(dbAndParams, "?", 2)
		dbName = tmp[0]
		values, _ := url.ParseQuery(tmp[1])
		for k, v := range values {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
	} else {
		dbName = dbAndParams
	}

	host := hostPort
	port := "3306"
	if strings.Contains(hostPort, ":") {
		parts := strings.SplitN(hostPort, ":", 2)
		host = parts[0]
		port = parts[1]
	}

	// 解析 maxOpenConns/maxIdleConns/connMaxLifetime/timeout
	maxOpen, _ := strconv.Atoi(params["maxOpenConns"])
	maxIdle, _ := strconv.Atoi(params["maxIdleConns"])
	connMaxLife, _ := time.ParseDuration(params["connMaxLifetime"])
	timeout, _ := time.ParseDuration(params["timeout"])

	return &DBConfig{
		Name:            name,
		DSN:             dsn,
		User:            user,
		Password:        pass,
		Host:            host,
		Port:            port,
		DBName:          dbName,
		Params:          params,
		MaxOpenConns:    maxOpen,
		MaxIdleConns:    maxIdle,
		ConnMaxLifetime: connMaxLife,
		Timeout:         timeout,
	}, nil
}
