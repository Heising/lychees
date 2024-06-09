<p align="center">
  <a href="https://redli.cn" target="_blank" rel="noopener noreferrer">
    <img width="360" src="lychees-web\public\lychees.svg" alt="Pinia logo">
  </a>
</p>

## 这是一个支持导入第三方 svg 的书签页全栈项目

<img src="assets\wallpaper-1717944582372.png" alt="wallpaper-1717944582372" style="width:50%;" /><img src="assets\wallpaper-1717945222743.png" alt="wallpaper-1717945222743"   style="width:50%;" />

<img src="assets\191938.png" alt="191938" />



通过 svg 的可以鼠标悬停切换logo中间的白色变为黑色 个性化你的标签页

<img src="assets\b9beae3986bf781abe93c5c78e6a3516.png" alt="b9beae3986bf781abe93c5c78e6a3516" />![屏幕截图 2024-06-09 223537](D:\root\lychees\assets\屏幕截图 2024-06-09 223537.png)



示例 svg  多色 svg 把需要控制变色的路径不要移除fill属性 通过css控制变色

```xml
<svg t="1717943897408" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="1872"
    width="200" height="200">
    <path
        d="M732.884976 369.62096a399.444957 399.444957 0 0 0 233.642975 75.093992v-168.319982c-16.511998 0-32.938996-1.707-49.066995-5.162999v132.436986a399.487957 399.487957 0 0 1-233.684974-75.049992V671.999928c0 171.775982-138.751985 310.996967-309.887967 310.996966a307.626967 307.626967 0 0 1-172.500982-52.607994A308.394967 308.394967 0 0 0 422.95501 1023.99989c171.135982 0 309.929967-139.220985 309.929966-311.039967V369.66396z m60.501994-169.727981a234.879975 234.879975 0 0 1-60.501994-137.300986V40.959996H686.379981a235.562975 235.562975 0 0 0 107.007989 158.932983zM309.632022 798.634914a141.994985 141.994985 0 0 1-28.927997-86.143991 141.994985 141.994985 0 0 1 184.74698-135.594985v-171.989981a311.509967 311.509967 0 0 0-49.023995-2.858v133.887985A141.994985 141.994985 0 0 0 231.68003 671.530928c0 55.551994 31.700997 103.679989 77.951992 127.103986z"
        fill="#FF004F" p-id="1873"></path>
    <path
        d="M683.775982 328.660965a399.529957 399.529957 0 0 0 233.684974 75.050992V271.274971a234.367975 234.367975 0 0 1-124.073986-71.381992A235.604975 235.604975 0 0 1 686.378981 40.959996H564.223994v671.999927a142.036985 142.036985 0 0 1-141.738984 141.823985 141.396985 141.396985 0 0 1-112.852988-56.149994 142.292985 142.292985 0 0 1-77.994992-127.102986 141.994985 141.994985 0 0 1 184.74698-135.594986V402.047957c-168.106982 3.499-303.359967 141.354985-303.359967 310.954966 0 84.649991 33.706996 161.364983 88.36299 217.428977a307.626967 307.626967 0 0 0 172.500982 52.607994c171.178982 0 309.887967-139.263985 309.887967-311.039966V328.703965z"
        p-id="1874"></path>
    <!-- 控制变色的路径不要给fill属性 通过css控制变色-->
    <path
        d="M917.460956 271.274971v-35.839996a232.959975 232.959975 0 0 1-124.073986-35.541996 234.111975 234.111975 0 0 0 124.073986 71.381992zM686.379981 40.959996a239.402974 239.402974 0 0 1-2.559999-19.327998V0H515.157v671.999928A141.994985 141.994985 0 0 1 373.420015 813.823913a140.756985 140.756985 0 0 1-63.786993-15.189999 141.396985 141.396985 0 0 0 112.852988 56.149994 142.036985 142.036985 0 0 0 141.738984-141.780985V40.959996h122.154987zM416.42701 402.047957v-38.100996a311.807967 311.807967 0 0 0-42.495995-2.902C202.752033 361.044961 64.000048 500.266946 64.000048 672.042928a310.996967 310.996967 0 0 0 137.386985 258.388972 310.527967 310.527967 0 0 1-88.31999-217.429977c0-169.556982 135.252985-307.454967 303.359967-310.953966z"
        fill="#00F2EA" p-id="1875"></path>
</svg>
```

### 导入 symbol 引用

