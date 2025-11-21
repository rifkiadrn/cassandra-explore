package config

import (
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewNoSQLDatabase(viper *viper.Viper, log *logrus.Logger) *gocql.Session {
	viper.AutomaticEnv()

	databaseKeyspace := viper.GetString("database.keyspace")
	if viper.GetString("DB_KEYSPACE") != "" {
		databaseKeyspace = viper.GetString("DB_KEYSPACE")
	}

	cluster := gocql.NewCluster("cassandra-seed", "cassandra-node2", "cassandra-node3")
	cluster.Port = 9042
	cluster.Keyspace = databaseKeyspace
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Fatal error cassandra setup: %v", err)
	}

	return session
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func (l *logrusWriter) Printf(message string, args ...interface{}) {
	l.Logger.Tracef(message, args...)
}
