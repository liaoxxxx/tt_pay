package tt_pay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/liaoxxxx/tt_pay/config"
	"github.com/liaoxxxx/tt_pay/consts"
	"github.com/liaoxxxx/tt_pay/util"
	"net/url"
	"time"
)

// 退款申请接口
func RefundCreate(ctx context.Context, req *RefundCreateRequest) (*RefundCreateResponse, error) {
	if err := req.checkParams(); err != nil {
		return nil, err
	}
	resp := NewRefundCreateResponse()
	err := Execute(ctx, req.TPClientTimeoutMs, req, resp)
	if err != nil {
		// 当出现请求失败错误时，不封装
		if _, ok := err.(*util.Error); ok {
			return nil, err
		}
		return nil, util.Wrap(err, "RefundCreate failed when [Execute()]")
	}
	return resp, nil
}

// 退款申请Request
type RefundCreateRequest struct {
	config.Config
	Method                string
	Format                string
	Charset               string
	Uid                   string
	OutOrderNo            string
	TradeNo               string
	SignType              string
	Version               string
	Timestamp             string
	OutRefundNo           string
	RefundAmount          int
	NotifyUrl             string
	RiskInfo              string
	path                  string
	SettlementProductCode string
	SettlementExt         string
	ProductCode           string
	PaymentType           string
	TransCode             string
	Reason                string
	ThridRefundAccount    string
	bizContent            *simplejson.Json
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
func NewRefundCreateRequest(config config.Config) *RefundCreateRequest {
	ret := new(RefundCreateRequest)
	ret.Config = config
	ret.Version = "1.0"
	ret.SignType = "MD5"
	ret.Format = "JSON"
	ret.Charset = "utf-8"
	ret.path = consts.TPPath
	if len(ret.Config.TPDomain) == 0 {
		ret.Config.TPDomain = consts.TPDomain
	}
	ret.Method = consts.MethodRefundCreate
	ret.Timestamp = fmt.Sprintf("%d", time.Now().Unix())
	ret.bizContent = simplejson.New()
	return ret
}

// 将Request编码成POST请求的Body
func (req *RefundCreateRequest) Encode() (string, error) {
	// 加签
	req.bizContent.Set("out_order_no", req.OutOrderNo)
	req.bizContent.Set("trade_no", req.TradeNo)
	req.bizContent.Set("merchant_id", req.Config.MerchantId)
	req.bizContent.Set("uid", req.Uid)
	req.bizContent.Set("out_refund_no", req.OutRefundNo)
	req.bizContent.Set("refund_amount", req.RefundAmount)
	req.bizContent.Set("notify_url", req.NotifyUrl)
	req.bizContent.Set("risk_info", req.RiskInfo)
	req.bizContent.Set("settlement_product_code", req.SettlementProductCode)
	req.bizContent.Set("settlement_ext", req.SettlementExt)
	req.bizContent.Set("product_code", req.ProductCode)
	req.bizContent.Set("payment_type", req.PaymentType)
	req.bizContent.Set("trans_code", req.TransCode)
	req.bizContent.Set("reason", req.Reason)
	req.bizContent.Set("third_refund_account", req.ThridRefundAccount)

	bizContentBytes, err := req.bizContent.Encode()
	if err != nil {
		util.Debug("RefundCreateRequest Encode bizContent.Encode err: %s, bizContent %s\n", err, *req.bizContent)
		return "", util.Wrap(err, "RefundCreateRequest Encode failed when [bizContent.Encode()]")
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
func (req *RefundCreateRequest) GetLogId() string {
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
func (req *RefundCreateRequest) GetUrl() string {
	return req.Config.TPDomain + "/" + req.path
}

// 提供该接口，方便业务方设置可选参数
// 比如product_code、payment_type等
func (req *RefundCreateRequest) SetBizContentKV(key string, val interface{}) {
	req.bizContent.Set(key, val)
}

// 退款接口响应
type RefundCreateResponse struct {
	Data         *simplejson.Json
	OutOrderNo   string `json:"out_order_no"`
	OutRefundNo  string `json:"out_refund_no"`
	RefundNo     string `json:"refund_no"`
	RefundAmount string `json:"refund_amount"`
}

// 初始化退款响应
func NewRefundCreateResponse() *RefundCreateResponse {
	ret := new(RefundCreateResponse)
	ret.Data = simplejson.New()
	return ret
}

// 将响应json数据反序列化为对应接口
func (resp *RefundCreateResponse) Decode() error {
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
func (resp *RefundCreateResponse) SetData(data *simplejson.Json) {
	resp.Data = data
}

// 目前只查验大写字母开头的参数(用户必传参数)
func (req *RefundCreateRequest) checkParams() error {

	if err := util.CheckAppId(req.AppId); err != nil {
		return err
	}

	if req.Method != consts.MethodRefundCreate {
		return fmt.Errorf(util.ErrorFormat, "Method", "must be tp.refund.create")
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

	if err := util.CheckOutRefundNo(req.OutRefundNo); err != nil {
		return err
	}

	if err := util.CheckRefundAmount(req.RefundAmount); err != nil {
		return err
	}

	if err := util.CheckNotifyUrl(req.NotifyUrl); err != nil {
		return err
	}

	if err := util.CheckRiskInfo(req.RiskInfo); err != nil {
		return err
	}

	return nil
}