应该说这才是未来的主流引用svg方式，[查看文章](https://www.iconfont.cn/help/detail?spm=a313x.help_detail.i1.d8d11a391.550b3a81IFTvPQ&helptype=code) 这种用法其实是做了一个svg的集合



# 使用哪个第三方呢？

#### 可以使用iconfont-阿里巴巴矢量图标库 https://www.iconfont.cn/

#### 搜索或者上传图标到阿里巴巴矢量图标库

#### 添加你需要的图标到项目

#### 拷贝项目下面生成的symbol导入链接：

```js
//at.alicdn.com/t/font_********.js
```





# 后端项目

## lychees-server 

需要安装Go语言编译 linux安装教程

```shell
wget https://go.dev/dl/go1.22.2.linux-amd64.tar.gz
#先移除旧的go安装包，再解压到/usr/local/go
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.2.linux-amd64.tar.gz
```

命令来安装 Go 语言。如果您已经安装了 Go 语言，请确保 $GOPATH 和 $GOROOT 环境变量已正确设置。

设置环境变量

```shell
export PATH=$PATH:/usr/local/go/bin
```

主要为[gin](https://gin-gonic.com/) + [gorm](https://gorm.io/) 请点击查阅相关文档

token存储需要 [protobuff](https://protobuf.dev/)  请到 proto 文件夹执行命令 `protoc --go_out=./ auth.proto`

邮箱用到 [templ](https://templ.guide/) HTML页面请到 utils 文件夹执行命令 `templ generate`



#### 编译程序，请在 linux 环境下编译

编译需要安装 cbrotli 如果不清楚 cbrotli 是什么，请到utils /website_info.go 注释掉 br 解码相关代码

lychees-server文件夹执行命令go build`编译



#### Prometheus

只统计qps，不需要 Prometheus 请到 middlewares 文件夹下 prometheus.go 注释代码



#### 推荐使用代理 nginx 或者 openresty 

示例

```nginx
server {
    # IP限流
    limit_req_zone $binary_remote_addr zone=ip_limit_api:10m rate=5r/s;
    limit_req zone=ip_limit_api burst=100 nodelay;

    listen 443 ssl;
    http2 on; # 新增这行来启用HTTP/2
    server_name api.lychees.com;
    
    # 修改为你的证书目录
    ssl_certificate /root/ssl/nginx/lychees.com_bundle.crt;
    ssl_certificate_key /root/ssl/nginx/lychees.com.key;

    ssl_session_cache shared:SSL:1m;
    ssl_session_timeout 5m;

    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;


    location / {
        proxy_pass http://127.0.0.0:8081/;
        proxy_set_header Host $host:$server_port;
        proxy_set_header X-Forwarded-For $remote_addr; # HTTP的请求端真实的IP
        # proxy_set_header X-Forwarded-Proto $scheme; # 为了正确地识别实际用户发出的协议是 http 还是 https
    }
}
```





## 使用方式

1. 安装并启用 redis mongodb postgresql ，请自行查看官方文档安装
2. 修改 `config.yaml`，见 [config.yaml](config.yaml)
3. 修改 Nginx 设置文件
4. 启用程序 (systemd/nohup)

#### systemd

```
vim /etc/systemd/system/lychees-server.service
```

- 创建文件 `/etc/systemd/system/lychees-server.service` (可以自选名字)
- 请替换为正确的路径

```
[Unit]
Description=lychees-server
After=network.target

[Service]
Type=simple
Restart=always
WorkingDirectory=/root/lychees-server/
ExecStart=/root/lychees-server/lychees-server

[Install]
WantedBy=multi-user.target
```



- 刷新后台程序 `systemctl daemon-reload`
- 启用后台程序 `systemctl enable lychees-server.service`
- 禁用后台程序 `systemctl disable lychees-server.service`
- 启动后台程序 `systemctl start lychees-server.service`
- 停止后台程序 `systemctl stop lychees-server.service`
- 检查后台程序状态 `systemctl status lychees-server.service`

#### nohup

- 程序路径执行 `nohup ./lychees-server &`
- 停止程序 `kill -9 1234`，1234 替换程序为 PID

# 前端页面

## lychees-web-account 

[前往查看实例页面](https://account.redli.cn/)

注册账号页面，包含重置密码，修改邮箱

主要是 react + shadcn 请查看官方文档

```shell
# 安装依赖
pnpm i
# 开发模式运行
pnpm run dev
# 打包
pnpm run build
```



## lychees-web

[前往查看实例页面](https://redli.cn/)

书签主页面

主要是 vue + elementPlus 请查看官方文档

```shell
# 安装依赖
pnpm i
# 开发模式运行
pnpm run dev
# 打包
pnpm run build
```

