图片服务
============================
#调试
```
#默认
go run main.go -d /tmp
#使用gin
gin --appPort 12345 -p 12346
```
#编译
```shell
env GOOS=linux GOARCH=amd64 go build -o MediaServer
```
#包依赖
```
go get github.com/astaxie/beego/logs
go get github.com/dchest/captcha
go get github.com/dgrijalva/jwt-go
go get github.com/garyburd/redigo/redis
go get github.com/jinzhu/gorm
go get github.com/go-sql-driver/mysql
go get github.com/pquerna/ffjson/ffjson
go get github.com/wendal/errors
go get gopkg.in/gin-gonic/gin.v1
go get github.com/codegangsta/gin
```