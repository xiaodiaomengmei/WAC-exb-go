package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
        "strconv"
	"math/rand"
)

func Hex2Dec(val string) int{
    n,err := strconv.ParseUint(val,16,32)
    if err!= nil {
        log.Println(err)
    }
    return int(n)
}


func GetSmsCode(phoneNumber string,Uid string) string {
	 smscode := RedisUtilget(phoneNumber)
	 if smscode == "" {
	 	if Uid == "" {
	 		i := 1
	 		for i <= 6 {
	 			myRand := rand.Intn(10) + 48
	 			newChar := string(byte(myRand))
	 			smscode = smscode + newChar
	 			i++
	 		}
	 	}else{
			w := md5.New()
			io.WriteString(w, Uid+phoneNumber)
			md5str := fmt.Sprintf("%x", w.Sum(nil))
                        md5int := Hex2Dec(md5str[len(md5str)-6:len(md5str)])
			smscode = strconv.Itoa(md5int%1000000)
		}
		RedisUtilset(phoneNumber, smscode)
	}
	return smscode
}

func GetToken(mac string) string {
	m := md5.New()
	io.WriteString(m, mac)
	arr := m.Sum(nil)
	return fmt.Sprintf("%x", arr)
}
