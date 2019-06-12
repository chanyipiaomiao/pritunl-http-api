package controllers

import (
	"github.com/astaxie/beego"
	"pritunl-http-api/custom/common"
	"pritunl-http-api/custom/ctype"
)

var (
	token     = beego.AppConfig.String("security::token")
	tokenName = beego.AppConfig.String("security::tokenName")
)

type BaseController struct {
	beego.Controller
}

func (b *BaseController) log(msg ctype.Data) ctype.Data {

	if _, ok := msg["clientIP"]; !ok {
		msg["clientIP"] = b.Data["RemoteIP"]
	}

	return msg
}

func (b *BaseController) LogInfo(entryType string, msg ctype.Data) {
	message := b.log(msg)
	if _, ok := msg["statuscode"]; !ok {
		message["statuscode"] = 0
	}
	common.GetLogger().Info(message, entryType)
}

func (b *BaseController) LogError(entryType string, msg ctype.Data) {
	message := b.log(msg)
	if _, ok := msg["statuscode"]; !ok {
		message["statuscode"] = 1
	}
	common.GetLogger().Error(message, entryType)
}

func (b *BaseController) json(entryType, errmsg string, statuscode int, data interface{}, isLog bool) {
	msg := ctype.Data{
		"entryType":  entryType,
		"error":      errmsg,
		"statusCode": statuscode,
		"data":       data,
	}
	b.Data["json"] = msg
	b.ServeJSON()

	msg["clientIP"] = b.Data["RemoteIP"]

	if isLog {
		go func() {
			if statuscode == 1 {
				b.LogError(entryType, msg)
			} else {
				b.LogInfo(entryType, msg)
			}
		}()
	}
}

func (b *BaseController) JsonError(entryType, errmsg string, data interface{}, isLog bool) {
	b.json(entryType, errmsg, 1, data, isLog)
}

func (b *BaseController) JsonOK(entryType string, data interface{}, isLog bool) {
	b.json(entryType, "", 0, data, isLog)
}

func (b *BaseController) Prepare() {

	// 获取客户端IP
	b.Data["RemoteIP"] = b.Ctx.Input.IP()

	// 获取 头部信息
	tokenInHeader := b.Ctx.Input.Header(tokenName)
	if tokenInHeader != token {
		tokenInGet := b.GetString(tokenName)
		if tokenInGet == "" {
			b.JsonError("token auth", "token auth error", "", true)
			b.Abort("403")
		}
	}

}
