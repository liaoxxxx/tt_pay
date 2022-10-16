package tt_pay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/liaoxxxx/tt_pay/config"
	"github.com/liaoxxxx/tt_pay/consts"
	"github.com/liaoxxxx/tt_pay/util"
)

// 提现下单接口
func WithdrawCreate(ctx context.Context, req *WithdrawCreateRequest) (*WithdrawCreateResponse, error) {
	if err := req.checkParams(); err != nil {
		return nil, err
	}
	resp := NewWithdrawCreateResponse(req)
	// 非登录态需要与财经后端通信
	if !req.WithLogin {
		err := Execute(ctx, req.TPClientTimeoutMs, req, resp)
		// 当出现请求失败错误时，不封装
		if _, ok := err.(*util.Error); ok {
			return nil, err
		}
		if err != nil {
			return nil, util.Wrap(err, "WithdrawCreate failed when [Execute()]")
		}
	}
	return resp, nil
}

// 提现下单Request
type WithdrawCreateRequest struct {
	config.Config
	WithLogin            bool // 此参数用来区分登录态及非登录态
	Method               string
	Format               string
	Charset              string
	SignType             string
	Timestamp            string
	Version              string
	bizContent           *simplejson.Json
	OutTradeNo           string
	Uid                  string
	TotalAmount          int
	Currency             string
	TradeName            string
	TradeDesc            string
	ProductCode          string
	PaymentType          string
	TradeTime            string
	ValidTime            string
	NotifyUrl            string
	ReturnUrl            string
	ExtParam             string
	SettlementExt        string
	RiskInfo             string
	AccountType          string
	SettlementProuctCode string
	TransCode            string
	Exts                 string
	path                 string
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
func NewWithdrawCreateRequest(config config.Config) *WithdrawCreateRequest {
	ret := new(WithdrawCreateRequest)
	ret.Config = config
	ret.Version = "1.0"
	ret.SignType = "MD5"
	ret.Format = "JSON"
	ret.Charset = "utf-8"
	ret.path = consts.TPPath
	if len(ret.Config.TPDomain) == 0 {
		ret.Config.TPDomain = consts.TPDomain
	}
	ret.Method = consts.MethodWithdrawCreate
	ret.Timestamp = fmt.Sprintf("%d", time.Now().Unix())
	ret.bizContent = simplejson.New()
	return ret
}

// 将Request编码成POST请求的Body
func (req *WithdrawCreateRequest) Encode() (string, error) {
	//加签
	req.bizContent.Set("out_trade_no", req.OutTradeNo)
	req.bizContent.Set("uid", req.Uid)
	req.bizContent.Set("merchant_id", req.MerchantId)
	req.bizContent.Set("amount", req.TotalAmount)
	req.bizContent.Set("currency", req.Currency)
	req.bizContent.Set("trade_name", req.TradeName)
	req.bizContent.Set("trade_desc", req.TradeDesc)
	req.bizContent.Set("product_code", req.ProductCode)
	req.bizContent.Set("payment_type", req.PaymentType)
	req.bizContent.Set("trade_time", req.TradeTime)
	req.bizContent.Set("valid_time", req.ValidTime)
	req.bizContent.Set("notify_url", req.NotifyUrl)
	req.bizContent.Set("return_url", req.ReturnUrl)
	req.bizContent.Set("ext_param", req.ExtParam)
	req.bizContent.Set("settlement_ext", req.SettlementExt)
	req.bizContent.Set("risk_info", req.RiskInfo)
	req.bizContent.Set("account_type", req.AccountType)
	req.bizContent.Set("settlement_product_code", req.SettlementProuctCode)

	// Json encode
	bizContentBytes, err := req.bizContent.Encode()
	if err != nil {
		util.Debug("WithdrawCreateRequest Encode bizContent.Encode err: %s, bizContent %s\n", err, *req.bizContent)
		return "", util.Wrap(err, "WithdrawCreateRequest Encode failed when [bizContent.Encode()]")
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
	// URL Encode
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

// GetLogId 生成该次请求logid
func (req *WithdrawCreateRequest) GetLogId() string {
	return fmt.Sprintf("%s_%s_%s_%s", req.AppId, req.MerchantId, req.OutTradeNo, req.Timestamp)
}

// GetUrl 获取请求url地址
func (req *WithdrawCreateRequest) GetUrl() string {
	return req.Config.TPDomain + "/" + req.path
}

// 提现下单响应
type WithdrawCreateResponse struct {
	Data            *simplejson.Json
	WithdrawTradeNo string                 `json:"withdraw_trade_no"`
	req             *WithdrawCreateRequest // 包含拉起收银台所需参数
}

// 初始化体现下单响应
func NewWithdrawCreateResponse(req *WithdrawCreateRequest) *WithdrawCreateResponse {
	ret := new(WithdrawCreateResponse)
	ret.Data = simplejson.New()
	ret.req = req
	return ret
}

// 返回拉起sdk收银台已经签名好的参数对
func (resp *WithdrawCreateResponse) GetCashdeskSdkParams() (string, error) {
	cashdeskParams, err := resp.getCashdeskSdkParams()
	if err != nil {
		return "", util.Wrap(err, "GetCashdeskSdkParams failed when [resp.getCashdeskSdkParams()]")
	}
	returnParams, err := util.JsonMarshal(cashdeskParams)
	if err != nil {
		return "", util.Wrap(err, "GetCashdeskSdkParams failed when [JsonMarshal()]")
	}
	return returnParams, nil
}

// 返回拉起H5收银台url
func (resp *WithdrawCreateResponse) GetCashdeskWithdrawH5() (string, error) {
	cashDeskParams, err := resp.getCashdeskSdkParams()
	if err != nil {
		return "", util.Wrap(err, "GetCashdeskWithdrawH5 failed when [resp.getCashdeskSdkParams()]")
	}
	// url encode cashDeskParams
	paramsForEncode := make(map[string][]string)
	for key, val := range cashDeskParams {
		paramsForEncode[key] = []string{val.(string)}
	}
	query := url.Values(paramsForEncode).Encode()
	return resp.req.TPDomain + "/redPacketWithdraw?" + query, nil
}

// 内部函数，对参数加签
func (resp *WithdrawCreateResponse) getCashdeskSdkParams() (map[string]interface{}, error) {
	cashDeskParams := make(map[string]interface{})

	cashDeskParams["app_id"] = resp.req.AppId
	if resp.req.OutTradeNo != "" {
		cashDeskParams["out_trade_no"] = resp.req.OutTradeNo
	}
	cashDeskParams["merchant_id"] = resp.req.MerchantId
	cashDeskParams["product_code"] = resp.req.ProductCode
	cashDeskParams["payment_type"] = resp.req.PaymentType
	if resp.req.Exts != "" {
		cashDeskParams["exts"] = resp.req.Exts
	}

	// 非登录态参数
	if !resp.req.WithLogin {
		cashDeskParams["withdraw_trade_no"] = resp.WithdrawTradeNo
		cashDeskParams["uid"] = resp.req.Uid

	} else {
		// 大于0代表商户传了该参数
		if resp.req.TotalAmount > 0 {
			cashDeskParams["total_amount"] = strconv.Itoa(resp.req.TotalAmount)
		}
		if resp.req.TransCode != "" {
			cashDeskParams["trans_code"] = resp.req.TransCode
		}
		if resp.req.NotifyUrl != "" {
			cashDeskParams["notify_url"] = resp.req.NotifyUrl
		}
	}

	// 在登录态，商户未传TotalAmount 且 exts参数里没有openid时，不需要加签
	// 反之，需要加签
	if !(resp.req.WithLogin && resp.req.TotalAmount == 0 && !strings.Contains(resp.req.Exts, "openid")) {
		cashDeskParams["sign_type"] = resp.req.SignType
		cashDeskParams["sign"] = util.BuildMd5WithSalt(cashDeskParams, resp.req.AppSecret)
	}

	if resp.req.ReturnUrl != "" {
		cashDeskParams["returnUrl"] = resp.req.ReturnUrl
	}
	if resp.req.RiskInfo != "" {
		cashDeskParams["risk_info"] = resp.req.RiskInfo
	}
	return cashDeskParams, nil
}

// 将响应json数据反序列化为对应接口
func (resp *WithdrawCreateResponse) Decode() error {
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
func (resp *WithdrawCreateResponse) SetData(data *simplejson.Json) {
	resp.Data = data
}

func (req *WithdrawCreateRequest) checkParams() error {
	if req.WithLogin {
		if err := req.checkParamsWithLogin(); err != nil {
			return err
		}
	} else {
		if err := req.checkParamsWithoutLogin(); err != nil {
			return err
		}
	}
	return nil
}

func (req *WithdrawCreateRequest) checkParamsWithLogin() error {
	if req.Method != consts.MethodWithdrawCreate {
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

	if req.ProductCode != "withdraw" {
		return errors.New("invalid param: ProductCode")
	}

	if err := util.CheckPaymentType(req.PaymentType); err != nil {
		return err
	}

	if req.TotalAmount < 0 {
		return errors.New("invalid param: TotalAmount")
	}

	if req.NotifyUrl != "" {
		if err := util.CheckNotifyUrl(req.NotifyUrl); err != nil {
			return err
		}
	}

	if req.RiskInfo != "" {
		if err := util.CheckRiskInfo(req.RiskInfo); err != nil {
			return err
		}
	}

	// 这里区分商户指定提现金额和商户未指定提现金额的参数查验
	if req.TotalAmount > 0 {
		if err := util.CheckOutTradeNo(req.OutTradeNo); err != nil {
			return err
		}

		if len(req.Exts) > 0 {
			if err := util.CheckExt(req.Exts); err != nil {
				return err
			}
		}
	} else {
		if len(req.Exts) > 0 {
			if err := util.CheckExt(req.Exts); err != nil {
				return err
			}
		}
	}

	return nil
}

func (req *WithdrawCreateRequest) checkParamsWithoutLogin() error {
	if req.Method != consts.MethodWithdrawCreate {
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

	if err := util.CheckTotalAmount(req.TotalAmount); err != nil {
		return err
	}

	if req.OutTradeNo != "" {
		if err := util.CheckOutTradeNo(req.OutTradeNo); err != nil {
			return err
		}
	}

	if err := util.CheckUid(req.Uid); err != nil {
		return err
	}

	if err := util.CheckCurrency(req.Currency); err != nil {
		return err
	}

	if err := util.CheckTradeName(req.TradeName); err != nil {
		return err
	}

	if err := util.CheckTradeDesc(req.TradeDesc); err != nil {
		return err
	}

	if err := util.CheckTradeTime(req.TradeTime); err != nil {
		return err
	}

	if err := util.CheckValidTime(req.ValidTime); err != nil {
		return err
	}

	if err := util.CheckNotifyUrl(req.NotifyUrl); err != nil {
		return err
	}

	if err := util.CheckRiskInfo(req.RiskInfo); err != nil {
		return err
	}

	if req.ProductCode != "withdraw" {
		return errors.New("invalid param: ProductCode")
	}

	if err := util.CheckPaymentType(req.PaymentType); err != nil {
		return err
	}

	if len(req.Exts) > 0 {
		if err := util.CheckExt(req.Exts); err != nil {
			return err
		}
	}

	if len(req.ExtParam) > 0 {
		if err := util.CheckExtParam(req.ExtParam); err != nil {
			return err
		}
	}

	if len(req.SettlementExt) > 0 {
		if err := util.CheckSettlementExt(req.SettlementExt); err != nil {
			return err
		}
	}

	return nil
}
