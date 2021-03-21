#!/bin/bash

if ! [ -x "$(command -v git)" ]; then
    echo 'Error: git is not installed.' >&2
    exit 1
fi

if ! [ -x "$(command -v go)" ]; then
    echo 'Error: go is not installed.' >&2
    exit 1
fi

user_name=$(git config user.name)
if [[ "$user_name" = "" ]]; then
    echo -n "输入用户名: "
    read user_name
    if [[ "$user_name" = "" ]]; then
        echo 'Error: user name is required.' >&2
        exit 1
    fi
fi

echo -n "输入项目名: "
read app_name
if [[ "$app_name" = "" ]]; then
    echo 'Error: app name is required.' >&2
    exit 1
fi


# 1. clone code ...
git clone https://github.com/lukedever/api-scaffold.git $app_name
cd $app_name

echo replace moudle name ...
find . -type f -not -path '*/\.*' | xargs sed -i "s/lukedever\/api/${user_name}\/${app_name}/g"
find . -type f -not -path '*/\.*' | xargs sed -i "s/api./${app_name}./g"
find . -maxdepth 1 -name "*.go" | xargs sed -i "s/api/${app_name}/g"

# 3.rename path
mv cmd/api cmd/${app_name}
rm -rf .git
go mod tidy
cp .env.example .env

echo done