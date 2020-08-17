#!/bin/bash

#默认git用户名
user=$(git config user.name)

#输入用户名
read -r -p "please enter your name [default $user]: " input_user
if [ "$input_user" != "" ]; then
    user=$input_user
fi
if [ "$user" == "" ]; then
    echo "user name is required"
    exit 1
fi

#输入项目名
read -p "please enter your project name: " name
if [ "$name" == "" ]; then
    echo "project name is required"
    exit 1
fi

#是否需要ui
read -p "do you need auth ui y/n? [default y]: " ui
if [ "$ui" == "n" ]; then
    ui=false
fi

#克隆项目
if git clone https://github.com/lukedever/gvue-scaffold $name; then
    echo "clone project success"
else
    echo "clone project failed"
    exit 1
fi

#替换代码
echo "replace code..."
sed -i "s/lukedever\/gvue-scaffold/$user\/$name/g" `grep -rl gvue-scaffold $name/`
sed -i "s/gvue-scaffold/$name/g" `grep -rl gvue-scaffold $name/`

#删除前端
if [ "$ui" = false ]; then
    echo "remove ui code..."
    # sed -i '/controllers"$/d' $name/main.go
    sed -i '/.*userController.*/d' $name/main.go
    sed -i '/.*auth.*/d' $name/main.go
    sed -i '/&User{}/d' $name/app/models/db.go
    rm $name/app/controllers/user.go $name/app/models/user.go
    rm -rf $name/resources/js/pages/auth
    sed -i '/.*pages\/auth.*/d' $name/resources/js/libs/routes.js
fi

#删除git
rm -rf $name/.git

#新建.env
cd $name && cp .env.example .env

echo "create project success! run yarn to install node_modules"