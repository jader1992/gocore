package foo

import (
	"fmt"

	"github.com/jader1992/gocore/framework/cobra"
)

var FooCommand = &cobra.Command{  // ./gocore command create 生成命令工具
	Use:   "foo",
	Short: "foo",
	RunE: func(c *cobra.Command, args []string) error {
        container := c.GetContainer()
		fmt.Println(container)
        fmt.Println("./gocore command create 生成命令工具")
		return nil
	},
}

