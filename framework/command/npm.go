package command

import (
	"github.com/jader1992/gocore/framework/cobra"
	"log"
	"os"
	"os/exec"
)

var npmCommand = &cobra.Command{
	Use:   "npm",
	Short: "运行 PATH/npm 的命令",
	RunE: func(c *cobra.Command, args []string) error {
		path, err := exec.LookPath("npm")
		if err != nil {
			log.Fatalln("hade npm: should install npm in your PATH")
		}

		cmd := exec.Command(path, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
		return nil
	},
}
