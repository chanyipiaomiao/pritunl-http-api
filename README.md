# Pritunl HTTP API

[Pritunl](https://pritunl.com/) 是一个基于OpenVPN协议实现的企业级VPN方案, 免费版本可以拥有一个服务器的创建, 付费版有更多的功能.

[安装配置文档](https://docs.pritunl.com/docs)

由于开源版不提供API调用，故有此项目, 通过模拟登陆的方式实现的一套HTTP API接口

## 安装运行

[下载](https://github.com/chanyipiaomiao/pritunl-http-api/releases)二进制文件即可

## 配置文件

```bash
# 配置文件中的${},首先从环境变量获取，获取不到则使用后边的默认值

appname = pritunl-http-api
autorender = false
copyrequestbody = true
EnableDocs = false

# 运行模式 dev | prod
runmode = "${PRITUNL_HTTP_API_RUNMODE||dev}"

# 监听地址
httpaddr = "${PRITUNL_HTTP_API_LISTEN_IP||127.0.0.1}"

# 监听端口
httpport = "${PRITUNL_HTTP_API_PORT||30080}"

[log]
# 日志路径
logPath = "${PRITUNL_HTTP_API_LOGPATH||logs/pritunl.log}"

[security]
# 是否启用token认证,启用会要求在请求头部加入token头

# true | false
enableToken = false

# 请求头中 token头名称
tokenName = "${PRITUNL_HTTP_API_TOKEN_NAME||PRITUNL-HTTP-API-TOKEN}"

# token 值
token = "${PRITUNL_HTTP_API_TOKEN||CyNINFvjdJTh4QTfqsPVJuNdRDUGvHnU}"


[common-vpn-wenba]
# common 是一个标识，代表的是那台服务器,可以随意定义, 比如有多台pritunl服务器,在调接口的时候会用到
# wenba 是pritunl里面一个organization的名称, pritunl里面可以包含多个organization, 在调接口的时候会用到
# common-vpn-wenba 三者结合起来指的就是要操作那台服务器的那个organization

# url是pritunl控制台地址
url = "${PRITUNL_COMMON_URL||https://x.x.x.x}"

# 登录的用户名和密码
username = "${PRITUNL_COMMON_USERNAME||wenba}"
password = "${PRITUNL_COMMON_PASSWORD||123456}"

# organization id 是对应上边wenba
organization = "${PRITUNL_COMMON_WENBA_ORG||123456}"

# debug 用于调试
# true | false
debug = "${PRITUNL_COMMON_DEBUG||false}"
```

### 如何获取 organization 的 名称 和 id？

先登录到pritunl web 控制台

chrome可以通过按F12或者右键检查菜单打开开发者工具,监控Network

切换到 Users 选项卡, 找到该接口的响应

```bash
/organization?page=0
```

里面有组织的名称和ID，在调用接口的时候会需要用到这2个值


## 接口

### 添加用户

```bash
POST /pritunl?vpn-name=VPN名称&org-name=组织名称&username=用户名

vpn-name 就是配置文件里面配置的名称, 如: common
org-name organization组织的名称, 如: wenba
username 用户名
email 可选

通过 vpn-name 和 org-name 在配置文件中找到对应的配置

```

返回

```bash
{
    "data": {
        "profileLink": "https://x.x.x.x/k/u3yUsEcL",
        "userId": "5d008e164ed924023806d388",
        "username": "test7777"
    },
    "entryType": "Pritunl OpenVPN Operation",
    "error": "",
    "statusCode": 0
}

profileLink 是下载配置文件的和配置2步验证的URL
```

### 删除用户

```bash
DELETE /pritunl?vpn-name=VPN名称&org-name=组织名称&username=用户名
```

返回

```bash
{
    "data": {
        "status": "deleted",
        "username": "test7777"
    },
    "entryType": "Pritunl OpenVPN Operation",
    "error": "",
    "statusCode": 0
}
```


### 禁用用户

```bash
PUT /pritunl?vpn-name=VPN名称&org-name=组织名称&username=用户名&status=disable
```

返回

```bash
{
    "data": {
        "status": "disabled",
        "username": "test7777"
    },
    "entryType": "Pritunl OpenVPN Operation",
    "error": "",
    "statusCode": 0
}
```

### 启用用户

```bash
PUT /pritunl?vpn-name=VPN名称&org-name=组织名称&username=用户名&status=enable
```

返回

```bash
{
    "data": {
        "status": "enabled",
        "username": "test7777"
    },
    "entryType": "Pritunl OpenVPN Operation",
    "error": "",
    "statusCode": 0
}
```

### 批量添加用户

```bash
POST /pritunl?vpn-name=VPN名称&org-name=组织名称&username=用户名&multi=yes
```

返回

```bash
{
    "data": [
        {
            "profileLink": "https://x.x.x.x/k/aP25mrUn",
            "userId": "5cff8b174ed92402380625f7",
            "username": "test5555"
        },
        {
            "profileLink": "https://x.x.x.x/k/x47p83KG",
            "userId": "5cff8b974ed924023806264f",
            "username": "test6666"
        }
    ],
    "entryType": "Pritunl OpenVPN Operation",
    "error": "",
    "statusCode": 0
}
```


