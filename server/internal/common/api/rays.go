package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mayfly-go/pkg/biz"
	"mayfly-go/pkg/req"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Rays struct{}

func (r *Rays) ReqConfs() *req.Confs {
	confs := []*req.Conf{
		req.NewGet("/overview", r.forward), req.NewGet("/accounts", r.forward), req.NewPost("/accounts", r.forward),
		req.NewPut("/accounts/:id", r.forward), req.NewDelete("/accounts/:id", r.forward),
		req.NewPost("/accounts/:id/revoke-sessions", r.forward), req.NewGet("/audit", r.forward), req.NewPost("/maintenance/expire-sessions", r.forward)}
	return req.NewConfs("/rays", confs...)
}
func operatorAllowed(id uint64, configured string) bool {
	for _, item := range strings.Split(configured, ",") {
		if strings.TrimSpace(item) == strconv.FormatUint(id, 10) {
			return id > 0
		}
	}
	return false
}
func (r *Rays) forward(rc *req.Ctx) {
	account := rc.GetLoginAccount()
	biz.IsTrue(operatorAllowed(account.Id, os.Getenv("RAYS_ADMIN_OPERATOR_IDS")), "当前账号未获 RAYS 管理授权")
	incoming := rc.GetRequest()
	path := strings.TrimPrefix(incoming.URL.Path, "/api/rays")
	method := incoming.Method
	if method == http.MethodPut {
		method = http.MethodPatch
	}
	data, err := bridgeRequest(incoming.Context(), method, path, incoming.URL.RawQuery, incoming.Body, account.Id)
	biz.ErrIsNil(err)
	rc.ResData = data
}
func bridgeRequest(ctx context.Context, method, path, query string, body io.Reader, actor uint64) (json.RawMessage, error) {
	base := strings.TrimRight(os.Getenv("RAYS_ADMIN_API_URL"), "/")
	secret := os.Getenv("RAYS_ADMIN_BRIDGE_TOKEN")
	u, err := url.Parse(base)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Path != "" || len(secret) < 32 {
		return nil, fmt.Errorf("RAYS 管理连接未配置")
	}
	payload, err := io.ReadAll(io.LimitReader(body, 65537))
	if err != nil || len(payload) > 65536 {
		return nil, fmt.Errorf("请求内容过大")
	}
	target := base + "/api/v1/admin" + path
	if query != "" {
		target += "?" + query
	}
	request, err := http.NewRequestWithContext(ctx, method, target, strings.NewReader(string(payload)))
	if err != nil {
		return nil, fmt.Errorf("管理请求无效")
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-RAYS-Admin-Token", secret)
	request.Header.Set("X-RAYS-Admin-Actor", fmt.Sprintf("mayfly:%d", actor))
	client := &http.Client{Timeout: 15 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("RAYS 服务暂时无法连接")
	}
	defer response.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(response.Body, 1048577))
	if err != nil || len(raw) > 1048576 {
		return nil, fmt.Errorf("管理响应过大")
	}
	var envelope struct {
		Success bool            `json:"success"`
		Data    json.RawMessage `json:"data"`
		Error   struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(raw, &envelope) != nil {
		return nil, fmt.Errorf("RAYS 返回无效响应")
	}
	if response.StatusCode >= 300 || !envelope.Success {
		return nil, fmt.Errorf("RAYS 操作失败 (%d): %s", response.StatusCode, envelope.Error.Message)
	}
	return envelope.Data, nil
}
