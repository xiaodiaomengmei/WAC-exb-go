package model

type UserHistoryTrace struct{
    Id uint64 `json:"id"`
    User_Name string `json:"user_name"`
    User_Human string `json:"user_humanname"`
    User_Phone string `json:"user_phnumber"`
    User_Email string `json:"user_email"`
    User_Department string `json:"user_department"`
    LoginDate string `json:"login_date"`
    Logintime string `json:"login_time"`
    Logouttime string `json:"logout_time"`
    Ap_Name string `json:"ap_name"`
    Ap_WanMac string `json:"ap_wanmac"`
    Ap_Address string `json:"ap_address"`
    Ap_Manufacture string `json:"ap_manufacture"`
    Ap_Model string `json:"ap_model"`
    Ap_Ens string `json:"ap_ens"`
    Client_mac string `json:"client_mac"`
    Client_ip string `json:"client_ip"`
    Token string `json:"token"`
}
type UserHistoryTraceShow struct{
    PageTotal int    `json:"pagetotal"` 
    UserHistoryTraceList []UserHistoryTrace `json:"list"`
}