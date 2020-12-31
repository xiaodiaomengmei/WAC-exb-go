package dao

import (
	"fmt"
	"log"
	"database/sql"
	"wifidog-server/model"
	"wifidog-server/framework"
	)

type UserNowTraceDao struct {
}
func SelectUserNowTraceTable(SqlStr string)[]model.UserNowTrace  {
	//fmt.Println("user now trace table sql:", SqlStr)
	rows,err:=framework.DB.Query(SqlStr)
	if err !=nil{
		fmt.Println("can not find user in table")
		return nil
	}
	var usernowtraces []model.UserNowTrace
	for rows.Next(){
		var usernowtrace model.UserNowTrace
		var state sql.NullString
		var uptime sql.NullString
		var ap_name sql.NullString
		var ap_wanmac sql.NullString 
		var ap_address sql.NullString 
		var ap_manufacture sql.NullString
		var ap_model sql.NullString
		var ap_ens sql.NullString 
		err:=rows.Scan(&usernowtrace.Id,
			&usernowtrace.User_Name,&usernowtrace.User_Human,
			&state,&uptime,
			&ap_name,&ap_wanmac,&ap_address,&ap_manufacture,
			&ap_model,&ap_ens,
			&usernowtrace.Client_mac,&usernowtrace.Client_ip,
			&usernowtrace.Token)
		usernowtrace.State=state.String
		usernowtrace.Uptime=uptime.String
		usernowtrace.Ap_Name=ap_name.String
		usernowtrace.Ap_WanMac=ap_wanmac.String
		usernowtrace.Ap_Address=ap_address.String
		usernowtrace.Ap_Manufacture=ap_manufacture.String
		usernowtrace.Ap_Model=ap_model.String
		usernowtrace.Ap_Ens=ap_ens.String
		if err !=nil{
			fmt.Println("get UserNowTraceDbElement error",err)
			continue
		}
		usernowtraces=append(usernowtraces,usernowtrace)
	}
	rows.Close()
	return usernowtraces
}
//get all user now trace info to display
func (p *UserNowTraceDao)GetUserNowTraceDbAll()[]model.UserNowTrace  {
	SqlStr := fmt.Sprintf(
		"SELECT id, user_name,human_name,state,uptime,ap_name,ap_wanmac,ap_address,ap_manufacture,ap_model,ap_ens,client_mac,client_ip,token FROM usernowtrace ORDER BY state;")
		
	return SelectUserNowTraceTable(SqlStr)
}

//get element by user name and client mac
func (p *UserNowTraceDao)GetUserNowTraceByUserNameAndClientMac(username string, clientmac string)[]model.UserNowTrace  {
	SqlStr := fmt.Sprintf(
		"SELECT id, user_name,human_name,state,uptime,ap_name,ap_wanmac,ap_address,ap_manufacture,ap_model,ap_ens,client_mac,client_ip,token FROM usernowtrace WHERE user_name='%s' and client_mac='%s';",username,clientmac)
	return SelectUserNowTraceTable(SqlStr)
}
//get element by ap
func (p *UserNowTraceDao)GetUserNowTraceByApWanmac(ap_wanmac string)[]model.UserNowTrace  {
	SqlStr := fmt.Sprintf(
		"SELECT id, user_name,human_name,state,uptime,ap_name,ap_wanmac,ap_address,ap_manufacture,ap_model,ap_ens,client_mac,client_ip,token FROM usernowtrace WHERE ap_wanmac='%s';",ap_wanmac)
	return SelectUserNowTraceTable(SqlStr)
}
//replace get token by client mac
func (p *UserNowTraceDao)GetUserNowTraceByClientMac(clientmac string)[]model.UserNowTrace  {
	SqlStr := fmt.Sprintf(
		"SELECT id, user_name,human_name,state,uptime,ap_name,ap_wanmac,ap_address,ap_manufacture,ap_model,ap_ens,client_mac,client_ip,token FROM usernowtrace WHERE client_mac='%s';",clientmac)
	return SelectUserNowTraceTable(SqlStr)
}
//delete user now trace
func (p *UserNowTraceDao)DeleteUserNowTraceById(id int)string  {
	_, err := framework.DB.Exec("DELETE FROM usernowtrace WHERE id=?;", id)
	if err != nil{
		log.Println("delete usernowtrace failed", err)
		return "error"
	}
	return "ok"
}

