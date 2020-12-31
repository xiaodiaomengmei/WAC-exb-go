package framework

import (
	"log"
	"fmt"
	"os"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	_"github.com/go-sql-driver/mysql"
	"strings"
)

type st_db_config struct{
	Db_type   string `json:"db_type"`
    Db_name   string `json:"db_name"`
    User      string `json:"user"`
    Password  string `json:"password"`
    Port      string `json:"port"`
    Host      string `json:"host"`
}

var DB *sql.DB
func InitDB(dbcfg_file string) {
	
	db_file,err:=os.Open(dbcfg_file)
	if err != nil{
        fmt.Println("open file failed:",err.Error())
        return
	}
	db_byte,_:=ioutil.ReadAll(db_file)
	defer db_file.Close()

	var db_cfg st_db_config
	err2 := json.Unmarshal(db_byte, &db_cfg);
	if  err2 != nil {
        fmt.Println(err2)
    }
	fmt.Println("database:",db_cfg)

	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=uft8"
	//注意：要想解析time.Time类型，必须要设置parseTime=True
	path := strings.Join([]string{db_cfg.User, ":", db_cfg.Password, "@tcp(", db_cfg.Host, ":", db_cfg.Port, ")/", db_cfg.Db_name, "?charset=utf8&parseTime=True&loc=Local"}, "")
	//打开数据库，前者是驱动名，所以要导入:_"github.com/go-sql-driver/mysql"
	DB, _ = sql.Open(db_cfg.Db_type, path)
	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		log.Panic(err)
	}
	log.Println("database connect success.")
}
