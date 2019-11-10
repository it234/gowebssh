package routers

import (
	"github.com/it234/gowebssh/internal/app/manageweb/controllers/webssh"

	"github.com/gin-gonic/gin"
)

func RegisterRouterWebssh(app *gin.RouterGroup) {
	serverInfo := webssh.ServerInfo{}
	app.GET("/serverinfo/list", serverInfo.List)
	app.GET("/serverinfo/detail", serverInfo.Detail)
	app.POST("/serverinfo/delete", serverInfo.Delete)
	app.POST("/serverinfo/update", serverInfo.Update)
	app.POST("/serverinfo/create", serverInfo.Create)
	app.GET("/xterm", serverInfo.Xterm)
	app.GET("/ws", serverInfo.Ws)
}
