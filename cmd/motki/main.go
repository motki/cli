// Command motki is a utility for interacting with the various data structures
// and processes in the overall motkid project.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/motki/motki-cli/app"
	motki "github.com/motki/motki/app"
	"github.com/motki/motki/log"
	"github.com/motki/motki/model"
	"github.com/motki/motki/proto/client"
)

var serverAddr = flag.String("server", "motki.org:18443", "Backend server host and port.")
var historyPath = flag.String("history-file", ".history", "Path to the CLI history file.")
var credentials = flag.String("credentials", "", "Username and password separated by a colon. (ie. \"frank:mypass\")")
var logLevel = flag.String("log-level", "warn", "Log level. Possible values: debug, info, warn, error.")
var insecureSkipVerify = flag.Bool("insecure-skip-verify", false, "INSECURE: Skip verification of server SSL cert.")
var version = flag.Bool("version", false, "Display the application version.")

// fatalf creates a default logger, writes the given message, and exits.
func fatalf(format string, vals ...interface{}) {
	logger := log.New(log.Config{})
	logger.Fatalf(format, vals...)
}

var Version = "dev"

func main() {
	flag.Parse()
	if *version {
		fmt.Printf("%s %s. %s\n", os.Args[0], Version, "https://github.com/motki/motki-cli")
		os.Exit(0)
	}

	appConf := &motki.Config{
		Backend: model.Config{
			Kind: model.BackendRemoteGRPC,
			RemoteGRPC: model.RemoteConfig{
				ServerAddr:         *serverAddr,
				InsecureSkipVerify: *insecureSkipVerify,
			},
		},
		Logging: log.Config{
			Level: *logLevel,
		},
	}

	// Writing to stderr offers a way to redirect the logger output to a file instead of
	// interrupting the user.
	appConf.Logging.OutputType = log.OutputStderr

	conf := app.NewCLIConfig(appConf)
	if *credentials != "" {
		parts := strings.Split(*credentials, ":")
		if len(parts) != 2 {
			fatalf("motki: invalid credentials, expected format \"username:password\"")
		}
		conf = conf.WithCredentials(parts[0], parts[1])
	} else {
		user, pass := os.Getenv("MOTKI_USERNAME"), os.Getenv("MOTKI_PASSWORD")
		if user != "" && pass != "" {
			conf = conf.WithCredentials(user, pass)
		}
	}

	env, err := app.NewCLIEnv(conf, *historyPath)
	if err != nil {
		if err == client.ErrBadCredentials {
			fmt.Println("Invalid username or password.")
		}
		fatalf("motki: error initializing application environment: %s", err.Error())
	}

	go env.LoopCLI()

	env.BlockUntilSignal(nil)
}
