package command

import (
	"fmt"
	"github.com/jader1992/gocore/framework/cobra"
	"github.com/jader1992/gocore/framework/contract"
)

// DemoCommand show current environment
var DemoCommand = &cobra.Command{
	Use:   "demo",
	Short: "demo for framework",
	Run: func(cmd *cobra.Command, args []string) {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)
		fmt.Println("app base folder:", appService.BaseFolder())
	},
}
