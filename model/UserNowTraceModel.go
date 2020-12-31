package model

type UserNowTrace struct{
    Id uint64 `json:"id"`
    User_Human string `json:"human_name"`
	User_Name string `json:"user_name"`
    State string `json:"state"`
    Uptime string `json:"uptime"`
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
type UserNowTraceShow struct{
    PageTotal int    `json:"pagetotal"` 
    UserNowTraceList []UserNowTrace `json:"list"`
}
