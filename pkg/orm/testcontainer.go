package orm

import (
	"context"
	"log"
	"path/filepath"
	"strconv"

	"github.com/docker/go-connections/nat"
	container "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DockerDatabaseConfig struct {
	MigrationQueryPath string
	SkipReaper         bool
	Reuse              bool
}

func NewDockerDatabase(config DockerDatabaseConfig) *Database {
	postgresPort := nat.Port("5432/tcp")
	postgres, err := container.GenericContainer(context.Background(),
		container.GenericContainerRequest{
			ContainerRequest: container.ContainerRequest{
				Image:        "postgres:14-alpine",
				ExposedPorts: []string{postgresPort.Port()},
				Env: map[string]string{
					"POSTGRES_PASSWORD": "POSTGRES_PASSWORD",
					"POSTGRES_USER":     "POSTGRES_USER",
					"POSTGRES_DB":       "POSTGRES_DB",
				},
				WaitingFor: wait.ForAll(
					wait.ForLog("database system is ready to accept connections"),
					wait.ForListeningPort(postgresPort),
				),
				Files: []container.ContainerFile{{
					HostFilePath:      config.MigrationQueryPath,
					ContainerFilePath: "/docker-entrypoint-initdb.d/" + filepath.Base(config.MigrationQueryPath),
					FileMode:          0o755,
				}},
				SkipReaper: config.SkipReaper,
			},
			Started: true,
			Reuse:   config.Reuse,
		})

	if err != nil {
		log.Fatalf("GenericContainer error: %v", err)
	}

	hostPort, err := postgres.MappedPort(context.Background(), postgresPort)
	if err != nil {
		log.Fatalf("Container host port map error: %v", err)
	}
	containerPort, _ := strconv.Atoi(hostPort.Port())

	dockerDatabase := NewDatabase(
		DatabaseConfig{
			Host:     "127.0.0.1",
			Port:     containerPort,
			User:     "POSTGRES_USER",
			Password: "POSTGRES_PASSWORD",
			Database: "POSTGRES_DB",
		})

	return dockerDatabase
}
