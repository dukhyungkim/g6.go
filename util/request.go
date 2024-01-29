package util

import (
	"fmt"
	"net/http"
)

type Request struct {
	State   State             `json:"state"`
	BaseURL string            `json:"base_url"`
	Cookies map[string]string `json:"cookies"`
	Session map[string]string `json:"session"`
}

func NewRequest(r *http.Request) Request {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s/", scheme, r.Host)

	return Request{
		BaseURL: baseURL,
	}
}

func FromMap(m map[string]any) Request {
	return MapToStruct[Request](m)
}

func (r Request) ToMap() map[string]any {
	return StructToMap(r)
}

type State struct {
	Config       Config `json:"config"`
	IsSuperAdmin bool   `json:"is_super_admin"`
	LoginMember  string `json:"login_member"`
	IsMobile     bool   `json:"is_mobile"`
	Editor       string `json:"editor"`
	Title        string `json:"title"`
	CookieDomain string `json:"cookie_domain"`
}

type Config struct {
	CfAddMeta string `json:"cf_add_meta"`
}
