package main

import (
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	r := gin.Default()
	r.GET("/appendString", AppendString)
	r.GET("/multi", Multi)
	r.Run(":9097")
}

func AppendString(c *gin.Context) {
	org := ""
	i := 0
	for ; i <= 1000; i++ {
		org += time.Now().String() + "__"
	}
	c.JSON(200, gin.H{"cnt": i})
}

func Multi(c *gin.Context) {
	i := 0
	org  := 1
	for i=1; i<=1000000; i++ {
		org += i
	}
	c.JSON(200, gin.H{"cnt": i})
}