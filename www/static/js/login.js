function login() {
    //jQuery写法
    //获取用户输入
    var username = $("#username").val();
    var password = hex_sha1($("#password").val());
    var client_mac = $("#username")[0].name;
    if (!username) {
        alert("请输入用户名");
        $("#username").focus();//获取焦点
        return;
    }
    if (!password) {
        alert("请输入密码");
        $("#password").focus();//获取焦点
        return;
    }
    $.ajax({
        //几个参数需要注意一下
        type: "post",//提交方法
        url: "/wifidog/login",//地址
        data: {
            "username": username,
            "password": password,
            "client_mac": client_mac
        },
        success: function (data) {
            var json = eval('(' + data + ')')
            if (json.code == "2000") {
                window.location.href = json.uri
            } else {
                alert(json.message)
            }
        },
    });
}
