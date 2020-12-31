package login
import(
    "log"
    "fmt"
    "net/http"
)
var portalPage = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>WiFi Portal</title>
    <meta name="viewport" content="width=device-width,minimum-scale=1.0,maximum-scale=1.0,user-scalable=no" />
</head>
<body>
    <h1>Welcome to WiFi Portal</h1>
</body>
</html>
`

func Portal(){
    
    http.HandleFunc("/wifidog/portal", func(w http.ResponseWriter, r *http.Request) {
        log.Println("portal", r.URL.RawQuery)
        fmt.Fprintf(w, portalPage)
    })

}