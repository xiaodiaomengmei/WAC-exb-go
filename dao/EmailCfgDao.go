package dao

import (
	"fmt"
	"log"
	"database/sql"
	"wifidog-server/model"
	"wifidog-server/framework"
	)

type EmailCfgDao struct {

}


func (p *EmailCfgDao)GetEmailCfg()[]model.EmailCfg  {
	ApSql := fmt.Sprintf("SELECT id,send_address,send_password,smtp_host,smtp_port,receiver_address,cc_address,subject,body,day_send,week_send,month_send FROM email_cfg")
	//fmt.Println("ap table sql:", ApSql)

	rows,err:=framework.DB.Query(ApSql)
	if err !=nil{
		fmt.Println("can not find email configuration in table")
		return nil
	}
	var EmailCfgs []model.EmailCfg
	for rows.Next(){
		var email_cfg model.EmailCfg
		var CCAddr sql.NullString 
		var Subject sql.NullString 
		var Body sql.NullString 
		var DaySend sql.NullString 
		var WeekSend sql.NullString 
		var MonthSend sql.NullString 
		
		err:=rows.Scan(&email_cfg.Id,&email_cfg.SendAddr,&email_cfg.SendPwd,&email_cfg.SmtpHost,&email_cfg.SmtpPort,
			&email_cfg.RcvAddr,&CCAddr,
			&Subject,&Body,
			&DaySend,&WeekSend,&MonthSend)
		
		email_cfg.CCAddr=CCAddr.String
		email_cfg.Subject=Subject.String
		email_cfg.Body=Body.String
		email_cfg.DaySend=DaySend.String
		email_cfg.WeekSend=WeekSend.String
		email_cfg.MonthSend=MonthSend.String
		//fmt.Println(ap)

		if err !=nil{
			fmt.Println("get email configuration error")
			continue
		}
		EmailCfgs=append(EmailCfgs,email_cfg)
	}
	rows.Close()
	return EmailCfgs
}


func (p *EmailCfgDao)SetEmailCfg(email_cfg model.EmailCfg)string  {
	SqlStr := fmt.Sprintf(
		"UPDATE email_cfg SET send_address='%s',send_password='%s',smtp_host='%s',smtp_port='%d',receiver_address='%s',cc_address='%s',subject='%s',body='%s',day_send='%s',week_send='%s',month_send='%s';",
		email_cfg.SendAddr,email_cfg.SendPwd,email_cfg.SmtpHost,email_cfg.SmtpPort,
		email_cfg.RcvAddr,email_cfg.CCAddr,email_cfg.Subject,email_cfg.Body,
		email_cfg.DaySend,email_cfg.WeekSend,email_cfg.MonthSend)
	fmt.Println("update email cfg table sql=",SqlStr)

	_, err := framework.DB.Exec(SqlStr)
	if err != nil{
		log.Println("modify email cfg table failed", err)
		return "error"
	}
	return "ok"
}