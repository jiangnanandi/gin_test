package main

import (
	"fmt"
	"github.com/didip/tollbooth/limiter"

	"github.com/didip/tollbooth"
	"runtime"
	"time"

	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func main() {
	defer log.Flush()
	testConfig := `
<seelog>
	<outputs formatid="main">
		<buffered  size="102400" flushperiod="2000">
			<file path="./main.log"/>
		</buffered>
	</outputs>
	<formats>
        <format id="main" format="[%LEVEL]: %Date %Time [%Func:%Line] %Msg%n"/>
    </formats>
</seelog>`

	logger, err := log.LoggerFromConfigAsBytes([]byte(testConfig))
	if err != nil {
		fmt.Println(err)
	}
	log.ReplaceLogger(logger)

	runtime.GOMAXPROCS(runtime.NumCPU())
	//r := gin.Default()
	r := gin.New()

	// Create a limiter struct.
	limiter := tollbooth.NewLimiter(1, nil)
	r.GET("/limiter", fLimiter(limiter), cLimiter)

	r.GET("/appendString", AppendString)
	r.GET("/multi", Multi)
	r.Run(":9097")
}

// 中间件，根据入参检查执行限流
func fLimiter(limiter *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		//判断参数
		a := c.Query("a")
		if a == "1" {
			httpError := tollbooth.LimitByRequest(limiter, c.Writer, c.Request)
			if httpError != nil {
				c.Data(httpError.StatusCode, limiter.GetMessageContentType(), []byte(httpError.Message))
				c.Abort()
			} else {
				c.Next()
			}
		} else {
			c.Next()
		}
	}
}

func cLimiter(c *gin.Context) {
	c.String(200, "Hello, world!")
	log.Trace("ok")
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
	/*i := 0
	org  := 1
	for i=1; i<=1000000; i++ {
		org += i
	}*/
	i := 0
	c.JSON(200, gin.H{"cnt": i})
	//c.JSON(200, gin.H{"n": n})
}
