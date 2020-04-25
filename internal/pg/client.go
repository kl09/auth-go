package pg

import (
	"database/sql"
	"io/ioutil"
	"time"

	"github.com/jinzhu/gorm"
	// postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/zerolog"
)

type Client struct {
	db                *gorm.DB
	logger            zerolog.Logger
	maxOpenCons       int
	maxIdleCons       int
	connectionTimeout time.Duration
}

// NewClient returns a new Client for DB connection.
func NewClient(options ...ConfigOption) *Client {
	c := Client{
		logger: zerolog.New(ioutil.Discard),
	}

	for _, opt := range options {
		opt(&c)
	}

	return &c
}

// ConfigOption configures the client.
type ConfigOption func(*Client)

// WithLogger configures a logger to debug interactions with Postgres.
func WithLogger(l zerolog.Logger) ConfigOption {
	return func(c *Client) {
		c.logger = l
	}
}

// WithMaxConnections configures a max connections to Postgres.
func WithMaxConnections(n int) ConfigOption {
	return func(c *Client) {
		c.maxOpenCons = n
	}
}

// WithMaxIdleConnections configures a max idle connections to Postgres.
func WithMaxIdleConnections(n int) ConfigOption {
	return func(c *Client) {
		c.maxIdleCons = n
	}
}

// WithConnectionTimeout configures a max connection timeout to Postgres.
func WithConnectionTimeout(t time.Duration) ConfigOption {
	return func(c *Client) {
		c.connectionTimeout = t
	}
}

// Open opens PostgreSQL connection.
func (c *Client) Open(source string) error {
	var err error

	c.logger.Debug().Msg("connecting to db")

	c.db, err = gorm.Open("postgres", source)
	if err != nil {
		c.logger.Err(err).Msg("sql open failed")
		return err
	}

	err = c.db.DB().Ping()
	if err != nil {
		c.logger.Err(err).Msg("sql ping failed")
		return err
	}

	c.logger.Debug().Msg("connected to db")

	c.db.SingularTable(true)
	c.db.DB().SetMaxOpenConns(c.maxOpenCons)
	c.db.DB().SetMaxIdleConns(c.maxIdleCons)
	c.db.DB().SetConnMaxLifetime(c.connectionTimeout)

	return nil
}

// Close closes PostgreSQL connection.
func (c *Client) Close() error {
	c.logger.Debug().Msg("connection to db closed")
	return c.db.Close()
}

// Schema sets up the initial schema.
func (c *Client) Schema() error {
	_, err := c.db.DB().Exec(Schema)
	return err
}

// Stats returns database statistics.
func (c *Client) Stats() sql.DBStats {
	return c.db.DB().Stats()
}
