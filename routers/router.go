package routers

import (
	"douyu/controllers"
	"douyu/ioclient"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/danmu", &controllers.MainController{}, "post:PostDanmu")
	beego.Router("/danmu", &controllers.MainController{}, "get:GetDanmu")
	beego.Router("/setcmd", &controllers.MainController{}, "get:SetCmd")

	beego.Handler("/socket.io", ioclient.GetSocketIO())
	beego.Router("/luck", &controllers.MainController{}, "get:Luck")
	beego.Router("/lucks", &controllers.MainController{}, "get:Lucks")
	beego.Router("/setluck", &controllers.MainController{}, "get:SetLuck")

	beego.Router("/process", &controllers.MainController{}, "get:Process")

	beego.Router("/statistics", &controllers.MainController{}, "get:Statistics")
}
