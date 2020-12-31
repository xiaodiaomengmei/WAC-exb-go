package utils

import (
    "fmt"
	"log"
	//"time"
	//"strconv"
	//"os"
	"gopkg.in/gomail.v2"
	"wifidog-server/dao"
	"wifidog-server/model"
)
var email_cfgDao=new(dao.EmailCfgDao)
func SendEmail(emailcfg model.EmailCfg, Subject string, AttFileName string){

    m := gomail.NewMessage()
	m.SetHeader("From",  m.FormatAddress(emailcfg.SendAddr, "拟态安全WiFi考勤系统"))
	m.SetHeader("To", emailcfg.RcvAddr)          //发送给多个用户
	m.SetHeader("Cc", emailcfg.CCAddr)           //抄送
	m.SetHeader("Subject", Subject)     //设置邮件主题
	bodyStr := fmt.Sprintf("请您查收考勤数据")
	m.SetBody("text/html", bodyStr)    //设置邮件正文
	file_name := fmt.Sprintf(AttFileName)
	m.Attach(file_name)

    d := gomail.NewDialer(emailcfg.SmtpHost, emailcfg.SmtpPort, emailcfg.SendAddr, emailcfg.SendPwd)
	err := d.DialAndSend(m)
    if err != nil {
        log.Println(err)
        fmt.Println("send fail")
    }
}

func SendMailForDay(Subject string, AttFileName string) {
	emailcfgs := email_cfgDao.GetEmailCfg()
	emailcfg := emailcfgs[0]
	if emailcfg.DaySend != "YES"{
		fmt.Println("the email config day send not set, not need send for day")
		return
	}
	SendEmail(emailcfg, Subject, AttFileName)
} 
func SendMailForWeek(Subject string, AttFileName string) {
	emailcfgs := email_cfgDao.GetEmailCfg()
	emailcfg := emailcfgs[0]
	if emailcfg.WeekSend != "YES"{
		fmt.Println("the email config week send not set, not need send for week")
		return
	}
	SendEmail(emailcfg, Subject, AttFileName)
}
func SendMailForMonth(Subject string, AttFileName string) {
	emailcfgs := email_cfgDao.GetEmailCfg()
	emailcfg := emailcfgs[0]
	if emailcfg.MonthSend != "YES"{
		fmt.Println("the email config month send not set, not need send for month")
		return
	}
	SendEmail(emailcfg, Subject, AttFileName)
}