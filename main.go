//go:generate fileb0x filebox.json
package main

import (
	"os"
	"path/filepath"

	"github.com/chinanf-boy/vue-with-go-example/server/cmd"
)

func main() {
	cmd.Paths, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	cmd.Execute()
}
