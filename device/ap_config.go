package device

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"wifidog-server/dao"
	"wifidog-server/model"
	"wifidog-server/utils"
)

var SsidDao = new(dao.SsidDao)
var apDao = new(dao.ApDao)

func Device_setting() {

	//ap configuration
	http.HandleFunc("/wifidog/north_api/apconfig", func(w http.ResponseWriter, r *http.Request) {
		log.Println("north api apconfig:", r.URL.RawQuery, r.Method)

		var returnData map[string]string
		returnData = make(map[string]string)
		var result string
		//请求方式为option
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Action, Module,Content-Type")
			return
		} else if r.Method == "POST" {
			r.ParseForm()
			wan_ip := r.PostFormValue("wan_ip")
			ssid := r.PostFormValue("ssid")
			lan_ip := r.PostFormValue("lan_ip")
			err := utils.ApConfig(wan_ip, ssid, lan_ip)
			if err != nil {
				log.Printf("ap configuration failed,err:%s.", err)
				result = "ok"
			}
		}

		if result != "ok" {
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

	//get ap device infomation from database
	http.HandleFunc("/wifidog/north_api/ap", func(w http.ResponseWriter, r *http.Request) {
		log.Println("north api ap:", r.URL.RawQuery)

		var ApsDisp model.ApDisp

		aps := apDao.GetApAll()
		ApsDisp.PageTotal = 0
		for _, v := range aps {
			ApsDisp.Ap = append(ApsDisp.Ap, v)
			ApsDisp.PageTotal++
		}
		jsonaps, err := json.Marshal(ApsDisp)
		if err != nil {
			fmt.Println("生成json字符串错误")
		}
		jsonstr := fmt.Sprintf("%s", string(jsonaps))
		//fmt.Println(jsonstr)

		w.Header().Set("content-type", "application/json")
		io.WriteString(w, jsonstr)
	})

	//add ap, delete ap, modify ap from database
	http.HandleFunc("/wifidog/north_api/ap_setting", func(w http.ResponseWriter, r *http.Request) {
		log.Println("north api ap_setting :", r.URL.RawQuery, r.Method)

		var returnData map[string]string
		returnData = make(map[string]string)
		var result string
		//请求方式为option
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Action, Module,Content-Type")
			return

		} else if r.Method == "POST" {

			r.ParseForm()
			flag := r.PostFormValue("flag")

			var ap model.ApDispList
			ap.Name = r.PostFormValue("name")
			ap.Address = r.PostFormValue("address")
			ap.Manufacture = r.PostFormValue("manufacture")
			ap.Model = r.PostFormValue("model")
			ap.Ens = r.PostFormValue("ens")
			ap.Wan_mac = r.PostFormValue("wan_mac")
			//flag == "1":add ap
			if flag == "1" {
				result = apDao.AddAp(ap)
			} else {
				//flag == "2":modify ap
				id, _ := strconv.Atoi(r.PostFormValue("id"))
				ap.Id = id
				result = apDao.ModifyAp(ap)
			}
		} else if r.Method == "DELETE" {
			//delete ap
			id, _ := strconv.Atoi(r.FormValue("id"))
			result = apDao.DeleteAp(id)
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

	http.HandleFunc("/wifidog/north_api/onlineap", func(w http.ResponseWriter, r *http.Request) {
		log.Println("north api ap count:", r.URL.RawQuery)
		counts := apDao.SelectcountAp()
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

	//add modify and delete ssid information
	http.HandleFunc("/wifidog/north_api/ssid_setting", func(w http.ResponseWriter, r *http.Request) {

		log.Println("north api ssid_setting :", r.URL.RawQuery, r.Method)

		var returnData map[string]string
		returnData = make(map[string]string)
		var result string
		//请求方式为option
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Action, Module,Content-Type")
			return

		} else if r.Method == "POST" {
			r.ParseForm()
			flag := r.PostFormValue("flag")

			var ssid model.Ssid
			ssid.Ssid = r.PostFormValue("ssid")
			ssid.Jurisdiction = r.PostFormValue("mode")
			//flag == "1":add ap
			if flag == "1" {
				result = SsidDao.AddSsid(ssid)
			} else {
				//flag == "2":modify ap
				ssid.Id, _ = strconv.Atoi(r.PostFormValue("id"))
				result = SsidDao.ModifySsid(ssid)
			}
		} else if r.Method == "DELETE" {
			//delete ap
			ssid_id := r.FormValue("ssid")
			result = SsidDao.DeleteSsid(ssid_id)
		}

		if result == "ok" {
			returnData["code"] = "2000"
			returnData["message"] = "操作成功"
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		} else if result == "1" {
			returnData["code"] = "2002"
			returnData["message"] = "数据库中已经存在该ssid，请重新设置ssid"
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		} else {
			returnData["code"] = "2001"
			returnData["message"] = "操作失败"
			json_byte, _ := json.Marshal(returnData)
			io.WriteString(w, string(json_byte))
		}
	})

	//get all ssid information
	http.HandleFunc("/wifidog/north_api/get_allssid", func(w http.ResponseWriter, r *http.Request) {
		log.Println(" get all ssid:", r.URL)

		var SsidDisp []model.Ssid
		Ssids := SsidDao.GetAllSsid()
		for _, v := range Ssids {
			//fmt.Printf("Id:%d,Ssid:%s,Jurisdiction:%s",v.Id,v.Ssid,v.Jurisdiction)
			SsidDisp = append(SsidDisp, v)
		}
		jsonaps, err := json.Marshal(SsidDisp)

		if err != nil {
			fmt.Println("生成json字符串错误")
		}
		jsonstr := fmt.Sprintf("%s", string(jsonaps))
		//fmt.Println(jsonstr)

		w.Header().Set("content-type", "application/json")
		io.WriteString(w, jsonstr)
	})
}
