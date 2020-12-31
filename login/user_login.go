package login

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
	"wifidog-server/model"
	"wifidog-server/utils"
)

func User_login() {

	http.HandleFunc("/wifidog/login", func(w http.ResponseWriter, r *http.Request) {
		log.Println("login handle process,method:", r.Method)

		if r.Method == "GET" {
			gw_port := r.URL.Query().Get("gw_port")
			gw_address := r.URL.Query().Get("gw_address")
			gw_id := r.URL.Query().Get("gw_id")
			client_mac := r.URL.Query().Get("mac")
			client_ip := r.URL.Query().Get("ip")
			ssid := r.URL.Query().Get("ssid")
			log.Println("login ap info:", gw_address, gw_id, gw_port, client_mac, ssid)

			//the client mac not in user now trace
			usertrace := usernowtraceDao.GetUserNowTraceByClientMac(client_mac)
			if usertrace == nil {
				if gw_id != "" {
					conn := utils.Pool.Get()
					conn.Do("HMSET", client_mac, "client_ip", client_ip, "gw_id", gw_id, "gw_address", gw_address, "gw_port", gw_port)
					conn.Close()
				}
				ssid_msg := utils.GetSSID(ssid)
				if ssid_msg == "0" {
					t, _ := template.ParseFiles("www/pml_guest_login.html")
					t.Execute(w, client_mac)
				} else {
					t, _ := template.ParseFiles("www/pml_account_login.html")
					t.Execute(w, client_mac)
				}
				log.Println("Not find mac in database,tx login html again.")
			} else {
				usertrace[0].Client_ip = client_ip
				apinfo := apDao.GetApDispInfoByLanMac(gw_id)
				if apinfo == nil {
					log.Println("the ap is invalid!!!")
				} else {
					//update  now trace table
					usertrace[0].Ap_Name = apinfo[0].Name
					usertrace[0].Ap_WanMac = apinfo[0].Wan_mac
					usertrace[0].Ap_Address = apinfo[0].Address
					usertrace[0].Ap_Manufacture = apinfo[0].Manufacture
					usertrace[0].Ap_Model = apinfo[0].Model
					usertrace[0].Ap_Ens = apinfo[0].Ens
					usertrace[0].State = "在线"
					//当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
					usertrace[0].Uptime = time.Now().Format("2006-01-02 15:04:05")
					ok := usernowtraceDao.UpdateUserNowTraceByUsernameAndClientmac(usertrace[0])
					log.Println("update user now trace database. action=", ok)
					io.WriteString(w, "Auth: 2")
				}
			}
		} else {
			var returnData map[string]string
			returnData = make(map[string]string)
			returnData["code"] = "2000"
			returnData["message"] = "success"

			r.ParseForm()
			username := r.PostFormValue("username")
			password := r.PostFormValue("password")
			client_mac := r.PostFormValue("client_mac")
			log.Println("[post form]username=", username, "password=", password, "client_mac", client_mac)
			if len(client_mac) == 0 {
				log.Println("client mac is null,return")
				return
			}

			users := userDao.SelectUserByUsername(username)
			log.Println("get user from database:", users)
			if len(users) == 0 {
				returnData["code"] = "2001"
				returnData["message"] = "用户不存在"
			} else if users[0].Password != password {
				returnData["code"] = "2002"
				returnData["message"] = "密码错误"
			} else {
				client_ip := utils.GetClientMes(client_mac, "client_ip")
				gw_id := utils.GetClientMes(client_mac, "gw_id")
				gw_address := utils.GetClientMes(client_mac, "gw_address")
				gw_port := utils.GetClientMes(client_mac, "gw_port")
				token := utils.GetToken(client_mac)
				log.Println("[redis info] client_mac=", client_mac, "clientip=:", client_ip, "gw_id=", gw_id, "gw_address", gw_address, "gw_port=", gw_port, "token=", token)
				//update user now trace database table
				var nowtraceTemp model.UserNowTrace
				nowtraceTemp.User_Name = username
				nowtraceTemp.User_Human = users[0].Human
				apinfo := apDao.GetApDispInfoByLanMac(gw_id)
				if apinfo == nil {
					log.Println("the ap is invalid!!!")
					return
				}
				//update  now trace table
				nowtraceTemp.Ap_Name = apinfo[0].Name
				nowtraceTemp.Ap_WanMac = apinfo[0].Wan_mac
				nowtraceTemp.Ap_Address = apinfo[0].Address
				nowtraceTemp.Ap_Manufacture = apinfo[0].Manufacture
				nowtraceTemp.Ap_Model = apinfo[0].Model
				nowtraceTemp.Ap_Ens = apinfo[0].Ens
				nowtraceTemp.Client_mac = client_mac
				nowtraceTemp.Client_ip = client_ip
				nowtraceTemp.Token = token
				nowtraceTemp.State = "在线"
				//当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
				nowtraceTemp.Uptime = time.Now().Format("2006-01-02 15:04:05")

				traceid := usernowtraceDao.InsertUserNowTraceByUsernameAndClientmac(nowtraceTemp)
				log.Println("insert user now trace to database. id=", traceid)
				log.Println("user and password check ok, login successed.")
				uri := fmt.Sprintf("http://%s:%s/wifidog/auth?token=%s&mac=%s", gw_address, gw_port, token, client_mac)
				log.Println("redirect url:", uri)
				returnData["uri"] = uri
			}
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		}
	})
}
