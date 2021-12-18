
package middleware

import (
	"fmt"
	"github.com/jader1992/gocore/framework/gin"
)

func Test1() gin.HandlerFunc  {
	return func(c *gin.Context) {
		fmt.Println("middleware pre test1")
		c.Next()
		fmt.Println("middleware after test1")
	}
}

func Test2() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("middleware pre test2")
		c.Next()
		fmt.Println("middleware after test2")
	}
}

func Test3() gin.HandlerFunc  {
	return func(c *gin.Context) {
		fmt.Println("middleware pre test3")
		c.Next()
		fmt.Println("middleware after test3")
	}
}
