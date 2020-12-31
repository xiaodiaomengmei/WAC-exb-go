package dao

import (
	"fmt"
	"log"
	"wifidog-server/framework"
	"wifidog-server/model"
)

type SsidDao struct {

}

//向数据库中增加ssid
func (s *SsidDao)AddSsid(ssid model.Ssid)string  {

	//判断是否在数据库中存在该ssid
	//fmt.Println(ssid.Ssid)
	rows, err := framework.DB.Query("SELECT * FROM ssid WHERE ssid=?;",ssid.Ssid)
	if err !=nil{
		log.Println("commend of checking ssid in dataabses failed")
		return "error"
	}
	//fmt.Println("*******************************")
	if rows.Next() != false {
		err := rows.Scan(&ssid.Id,&ssid.Ssid, &ssid.Jurisdiction)
		if err != nil {
			fmt.Println("get ssid error")
		}else{
			fmt.Printf("id:%d,Ssid:%s,Jurisdiction:%s",ssid.Id, ssid.Ssid, ssid.Jurisdiction)
		}
		log.Println("existence in the database")
		return "1"
	}else{
		//在不重复的情况下，向数据库中添加该ssid及模式
		_, err = framework.DB.Exec("INSERT INTO ssid (ssid, jurisdiction) VALUES (?, ?);",ssid.Ssid, ssid.Jurisdiction)
		if err != nil{
			log.Println("add ssid failed", err)
			return "error"
		}else{
			log.Printf("add ssid: %s ,complement\n",ssid.Ssid)
			return "ok"
		}
	}
}

//修改数据库中的ssid信息
func (s *SsidDao)ModifySsid(ssid model.Ssid)string  {
	//1.判断修改的是ssid还是模式
	//fmt.Println(ssid.Id)
	rows, err := framework.DB.Query("SELECT * FROM ssid WHERE id=?;",ssid.Id)
	if err !=nil{
		log.Println("commend of checking ssid in dataabses failed")
		return "error"
	}
	tag := 0    //标志位
	rows.Next()
	var ssid_database model.Ssid
	err = rows.Scan(&ssid_database.Id, &ssid_database.Ssid, &ssid_database.Jurisdiction)
	if err != nil {
		fmt.Println("get ssid error")
	}else{
		fmt.Printf("Id:%d,Ssid:%s,Jurisdiction:%s", ssid_database.Id, ssid_database.Ssid, ssid_database.Jurisdiction)
	}
	if ssid_database.Ssid == ssid.Ssid{
		tag =1
	}

	fmt.Printf("tag =%d\n",tag)
	//2.如果修改的是ssid
	if tag == 0 {
		//fmt.Println(ssid.Ssid)
		rows, err := framework.DB.Query("SELECT * FROM ssid WHERE ssid=?;", ssid.Ssid)
		if err != nil {
			log.Println("commend of checking ssid in dataabses failed")
			return "error"
		}
		fmt.Println("*******************************")
		if rows.Next() != false {
			err := rows.Scan(&ssid.Id, &ssid.Ssid, &ssid.Jurisdiction)
			if err != nil {
				fmt.Println("get ssid error")
			} else {
				fmt.Printf("id:%d,Ssid:%s,Jurisdiction:%s", ssid.Id, ssid.Ssid, ssid.Jurisdiction)
				log.Println("existence in the database")
			}
			return "1"
		} else {
			_, err := framework.DB.Exec("UPDATE ssid SET ssid=?,jurisdiction=? WHERE id=?;", ssid.Ssid, ssid.Jurisdiction, ssid.Id)
			//fmt.Println(ssid.Ssid)
			if err != nil {
				log.Println("modify ssid failed", err)
				return "error"
			} else {
				log.Printf("modify ssid: %s ,success\n", ssid.Ssid)
				return "ok"
			}
		}
	}else{
		//2.如果没有修改的是ssid
		_, err := framework.DB.Exec("UPDATE ssid SET ssid=?,jurisdiction=? WHERE id=?;", ssid.Ssid, ssid.Jurisdiction, ssid.Id)
		//fmt.Println(ssid.Ssid)
		if err != nil {
			log.Println("modify ssid failed", err)
			return "error"
		} else {
			log.Printf("modify ssid: %s ,success\n", ssid.Ssid)
			return "ok"
		}
	}

}

//删除数据库中的ssid
func (s *SsidDao)DeleteSsid(ssid_id string)string  {
	//fmt.Println(id)
	_, err := framework.DB.Exec("DELETE FROM ssid WHERE ssid=?;", ssid_id)
	if err != nil{
		log.Println("delete ap failed", err)
		return "error"
	}else{
		log.Printf("delete ssid: %s ,complement\n",ssid_id)
		return "ok"
	}
}

//查询所有的ssid信息
func (s *SsidDao)GetAllSsid()[]model.Ssid{
	SsidSql := fmt.Sprintf("SELECT * FROM ssid;")
	//fmt.Println("ap table sql:", ApSql)

	rows,err:=framework.DB.Query(SsidSql)
	if err !=nil{
		fmt.Println("can not find ssid in table")
		return nil
	}else{
		var ssids []model.Ssid
		for rows.Next(){
			var ssid model.Ssid
			err:=rows.Scan(&ssid.Id,&ssid.Ssid,&ssid.Jurisdiction)
			//fmt.Printf("Id:%d,Ssid:%s,Jurisdiction:%s",ssid.Id,ssid.Ssid,ssid.Jurisdiction)
			if err !=nil{
				fmt.Println("get ssid error")
				continue
			}
			//fmt.Printf("Id:%s,Ssid:%s,Jurisdiction:%s",ssid.Id,ssid.Ssid,ssid.Jurisdiction)
			ssids=append(ssids,ssid)
		}
		rows.Close()
		log.Println("get all ssid complement")
		return ssids
	}
}

