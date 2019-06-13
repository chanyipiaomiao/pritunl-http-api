package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/imroc/req"
	"pritunl-http-api/custom/ctype"
)

const (
	authSessionAPI         = "/auth/session"
	userAPI                = "/user"
	downloadProfileLinkAPI = "/key"
	getCsrfTokenAPI        = "/state"
)

var (
	pritunlUrl      string
	authUrl         string
	getCsrfTokenUrl string
	password        string
	username        string
	organization    string
)

type PritunlLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 创建用户请求body
type PritunlCreateUserReqBody struct {
	Organization    string   `json:"organization"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	AuthType        string   `json:"auth_type"`
	BypassSecondary bool     `json:"bypass_secondary"`
	ClientToClient  bool     `json:"client_to_client"`
	Groups          []string `json:"groups"`
	Pin             *string  `json:"pin"`
	NetworkLinks    []string `json:"network_links"`
	DnsServers      []string `json:"dns_servers"`
	DnsSuffix       string   `json:"dns_suffix"`
	PortForwarding  []string `json:"port_forwarding"`
}

// 批量添加用户时请求体
type PritunlMultiCreateUserReqBody []struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// 用户信息
type pritunlUserPayload struct {
	Organization     string   `json:"organization"`
	OrganizationName string   `json:"organization_name"`
	Name             string   `json:"name"`
	Email            string   `json:"email"`
	AuthType         string   `json:"auth_type"`
	BypassSecondary  bool     `json:"bypass_secondary"`
	ClientToClient   bool     `json:"client_to_client"`
	Groups           []string `json:"groups"`
	Pin              bool     `json:"pin"`
	NetworkLinks     []string `json:"network_links,omitempty"`
	DnsServers       []string `json:"dns_servers"`
	DnsSuffix        string   `json:"dns_suffix"`
	PortForwarding   []string `json:"port_forwarding"`
	Disabled         bool     `json:"disabled"`
	OtpSecret        string   `json:"otp_secret"`
	Type             string   `json:"type"`
	Id               string   `json:"id"`
}

// 创建用户响应body
type PritunlCreateUserRespBody []pritunlUserPayload

// 获取csrf_token接口响应body
type PritunlGetCSRFTokenRespBody struct {
	CsrfToken string  `json:"csrf_token"`
	SSO       *string `json:"sso,omitempty"`
	SuperUser bool    `json:"super_user"`
	Version   int64   `json:"version"`
	Plan      *string `json:"plan,omitempty"`
	Active    bool    `json:"active"`
	Theme     string  `json:"theme"`
}

// 获取profile下载链接响应body
type PritunlGetProfileLinkRespBody struct {
	ViewUrl   string `json:"view_url"`
	KeyUrl    string `json:"key_url"`
	UriUrl    string `json:"uri_url"`
	KeyZipUrl string `json:"key_zip_url"`
	KeyOncUrl string `json:"key_onc_url"`
	Id        string `json:"id"`
}

// 搜索用户返回body
type pritunlSearchUserRespUserBody struct {
	Hidden     string      `json:"hidden,omitempty"`
	AuthType   string      `json:"auth_type"`
	DNSServers interface{} `json:"dns_servers"`
	Pin        bool        `json:"pin"`
	DNSSuffix  interface{} `json:"dns_suffix"`
	Servers    []struct {
		Status         bool        `json:"status"`
		Platform       interface{} `json:"platform"`
		ServerID       string      `json:"server_id"`
		VirtAddress6   string      `json:"virt_address6"`
		VirtAddress    string      `json:"virt_address"`
		Name           string      `json:"name"`
		RealAddress    interface{} `json:"real_address"`
		ConnectedSince interface{} `json:"connected_since"`
		ID             string      `json:"id"`
		DeviceName     interface{} `json:"device_name"`
	} `json:"servers"`
	Disabled         bool          `json:"disabled"`
	NetworkLinks     []interface{} `json:"network_links"`
	PortForwarding   []interface{} `json:"port_forwarding"`
	ID               string        `json:"id"`
	OrganizationName string        `json:"organization_name"`
	Type             string        `json:"type"`
	Email            string        `json:"email"`
	Status           bool          `json:"status"`
	DNSMapping       interface{}   `json:"dns_mapping"`
	OtpSecret        string        `json:"otp_secret"`
	ClientToClient   bool          `json:"client_to_client"`
	Sso              interface{}   `json:"sso"`
	BypassSecondary  bool          `json:"bypass_secondary"`
	Groups           []interface{} `json:"groups"`
	Audit            bool          `json:"audit"`
	Name             string        `json:"name"`
	Gravatar         bool          `json:"gravatar"`
	OtpAuth          bool          `json:"otp_auth"`
	Organization     string        `json:"organization"`
	Haskey           bool          `json:"has_key,omitempty"`
}

// 搜索用户返回body
type PritunlSearchUserRespBody struct {
	Search      string                          `json:"search"`
	Users       []pritunlSearchUserRespUserBody `json:"users"`
	SearchMore  bool                            `json:"search_more"`
	SearchLimit int                             `json:"search_limit"`
	SearchTime  float64                         `json:"search_time"`
	SearchCount int                             `json:"search_count"`
	ServerCount int                             `json:"server_count"`
}

// Pritunl 客户端
type Pritunl struct {
	Request   *req.Req
	CsrfToken string
}

func NewPritunlClient(vpnName, orgName string) (*Pritunl, error) {
	var (
		err           error
		resp          *req.Resp
		doReq         = req.New()
		debug         bool
		csrfTokenBody = &PritunlGetCSRFTokenRespBody{}
	)

	pritunlUrl = beego.AppConfig.String(fmt.Sprintf("%s-vpn-%s::url", vpnName, orgName))
	username = beego.AppConfig.String(fmt.Sprintf("%s-vpn-%s::username", vpnName, orgName))
	password = beego.AppConfig.String(fmt.Sprintf("%s-vpn-%s::password", vpnName, orgName))
	organization = beego.AppConfig.String(fmt.Sprintf("%s-vpn-%s::organization", vpnName, orgName))
	debug, _ = beego.AppConfig.Bool(fmt.Sprintf("%s-vpn-%s::debug", vpnName, orgName))

	authUrl = fmt.Sprintf("%s%s", pritunlUrl, authSessionAPI)
	getCsrfTokenUrl = fmt.Sprintf("%s%s", pritunlUrl, getCsrfTokenAPI)

	if debug {
		req.Debug = true
	}

	// 登录
	doReq.EnableInsecureTLS(true)
	doReq.EnableCookie(true)
	resp, err = doReq.Post(authUrl, req.BodyJSON(&PritunlLogin{Username: username, Password: password}))
	if err != nil {
		return nil, err
	}

	// 获取csrf_token
	resp, err = doReq.Get(getCsrfTokenUrl)
	if err != nil {
		return nil, err
	}

	if err = resp.ToJSON(csrfTokenBody); err != nil {
		return nil, err
	}

	return &Pritunl{Request: doReq, CsrfToken: csrfTokenBody.CsrfToken}, nil
}

// 根据用户名搜索账号
func (p *Pritunl) SearchUser(username string) (*PritunlSearchUserRespBody, error) {

	var (
		searchUserUrl      = fmt.Sprintf("%s%s/%s", pritunlUrl, userAPI, organization)
		resp               *req.Resp
		err                error
		header             = req.Header{"Csrf-Token": p.CsrfToken}
		searchUserRespBody = &PritunlSearchUserRespBody{} //搜索用户时返回的body
	)

	if resp, err = p.Request.Get(searchUserUrl, header, req.Param{"search": username}); err != nil {
		return nil, err
	}

	if err = resp.ToJSON(searchUserRespBody); err != nil {
		return nil, err
	}

	return searchUserRespBody, nil
}

// 创建账号
func (p *Pritunl) CreateUser(name, email string) (ctype.Data, error) {
	var (
		createUserUrl      = fmt.Sprintf("%s%s/%s", pritunlUrl, userAPI, organization)
		header             = req.Header{"Csrf-Token": p.CsrfToken}
		createUserReqBody  *PritunlCreateUserReqBody
		createUserRespBody = PritunlCreateUserRespBody{}
		getProfileLinkUrl  string
		resp               *req.Resp
		err                error
		userId             string
		username           string
		searchUserResult   *PritunlSearchUserRespBody
	)

	// 检查该用户名是否存在
	if searchUserResult, err = p.SearchUser(name); err != nil {
		return nil, err
	}

	if len(searchUserResult.Users) != 0 {
		return nil, fmt.Errorf("already exist: %s", name)
	}

	// 创建用户
	createUserReqBody = &PritunlCreateUserReqBody{
		Organization:    organization,
		Name:            name,
		Email:           email,
		AuthType:        "local",
		BypassSecondary: false,
		ClientToClient:  false,
		Groups:          []string{},
		Pin:             nil,
		NetworkLinks:    []string{},
		DnsServers:      []string{},
		DnsSuffix:       "",
		PortForwarding:  []string{},
	}

	resp, err = p.Request.Post(createUserUrl, req.BodyJSON(createUserReqBody), header)
	if err != nil {
		return nil, err
	}

	if err = resp.ToJSON(&createUserRespBody); err != nil {
		return nil, err
	}

	userId = createUserRespBody[0].Id
	username = createUserRespBody[0].Name

	// 获取profile下载链接
	if getProfileLinkUrl, err = p.GetProfileLink(userId); err != nil {
		return nil, err
	}

	return ctype.Data{"userId": userId, "username": username,
		"profileLink": fmt.Sprintf("%s%s", pritunlUrl, getProfileLinkUrl),
	}, nil
}

// 获取某个用户的profile文件临时链接
func (p *Pritunl) GetProfileLink(userId string) (string, error) {
	var (
		getProfileLinkRespBody = PritunlGetProfileLinkRespBody{}
		header                 = req.Header{"Csrf-Token": p.CsrfToken}
		getProfileLinkUrl      string
		resp                   *req.Resp
		err                    error
	)

	// 获取profile下载链接
	getProfileLinkUrl = fmt.Sprintf("%s%s/%s/%s", pritunlUrl, downloadProfileLinkAPI, organization, userId)
	resp, err = p.Request.Get(getProfileLinkUrl, header)
	if err != nil {
		return "", err
	}

	if err = resp.ToJSON(&getProfileLinkRespBody); err != nil {
		return "", err
	}

	return getProfileLinkRespBody.ViewUrl, nil
}

// 批量添加用户
func (p *Pritunl) MultiCreateUser(userinfo []byte) ([]ctype.Data, error) {
	var (
		resp               *req.Resp
		err                error
		header             = req.Header{"Csrf-Token": p.CsrfToken}
		multiCreateUserReq = PritunlMultiCreateUserReqBody{}
		multiCreateUserUrl = fmt.Sprintf("%s%s/%s/multi", pritunlUrl, userAPI, organization)
		createResult       = PritunlCreateUserRespBody{}
		data               = []ctype.Data{}
	)

	if err = json.Unmarshal(userinfo, &multiCreateUserReq); err != nil {
		return nil, err
	}

	resp, err = p.Request.Post(multiCreateUserUrl, req.BodyJSON(multiCreateUserReq), header)
	if err != nil {
		return nil, err
	}

	if err = resp.ToJSON(&createResult); err != nil {
		return nil, err
	}

	for _, user := range createResult {
		profileUrl, _ := p.GetProfileLink(user.Id)
		data = append(data, ctype.Data{"username": user.Name, "userId": user.Id,
			"profileLink": fmt.Sprintf("%s%s", pritunlUrl, profileUrl)})
	}

	return data, nil
}

// 删除用户
func (p *Pritunl) DeleteUser(username string) (ctype.Data, error) {
	var (
		resp             *req.Resp
		err              error
		header           = req.Header{"Csrf-Token": p.CsrfToken}
		searchUserResult *PritunlSearchUserRespBody
		userResult       pritunlSearchUserRespUserBody
		deleteUserUrl    string
		length           int
	)

	if searchUserResult, err = p.SearchUser(username); err != nil {
		return nil, err
	}

	length = len(searchUserResult.Users)

	if length == 0 {
		return nil, fmt.Errorf("not found user: %s", username)
	}

	if length == 1 {
		userResult = searchUserResult.Users[0]
		deleteUserUrl = fmt.Sprintf("%s%s/%s/%s", pritunlUrl, userAPI, organization, userResult.ID)
		if resp, err = p.Request.Delete(deleteUserUrl, header); err != nil {
			return nil, err
		}

		if resp.Response().StatusCode != 200 {
			return nil, fmt.Errorf("call pritunl api error, statuscode: %d", resp.Response().StatusCode)
		}

		return ctype.Data{"username": userResult.Name, "status": "deleted"}, nil
	}

	return nil, fmt.Errorf("found %d user as same username, please use admin web console delete", length)
}

// 启用禁用用户
func (p *Pritunl) EnableDisableUser(username, status string) (ctype.Data, error) {
	var (
		resp             *req.Resp
		err              error
		header           = req.Header{"Csrf-Token": p.CsrfToken}
		searchUserResult *PritunlSearchUserRespBody
		userResult       pritunlSearchUserRespUserBody
		disableUserUrl   string
		length           int
	)

	if searchUserResult, err = p.SearchUser(username); err != nil {
		return nil, err
	}

	length = len(searchUserResult.Users)

	if length == 0 {
		return nil, fmt.Errorf("not found user: %s", username)
	}

	if length == 1 {
		userResult = searchUserResult.Users[0]
		disableUserUrl = fmt.Sprintf("%s%s/%s/%s", pritunlUrl, userAPI, organization, userResult.ID)

		switch status {
		case "disable":
			userResult.Disabled = true
			if resp, err = p.Request.Put(disableUserUrl, req.BodyJSON(userResult), header); err != nil {
				return nil, err
			}

			if resp.Response().StatusCode == 200 {
				return ctype.Data{"username": username, "status": "disabled"}, nil
			}
		case "enable":
			userResult.Disabled = false
			if resp, err = p.Request.Put(disableUserUrl, req.BodyJSON(userResult), header); err != nil {
				return nil, err
			}

			if resp.Response().StatusCode == 200 {
				return ctype.Data{"username": username, "status": "enabled"}, nil
			}
		default:
			return nil, fmt.Errorf("status fields can only be enable/disable")
		}
	}

	return nil, fmt.Errorf("found %d user as same username, please use admin web console disable", length)

}
