<h1>GOAPP</h1>

<div>
 基于goapp+xterm实现webssh-网页上的SSH终端 <br/>
</div>
<br/>

### 特性

- 在网页上实现一个SSH终端。从而无需Xshell之类的模拟终端工具进行SSH连接，方便管理你的服务器。

- 可以对交互命令进行审计、记录

- 在页面上按一个键，就能打开一个webssh，并且自动登录 ，方便地管理各个服务器

- 可以集成到自有的后台管理体系中

## 下载并运行

### 获取代码

```
go get -v github.com/it234/gowebssh
```

### 运行

- 可以直接下载体验，下载地址: https://pan.baidu.com/s/1cgMF0rXf5hlx0DF3N7nVUw 提取码: have ，下载后直接运行gowebssh.exe，然后添加你的服务器信息即可。

- 运行服务端：cd cmd/manageweb，go run main.go，运行成功后打开 127.0.0.1:8080，如果是在windows下操作，需要提前安装并配置好mingw（sqlite的操作库用到），安装方式请自行百度/谷歌。
- 配置文件在(`cmd/manageweb/config.yaml`)中，用户默认为：admin/123456


## 未完待续部分

- 证书登录 
- 前端优化
- 端面客户端
- SSH命令审核
- 其他

## Donate

- If you find this project useful, you can buy author a glass of juice 
- alipay
- <img src="./img/alipay.jpg" width="256" height="256" />
- wechat
- <img src="./img/wxpay.jpg" width="256" height="256" />
- [Buy me a coffee](https://www.buymeacoffee.com/it234)
- bitcoin address : 1LwTcCZ1p5kq8UokZGUBVy3BL1wRa3q5Wn
- eth address : 0x68ca43651529D12996183d09a052a654F845cB89
- eos address : 123451234534


## MIT License

    Copyright (c) 2019 it234

## 与作者对话

> 微信号：it23456789，微信二维码：

<img src="./wechat.jpeg" width="256" height="256" />




