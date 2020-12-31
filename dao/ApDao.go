package dao

import (
	"fmt"
	"log"
	"time"
	"database/sql"
	"wifidog-server/model"
	"wifidog-server/framework"
	)

type ApDao struct {

}

//向数据库中增加ap
func (p *ApDao)AddAp(ap model.ApDispList)string  {
	_, err := framework.DB.Exec("INSERT INTO ap (name, wan_mac, address,manufacture, model, ens) VALUES (?, ?, ?, ?, ?,?);",
		ap.Name, ap.Wan_mac, ap.Address, ap.Manufacture, ap.Model,ap.Ens)
	if err != nil{
		log.Println("add ap failed", err)
		return "error"
	}
	return "ok"
}

//修改数据库中的ap信息
func (p *ApDao)ModifyAp(ap model.ApDispList)string  {
	_, err:=framework.DB.Exec("UPDATE ap SET name=?, address=?,manufacture=?, model=?, ens=? ,wan_mac=? WHERE id=?;", ap.Name, ap.Address, ap.Manufacture, ap.Model,ap.Ens,ap.Wan_mac,ap.Id)


	if err != nil{
		log.Println("modify ap failed", err)
		return "error"
	}
	return "ok"
}

//删除数据库中的ap
func (p *ApDao)DeleteAp(id int)string  {
	_, err := framework.DB.Exec("DELETE FROM ap WHERE id=?;", id)
	if err != nil{
		log.Println("delete ap failed", err)
		return "error"
	}
	return "ok"
}

func (p *ApDao)UpdateApByWanMac(ap model.ApDbInfo)string  {
	ApSql := fmt.Sprintf(
		"UPDATE ap SET wan_ip='%s',lan_mac='%s',lan_ip='%s', uptime=%d, state='%s' WHERE wan_mac='%s';", 
		ap.Wan_ip, ap.Lan_mac,ap.Lan_ip,ap.Uptime,ap.State,ap.Wan_mac)
	//fmt.Println("ap table sql:", ApSql)

	_, err:=framework.DB.Exec(ApSql)
	if err != nil{
		fmt.Println("update ap info failed!",err)
		return "error"
	}
	return "ok"
}
	
//get ap by gw_id(lan_mac).only one
func (p *ApDao)GetApDispInfoByLanMac(lan_mac string)[]model.ApDispList  {
	ApSql := fmt.Sprintf(
		"SELECT id,name,address,manufacture,model,ens,state,wan_mac,wan_ip,lan_mac,lan_ip,uptime FROM ap WHERE lan_mac='%s';", 
		lan_mac)
	//fmt.Println("ap table sql:", ApSql)
	rows, err:=framework.DB.Query(ApSql)
	if err != nil{
		fmt.Println("find ap info failed!",err)
		return nil
	}
	var aps []model.ApDispList
	for rows.Next(){
		var ap model.ApDispList
		var Manufacture sql.NullString 
		var Model sql.NullString 
		var Ens sql.NullString 
		var State sql.NullString 
		var Wan_mac sql.NullString
		var Wan_ip sql.NullString
		var Lan_mac sql.NullString 
		var Lan_ip sql.NullString 
		var Uptime sql.NullString 

		err:=rows.Scan(&ap.Id,&ap.Name,&ap.Address,
			&Manufacture,&Model,&Ens,
			&State,
			&Wan_mac,&Wan_ip,&Lan_mac,&Lan_ip,&Uptime)
		if State.String == "on"{
			ap.State = "在线"
		}else{
			ap.State = "离线"
		}
		ap.Manufacture=	Manufacture.String
		ap.Model=	Model.String
		ap.Ens=	Ens.String
		ap.Wan_mac=	Wan_mac.String
		ap.Wan_ip=	Wan_ip.String
		ap.Lan_mac=	Lan_mac.String
		ap.Lan_ip=	Lan_ip.String
		ap.Uptime=	Uptime.String
		//fmt.Println(ap)

		if err !=nil{
			fmt.Println("get ap error")
			continue
		}
		aps=append(aps,ap)
	}
	rows.Close()
	return aps
}

