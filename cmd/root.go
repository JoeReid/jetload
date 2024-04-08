package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/JoeReid/jetload/nats"
	"github.com/JoeReid/jetload/specfile"
	"github.com/spf13/cobra"
)

func Execute() error {
	return (&rootCmd{fsys: &defaultFS{}}).cmd().Execute()
}

type rootCmd struct {
	fsys fs.FS

	wait    bool
	timeout time.Duration
	natsURL string
}

func (r *rootCmd) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jetload",
		Short: " A simple CLI tool to load test data into a NATS: JetStream subject",
		Long:  " A simple CLI tool to load test data into a NATS: JetStream subject",
		RunE:  r.run,
	}

	cmd.Flags().BoolVarP(&r.wait, "wait", "w", false, "Wait until all messages are consumed before exiting")
	cmd.Flags().DurationVarP(&r.timeout, "timeout", "t", time.Minute, "Timeout for waiting for messages to be consumed")
	cmd.Flags().StringVarP(&r.natsURL, "nats-url", "n", "nats://127.0.0.1:4222", "URL of the NATS server")

	return cmd
}

func (r *rootCmd) run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	if len(args) == 0 {
		return errors.New("no files provided")
	}

	files, err := specfile.LoadPaths(r.fsys, args...)
	if err != nil {
		return err
	}

	loader, err := nats.NewLoader(r.natsURL)
	if err != nil {
		return err
	}
	defer loader.Close()

	for _, file := range files {
		if err := loader.Load(ctx, file, r.wait); err != nil {
			return fmt.Errorf("error loading file: %s", err.Error())
		}
	}

	return nil
}

type defaultFS struct{}

func (d defaultFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}
