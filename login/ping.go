package login

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wifidog-server/dao"
	"wifidog-server/model"
	"wifidog-server/utils"
)

func Ping() {
	http.HandleFunc("/wifidog/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Println("ping message : ", r)
		//save ap database
		var ApDbInfo model.ApDbInfo
		ApDbInfo.Wan_mac = r.URL.Query().Get("wan_mac")
		ApDbInfo.Wan_ip = r.URL.Query().Get("wan_ip")
		ApDbInfo.Lan_mac = r.URL.Query().Get("gw_id")
		ApDbInfo.Lan_ip = r.URL.Query().Get("gw_address")
		uptime, err := strconv.Atoi(r.URL.Query().Get("sys_uptime"))
		if err != nil {
			log.Println("字符串转换成整数失败")
			uptime = 0
		}
		ApDbInfo.Uptime = uptime
		ApDbInfo.Memory = r.URL.Query().Get("sys_memfree")
		ApDbInfo.State = "on"
		ret := apDao.UpdateApLastedTimeByWanmac(ApDbInfo.Wan_mac)
		if ret != "ok" {
			fmt.Println("更新ap最新在线时间错误")
			return
		}

		if len(ApDbInfo.Wan_mac) == 0 {
			log.Println("the ap deivce not register, Invalid!")
			fmt.Fprintf(w, "AP Invalid,not register!")
			return
		}

		//check attach client
		clientmac := r.URL.Query().Get("clientmac")
		usernowtraces := usernowtraceDao.GetUserNowTraceByApWanmac(ApDbInfo.Wan_mac)
		for _, v := range usernowtraces {
			user_clientmac := v.Client_mac
			//the client leave the ap
			if true != strings.Contains(clientmac, user_clientmac) {
				usertraceByclients := usernowtraceDao.GetUserNowTraceByClientMac(user_clientmac)
				if usertraceByclients != nil {
					//write user history trace table
					userhistorytraceDao.AddUserHistoryTrace(usertraceByclients[0])
					log.Println("AddUserHistoryTrace case1 : user offline")
					//update user now trace table
					usertraceByclients[0].State = "离线"
					usertraceByclients[0].Uptime = ""
					usertraceByclients[0].Ap_Name = ""
					usertraceByclients[0].Ap_WanMac = ""
					usertraceByclients[0].Ap_Address = ""
					usertraceByclients[0].Ap_Manufacture = ""
					usertraceByclients[0].Ap_Model = ""
					usertraceByclients[0].Ap_Ens = ""
					usertraceByclients[0].Client_ip = ""
					usernowtraceDao.UpdateUserNowTraceByUsernameAndClientmac(usertraceByclients[0])
				}
			}
		}

		//update ap table
		ok := apDao.UpdateApByWanMac(ApDbInfo)
		if ok != "ok" {
			log.Println("AP infomation update failed!")
		}
		//kick off the client in ap

		ackStr := fmt.Sprintf("Pong")
		fmt.Fprintf(w, ackStr)
	})
}

//定时器每个3分钟去比较一下数据库中最新的时间
func ApState(){
	//获取ap数据库中所有的ap信息
	var apDao = new(dao.ApDao)
	aps := apDao.GetApAll()
	for _,v := range aps {
		if strings.EqualFold(v.State,"离线"){
			continue
		}
		time_part, _ := time.ParseDuration(utils.Str2Duration(v.Latestonline_time,utils.Time2Str()))
		if time_part.Minutes() >3 {
			apDao.UpdateApState(v)
		}
	}
}
