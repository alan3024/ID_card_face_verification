package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ValidationResult 保存从 API 解析的统一结果.
type ValidationResult struct {
	Success     bool
	ResultCode  int
	Message     string
	Score       float64
	Sex         string
	Birthday    string
	Address     string
	RawResponse string // 保存原始的 API 响应字符串
}

// Client 是身份验证 API 的接口.
// 用户可以实现此接口以支持不同的 API 提供商.
type Client interface {
	Validate(name, idCardNo, imageBase64 string) (*ValidationResult, error)
	SetAppCode(appCode string)
}

// aliyunResponseData 是阿里云响应中 'data' 字段的结构.
type aliyunResponseData struct {
	Result   int     `json:"result"`
	Msg      string  `json:"msg"`
	Score    float64 `json:"score"`
	Name     string  `json:"name"`
	Sex      string  `json:"sex"`
	Birthday string  `json:"birthday"`
	Address  string  `json:"address"`
	IdCardNo string  `json:"idCardNo"`
}

// aliyunAPIResponse 是阿里云 JSON 响应的顶层结构.
type aliyunAPIResponse struct {
	Success bool               `json:"success"`
	Data    aliyunResponseData `json:"data"`
	Message string             `json:"msg"` // 有时错误信息会在这里
}

// AliyunClient 是阿里云人脸识别 API 的默认实现.
type AliyunClient struct {
	URL     string
	AppCode string
	client  *http.Client
}

// NewAliyunClient 创建一个新的 AliyunClient.
func NewAliyunClient(appCode string) *AliyunClient {
	return &AliyunClient{
		URL:     "https://jmfaceid.market.alicloudapi.com/idcard-face/validate",
		AppCode: appCode,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// SetAppCode 允许在初始化后更改 AppCode.
func (c *AliyunClient) SetAppCode(appCode string) {
	c.AppCode = appCode
}

// Validate 为阿里云 API 实现了 Client 接口.
func (c *AliyunClient) Validate(name, idCardNo, imageBase64 string) (*ValidationResult, error) {
	if c.AppCode == "" || c.AppCode == "你自己的AppCode" {
		return nil, fmt.Errorf("请提供有效的 AppCode")
	}

	params := url.Values{}
	params.Add("name", name)
	params.Add("idCardNo", idCardNo)
	params.Add("facePhotoBase64", imageBase64)

	req, err := http.NewRequest("POST", c.URL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "APPCODE "+c.AppCode)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}
	rawResponse := string(body)

	if resp.StatusCode != http.StatusOK {
		// 尝试解析可能存在的 JSON 错误信息
		var errResp aliyunAPIResponse
		if json.Unmarshal(body, &errResp) == nil && errResp.Message != "" {
			return nil, fmt.Errorf("API 错误, 状态码 %d: %s", resp.StatusCode, errResp.Message)
		}
		return nil, fmt.Errorf("API 错误, 状态码 %d: %s", resp.StatusCode, rawResponse)
	}

	var apiResp aliyunAPIResponse
	// 使用 bytes.NewReader 因为原始的 resp.Body 已经被读取
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&apiResp); err != nil {
		return &ValidationResult{
			Success:     false,
			RawResponse: rawResponse,
		}, fmt.Errorf("解析JSON响应失败: %w", err)
	}

	// 即使 API success=false，也返回结构体，让调用者决定如何处理
	return &ValidationResult{
		Success:     apiResp.Success,
		ResultCode:  apiResp.Data.Result,
		Message:     apiResp.Data.Msg,
		Score:       apiResp.Data.Score,
		Sex:         apiResp.Data.Sex,
		Birthday:    apiResp.Data.Birthday,
		Address:     apiResp.Data.Address,
		RawResponse: rawResponse,
	}, nil
}
