package controllers

import (
	"douyu/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.html"
}

func (c *MainController) PostDanmu() {
	rec := c.Ctx.Input.RequestBody
	ds := make([]models.Danmu, 0, 1000)
	json.Unmarshal(rec, &ds)
	fmt.Println(ds)
	models.RecDanmu(ds)
	c.Ctx.Output.Context.WriteString("0")
}

func (c *MainController) GetDanmu() {
	from := c.GetString("from")
	res := models.GetDanmu(from)
	c.Ctx.Output.JSON(res, false, false)
}

func (c *MainController) SetCmd() {
	cmd := c.GetString("cmd")
	log.Println(cmd)
	models.SetCommand(cmd)
	c.Ctx.Output.Context.WriteString("success")
}

func (c *MainController) Luck() {
	from := c.GetString("from")
	user := c.GetString("user")
	models.WriteLuck(from, user)
	c.Ctx.Output.Context.WriteString("0")
}

func (c *MainController) Lucks() {
	f := path.Join("static", "file", "luck.json")
	res, _ := ioutil.ReadFile(f)
	c.Ctx.Output.Context.WriteString(string(res))
}

func (c *MainController) SetLuck() {
	from := c.GetString("from")
	user := c.GetString("user")
	text := c.GetString("text")
	models.SetLuck(from, user, text)
	c.Ctx.Output.Context.WriteString("0")
}

func (c *MainController) Process() {
	c.TplName = "process.html"
}

func (c *MainController) Statistics() {
	keys := c.GetString("keyword")
	keyArray := strings.Split(keys, ",")
	file := models.StatisticsDanmu(keyArray)
	c.Ctx.Output.Download(file)
}
