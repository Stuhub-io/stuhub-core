package scylla

import (
	"strconv"

	"github.com/Stuhub-io/logger"
	gocql "github.com/gocql/gocql"
)

// CQLLogger implements gocql.Trace for logging queries
type CQLLogger struct {
	logger logger.Logger
}

// NewCQLLogger initializes a new CQLLogger with Logrus
func NewCQLLogger(logger logger.Logger) *CQLLogger {
	return &CQLLogger{logger: logger}
}

// Trace method logs query execution details
func (c *CQLLogger) Trace(traceId []byte) {

}

func Must(hosts []string, port string, keyspace string, isDebug bool, logger logger.Logger) *gocql.Session {
	nPort, err := strconv.Atoi(port)
	if err != nil {
		logger.Fatalf(err, "failed to convert port to int: %v")
		panic(err)
	}
	cluster := gocql.NewCluster(hosts...)
	cluster.Port = nPort
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum

	// FIXME: Possible enable SSL
	// cluster.SslOpts = &gocql.SslOptions{
	// 	EnableHostVerification: true,
	// }

	session, err := cluster.CreateSession()
	session.SetTrace(NewCQLLogger(logger))

	if err != nil {
		logger.Fatalf(err, "failed to create Cassandra session: %v")
		panic(err)
	}

	if isDebug {
		session.SetConsistency(gocql.Quorum)
	}

	logger.Info("SCylla connected")
	return session
}
