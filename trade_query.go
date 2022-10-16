package tt_pay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/zoooozz/tt_pay/config"
	"github.com/zoooozz/tt_pay/consts"
	"github.com/zoooozz/tt_pay/util"

	"github.com/bitly/go-simplejson"
)

// 订单查询接口
func TradeQuery(ctx context.Context, req *TradeQueryRequest) (*TradeQueryResponse, error) {
	if err := req.checkParams(); err != nil {
		return nil, err
	}
	resp := NewTradeQueryResponse()
	err := Execute(ctx, req.TPClientTimeoutMs, req, resp)
	// 当出现请求失败错误时，不封装
	if _, ok := err.(*util.Error); ok {
		return nil, err
	}
	if err != nil {
		return nil, util.Wrap(err, "TradeQuery failed when [Execute()]")
	}
	return resp, nil
}

// 订单查询Request
type TradeQueryRequest struct {
	config.Config
	Method     string
	Format     string
	Charset    string
	SignType   string
	Version    string
	Timestamp  string
	Uid        string
	UidType    string
	OutOrderNo string
	TradeNo    string
	path       string
	bizContent *simplejson.Json
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
func NewTradeQueryRequest(config config.Config) *TradeQueryRequest {
	ret := new(TradeQueryRequest)
	ret.Config = config
	ret.Version = "1.0"
	ret.SignType = "MD5"
	ret.Format = "JSON"
	ret.Charset = "utf-8"
	ret.path = consts.TPPath
	if len(ret.Config.TPDomain) == 0 {
		ret.Config.TPDomain = consts.TPDomain
	}
	ret.Method = consts.MethodTradeQuery
	ret.Timestamp = fmt.Sprintf("%d", time.Now().Unix())
	ret.bizContent = simplejson.New()
	return ret
}

// 将Request编码成POST请求的Body
func (req *TradeQueryRequest) Encode() (string, error) {
	// 加签
	req.bizContent.Set("merchant_id", req.MerchantId)
	req.bizContent.Set("uid", req.Uid)
	req.bizContent.Set("uid_type", req.UidType)
	req.bizContent.Set("out_order_no", req.OutOrderNo)
	req.bizContent.Set("trade_no", req.TradeNo)

	bizContentBytes, err := req.bizContent.Encode()
	if err != nil {
		util.Debug("TradeQueryRequest Encode bizContent.Encode err: %s, bizContent %s\n", err, *req.bizContent)
		return "", util.Wrap(err, "TradeQueryRequest Encode failed when [bizContent.Encode()]")
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
func (req *TradeQueryRequest) GetLogId() string {
	id := ""
	if len(req.TradeNo) != 0 {
		id = req.TradeNo
	}
	if len(req.OutOrderNo) != 0 {
		id = req.OutOrderNo
	}
	return fmt.Sprintf("%s_%s_%s_%s", req.Config.AppId, req.Config.MerchantId, id, req.Timestamp)
}

// 获取请求url地址
func (req *TradeQueryRequest) GetUrl() string {
	return req.Config.TPDomain + "/" + req.path
}

// 提供该接口，方便业务方设置可选参数，比如product_code、payment_type等
func (req *TradeQueryRequest) SetBizContentKV(key string, val interface{}) {
	req.bizContent.Set(key, val)
}

// 下单查询响应
type TradeQueryResponse struct {
	Data        *simplejson.Json
	TradeNo     string `json:"trade_no"`
	OutOrderNo  string `json:"out_order_no"`
	MerchantId  string `json:"merchant_id"`
	Uid         string `json:"uid"`
	Mid         string `json:"m_id"`
	CreateTime  string `json:"create_time"`
	PayTime     string `json:"pay_time"`
	TradeTime   string `json:"trade_time"`
	ExpireTime  string `json:"expire_time"`
	TradeStatus string `json:"trade_status"`
	TradeName   string `json:"trade_name"`
	TradeDesc   string `json:"trade_desc"`
	TotalAmount string `json:"total_amount"`
	Currency    string `json:"currency"`
	PayChannel  string `json:"pay_channel"`
	CouponNo    string `json:"coupon_no"`
	RealAmount  string `json:"real_amount"`
	ChannelExt  string `json:"channel_ext"`
}

// 初始化订单查询响应
func NewTradeQueryResponse() *TradeQueryResponse {
	ret := new(TradeQueryResponse)
	ret.Data = simplejson.New()
	return ret
}

// 将响应json数据反序列化为对应接口
func (resp *TradeQueryResponse) Decode() error {
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
func (resp *TradeQueryResponse) SetData(data *simplejson.Json) {
	resp.Data = data
}

// 参数查验
func (req *TradeQueryRequest) checkParams() error {
	if req.Method != consts.MethodTradeQuery {
		return fmt.Errorf(util.ErrorFormat, "Method", "must be tp.trade.query")
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

	if err := util.CheckUid(req.Uid); err != nil {
		return err
	}

	// 二选一参数判断
	if req.OutOrderNo == "" && req.TradeNo == "" {
		return errors.New("OutOrderNo and TradeNo can't both be blank")
	}

	if req.OutOrderNo != "" {
		if err := util.CheckOutOrderNo(req.OutOrderNo); err != nil {
			return err
		}
	}

	if req.TradeNo != "" {
		if err := util.CheckTradeNo(req.TradeNo); err != nil {
			return err
		}
	}

	return nil
}
