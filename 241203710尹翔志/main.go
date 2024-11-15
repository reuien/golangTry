package main

//实现了gin和html的应用吧，在浏览器可以看到用户名
import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testmy.com/gocode/241203710尹翔志/routers"
)

type Test struct {
	ID   string
	Name string
}

func main() {
	r := routers.Routers()
	r.LoadHTMLFiles("testProgramme/test.html")

	test := Test{
		ID:   "123456",
		Name: "dijkstra",
	}

	r.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test.html", gin.H{
			"UserID": test.ID,
			"Name":   test.Name,
		})
	})

	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Run("localhost:8082")
}
