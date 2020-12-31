package login

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"wifidog-server/dao"
	"wifidog-server/files"
	"wifidog-server/model"
	"wifidog-server/utils"
)

var userDao = new(dao.UserDao)
var apDao = new(dao.ApDao)
var email_cfgDao = new(dao.EmailCfgDao)
var SsidDao = new(dao.SsidDao)

func North_api() {

	//admin login
	http.HandleFunc("/wifidog/north_api/admin_login", func(w http.ResponseWriter, r *http.Request) {
		log.Println("admin login method:", r.Method)

		var returnData map[string]string
		returnData = make(map[string]string)
		returnData["code"] = "2000"
		returnData["message"] = "success"
		returnData["role"] = "1"
		if r.Method == "POST" {
			r.ParseForm()
			username := r.PostFormValue("username")
			password := r.PostFormValue("password")
			users := userDao.SelectUserByUsername(username)
			log.Println("get user from database:", users)
			newpassword := utils.Sha1String(password)
			if len(users) == 0 {
				returnData["code"] = "2001"
				returnData["message"] = "用户不存在"
			} else if users[0].Password != newpassword {
				returnData["code"] = "2002"
				returnData["message"] = "密码错误"
			} else if users[0].Role != "1" {
				returnData["code"] = "0"
				returnData["username"] = users[0].Name
			}
		}
		json_byte, _ := json.Marshal(returnData)
		io.WriteString(w, string(json_byte))
	})

	//get user account infomation from database
	http.HandleFunc("/wifidog/north_api/user_account", func(w http.ResponseWriter, r *http.Request) {
		log.Println("north api user_account :", r.URL.RawQuery)

		var userinfos model.UserDisp

		users := userDao.GetUserAll()
		userinfos.PageTotal = 0
		for _, v := range users {
			userinfos.User = append(userinfos.User, v)
			userinfos.PageTotal++
		}

		jsonusers, err := json.Marshal(userinfos)
		if err != nil {
			fmt.Println("生成json字符串错误")
		}
		jsonstr := fmt.Sprintf("%s", string(jsonusers))
		//log.Println("users : ", jsonstr)
		w.Header().Set("content-type", "application/json")
		io.WriteString(w, jsonstr)
	})

	http.HandleFunc("/wifidog/north_api/setpassword", func(w http.ResponseWriter, r *http.Request) {
		var returnData map[string]string
		returnData = make(map[string]string)
		returnData["code"] = "200"
		returnData["message"] = "success"
		password := r.PostFormValue("password")
		fmt.Println(password)
		username := r.PostFormValue("username")
		newpassword := utils.Sha1String(password)
		var useraccount model.Useraccount
		useraccount.Name = username
		useraccount.Password = newpassword
		result := userDao.ModifyUserByUsername(useraccount)
		if result == "ok" {
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		} else {
			returnData["code"] = "201"
			returnData["message"] = "修改密码失败"
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		}
	})

	//add user, delete user, modify user from database
	http.HandleFunc("/wifidog/north_api/user_setting", func(w http.ResponseWriter, r *http.Request) {
		log.Println("north api user_setting :", r.URL.RawQuery, r.Method)

		var returnData map[string]string
		returnData = make(map[string]string)
		var result string
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Action, Module")
			return
		} else if r.Method == "GET" {
			//get user
			username := r.URL.Query()["username"]
			users := userDao.SelectUserByUsername(username[0])
			jsonusers, _ := json.Marshal(users[0])
			jsonstr := fmt.Sprintf("%s", string(jsonusers))
			io.WriteString(w, jsonstr)
			return
		} else if r.Method == "POST" {
			r.ParseForm()
			flag := r.PostFormValue("flag")
			var user model.Useraccount
			user.Name = r.PostFormValue("username")
			user.Human = r.PostFormValue("realname")
			user.PhoneNumber = r.PostFormValue("phoneNumber")
			user.Email = r.PostFormValue("email")
			user.Department = r.PostFormValue("department")
			//flag == "1":add user
			if flag == "1" {
				token := utils.GetToken(user.Name)
				pwd := token[len(token)-8 : len(token)]
				sha1 := sha1.New()
				sha1.Write([]byte(pwd))
				user.Password = hex.EncodeToString(sha1.Sum([]byte(nil)))
				result = userDao.AddUser(user)
			} else {
				//flag == "2":modify user
				id, _ := strconv.Atoi(r.PostFormValue("id"))
				role := r.PostFormValue("role")
				user.Id = id
				user.Role = role
				result = userDao.ModifyUser(user)
			}
		} else if r.Method == "DELETE" {
			//delete user
			id, _ := strconv.Atoi(r.FormValue("id"))
			result = userDao.DeleteUser(id)
		}

		if result == "ok" {
			returnData["code"] = "2000"
			returnData["message"] = "操作成功"
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		} else {
			returnData["code"] = "2001"
			returnData["message"] = "操作失败"
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		}
	})

	//get user now trace, delete user now trace from database
	http.HandleFunc("/wifidog/north_api/usernowtrace", func(w http.ResponseWriter, r *http.Request) {
		log.Println("north api usernowtrace:", r.URL.RawQuery, r.Method)

		var returnData map[string]string
		returnData = make(map[string]string)
		var result string
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Action, Module")
			return
		} else if r.Method == "GET" {
			var UsernowtraceDisp model.UserNowTraceShow
			usenowtraces := usernowtraceDao.GetUserNowTraceDbAll()
			UsernowtraceDisp.PageTotal = 0
			for _, v := range usenowtraces {
				UsernowtraceDisp.UserNowTraceList = append(UsernowtraceDisp.UserNowTraceList, v)
				UsernowtraceDisp.PageTotal++
			}
			jsonUsernowtraces, err := json.Marshal(UsernowtraceDisp)
			if err != nil {
				fmt.Println("生成json字符串错误")
			}
			jsonstr := fmt.Sprintf("%s", string(jsonUsernowtraces))
			//fmt.Println(jsonstr)
			w.Header().Set("content-type", "application/json")
			io.WriteString(w, jsonstr)
			return
		} else if r.Method == "DELETE" {
			//delete user now trace from database
			id, _ := strconv.Atoi(r.FormValue("id"))
			result = usernowtraceDao.DeleteUserNowTraceById(id)
		}

		if result == "ok" {
			returnData["code"] = "2000"
			returnData["message"] = "操作成功"
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		} else {
			returnData["code"] = "2001"
			returnData["message"] = "操作失败"
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		}
	})

	//get user history trace from database
	http.HandleFunc("/wifidog/north_api/userhistorytracebyhumannameandclientmac", func(w http.ResponseWriter, r *http.Request) {
		log.Println("north api userhistorytrace:", r.URL.RawQuery)

		var UserHistorytraceDisp model.UserHistoryTraceShow

		humanname := r.URL.Query().Get("user_humanname")
		client_mac := r.URL.Query().Get("client_mac")
		log.Println(humanname)
		log.Println(client_mac)

		usehistorytraces := userhistorytraceDao.GetUserHistoryTraceByHumanNameAndClientMac(humanname, client_mac)
		UserHistorytraceDisp.PageTotal = 0
		for _, v := range usehistorytraces {
			UserHistorytraceDisp.UserHistoryTraceList = append(UserHistorytraceDisp.UserHistoryTraceList, v)
			UserHistorytraceDisp.PageTotal++
		}
		jsonUserhistorytraces, err := json.Marshal(UserHistorytraceDisp)
		if err != nil {
			fmt.Println("生成json字符串错误")
		}
		jsonstr := fmt.Sprintf("%s", string(jsonUserhistorytraces))
		//fmt.Println(jsonstr)

		w.Header().Set("content-type", "application/json")
		io.WriteString(w, jsonstr)
	})

	//get /post the email configuration
	http.HandleFunc("/wifidog/north_api/email_cfg", func(w http.ResponseWriter, r *http.Request) {
		log.Println("north api email configuration:", r.URL.RawQuery)

		emailcfgs := email_cfgDao.GetEmailCfg()
		emailcfg := emailcfgs[0]
		jsonEmailCfg, err := json.Marshal(emailcfg)
		if err != nil {
			fmt.Println("生成json字符串错误")
		}
		jsonstr := fmt.Sprintf("%s", string(jsonEmailCfg))
		//fmt.Println(jsonstr)

		w.Header().Set("content-type", "application/json")
		io.WriteString(w, jsonstr)
	})

	http.HandleFunc("/wifidog/north_api/email_setting", func(w http.ResponseWriter, r *http.Request) {
		log.Println("north api email_setting :", r.URL.RawQuery, r.Method)
		var returnData map[string]string
		returnData = make(map[string]string)
		var result string
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Action, Module")
			return
		} else if r.Method == "GET" {
			//get user
			//query := r.URL.Query()["id"]
			//id, _ := strconv.Atoi(query[0])
			//result = userDao.DeleteUser(id)
		} else if r.Method == "POST" {
			r.ParseForm()
			var email model.EmailCfg
			var timer []string
			email.SendAddr = r.PostFormValue("send_address")
			email.SendPwd = r.PostFormValue("send_password")
			email.SmtpHost = r.PostFormValue("smtpaddress")
			port, _ := strconv.Atoi(r.PostFormValue("port"))
			email.RcvAddr = r.PostFormValue("recipient")
			email.CCAddr = r.PostFormValue("copyperson")
			email.Subject = r.PostFormValue("headline")
			email.Body = r.PostFormValue("body")
			time1 := r.PostFormValue("time_type[0]")
			time2 := r.PostFormValue("time_type[1]")
			time3 := r.PostFormValue("time_type[2]")
			if time1 != "" {
				timer = append(timer, time1)
			}
			if time2 != "" {
				timer = append(timer, time2)
			}
			if time3 != "" {
				timer = append(timer, time3)
			}
			for i := 0; i < len(timer); i++ {
				fmt.Println(timer[i])
				switch timer[i] {
				case "每天":
					email.DaySend = "YES"
				case "每周":
					email.WeekSend = "YES"
				case "每月":
					email.MonthSend = "YES"
				}
			}

			email.SmtpPort = port
			result = email_cfgDao.SetEmailCfg(email)
		}
		if result == "ok" {
			returnData["code"] = "2000"
			returnData["message"] = "操作成功"
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		} else {
			returnData["code"] = "2001"
			returnData["message"] = "操作失败"
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		}
	})

	http.HandleFunc("/wifidog/north_api/onlineuser", func(w http.ResponseWriter, r *http.Request) {
		log.Println("north api user count:", r.URL.RawQuery)
		counts := usernowtraceDao.SelectcountUser()
		fmt.Println(counts)
		jsoncounts, err := json.Marshal(counts)
		if err != nil {
			fmt.Println("生成json字符串错误")
		}
		jsonstr := fmt.Sprintf("%s", string(jsoncounts))
		//fmt.Println(jsonstr)

		w.Header().Set("content-type", "application/json")
		io.WriteString(w, jsonstr)
	})

	http.HandleFunc("/wifidog/north_api/delbatch", func(w http.ResponseWriter, r *http.Request) {
		var returnData map[string]string
		var result string
		returnData = make(map[string]string)
		delList := r.URL.Query()
		for _, value := range delList {
			result = usernowtraceDao.DeleteUserNowTraceBatch(value[0])
		}
		if result == "ok" {
			returnData["code"] = "2000"
			returnData["message"] = "操作成功"
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		} else {
			returnData["code"] = "2001"
			returnData["message"] = "操作失败"
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		}
	})

	http.HandleFunc("/wifidog/north_api/registerforexcel", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		file, handler, err := r.FormFile("excelFile")
		if err != nil {
			fmt.Println("form file err: ", err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		//创建上传的目的文件
		f, err := os.OpenFile("./files/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("open file err: ", err)
			return
		}
		defer f.Close()
		//拷贝文件
		io.Copy(f, file)
		pathing := files.Path()
		strs := ReadExcel(pathing)
		fmt.Println(strs)
		strs = strs[5:]
		length := len(strs)
		index := length / 5
		var user model.Useraccount
		var result string
		var returnData map[string]string
		returnData = make(map[string]string)
		for i := 0; i < index; i++ {
			list := strs[:5]
			user.Name = list[0]
			user.Human = list[4]
			user.PhoneNumber = list[2]
			user.Email = list[1]
			user.Department = list[3]
			token := utils.GetToken(user.Name)
			pwd := token[len(token)-8 : len(token)]
			sha1 := sha1.New()
			sha1.Write([]byte(pwd))
			user.Password = hex.EncodeToString(sha1.Sum([]byte(nil)))
			result = userDao.AddUser(user)
			fmt.Println("导入成功")
			strs = strs[5:]
		}
		if result == "ok" {
			returnData["code"] = "2000"
			returnData["message"] = "添加成功"
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		} else {
			returnData["code"] = "2001"
			returnData["message"] = "添加失败"
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		}
	})

}
