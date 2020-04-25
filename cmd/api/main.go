package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oklog/run"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/kl09/auth-go/internal/api"
	"github.com/kl09/auth-go/internal/generator"
	"github.com/kl09/auth-go/internal/pg"
)

func main() {
	var err error

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logger.Info().Msg("starting app")

	fs := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	{
		fs.String(
			"pg.conn-string",
			"user=auth password=auth host=localhost port=5432 dbname=auth_test connect_timeout=3 sslmode=disable",
			"Postgresql connection string",
		)
		fs.Int("pg.max-cons", 5, "Max connections to Postgres.")
		fs.Int("pg.max-idle-cons", 2, "Max idle connections to Postgres.")
		fs.Duration("pg.connection-timeout", time.Minute, "Max connection timeout to Postgres.")

		fs.String("http-addr", ":8080", "Address to listen for System API")

		fs.String("log-lvl", "info", "Log level.")
	}

	if err = viper.BindPFlags(fs); err != nil {
		logger.Fatal().Err(err).Msg("failed bind pflags")
		os.Exit(1)
	}

	logLvl, err := zerolog.ParseLevel(viper.GetString("log-lvl"))
	if err != nil {
		logger.Err(err).Msg("Couldn't parse log lvl.")
	} else {
		logger.Info().Msgf("Set log level: %s", viper.GetString("log-lvl"))
		logger = logger.Level(logLvl)
	}

	pgClient := pg.NewClient(
		pg.WithLogger(logger),
		pg.WithMaxConnections(viper.GetInt("pg.max-cons")),
		pg.WithMaxIdleConnections(viper.GetInt("pg.max-idle-cons")),
		pg.WithConnectionTimeout(viper.GetDuration("pg.connection-timeout")),
	)
	if err = pgClient.Open(viper.GetString("pg.conn-string")); err != nil {
		logger.Fatal().Err(err).Msg("db connection failed")
		os.Exit(1)
	}

	defer func() {
		if err = pgClient.Close(); err != nil {
			logger.Error().Err(err).Msg("db close failed")
		}
	}()

	credRepository := pg.NewCredentialRepository(pgClient)

	r := api.NewRouter(
		api.NewCredentialService(
			credRepository,
			func() time.Time {
				return time.Now().UTC()
			},
			generator.GenerateRandomString,
		),
	)

	apiServer := &http.Server{
		Addr:    viper.GetString("http-addr"),
		Handler: r.Handler(),
	}

	ctx, cancel := context.WithCancel(context.Background())

	var g run.Group
	{
		g.Add(
			func() error {
				sig := make(chan os.Signal, 1)
				signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
				select {
				case <-ctx.Done():
					return nil
				case si := <-sig:
					return fmt.Errorf("signal received: %v", si)
				}
			}, func(err error) {
				logger.Fatal().Err(err).Msg("app was interrupted")
				cancel()
			},
		)
	}
	{
		g.Add(func() error {
			logger.Info().Msgf("started server for addr: %s", apiServer.Addr)
			return apiServer.ListenAndServe()
		}, func(err error) {
			logger.Info().Err(err).Msg("server was stopped")
		})
	}
	{
		g.Add(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case <-time.After(5 * time.Second):
					b, _ := json.Marshal(pgClient.Stats())
					logger.Debug().Msgf("pgStats :%v", string(b))
				}
			}
		}, func(err error) {

		})
	}

	err = g.Run()
	logger.Info().Err(err).Msg("app was stopped")
}
