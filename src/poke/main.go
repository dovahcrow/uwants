package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"poke/controllers"
	_ "poke/routers"
	"time"
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
	wait, err := cfg.Int(`waittime`)
	if err != nil {
		fmt.Println(`parse wait time error`)
		return
	}
	controllers.Waittime = time.Duration(wait)
	beego.Run()
}
