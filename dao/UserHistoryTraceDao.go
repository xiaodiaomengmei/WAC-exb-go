package dao

import (
	"fmt"
	"log"
	"time"
	"wifidog-server/framework"
	//"database/sql"
	"wifidog-server/model"
)

type UserHistoryTraceDao struct {
}

var userDao = new(UserDao)

func SelectUserHistoryTraceTable(SqlStr string) []model.UserHistoryTrace {
	//fmt.Println("user history trace table sql:", SqlStr)
	rows, err := framework.DB.Query(SqlStr)
	if err != nil {
		fmt.Println("can not find user in table")
		return nil
	}
	//all table element must be not null
	var userhistorytraces []model.UserHistoryTrace
	for rows.Next() {
		var userhistorytrace model.UserHistoryTrace

		err := rows.Scan(&userhistorytrace.Id,
			&userhistorytrace.User_Name, &userhistorytrace.User_Human, &userhistorytrace.User_Phone, &userhistorytrace.User_Email, &userhistorytrace.User_Department,
			&userhistorytrace.LoginDate, &userhistorytrace.Logintime, &userhistorytrace.Logouttime,
			&userhistorytrace.Ap_Name, &userhistorytrace.Ap_WanMac, &userhistorytrace.Ap_Address, &userhistorytrace.Ap_Manufacture, &userhistorytrace.Ap_Model, &userhistorytrace.Ap_Ens,
			&userhistorytrace.Client_mac, &userhistorytrace.Client_ip)

		if err != nil {
			fmt.Println("get UserHistoryTraceDbElement error", err)
			continue
		}
		userhistorytraces = append(userhistorytraces, userhistorytrace)
	}
	rows.Close()
	return userhistorytraces
}

//get all user history trace info to display
func (p *UserHistoryTraceDao) GetUserHistoryTraceDbAll() []model.UserHistoryTrace {
	SqlStr := fmt.Sprintf(
		"SELECT id, user_name,user_humanname,user_phnumber,user_email,user_department,login_date,login_time,logout_time,ap_name,ap_wanmac,ap_address,ap_manufacture,ap_model,ap_ens,client_mac,client_ip FROM userhistorytrace;")

	return SelectUserHistoryTraceTable(SqlStr)
}

//get element by user name,used for display
func (p *UserHistoryTraceDao) GetUserHistoryTraceByHumanNameAndClientMac(username, client_mac string) []model.UserHistoryTrace {
	SqlStr := fmt.Sprintf(
		"SELECT id, user_name,user_humanname,user_phnumber,user_email,user_department,login_date,login_time,logout_time,ap_name,ap_wanmac,ap_address,ap_manufacture,ap_model,ap_ens,client_mac,client_ip FROM userhistorytrace WHERE user_humanname='%s' And client_mac='%s';", username, client_mac)
	return SelectUserHistoryTraceTable(SqlStr)
}

//get element by user name and login date,used for report xls
func (p *UserHistoryTraceDao) GetUserHistoryTraceByUserNameAndLoginDate(username string, logindate string) []model.UserHistoryTrace {
	SqlStr := fmt.Sprintf(
		"SELECT id, user_name,user_humanname,user_phnumber,user_email,user_department,login_date,login_time,logout_time,ap_name,ap_wanmac,ap_address,ap_manufacture,ap_model,ap_ens,client_mac,client_ip FROM userhistorytrace WHERE user_name='%s' and login_date='%s';", username, logindate)
	return SelectUserHistoryTraceTable(SqlStr)
}

//insert element
func InsertUserHistoryTrace(historytrace model.UserHistoryTrace) int64 {
	SqlStr := fmt.Sprintf(
		"INSERT INTO userhistorytrace(user_name,user_humanname,user_phnumber,user_email,user_department,login_date,login_time,logout_time,ap_name,ap_wanmac,ap_address,ap_manufacture,ap_model,ap_ens,client_mac,client_ip)VALUE('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s');",
		historytrace.User_Name, historytrace.User_Human, historytrace.User_Phone, historytrace.User_Email, historytrace.User_Department,
		historytrace.LoginDate, historytrace.Logintime, historytrace.Logouttime,
		historytrace.Ap_Name, historytrace.Ap_WanMac, historytrace.Ap_Address, historytrace.Ap_Manufacture, historytrace.Ap_Model, historytrace.Ap_Ens,
		historytrace.Client_mac, historytrace.Client_ip)
	fmt.Println("inset user history trace table sql=", SqlStr)
	result, err := framework.DB.Exec(SqlStr)
	if err != nil {
		log.Println("insert user history trace error")
		return 0
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("insert user history trace error")
		return 0
	}
	return id
}

func (p *UserHistoryTraceDao) AddUserHistoryTrace(nowtrace model.UserNowTrace) {
	//if the now trace ap name is not null, roam, need update history trace table
	var userhistorytrace model.UserHistoryTrace
	userAccounts := userDao.SelectUserByUsername(nowtrace.User_Name)
	if userAccounts != nil {
		userhistorytrace.User_Name = userAccounts[0].Name
		userhistorytrace.User_Human = userAccounts[0].Human
		userhistorytrace.User_Phone = userAccounts[0].PhoneNumber
		userhistorytrace.User_Email = userAccounts[0].Email
		userhistorytrace.User_Department = userAccounts[0].Department
	} else {
		userhistorytrace.User_Name = "visitor"
		userhistorytrace.User_Human = "visitor"
		userhistorytrace.User_Phone = nowtrace.User_Name
		userhistorytrace.User_Email = ""
		userhistorytrace.User_Department = ""
	}
	userhistorytrace.Ap_Name = nowtrace.Ap_Name
	userhistorytrace.Ap_WanMac = nowtrace.Ap_WanMac
	userhistorytrace.Ap_Address = nowtrace.Ap_Address
	userhistorytrace.Ap_Manufacture = nowtrace.Ap_Manufacture
	userhistorytrace.Ap_Model = nowtrace.Ap_Model
	userhistorytrace.Ap_Ens = nowtrace.Ap_Ens
	userhistorytrace.Client_mac = nowtrace.Client_mac
	userhistorytrace.Client_ip = nowtrace.Client_ip
	userhistorytrace.LoginDate = time.Now().Format("2006-01-02")
	userhistorytrace.Logintime = nowtrace.Uptime
	userhistorytrace.Logouttime = time.Now().Format("2006-01-02 15:04:05")
	traceid := InsertUserHistoryTrace(userhistorytrace)
	log.Println("insert user history trace to database. id=", traceid)
	return
}
