package pg

import (
	"io/ioutil"

	"github.com/jinzhu/gorm"
	// postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/zerolog"
)

type Client struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func NewClient() *Client {
	return &Client{
		logger: zerolog.New(ioutil.Discard),
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
