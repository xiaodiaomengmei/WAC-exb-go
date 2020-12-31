package login
import(
    "log"
    "fmt"
    "time"
    "strings"
    "io/ioutil"
    "net/http"
    "wifidog-server/dao"
)

const roam_token="1234567"
const no_token="no"
var usernowtraceDao=new(dao.UserNowTraceDao)
var userhistorytraceDao=new(dao.UserHistoryTraceDao)

func Authen(){
    http.HandleFunc("/wifidog/auth", func(w http.ResponseWriter, r *http.Request) {
        stage := r.URL.Query().Get("stage")
        url_clientmac:=r.URL.Query().Get("mac")
        url_clientip:=r.URL.Query().Get("ip")
        url_token:=r.URL.Query().Get("token")
        url_token =strings.Replace(url_token," ","", -1)
        var db_token string
        usertrace := usernowtraceDao.GetUserNowTraceByClientMac(url_clientmac)
        if usertrace == nil{
            fmt.Println("can not find token in db")
            db_token = no_token
        }else{
            db_token = usertrace[0].Token
        }

        log.Println("auth", stage, r.URL.RawQuery)
        log.Println("db token=", db_token)
        log.Println("packet token=", url_token)
    
        if stage == "login" {
            if url_token==db_token||url_token==roam_token{
               fmt.Fprintf(w, "Auth: 1")
               log.Println("auth url,stage=login,auth=1")
            }else {
               fmt.Fprintf(w, "Auth: 0")
               log.Println("auth url,stage=login,auth=0")
            }
        } else if stage == "counters" {
            body, _ := ioutil.ReadAll(r.Body)
            r.Body.Close()
            log.Println(string(body))
            log.Println("auth url,stage=counter")
            fmt.Fprintf(w, "{\"resp\":[]}")
        } else if stage == "roam" {
            if db_token != no_token {
                fmt.Fprintf(w, "token=%s",roam_token)
                gw_id := r.URL.Query().Get("gw_id")
                usertrace[0].Client_ip = url_clientip
                apinfo := apDao.GetApDispInfoByLanMac(gw_id)
                if apinfo == nil {
                    log.Println("the ap is invalid!!!")
                }
                //update  now trace table
                usertrace[0].State = "在线"
                //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
                usertrace[0].Uptime = time.Now().Format("2006-01-02 15:04:05")
                usertrace[0].Ap_Name = apinfo[0].Name
                usertrace[0].Ap_WanMac = apinfo[0].Wan_mac
                usertrace[0].Ap_Address = apinfo[0].Address
                usertrace[0].Ap_Manufacture = apinfo[0].Manufacture
                usertrace[0].Ap_Model = apinfo[0].Model
                usertrace[0].Ap_Ens = apinfo[0].Ens
                ok:=usernowtraceDao.UpdateUserNowTraceByUsernameAndClientmac(usertrace[0])
                log.Println("update user now trace database. action=",ok)
                log.Println("auth url,stage=roam, client already register in database.")
            } else {
                fmt.Fprintf(w, "deny")
                log.Println("auth url,stage=roam,client is not find in database,deny.turn to login")
            }
        } else {
            fmt.Fprintf(w, "OK")
            log.Println("auth url,stage=not handle,return ok.")
        }
    })
}