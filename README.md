# 媒体服务 V0.1
---
## 功能
- 图片
  - 上传,下载,压缩 `jpg,png,gif`
  - 封面图 `gif`
- 视频
  - 上传,下载,压缩,封面图 `支持主流媒体格式`
  
## 包&环境依赖
- golang1.8+
- [gin](https://github.com/gin-gonic/gin) `高性能,轻量级 golang 开发框架`
- [ffmpeg](https://ffmpeg.org/) `C视频处理库`
- [imageserver](https://github.com/GxlZ/imageserver) `golang图片处理库` 基于[pierrre/imageserver](https://github.com/pierrre/imageserver) fork 后进行部分改造

## 参数配置
- 启动参数
  - `-p 启动端口`
  - `-d <文件保存目录>` **注:`如果输入的目录不存在不会自动创建,并会提示错误`,`会在指定目录下自动生成videos,imgs,tmp文件夹`**
- 环境参数
  - 在系统环境变量中设置 `MEDIA_SERVER_ENV` 指定运行环境 
    - `dev` 开发环境
    - `pro 或 不设置` 生产环境   

## 运行
- 调试
  - 直接运行 `go run main.go -p 12345 -d /tmp/data`
  - 通过gin代理端口,监听文件自动编译 `gin --appPort 12345 -p 12346` 启动后使用代理端口访问程序 
    - `-p <代理端口>`
    - `--appPort <MediaServer启动端口>`
    - **注:此gin`是本地代理,用于实时编译代码`,非golang开发框架的gin。**
- 编译
  - `go build && ./MediaServer -d /tmp/data` 本地版本
  - `env GOOS=linux GOARCH=amd64 go build && ./MediaServer -d /tmp/data` linux amd64版本

## TODOS
- [ ] 接入CDN配置
- [ ] 支持GIF尺寸调整
- [ ] 使用CGO调用ffmpeg
- [ ] 支持参数自定义尺寸
