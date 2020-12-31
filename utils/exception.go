package utils

import (
    "log"
	//"encoding/json"
)


func checkErr(err error){
    if err !=nil{
        log.Println(err)
    }
}
type ExceptionStruct struct {
    Try     func()
    Catch   func(Exception)
}
type Exception interface{}
func Throw(up Exception) {
    panic(up)
}
func (this ExceptionStruct) Do() {
    if this.Catch != nil {
        defer func() {
            if e := recover(); e != nil {
                this.Catch(e)
            }
        }()
    }
    this.Try()
}

type StatusCode struct{
    Code int `json:"code,string"`

}