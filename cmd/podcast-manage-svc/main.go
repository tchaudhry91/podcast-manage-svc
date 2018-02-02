package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/tchaudhry91/podcast-manage-svc/service"
	"net/http"
	"os"
)

func main() {
	var (
		httpAddr   = flag.String("http.addr", ":8080", "HTTP listen address")
		dbDialect  = flag.String("db.dialect", "postgres", "Dialect of the Database to talk to")
		dbHostname = flag.String("db.hostname", "localhost", "Location of the database host")
		dbUser     = flag.String("db.user", "test", "User to connect to database")
		dbPassword = flag.String("db.password", "", "Password to connect to the database")
		dbName     = flag.String("db.name", "podcastmg", "Name of the database to connect to")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	dbConnString := BuildDBConnString(*dbDialect, *dbHostname, *dbUser, *dbPassword, *dbName)

	var svc service.PodcastManageService
	{
		var err error
		svc, err = service.NewSQLStorePodcastManageService(*dbDialect, dbConnString)
		if err != nil {
			logger.Log("err", err.Error())
			panic("Could not create service")
		}
	}

	var h http.Handler
	{
		h = service.MakeHTTPHandler(svc)
	}

	http.ListenAndServe(*httpAddr, h)
}

func BuildDBConnString(dialect, hostname, user, password, name string) (connString string) {
	switch dialect {
	case "postgres":
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", hostname, user, password, name)
	default:
		return ""
	}
}