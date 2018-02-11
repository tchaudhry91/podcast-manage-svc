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
		httpAddr         = flag.String("http.addr", ":8080", "HTTP listen address")
		dbDialect        = flag.String("db.dialect", "postgres", "Dialect of the Database to talk to")
		dbHostname       = flag.String("db.hostname", "localhost", "Location of the database host")
		dbUser           = flag.String("db.user", "test", "User to connect to database")
		dbPassword       = flag.String("db.password", "", "Password to connect to the database")
		dbName           = flag.String("db.name", "podcastmg", "Name of the database to connect to")
		dbSSLMode        = flag.String("db.sslmode", "disable", "SSLMode enable/disable when applicable")
		svcSigningSecret = flag.String("svc.signingSharedSecret", "", "Token Signing Secret for the service")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	dbConnString := BuildDBConnString(*dbDialect, *dbHostname, *dbUser, *dbPassword, *dbName, *dbSSLMode)

	// Base Service
	var svc service.PodcastManageService
	{
		var err error
		svc, err = service.NewSQLStorePodcastManageService(*dbDialect, dbConnString, *svcSigningSecret, logger)
		if err != nil {
			logger.Log("err", err.Error())
			panic("Could not create service")
		}
	}

	// Middlewares
	svc = service.MakeNewLoggingMiddleware(logger, svc)

	var h http.Handler
	{
		h = service.MakeHTTPHandler(svc, *svcSigningSecret)
	}

	http.ListenAndServe(*httpAddr, h)
}

// BuildDBConnString returns a GORM connection string from the given parameters
func BuildDBConnString(dialect, hostname, user, password, name, sslmode string) (connString string) {
	switch dialect {
	case "postgres":
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", hostname, user, password, name, sslmode)
	default:
		return ""
	}
}
