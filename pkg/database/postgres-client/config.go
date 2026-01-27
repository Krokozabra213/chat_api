package postgresclient

import "time"

// PGConfig конфигурация подключения
type PGConfig struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
	sslMode  string

	// Pool settings
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
}

func NewPGConfig(host, port, user, password, dbName, sslMode string,
	maxOpenConns, maxIdleConns int, connMaxLifetime time.Duration,
) PGConfig {
	return PGConfig{
		host:            host,
		port:            port,
		user:            user,
		password:        password,
		dbName:          dbName,
		sslMode:         sslMode,
		maxOpenConns:    maxOpenConns,
		maxIdleConns:    maxIdleConns,
		connMaxLifetime: connMaxLifetime,
	}
}
