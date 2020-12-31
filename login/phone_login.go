package login

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
	"wifidog-server/model"
	"wifidog-server/utils"
)

func Phone_login() {

	http.HandleFunc("/wifidog/sendSmsValidate", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		phoneNumber := r.PostFormValue("phoneNumber")
		//mimic frame,add timestamp in http header
		Uid := r.Header.Get("Uid")
		var randomNumber string
		random := utils.GetSmsCode(phoneNumber, Uid)
		randomNumber = random[:6]
		client, _ := dysmsapi.NewClientWithAccessKey("cn-beijing", "LTAI4FhvWqFjMgDCLHRoWYff", "abJ29RvnOnWHDh87sFq58OHQi1GRCx")
		request := dysmsapi.CreateSendSmsRequest()
		request.Method = "POST"
		request.Scheme = "http" // https | http
		request.Domain = "dysmsapi.aliyuncs.com"
		request.QueryParams["RegionId"] = "cn-beijing"
		request.QueryParams["PhoneNumbers"] = phoneNumber     //手机号
		request.QueryParams["SignName"] = "拟态安全WAC"            //阿里云验证过的项目名 自己设置
		request.QueryParams["TemplateCode"] = "SMS_203185432" //阿里云的短信模板号 自己设置
		request.QueryParams["TemplateParam"] = "{\"code\":\"" + randomNumber + "\"}"
		response, _ := client.SendSms(request)
		log.Println("response : ", response)
		io.WriteString(w, "success")
	})

	http.HandleFunc("/wifidog/loginByPhone", func(w http.ResponseWriter, r *http.Request) {
		log.Println("loginByPhone,method=", r.Method)
		if r.Method == "GET" {
			client_mac := r.FormValue("client_mac")
			log.Println("client_mac=", client_mac)
			t, _ := template.ParseFiles("www/pml_phone_login.html")
			t.Execute(w, client_mac)
		} else {
			var returnData map[string]string
			returnData = make(map[string]string)
			returnData["code"] = "2000"
			returnData["message"] = "success"

			r.ParseForm()
			phoneNumber := r.PostFormValue("phoneNumber")
			verifyCode := r.PostFormValue("verifyCode")
			client_mac := r.PostFormValue("client_mac")
			flag := r.PostFormValue("flag")
			token := utils.GetToken(client_mac)
			log.Println("[post form]phoneNumber=", phoneNumber, "verifyCode=", verifyCode, "client_mac", client_mac)
			log.Println("get token:", token)
			if len(client_mac) == 0 {
				fmt.Println("client mac is null,return")
				return
			}
			users := userDao.SelectUserByPhoneNumber(phoneNumber)
			if flag == "1" {
				if users == nil {
					returnData["code"] = "2001"
					returnData["message"] = "用户不存在"
					json_byte, _ := json.Marshal(returnData)
					io.WriteString(w, string(json_byte))
					return
				}
			}
			if utils.RedisUtilget(phoneNumber) == "" {
				returnData["code"] = "2002"
				returnData["message"] = "验证码发送失败"
			} else if utils.RedisUtilget(phoneNumber) != verifyCode {
				returnData["code"] = "2002"
				returnData["message"] = "验证码错误"
			} else {
				client_ip := utils.GetClientMes(client_mac, "client_ip")
				gw_id := utils.GetClientMes(client_mac, "gw_id")
				gw_address := utils.GetClientMes(client_mac, "gw_address")
				gw_port := utils.GetClientMes(client_mac, "gw_port")
				token := utils.GetToken(client_mac)
				log.Println("[redis info]client_mac=", client_mac, "clientip=:", client_ip, "gw_id=", gw_id, "gw_address", gw_address, "gw_port=", gw_port, "token=", token)
				//update user now trace database table
				var nowtraceTemp model.UserNowTrace
				if flag == "0" {
					nowtraceTemp.User_Name = phoneNumber
					nowtraceTemp.User_Human = "访客"
				} else {
					nowtraceTemp.User_Name = users[0].Name
					nowtraceTemp.User_Human = users[0].Human
				}
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
				log.Println("user phone number and verify code check ok, login successed.")
				uri := fmt.Sprintf("http://%s:%s/wifidog/auth?token=%s&mac=%s", gw_address, gw_port, token, client_mac)
				returnData["uri"] = uri
			}
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		}
	})
}
