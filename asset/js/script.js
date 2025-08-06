document.getElementById('loginForm').addEventListener('submit', function (event) {

    // 阻止表单默认的提交行为
    event.preventDefault();

    // 获取表单输入的值
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const selectedRole =document.getElementById('role').value;


    // 创建一个对象来存储登录信息
    const loginData = {
        name: username,
        pwd: password,
        role: selectedRole
    };

    // 使用fetch API调用登录接口
    fetch('/user/login', {
        method: 'POST', // 假设你的登录接口使用POST方法
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(loginData) // 将登录信息转换为JSON字符串并发送
    })
        .then(response => response.json()) // 解析响应为JSON
        .then(data => {
            // 根据接口返回的数据处理登录结果
            console.log(data);
            var code = data.code;
            if (data.code === "20000" && data.data.redirectURL) {
                console.log("即将跳转到:", data.data.redirectURL);
                // 将用户信息存储在 localStorage
                localStorage.setItem("username", data.data.name);
                localStorage.setItem("role", data.data.role);
                localStorage.setItem("avatar", data.data.avatar);
                window.location.href = data.data.redirectURL; // 跳转页面
            } else {
                console.error("登录失败或跳转路径不存在:", data.msg);
                alert(data.msg);
            }
        })
        .catch(error => {
            console.error('登录时发生错误:', error);
            // 在这里你可以处理错误情况，比如显示一个通用的错误消息给用户
        });
});

