## 关于

GVUE = gin + vue

前端使用 vue + tailwindcss + laravel mix

后端使用 gin + gorm

.env 中的配置优先于 config.yml 文件中的配置，比如 .env 文件中的 `GVUE_APP_ADDR` 优先于 config.yml 文件中的 `app.addr`，`GVUE_` 变量作用于 go，`MIX_` 变量作用于 js

预览地址 [https://gvue.ideaclips.com/](https://gvue.ideaclips.com/) 

## Todo

- [x] 基于 jwt 的用户认证(注册、登录、忘记密码)
- [x] 发送邮件(使用阿里云邮件，每日200封免费)
- [ ] 接口限流
- [ ] 日志
- [ ] 数据库迁移
- [ ] ...

## 起步

```sh
$ git clone https://github.com/lukedever/gvue-scaffold
$ cd gvue-scaffold
$ cp .env.example .env                  #修改相关配置
$ yarn                                  #安装前端依赖
$ yarn prod                             #开发时npm run watch
$ make                                  #编译
$ bin/server -f config.yml              #或make run运行
```

打开浏览器 [http://localhost:3000](http://localhost:3000)

## 如何修改项目名

克隆到指定目录

```sh
$ git clone https://github.com/lukedever/gvue-scaffold yourpath
```

sed 批量替换

```sh
$ sed -i 's/lukedever\/gvue-scaffold/yourname\/yourpath/g' `grep -rl gvue-scaffold yourpath/`
```
