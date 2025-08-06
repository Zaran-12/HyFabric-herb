document.getElementById("upload-btn").addEventListener("click", function () {
    const avatarInput = document.getElementById("avatar");
    const previewImage = document.getElementById("avatar-preview");

    if (avatarInput.files && avatarInput.files[0]) {
        const file = avatarInput.files[0];
        const reader = new FileReader();

        reader.onload = function (e) {
            previewImage.src = e.target.result; // 设置预览图片的 src
            previewImage.style.display = "block"; // 显示预览图片
        };

        reader.readAsDataURL(file); // 读取文件内容并触发 onload
    } else {
        alert("请先选择一张图片！");
    }
});
document.getElementById('loginForm').addEventListener('submit', function (event) {

    // 阻止表单默认的提交行为
    event.preventDefault();

    // 获取表单输入的值
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const repeatpwd = document.getElementById('repeatpwd').value;
    const selectedRole =document.getElementById('role').value;
    const avatar = document.getElementById('avatar').files[0];
    // 验证头像是否已选择
    if (!avatar) {
        alert("请上传头像！");
        return;
    }

    if (avatar) {
        const reader = new FileReader();
        reader.onload = function (e) {
            const preview = document.getElementById("avatarPreview");
            preview.src = e.target.result;
            preview.style.display = "block";
        };

        reader.readAsDataURL(avatar); // 读取文件内容
    }


    // 创建一个对象来存储注册信息
    // 创建 FormData 对象
    const formData = new FormData();
    formData.append("user", username);
    formData.append("pwd", password);
    formData.append("repeatpwd", repeatpwd);
    formData.append("role", selectedRole);
    formData.append("avatar", avatar); // 文件数据

    // 使用 fetch API 发送数据
    fetch("/user/register", {
        method: "POST",
        body: formData
    })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log("注册结果：", data);
            if (data.code === "10000") {
                alert("注册成功！");
                window.location.href = "/index.html";
            } else {
                alert(`注册失败：${data.msg}`);
            }
        })
        .catch(error => {
            console.error("注册时发生错误：", error);
            alert("注册失败，请稍后重试！");
        });
});

