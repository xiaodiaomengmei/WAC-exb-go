package model

//ap info need update word
type ApDbInfo struct {
    Id int `json:"id"`
    Wan_mac string `json:"wan_mac"`
    Wan_ip string `json:"wan_ip"`
    Lan_mac string `json:"lan_mac"`
    Lan_ip string `json:"lan_ip"`
    Uptime int `json:"uptime"`    //startup second 
    Cpu string `json:"cpu"`
    Memory string `json:"memory"`    
    State string `json:"state"`
}

//ap info need to display word
type ApDispList struct{
    Id int `json:"id"`
	Name string `json:"name"`
    Address string `json:"address"`
    Manufacture string `json:"manufacture"`
    Model string `json:"model"`
    Ens string `json:"ens"`
    State string `json:"state"`
    Wan_mac string `json:"wan_mac"`
    Wan_ip string `json:"wan_ip"`
    Lan_mac string `json:"lan_mac"`
    Lan_ip string `json:"lan_ip"`
    Uptime string `json:"uptime"`
    Latestonline_time string  `json:"latestonline_time"`
}
type ApDisp struct{
    PageTotal int    `json:"pagetotal"` 
    Ap []ApDispList `json:"list"`
}