// GoogleMapsCLI — query Google Maps from your terminal.
// Built by Krishna Shahane.
package main

import (
	"io"
	"os"

	"github.com/krishnashahane/googlemapscli/internal/terminal"
)

var quit = os.Exit

func main() {
	quit(start(os.Args[1:], os.Stdout, os.Stderr))
}

func start(args []string, stdout io.Writer, stderr io.Writer) int {
	return terminal.Execute(args, stdout, stderr)
}
