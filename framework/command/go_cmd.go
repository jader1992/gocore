package command

import (
	"github.com/jader1992/gocore/framework/cobra"
	"log"
	"os"
	"os/exec"
)

var goCommand = &cobra.Command{
	Use:   "go",
	Short: "运行path/go程序, 要求go 必须安装",
	RunE: func(c *cobra.Command, args []string) error {
		path, err := exec.LookPath("go")
		if err != nil {
			log.Fatalln("core go: should install go in your PATH")
		}

		cmd := exec.Command(path, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
		return nil
	},
}
