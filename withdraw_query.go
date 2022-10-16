package tt_pay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"net/url"
	"time"

	"github.com/liaoxxxx/tt_pay/config"
	"github.com/liaoxxxx/tt_pay/consts"
	"github.com/liaoxxxx/tt_pay/util"
)

// 订单查询接口
func WithdrawQuery(ctx context.Context, req *WithdrawQueryRequest) (*WithdrawQueryResponse, error) {
	if err := req.checkParams(); err != nil {
		return nil, err
	}
	resp := NewWithdrawQueryResponse()
	err := Execute(ctx, req.TPClientTimeoutMs, req, resp)
	// 当出现请求失败错误时，不封装
	if _, ok := err.(*util.Error); ok {
		return nil, err
	}
	if err != nil {
		return nil, util.Wrap(err, "WithdrawQuery failed when [Execute()]")
	}
	return resp, nil
}

// 订单查询Request
type WithdrawQueryRequest struct {
	config.Config
	Method          string
	Format          string
	Charset         string
	SignType        string
	Timestamp       string
	Version         string
	bizContent      *simplejson.Json
	OutTradeNo      string
	WithdrawTradeNo string
	path            string
}

// New函数内赋默认值，目前含默认值（或仅支持一个值的）参数包括：
// Version = "1.0"
// SignType = "MD5"
// Format = "JSON"
// Charset = "utf-8"
// Path = "gateway"
// Config.TPDomain = "https://tp-pay.snssdk.com"
// Method 根据不同接口设置
// Timestamp 自动设置Unix时间戳
// 另外，注意初始化bizContent，以免出现nil指针错误
func NewWithdrawQueryRequest(config config.Config) *WithdrawQueryRequest {
	ret := new(WithdrawQueryRequest)
	ret.Config = config
	ret.Version = "1.0"
	ret.SignType = "MD5"
	ret.Format = "JSON"
	ret.Charset = "utf-8"
	ret.path = consts.TPPath
	if len(ret.Config.TPDomain) == 0 {
		ret.Config.TPDomain = consts.TPDomain
	}
	ret.Method = consts.MethodWithdrawQuery
	ret.Timestamp = fmt.Sprintf("%d", time.Now().Unix())
	ret.bizContent = simplejson.New()
	return ret
}

// 将Request编码成POST请求的Body
func (req *WithdrawQueryRequest) Encode() (string, error) {
	// 加签
	req.bizContent.Set("merchant_id", req.MerchantId)
	req.bizContent.Set("out_trade_no", req.OutTradeNo)
	req.bizContent.Set("withdraw_trade_no", req.WithdrawTradeNo)

	bizContentBytes, err := req.bizContent.Encode()
	if err != nil {
		util.Debug("WithdrawQueryRequest Encode bizContent.Encode err: %s, bizContent %s\n", err, *req.bizContent)
		return "", util.Wrap(err, "WithdrawQueryRequest Encode failed when [bizContent.Encode()]")
	}

	signParams := make(map[string]interface{})
	signParams["app_id"] = req.Config.AppId
	signParams["method"] = req.Method
	signParams["format"] = req.Format
	signParams["charset"] = req.Charset
	signParams["sign_type"] = req.SignType
	signParams["timestamp"] = req.Timestamp
	signParams["version"] = req.Version
	signParams["biz_content"] = string(bizContentBytes)

	sign := util.BuildMd5WithSalt(signParams, req.Config.AppSecret)
	// 序列化
	values := url.Values{}
	values.Set("app_id", req.Config.AppId)
	values.Set("method", req.Method)
	values.Set("format", req.Format)
	values.Set("charset", req.Charset)
	values.Set("sign_type", req.SignType)
	values.Set("sign", sign)
	values.Set("timestamp", req.Timestamp)
	values.Set("version", req.Version)
	values.Set("biz_content", string(bizContentBytes))

	return values.Encode(), nil
}

// 生成该次请求logid
// out_order_no 和 trade_no 哪个不空用哪个，都不空优先用out_order_no
func (req *WithdrawQueryRequest) GetLogId() string {
	id := ""
	if len(req.OutTradeNo) != 0 {
		id = req.OutTradeNo
	}
	if len(req.WithdrawTradeNo) != 0 {
		id = req.WithdrawTradeNo
	}
	return fmt.Sprintf("%s_%s_%s_%s", req.Config.AppId, req.Config.MerchantId, id, req.Timestamp)
}

// 获取请求url地址
func (req *WithdrawQueryRequest) GetUrl() string {
	return req.Config.TPDomain + "/" + req.path
}

// 提供该接口，方便业务方设置可选参数，比如product_code、payment_type等
func (req *WithdrawQueryRequest) SetBizContentKV(key string, val interface{}) {
	req.bizContent.Set(key, val)
}

// 提现查询响应
type WithdrawQueryResponse struct {
	Data            *simplejson.Json
	WithdrawTradeNo string `json:"withdraw_trade_no"`
	OutTradeNo      string `json:"out_trade_no"`
	MerchantId      string `json:"merchant_id"`
	Uid             string `json:"uid"`
	CreateTime      string `json:"create_time"`
	TradeTime       string `json:"trade_time"`
	Status          string `json:"status"`
	TradeName       string `json:"trade_name"`
	TradeDesc       string `json:"trade_desc"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	WithdrawType    string `json:"withdraw_type"`
	Account         string `json:"account"`
	Name            string `json:"name"`
	ValiditySeconds string `json:"validity_seconds"`
	ErrorCode       string `json:"err_code"`
	ErrMsg          string `json:"err_msg"`
}

// 初始化提现查询响应
func NewWithdrawQueryResponse() *WithdrawQueryResponse {
	ret := new(WithdrawQueryResponse)
	ret.Data = simplejson.New()
	return ret
}

// 将响应json数据反序列化为对应接口
func (resp *WithdrawQueryResponse) Decode() error {
	respBytes, err := resp.Data.Get("response").Encode()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(respBytes, resp); err != nil {
		return err
	}
	return nil
}

// 设置原始响应
func (resp *WithdrawQueryResponse) SetData(data *simplejson.Json) {
	resp.Data = data
}

// 参数查验
func (req *WithdrawQueryRequest) checkParams() error {
	if err := util.CheckVersion(req.Version); err != nil {
		return err
	}

	if err := util.CheckAppId(req.AppId); err != nil {
		return err
	}

	if req.Method != consts.MethodWithdrawQuery {
		return fmt.Errorf(util.ErrorFormat, "Method", "must be tp.withdraw.create")
	}

	if err := util.CheckFormat(req.Format); err != nil {
		return err
	}

	if err := util.CheckCharset(req.Charset); err != nil {
		return err
	}

	if err := util.CheckSignType(req.SignType); err != nil {
		return err
	}

	if err := util.CheckTimeStamp(req.Timestamp); err != nil {
		return err
	}

	if err := util.CheckVersion(req.Version); err != nil {
		return err
	}

	if err := util.CheckBizContent(req.bizContent); err != nil {
		return err
	}

	if err := util.CheckAppId(req.AppId); err != nil {
		return err
	}

	if err := util.CheckMerchantId(req.MerchantId); err != nil {
		return err
	}

	// 二选一参数判断
	if req.OutTradeNo == "" && req.WithdrawTradeNo == "" {
		return errors.New("OutOrderNo and TradeNo can't both be blank")
	}

	return nil
}
