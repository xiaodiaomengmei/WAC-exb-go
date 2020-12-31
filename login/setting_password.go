package login
import(
    "io"
    "log"
    //"fmt"
    "net/http"
    "html/template"
    "encoding/json"
    "wifidog-server/model"
    "wifidog-server/utils"
)

func Setting_pwd() {

	http.HandleFunc("/wifidog/modifyPwd", func(w http.ResponseWriter, r *http.Request) {
		log.Println("login handle process,method:", r.Method)

		if r.Method == "GET" {
			client_mac := r.FormValue("client_mac")
			log.Println("client_mac=", client_mac)
			t, _ := template.ParseFiles("www/modifyPwd.html")
			t.Execute(w, client_mac)
		} else {
			r.ParseForm()
			phoneNumber := r.PostFormValue("phoneNumber")
			verifyCode := r.PostFormValue("verifyCode")
			password := r.PostFormValue("password")
			if utils.RedisUtilget(phoneNumber) == "" {
				code := utils.StatusCode{2003}
				json_byte, _ := json.MarshalIndent(code, "", "")
				io.WriteString(w, string(json_byte))
			} else if utils.RedisUtilget(phoneNumber) == verifyCode {
				var user model.Useraccount
				user.PhoneNumber = phoneNumber
				user.Password = password
				result := userDao.ModifyUserByPhoneNumber(user)
				if result == "ok" {
					code := utils.StatusCode{2001}
					json_byte, _ := json.MarshalIndent(code, "", "")
					io.WriteString(w, string(json_byte))
				} else {
					code := utils.StatusCode{2004}
					json_byte, _ := json.MarshalIndent(code, "", "")
					io.WriteString(w, string(json_byte))
				}
			} else {
				code := utils.StatusCode{2002}
				json_byte, _ := json.MarshalIndent(code, "", "")
				io.WriteString(w, string(json_byte))
			}
		}
	})
}