func (p *ApDao)GetApAll()[]model.ApDispList  {
	ApSql := fmt.Sprintf("SELECT id,name,address,manufacture,model,ens,state,wan_mac,wan_ip,lan_mac,lan_ip,uptime,latestonline_time FROM ap ORDER BY state desc;")
	//fmt.Println("ap table sql:", ApSql)

	rows,err:=framework.DB.Query(ApSql)
	if err !=nil{
		fmt.Println("can not find ap in table")
		return nil
	}
	var aps []model.ApDispList
	for rows.Next(){
		var ap model.ApDispList
		var Manufacture sql.NullString 
		var Model sql.NullString 
		var Ens sql.NullString 
		var State sql.NullString 
		var Wan_mac sql.NullString
		var Wan_ip sql.NullString
		var Lan_mac sql.NullString 
		var Lan_ip sql.NullString 
		var Uptime sql.NullString
		var Latestonline_time sql.NullString

		err:=rows.Scan(&ap.Id,&ap.Name,&ap.Address,
			&Manufacture,&Model,&Ens,
			&State,
			&Wan_mac,&Wan_ip,&Lan_mac,&Lan_ip,&Uptime,&Latestonline_time)
		if State.String == "on"{
			ap.State = "在线"
		}else{
			ap.State = "离线"
		}
		ap.Manufacture=	Manufacture.String
		ap.Model=	Model.String
		ap.Ens=	Ens.String
		ap.Wan_mac=	Wan_mac.String
		ap.Wan_ip=	Wan_ip.String
		ap.Lan_mac=	Lan_mac.String
		ap.Lan_ip=	Lan_ip.String
		if Uptime.String == "0"{
			ap.Uptime=""
		}else{
			ap.Uptime=	Uptime.String
		}
		ap.Latestonline_time=Latestonline_time.String
		//fmt.Println(ap)

		if err !=nil{
			fmt.Println("get ap error")
			continue
		}
		aps=append(aps,ap)
	}
	rows.Close()
	return aps
}

func (p *ApDao)SelectcountAp()string  {
	rows, err := framework.DB.Query("SELECT count(*) FROM ap")
	var count string
	if err != nil{
		log.Println("select ap count failed", err)
		return "error"
	}
	for rows.Next(){
		err=rows.Scan(&count)
		fmt.Println(count)
	}
	rows.Close()
	return count
}


//根据gw_id将当前时间写入到ap表中对应的online时间
func (p *ApDao)UpdateApLastedTimeByWanmac(wan_mac string)string  {
	ApSql := fmt.Sprintf(
		"UPDATE ap SET latestonline_time='%s' WHERE wan_mac='%s';",time.Now().Format("2006-01-02 15:04:05"),wan_mac)
	//fmt.Println("ap table sql:", ApSql)

	_, err:=framework.DB.Exec(ApSql)
	if err != nil{
		fmt.Println("update apLatestOlineTime  failed!",err)
		return "error"
	}
	return "ok"
}


func (p *ApDao)UpdateApState(ap model.ApDispList)string  {
	ApSql := fmt.Sprintf(
		"UPDATE ap SET wan_ip='%s',lan_mac='%s',lan_ip='%s',uptime=%d,state='%s' WHERE wan_mac='%s';","","","",0,"off",ap.Wan_mac)
	fmt.Println("ap table sql:", ApSql)

	_, err:=framework.DB.Exec(ApSql)
	if err != nil{
		fmt.Println("update ap State and clear some field  failed!",err)
		return "error"
	}
	//fmt.Println("success")
	log.Printf("ap is offline.name:%s,wan_ip:%s,lan_mac:%s,lan_ip:%s,address:%s,ens:%s;\n",ap.Name,ap.Wan_ip,ap.Lan_mac,ap.Lan_ip,ap.Address,ap.Ens)
	return "ok"
}


//从数据库中获取指定hostname的AP的hostkey
func (p *ApDao)GetApHostkeyDB(host,port string)string{

	var hostkey string
	err := framework.DB.QueryRow("SELECT hostkey FROM ap WHERE wan_ip=?", host).Scan(&hostkey)
	if err !=nil{
		fmt.Println("can not find aphostkey in table")
		return "nil"
	}
	hostkey ="[localhost]:"+port+" "+hostkey
	log.Println(hostkey)
	return hostkey
}

//从数据库中获取指定hostname的AP的hostkey
func (p *ApDao)GetForwardPort(host string)string{

	var forward_port string
	err := framework.DB.QueryRow("SELECT forward_port FROM ap WHERE wan_ip=?", host).Scan(&forward_port)
	if err !=nil{
		fmt.Println("can not find forward port in table")
		return "err"
	}
	log.Println(forward_port)
	return forward_port
}