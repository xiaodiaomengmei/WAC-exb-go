package main

import (
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	"wifidog-server/app"
	"wifidog-server/framework"
	"wifidog-server/login"
	"wifidog-server/device"
)

type st_sys_config struct {
	WifidogserverHost string `json:"wifidogserver-host"`
	WifidogserverPort string `json:"wifidogserver-port"`
	StaticPath        string `json:"static_path"`
	SslCert           string `json:"ssl_cert"`
	SslKey            string `json:"ssl_key"`
}

const wifidog_server_ver = "1.0.6"

func main() {
	var syscfg_file string
	var dbcfg_file string

	flag.StringVar(&syscfg_file, "syscfg", "config_sys.json", "system config file name")
	flag.StringVar(&dbcfg_file, "dbcfg", "config_db.json", "database config file name")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	log.Println("WAC version: ", wifidog_server_ver)

	sys_file, err := os.Open(syscfg_file)
	if err != nil {
		fmt.Println("open file failed:", err.Error())
		return
	}
	sys_byte, _ := ioutil.ReadAll(sys_file)
	defer sys_file.Close()
	var sys_cfg st_sys_config
	err2 := json.Unmarshal(sys_byte, &sys_cfg)
	if err2 != nil {
		fmt.Println(err2)
	}
	log.Println("sys config:", sys_cfg)
	currentPath, _ := os.Getwd()
	h := http.FileServer(http.Dir(currentPath + "/www/static"))
	http.Handle("/wifidog/static/", http.StripPrefix("/wifidog/static/", h))

	framework.InitDB(dbcfg_file)

	c := cron.New()
	//every day 23:00 clock timer: clear user now trace, summary attendace for day
	c.AddFunc("0 0 23 * * *", app.AttendanceTmrDay)
	//every week Saturday 23:00 clock timer: summary attendace for last weekend
	c.AddFunc("0 0 23 * * 0", app.AttendanceTmrWeekend)
	//every month first day 23:00 clock timer: summary attendace for last month
	c.AddFunc("0 0 23 1 * *", app.AttendanceTmrMonth)
	//half minutes to watch out the ap state
	c.AddFunc("* */5 * * * *", login.ApState)
	//c.AddFunc("*/5 * * * *", app.AttendanceTmrDay)
	c.Start()
	//defer c.Stop()

	login.User_login()
	login.Phone_login()
	login.Authen()
	login.Portal()
	login.Ping()
	login.Setting_pwd()
	login.North_api()
	device.Device_setting()

	wifidogserver := fmt.Sprintf("%s:%s", sys_cfg.WifidogserverHost, sys_cfg.WifidogserverPort)
	ssl_cert := fmt.Sprintf("%s", sys_cfg.SslCert)
	ssl_key := fmt.Sprintf("%s", sys_cfg.SslKey)

	if ssl_cert != "" && ssl_key != "" {
		log.Println("Listen on: ", wifidogserver, "SSL on")
		log.Fatal(http.ListenAndServeTLS(wifidogserver, ssl_cert, ssl_key, nil))
	} else {
		log.Println("Listen on: ", wifidogserver, "SSL off")
		log.Fatal(http.ListenAndServe(wifidogserver, nil))
	}

}
