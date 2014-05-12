package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	_ "poke/controllers"
	_ "poke/routers"

	"uwants"
)

func main() {
	cfg, err := config.NewConfig(`ini`, `config.conf`)
	if err != nil {
		fmt.Println(`open config error`)
		return
	}
	proxy := cfg.String(`proxy`)
	uwants.Proxy = proxy

	beego.SessionOn = true
	beego.SessionCookieLifeTime = 3600
	beego.SessionName = `boss!`
	beego.SessionProvider = `memory`
	beego.Run()
}
