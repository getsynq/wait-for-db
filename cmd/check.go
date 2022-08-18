package cmd

import (
	"database/sql"
	"github.com/mxssl/wait-for-pg/util"
	"github.com/spf13/cobra"
	"log"
	nurl "net/url"
	"os"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/lib/pq"
)

type config struct {
	database string
	retry    int
	sleep    int
}

var c config

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if db is ready",
	Long:  `Check if db is ready`,
	Run: func(cmd *cobra.Command, args []string) {
		check(c)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringVar(&c.database,
		"database",
		c.database,
		"database")

	checkCmd.Flags().IntVar(&c.retry,
		"retry",
		1,
		"retry")

	checkCmd.Flags().IntVar(&c.sleep,
		"sleep",
		1,
		"sleep")
}

type ClickHouse struct {
	conn *sql.DB
}

type Pg struct {
	conn *sql.DB
}

func check(c config) {
	purl, err := nurl.Parse(c.database)

	if err != nil {
		log.Fatal("Could not parse database url")
	}

	scheme, err := util.SchemeFromURL(c.database)
	if err != nil {
		log.Fatal(err)
	}

	q := util.FilterCustomQuery(purl)

	// Make this pretty
	if scheme == "clickhouse" {
		q.Scheme = "tcp"
	}

	log.Println(q.String())

	for i := 0; i < c.retry; i++ {
		time.Sleep(time.Duration(c.sleep) * time.Second)
		db, err := sql.Open(scheme, q.String())
		if db != nil {
			err := db.Ping()
			defer db.Close()
			if err != nil {
				log.Printf("Error: %s", err.Error())
				continue
			}
		}
		if db == nil {
			log.Printf("Error: %s", err.Error())
			continue
		}
		log.Printf("DB is ready!")
		os.Exit(0)
	}
	log.Printf("DB isn't ready! Retry counter exceeded")
	os.Exit(1)
}
