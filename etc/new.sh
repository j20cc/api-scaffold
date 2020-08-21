#!/bin/bash

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
sed -i "s/gvue-scaffold/$name/g" `grep -rl gvue-scaffold $name/`

#删除前端
if [ "$ui" = false ]; then
    echo "remove ui code..."
    rm $name/app/controllers/user.go
    rm $name/app/models/user.go
    sed -i '/mysqlCli.AutoMigrate(&User{})/d' $name/app/models/db.go
    sed -i '/\/\/auth-route-start/,/\/\/auth-route-end/d' $name/main.go

    rm -rf $name/resources/js/pages/auth
    sed -i '/<!-- auth-btn-start -->/,/<!-- auth-btn-end -->/d' $name/resources/js/components/MyHeader.vue
    sed -i '/\/\/auth-route-start/,/\/\/auth-route-end/p' $name/resources/js/libs/routes.js
fi

#删除git
rm -rf $name/.git

#新建.env
cd $name && cp .env.example .env

echo "create project success! run yarn to install node_modules"