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

/**
 *输入验证码后提交到后台是否输入合法
 */
function dologin(flag) {
    var fverifyCode = $("#fverifyCode");
    var phoneNumber = $("#phoneNumber");
    if ($.trim(fverifyCode.val()) == "") {
        alert("请输入验证码");
        fverifyCode.val("");
        fverifyCode.focus();
    } else {
        $.ajax({
            type: "POST",
            url: "/wifidog/loginByPhone",
            data: {
                "phoneNumber": phoneNumber.val(),
                "verifyCode": fverifyCode.val(),
                "client_mac": $("#phoneNumber")[0].name,
                "flag": flag
            },
            success: function (data) {
                var json = eval('(' + data + ')')
                if (json.code == "2000") {
                    window.location.href = json.uri
                } else {
                    alert(json.message)
                }
            }
        })
    }
}
