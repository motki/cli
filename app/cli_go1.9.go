//+build go1.9

package app

import (
	"github.com/motki/motki/app"
)

// Config represents the configuration of a MOTKI Env.
//
// This exists as a type alias so that packages depending on this package do not
// need to import both "github.com/motki/motki/app" and "github.com/motki/motki-cli/app".
type Config = app.Config
