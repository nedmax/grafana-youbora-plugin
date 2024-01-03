//go:build mage
// +build mage

package main

import (
	"fmt"
	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
)

// Build build for all archs
func Build() {
	fmt.Println("Building plugin")
	build.BuildAll()
}

// Default configures the default target.
var Default = Build
