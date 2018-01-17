// Package app contains functionality related to creating an interactive
// command-line interface environment with all the necessary dependencies.
//
// The goal of the app package is to provide a single, reusable base for
// building client-only MOTKI CLI applications.
//
// This package imports every other motki/cli package. As such, it cannot
// be imported from the "library" portion of the project. It is intended to be
// used from external packages only. For a real example of this, check the motki
// command source code.
package app // import "github.com/motki/cli/app"

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/motki/core/app"

	"github.com/motki/cli"
	"github.com/motki/cli/command"
	"github.com/motki/cli/text"
)

// A CLIEnv wraps a ClientEnv, providing CLI specific facilities.
type CLIEnv struct {
	*app.ClientEnv
	CLI      *cli.Server
	Prompter *cli.Prompter

	// historyPath is the path to the CLI history file.
	historyPath string

	// signals is used to shutdown the program.
	//
	// Write any os.Signal to this channel and the program will attempt to
	// shutdown gracefully.
	//
	// The implication is os.Exit() should not be used; stick solely to exiting
	// by writing to this channel to exit.
	signals chan os.Signal
}

// CLIConfig wraps an app.Config and contains optional credentials.
type CLIConfig struct {
	*app.Config

	username string
	password string
}

// NewCLIConfig creates a CLI-specific configuration using the given conf.
func NewCLIConfig(base *Config) CLIConfig {
	// TODO: Need this for compatibility with Go 1.8 and earlier.
	conf := app.Config(*base)
	return CLIConfig{Config: &conf}
}

// WithCredentials returns a copy of the CLIConfig with the given credentials embedded.
func (c CLIConfig) WithCredentials(username, password string) CLIConfig {
	return CLIConfig{c.Config, username, password}
}

// NewCLIEnv initializes a new CLI environment.
//
// If the given CLIConfig contains a username or password, authentication
// will be attempted. If authentication fails, an error is returned.
func NewCLIEnv(conf CLIConfig, historyPath string) (*CLIEnv, error) {
	appEnv, err := app.NewClientEnv(conf.Config)
	if err != nil {
		return nil, err
	}
	if !filepath.IsAbs(historyPath) {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		historyPath = filepath.Join(cwd, historyPath)
	}
	if conf.username != "" || conf.password != "" {
		if err = appEnv.Client.Authenticate(conf.username, conf.password); err != nil {
			return nil, err
		} else {
			fmt.Println("Welcome, " + text.Boldf(conf.username) + "!")
		}
	} else {
		fmt.Printf("Welcome to the %s command line interface.\n", text.Boldf("app"))
		fmt.Println()
		fmt.Printf("Enter \"%s\" in the prompt for detailed help information.\n", text.Boldf("help"))
	}

	srv := cli.NewServer(appEnv.Logger)
	if f, err := os.Open(historyPath); err == nil {
		srv.ReadHistory(f)
		f.Close()
	}
	prompter := cli.NewPrompter(srv, appEnv.Client, appEnv.Logger)

	srv.SetCommands(
		command.NewEVETypesCommand(prompter),
		command.NewProductCommand(appEnv.Client, prompter, appEnv.Logger))
	srv.SetCtrlCAborts(true)
	env := &CLIEnv{
		ClientEnv: appEnv,
		CLI:       srv,
		Prompter:  prompter,

		historyPath: historyPath,

		signals: make(chan os.Signal, 1),
	}

	return env, nil
}

// LoopCLI starts an endless loop to perform commands read from stdin.
//
// This function is intended to be started in a goroutine. The loop will end when
// the application exits or when the associated Env receives an abort signal. Use this
// method in combination with env.BlockUntilSignal() for clean exits.
//   env, err := app.NewCLIEnv(conf, "")
//   // handle err...
//
//   go env.LoopCLI()
//
//   env.BlockUntilSignal()
func (env *CLIEnv) LoopCLI() {
	env.CLI.LoopCLI()
	env.signals <- os.Interrupt
}

// BlockUntilSignal will block until it receives the signals signal.
//
// This function attempts to perform a graceful shutdown, shutting
// down all related services and doing whatever clean up processes are
// necessary.
func (env *CLIEnv) BlockUntilSignal(s chan os.Signal) {
	if s == nil {
		s = env.signals
	}
	abortFuncs := []app.ShutdownFunc{
		func() {
			if err := env.Scheduler.Shutdown(); err != nil {
				env.Logger.Warnf("app: error shutting down scheduler: %s", err.Error())
			}
		},
		func() {
			if f, err := os.Create(env.historyPath); err == nil {
				env.CLI.WriteHistory(f)
				f.Close()
			} else {
				env.Logger.Warnf("unable to write CLI history: %s", err.Error())
			}
			env.CLI.Close()
		}}
	env.BlockUntilSignalWith(s, abortFuncs...)
}
