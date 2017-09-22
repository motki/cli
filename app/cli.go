package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/motki/motki-cli/cli"
	"github.com/motki/motki-cli/cli/command"
	"github.com/motki/motki-cli/cli/text"
	"github.com/motki/motki/app"
)

// A CLIEnv wraps an *app.Env, providing CLI specific facilities.
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

// CLIConfig wraps a *Config and contains optional credentials.
type CLIConfig struct {
	*app.Config

	username string
	password string
}

// NewCLIConfig creates a new CLI-specific configuration using the given conf.
func NewCLIConfig(appConf *app.Config) CLIConfig {
	return CLIConfig{Config: appConf}
}

// WithCredentials returns a copy of the CLIConfig with the given credentials.
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
	srv := cli.NewServer(appEnv.Logger)
	prompter := cli.NewPrompter(srv, appEnv.Client, appEnv.Logger)
	if conf.username != "" || conf.password != "" {
		if _, err := appEnv.Client.Authenticate(conf.username, conf.password); err != nil {
			return nil, err
		} else {
			fmt.Println("Welcome, " + text.Boldf(conf.username) + "!")
		}
	} else {
		fmt.Printf("Welcome to the %s command line interface.\n", text.Boldf("motki"))
		fmt.Println()
		fmt.Printf("Enter \"%s\" in the prompt for detailed help information.\n", text.Boldf("help"))
	}
	srv.SetCommands(
		command.NewEVETypesCommand(prompter),
		command.NewProductCommand(appEnv.Client, prompter, appEnv.Logger))
	if f, err := os.Open(historyPath); err == nil {
		srv.ReadHistory(f)
		f.Close()
	}
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

func (env *CLIEnv) LoopCLI() {
	env.CLI.LoopCLI()
	env.signals <- os.Interrupt
}

// BlockUntilSignal will block until it receives the signals signal.
//
// This function attempts to perform a graceful shutdown, shutting
// down all related services and doing whatever clean up processes are
// necessary.
func (env *CLIEnv) BlockUntilSignal(signals chan os.Signal) {
	if signals == nil {
		signals = env.signals
	}
	abortFuncs := []app.ShutdownFunc{
		func() {
			env.Shutdown()
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
	env.BlockUntilSignalWith(signals, abortFuncs...)
}
