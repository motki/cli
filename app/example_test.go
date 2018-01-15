package app_test

import (
	"fmt"
	"time"

	"github.com/motki/motki-cli/app"
	"github.com/motki/motki/log"
	"github.com/motki/motki/proto"
	"github.com/motki/motki/proto/client"
)

// ExampleNewCLIEnv_LoopCLI shows the bare-minimum to connect to the public MOTKI
// application and start an interactive CLI session.
func ExampleNewCLIEnv_LoopCLI() {
	conf := app.NewCLIConfig(&app.Config{
		Backend: proto.Config{
			Kind:       proto.BackendRemoteGRPC,
			RemoteGRPC: proto.RemoteConfig{ServerAddr: "motki.org"},
		},
		Logging: log.Config{
			Level: "debug",
		},
	})

	// A blank history path puts the history file in the current working directory.
	historyPath := ""

	// Create the CLI environment.
	env, err := app.NewCLIEnv(conf, historyPath)
	if err != nil {
		panic("motki: error initializing application environment: " + err.Error())
	}

	// Begin reading standard input and executing commands.
	go env.LoopCLI()

	// Since this is a proper test and example, exit after 1 second.
	//
	// Remove this select statement if you want to actually run the CLI.
	select {
	case <-time.After(1 * time.Second):
		break
	}

	// Output:
	// Welcome to the app command line interface.
	//
	// Enter "help" in the prompt for detailed help information.
}

// ExampleCLIEnv_WithCredentials shows how to use the authenticate method on a CLIEnv.
//
// Credentials are set before creating the CLIEnv. If authentication fails, for any reason,
// an error is returned and a CLIEnv does not initialize.
//
// Since this test runs in a sandbox with no network access, this test always fails.
func ExampleCLIEnv_WithCredentials() {
	conf := app.NewCLIConfig(&app.Config{
		Backend: proto.Config{
			Kind:       proto.BackendRemoteGRPC,
			RemoteGRPC: proto.RemoteConfig{ServerAddr: "motki.org"},
		},
		Logging: log.Config{
			Level: "debug",
		},
	})

	// MOTKI.org credentials. These are invalid.
	conf = conf.WithCredentials("frank", "mypass")

	// A blank history path puts the history file in the current working directory.
	historyPath := ""

	// Create the CLI environment.
	//
	// We don't use it in this test since we always fail authentication, so we have to leave
	// it un-named for the example to compile.
	_, err := app.NewCLIEnv(conf, historyPath)
	if err != nil {
		if err == client.ErrBadCredentials {
			fmt.Println("Invalid username or password.")
		}
		fmt.Println("motki: error initializing application environment: " + err.Error())
	}

	// Output:
	// motki: error initializing application environment: rpc error: code = Unavailable desc = all SubConns are in TransientFailure
}
