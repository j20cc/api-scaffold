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

#克隆项目
if git clone https://github.com/lukedever/gvue-scaffold $name; then
    echo "clone project success"
else
    echo "clone project failed"
    exit 1
fi

#替换代码
echo "replace code ..."
sed -i "s/lukedever\/gvue-scaffold/$user\/$name/g" `grep -rl gvue-scaffold $name/`
sed -i "s/gvue-scaffold/$name/g" `grep -rl gvue-scaffold $name/`

#删除git
rm -rf $name/.git

echo "create project success!"