## 关于

GVUE = gin + vue

前端使用 vue + tailwindcss + laravel mix

后端使用 gin + gorm

.env 中的配置优先于 config.yml 文件中的配置，比如 .env 文件中的 `GVUE_APP_ADDR` 优先于 config.yml 文件中的 `app.addr`，`GVUE_` 变量作用于 go，`MIX_` 变量作用于 js

预览地址 [https://gvue.ideaclips.com/](https://gvue.ideaclips.com/) 

## Todo

- [x] 基于 jwt 的用户认证(注册、登录、忘记密码)
- [x] 发送邮件(使用阿里云邮件，每日200封免费)
- [x] 多语言表单验证
- [x] 日志
- [x] 数据库迁移
- [ ] 接口限流
- [ ] i18n
- [ ] ...

## 起步

```sh
$ git clone https://github.com/lukedever/gvue-scaffold
$ cd gvue-scaffold
$ cp .env.example .env                  #修改相关配置
$ yarn                                  #安装前端依赖
$ yarn prod                             #开发时npm run watch
$ make build                            #编译
$ bin/server run                        #运行
```

打开浏览器 [http://localhost:3000](http://localhost:3000)

## 数据库迁移

新建迁移

```sh
$ bin/server migrate create create_foos_table
```

执行/回滚 迁移

```sh
$ bin/server migrate up/down
```

## 使用脚本创建

```sh
$ wget https://raw.githubusercontent.com/lukedever/gvue-scaffold/master/etc/new.sh
$ bash new.sh
please enter your project name: project_name
do you need auth ui y/n? [default y]: y
```
