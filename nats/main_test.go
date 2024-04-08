package nats

import (
	"fmt"
	"os"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/ory/dockertest/v3"
)

var natsURL string

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		fmt.Println("Could not connect to docker: ", err)
		os.Exit(1)
	}

	container, err := pool.Run("nats", "latest", nil)
	if err != nil {
		fmt.Println("Could not start nats container: ", err)
		os.Exit(1)
	}

	// Tell docker to clean up the container in 5 minutes.
	// This is mainly a failsafe against crashing or being killed leaving
	// orphaned containers running.
	if err := container.Expire(5 * 60); err != nil {
		fmt.Println("Could not set container expiration: ", err)
		os.Exit(1)
	}

	natsURL = fmt.Sprintf("nats://localhost:%s", container.GetPort("4222/tcp"))

	if err := pool.Retry(func() error {
		nc, err := nats.Connect(natsURL)
		if err != nil {
			return err
		}
		defer nc.Close()

		return nil
	}); err != nil {
		fmt.Println("Could not connect to nats: ", err)
		os.Exit(1)
	}

	exitCode := m.Run()

	if err := pool.Purge(container); err != nil {
		fmt.Println("Could not purge container: ", err)
	}

	os.Exit(exitCode)
}
