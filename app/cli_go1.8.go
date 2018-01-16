//+build !go1.9

package app

import (
	"github.com/motki/core/app"
)

// Config represents the configuration of a MOTKI Env.
//
// This exists so that packages depending on this package do not need to import
// both "github.com/motki/core/app" and "github.com/motki/motki-cli/app".
type Config app.Config
