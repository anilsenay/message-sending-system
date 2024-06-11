package repository_test

import (
	"os"
	"testing"

	"github.com/anilsenay/message-sending-system/pkg/orm"
)

var dockerDatabase *orm.Database

func TestMain(m *testing.M) {
	dockerDatabase = orm.NewDockerDatabase(orm.DockerDatabaseConfig{
		MigrationQueryPath: "../../db.sql",
	})

	os.Exit(m.Run())
}
