package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"pritunl-http-api/custom/ctype"
	"pritunl-http-api/models"
)

const (
	pritunlEntry = "Pritunl OpenVPN Operation"
)

var (
	token          = beego.AppConfig.String("security::token")
	tokenName      = beego.AppConfig.String("security::tokenName")
	enableToken, _ = beego.AppConfig.Bool("security::enableToken")
)

// Pritunl 控制器
type PritunlController struct {
	BaseController
}

func (p *PritunlController) Prepare() {
	var (
		vpnName       string
		orgName       string
		tokenInHeader string
		tokenInGet    string
	)

	if enableToken {
		// 获取 头部信息
		tokenInHeader = p.Ctx.Input.Header(tokenName)
		if tokenInHeader == "" {
			p.JsonError("token auth", "token auth error", "", true)
			p.Abort("403")
		}else if tokenInHeader != token {
			p.JsonError("token auth", "invalid token", "", true)
			p.Abort("403")
		}
	}


	vpnName = p.GetString("vpn-name")
	orgName = p.GetString("org-name")

	if vpnName == "" || orgName == "" {
		p.JsonError(pritunlEntry, "need vpn-name/org-name fields", "", true)
		p.StopRun()
	}

}

// 搜索用户
func (p *PritunlController) Get() {
	var (
		err           error
		pritunlClient *models.Pritunl
		searchResult  *models.PritunlSearchUserRespBody
		vpnName       string
		orgName       string
		username      string
	)

	vpnName = p.GetString("vpn-name")
	orgName = p.GetString("org-name")
	username = p.GetString("username")
	if username == "" {
		p.JsonError(pritunlEntry, "need username fields", "", true)
		return
	}

	pritunlClient, err = models.NewPritunlClient(vpnName, orgName)
	if err != nil {
		p.JsonError(pritunlEntry, fmt.Sprintf("pritunl  init client operation error: %s", err), "", true)
		return
	}

	if searchResult, err = pritunlClient.SearchUser(username); err != nil {
		p.JsonError(pritunlEntry, fmt.Sprintf("pritunl search user error: %s", err), "", true)
		return
	}

	p.JsonOK(pritunlEntry, searchResult, true)
}

// 添加用户
func (p *PritunlController) Post() {
	var (
		err           error
		pritunlClient *models.Pritunl
		vpnName       string
		orgName       string
		username      string
		multi         string
		email         string
		data          ctype.Data
		multiData     []ctype.Data
	)
	vpnName = p.GetString("vpn-name")
	orgName = p.GetString("org-name")
	multi = p.GetString("multi")

	switch multi {
	case "":
		username = p.GetString("username")
		email = p.GetString("email")
		if username == "" {
			p.JsonError(pritunlEntry, "need username fields", "", true)
			return
		}

		pritunlClient, err = models.NewPritunlClient(vpnName, orgName)
		if err != nil {
			p.JsonError(pritunlEntry, fmt.Sprintf("pritunl init client operation error: %s", err), "", true)
			return
		}

		if data, err = pritunlClient.CreateUser(username, email); err != nil {
			p.JsonError(pritunlEntry, fmt.Sprintf("pritunl create user error: %s", err), "", true)
			return
		}

		p.JsonOK(pritunlEntry, data, true)
		return
	case "yes":
		pritunlClient, err = models.NewPritunlClient(vpnName, orgName)
		if err != nil {
			p.JsonError(pritunlEntry, fmt.Sprintf("pritunl init client operation error: %s", err), "", true)
			return
		}

		if multiData, err = pritunlClient.MultiCreateUser(p.Ctx.Input.RequestBody); err != nil {
			p.JsonError(pritunlEntry, fmt.Sprintf("pritunl multi create user error: %s", err), "", true)
			return
		}
		p.JsonOK(pritunlEntry, multiData, true)
	default:
		p.JsonError(pritunlEntry, "pritunl create user error: multi fields can be yes ", "", true)
	}

}

// 删除用户
func (p *PritunlController) Delete() {
	var (
		err           error
		pritunlClient *models.Pritunl
		vpnName       string
		orgName       string
		username      string
		data          ctype.Data
	)

	vpnName = p.GetString("vpn-name")
	orgName = p.GetString("org-name")
	username = p.GetString("username")
	if username == "" {
		p.JsonError(pritunlEntry, "need username fields", "", true)
		return
	}

	pritunlClient, err = models.NewPritunlClient(vpnName, orgName)
	if err != nil {
		p.JsonError(pritunlEntry, fmt.Sprintf("pritunl init client operation error: %s", err), "", true)
		return
	}

	if data, err = pritunlClient.DeleteUser(username); err != nil {
		p.JsonError(pritunlEntry, fmt.Sprintf("pritunl delete user operation error: %s", err), "", true)
		return
	}

	p.JsonOK(pritunlEntry, data, true)
}

// 禁用用户
func (p *PritunlController) Put() {
	var (
		err           error
		pritunlClient *models.Pritunl
		vpnName       string
		orgName       string
		username      string
		status        string
		data          ctype.Data
	)

	vpnName = p.GetString("vpn-name")
	orgName = p.GetString("org-name")
	username = p.GetString("username")
	status = p.GetString("status")

	if username == "" || status == "" {
		p.JsonError(pritunlEntry, "need username,status fields", "", true)
		return
	}

	pritunlClient, err = models.NewPritunlClient(vpnName, orgName)
	if err != nil {
		p.JsonError(pritunlEntry, fmt.Sprintf("pritunl init client operation error: %s", err), "", true)
		return
	}

	if data, err = pritunlClient.EnableDisableUser(username, status); err != nil {
		p.JsonError(pritunlEntry, fmt.Sprintf("pritunl enable/disable user operation error: %s", err), "", true)
		return
	}

	p.JsonOK(pritunlEntry, data, true)
}