//update element
func (p *UserNowTraceDao)UpdateUserNowTraceByUsernameAndClientmac(nowtrace model.UserNowTrace)string  {
	SqlStr := fmt.Sprintf(
		"UPDATE usernowtrace SET user_name='%s',human_name='%s',state='%s',uptime='%s',ap_name='%s',ap_wanmac='%s',ap_address='%s',ap_manufacture='%s',ap_model='%s',ap_ens='%s',client_mac='%s',client_ip='%s',token='%s' WHERE user_name='%s' and client_mac='%s';",
		nowtrace.User_Name,nowtrace.User_Human,nowtrace.State,nowtrace.Uptime,
		nowtrace.Ap_Name,nowtrace.Ap_WanMac,nowtrace.Ap_Address,nowtrace.Ap_Manufacture,nowtrace.Ap_Model,nowtrace.Ap_Ens,
		nowtrace.Client_mac,nowtrace.Client_ip,nowtrace.Token,
		nowtrace.User_Name,nowtrace.Client_mac)
	fmt.Println("update user now trace table sql=",SqlStr)
	_, err:=framework.DB.Exec(SqlStr)
	if err != nil{
		log.Println("set user now trace failed",err)
		return "error"
	}
	return "ok"
}

//insert element
func (p *UserNowTraceDao)InsertUserNowTraceByUsernameAndClientmac(nowtrace model.UserNowTrace)int64{
	SqlStr := fmt.Sprintf(
		"INSERT INTO usernowtrace(user_name,human_name,state,uptime,ap_name,ap_wanmac,ap_address,ap_manufacture,ap_model,ap_ens,client_mac,client_ip,token)VALUE('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s');",
		nowtrace.User_Name,nowtrace.User_Human,nowtrace.State,nowtrace.Uptime,nowtrace.Ap_Name,nowtrace.Ap_WanMac,nowtrace.Ap_Address,nowtrace.Ap_Manufacture,
		nowtrace.Ap_Model,nowtrace.Ap_Ens,nowtrace.Client_mac,nowtrace.Client_ip,nowtrace.Token )
	fmt.Println("inset user now trace table sql=",SqlStr)
	result,err:=framework.DB.Exec(SqlStr)
	if err!=nil{
		log.Println("insert user now trace error")
		return 0
	}
	id,err:=result.LastInsertId()
	if err!=nil{
		log.Println("insert user now trace error")
		return 0
	}
	return id
}
func (p *UserNowTraceDao)ClearAllUserNowTrace(){
	SqlStr := fmt.Sprintf("truncate usernowtrace;")
	fmt.Println("operate user now trace table sql=",SqlStr)
	framework.DB.Exec(SqlStr)
}

func (p *UserNowTraceDao)SelectcountUser()string  {
	rows, err := framework.DB.Query("SELECT count(*) FROM usernowtrace where state='在线'")
	var count string
	if err != nil{
		log.Println("select user count failed", err)
		return "error"
	}
	for rows.Next(){
		err=rows.Scan(&count)
		fmt.Println(count)
	}
	rows.Close()
	return count
}

//delete user now trace batch
func (p *UserNowTraceDao)DeleteUserNowTraceBatch(id string)string  {
	_, err := framework.DB.Exec("DELETE FROM usernowtrace WHERE id=?;", id)
	if err != nil{
		log.Println("delete usernowtracebatch failed", err)
		return "error"
	}
	return "ok"
}
