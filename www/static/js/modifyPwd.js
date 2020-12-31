var wait = 60;
$(".get-captcha_button")[0].onclick = function () {
    var phoneNumber = $("#phoneNumber");
    if ($.trim(phoneNumber.val()) == "") {
        phoneNumber.focus();
        alert("请输入手机号");
        return
    }
    if (this.innerHTML == "获取") {
        getVerifyCode(phoneNumber.val())
    }
    time(this);
}

function time(o) {
    if (wait == 0) {
        o.removeAttribute("disabled");
        o.innerHTML = "获取";
        wait = 60;
    } else {
        o.setAttribute("disabled", true);
        o.innerHTML = wait + "秒";
        wait--;
        setTimeout(function () {
            time(o)
        }, 1000)
    }
}

/**
 *点击获取验证码按钮后将手机号传到后台获取验证码
 */
function getVerifyCode(phoneNumber) {
    $.ajax({
        url: "/wifidog/sendSmsValidate",
        type: "POST",
        data: {"phoneNumber": phoneNumber},
    });
}

function modifyPwd() {
    var phoneNumber = $("#phoneNumber").val();
    var verifyCode = $("#verifyCode").val();
    var password = hex_sha1($("#password").val());
    var password2 = hex_sha1($("#password2").val());
    if (!$.trim(phoneNumber)) {
        alert("请输入手机号！");
        $("#phoneNumber").focus();//获取焦点
        return;
    }
    if (!$.trim(verifyCode)) {
        alert("请输入验证码！");
        $("#verifyCode").focus();//获取焦点
        return;
    }
    if (!$.trim(password)) {
        alert("请输入新密码！");
        $("#password").focus();//获取焦点
        return;
    }
    if (!$.trim(password2)) {
        alert("请再次输入新密码！");
        $("#password2").focus();//获取焦点
        return;
    }
    if (password != password2) {
        alert("两次输入的密码不一致！");
        $("#password").focus();//获取焦点
        return;
    }
    $.ajax({
        //几个参数需要注意一下
        type: "post",//提交方法
        url: "/wifidog/modifyPwd",//地址
        data: {
            "phoneNumber": phoneNumber,
            "verifyCode": verifyCode,
            "password": password
        },
        dataType: "json",//预期的服务器响应的类型
        success: function (data) {
            //返回数据为data
            if (data["code"] == "2001") {
                alert("修改成功");
                window.location.href = "/wifidog/login";
            } else if (data["code"] === "2004") {
                alert("修改失败");
            } else {
                alert("验证码错误");
            }
        }
    });
}