# motki-cli

Command `motki` contains interactive command-line tools for managing EVE Online character and corporation assets and industrial processes.

## Getting started

Install `motki` using `go get`.

```bash
go get -u github.com/motki/motki-cli/...
```

Run the program.

```bash
motki
```

> Note that the default configuration connects to `motki.org:18443` using SSL.

### Command-line options

```
Usage of motki:
  -credentials string
    	Username and password separated by a colon. (ie. "frank:mypass")
  -history-file string
    	Path to the CLI history file. (default ".history")
  -insecure-skip-verify
    	INSECURE: Skip verification of server SSL cert.
  -log-level string
    	Log level. Possible values: debug, info, warn, error. (default "warn")
  -server string
    	Backend server host and port. (default "motki.org:18443")
  -version
    	Display the application version.
```

## Authenticating

Some functionality in the application requires authenticating with the remote motkid installation (by default, the [Moritake Industries](https://moritakeindustries.com) website).

To authenticate:

1. Ensure you have a valid account with characters linked on the remote motkid installation.
2. Configure `motki` to use your credentials.
   1. Pass them via command-line option:
   ```
   motki -credentials username:password
   ```
       
   2. Pass them via environment variables:
   ```
   MOTKI_USERNAME=username MOTKI_PASSWORD=password motki
   ```