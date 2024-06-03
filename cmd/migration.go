package cmd

import (
	"github.com/pressly/goose"
	"github.com/spf13/cobra"
	"pcstakehometest/database/postgres"
	"pcstakehometest/package/logger"
)

var migration = &cobra.Command{
	Use:   "migration",
	Short: "Migration database",
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.NewLogRus()
		postgres := postgres.NewPostgres(log)
		if len(args) == 0 {
			log.Error("please insert argument up or down")
			return
		} else if args[0] == "up" {
			if err := Up(postgres); err != nil {
				log.Error(err)
				return
			}
		} else if args[0] == "down" {
			if err := Down(postgres); err != nil {
				log.Error(err)
				return
			}
		} else {
			log.Error("migration argument not found")
			return
		}
	},
}

func Up(db *postgres.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db.Sql, "sql"); err != nil {
		return err
	}

	return nil
}

func Down(db *postgres.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Down(db.Sql, "sql"); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(migration)
}
