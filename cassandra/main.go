package cassandra

import (
	"fmt"
	"os"

	"github.com/gocql/gocql"
)

var (
	Session               *gocql.Session
	KEY_SPACE             = os.Getenv("KEY_SPACE")             //accountsapi
	CASSANDRA_CLUSTER_URL = os.Getenv("CASSANDRA_CLUSTER_URL") // cassandra.default.svc.cluster.local
)

func init() {
	var err error

	cluster := gocql.NewCluster(CASSANDRA_CLUSTER_URL)
	cluster.Keyspace = KEY_SPACE
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("cassandra init done")
}
