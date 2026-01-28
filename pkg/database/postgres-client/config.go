// Package postgresclient provides PostgreSQL connection configuration and management.
package postgresclient

import "time"

// PGConfig holds PostgreSQL connection parameters and pool settings.
type PGConfig struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
	sslMode  string

	// Connection pool settings
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
}

// NewPGConfig creates new PostgreSQL configuration.
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
