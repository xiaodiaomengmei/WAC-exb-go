package model

type EmailCfg struct{
    Id uint64 `json:"id"`
    SendAddr string `json:"send_address"`
	SendPwd string `json:"send_password"`
    SmtpHost string `json:"smtp_host"`
    SmtpPort int `json:"smtp_port"`
    RcvAddr string `json:"receiver_address"`
    CCAddr string `json:"cc_address"`
    Subject string `json:"subject"`
    Body string `json:"body"`
    DaySend string `json:"day_send"`
    WeekSend string `json:"week_send"`
    MonthSend string `json:"month_send"`
}
