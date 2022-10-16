package tt_pay

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"net/url"
	"strconv"
	"time"

	"github.com/liaoxxxx/tt_pay/config"
	"github.com/liaoxxxx/tt_pay/consts"
	"github.com/liaoxxxx/tt_pay/util"
)

// 预下单接口
func TradeCreate(ctx context.Context, req *TradeCreateRequest) (*TradeCreateResponse, error) {
	resp := NewTradeCreateResponse(req)
	// 1.0需要与财经后端通信取得"trade_no"
	if req.AppletVersion == "2.0+" || req.AppletVersion == "1.0" {
		// 查验1.0参数
		if err := req.checkParams1_0(); err != nil {
			return nil, err
		}

		// 2019/08/06
		// 现在不需要从交易获取trade_no了，可以直接用out_order_no代替trade_no

		//err := Execute(ctx, req.TPClientTimeoutMs, req, resp)
		//// 当出现请求失败错误时，不封装
		//if _, ok := err.(*util.Error); ok {
		//	return nil, err
		//}
		//if err != nil {
		//	return nil, util.Wrap(err, "TradeCreate failed when [Execute]")
		//}
	}
	if req.AppletVersion == "2.0" || req.AppletVersion == "2.0+" {
		// 查验2.0参数
		if err := req.checkParams2_0(); err != nil {
			return nil, err
		}
	}
	return resp, nil
}

// 预下单Request
type TradeCreateRequest struct {
	config.Config
	Method         string
	Format         string
	Charset        string
	SignType       string
	Timestamp      string
	Version        string
	AppletVersion  string
	bizContent     *simplejson.Json
	OutOrderNo     string
	Uid            string
	UidType        string
	TotalAmount    int
	Currency       string
	TradeType      string
	Subject        string
	Body           string
	ProductCode    string
	PaymentType    string
	PaymentType1_0 string
	TradeTime      string
	ValidTime      string
	NotifyUrl      string
	RiskInfo       string
	Params         string
	ProductId      string
	PayChannel     string
	PayDiscount    string
	ServiceFee     string
	LimitPay       string
	Path           string
	AlipayUrl      string
	WxUrl          string
	WxType         string
	ExtParam       string
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
func NewTradeCreateRequest(config config.Config) *TradeCreateRequest {
	ret := new(TradeCreateRequest)
	ret.Config = config
	ret.Version = "1.0"
	ret.SignType = "MD5"
	ret.Format = "JSON"
	ret.Charset = "utf-8"
	ret.Path = consts.TPPath
	if len(ret.Config.TPDomain) == 0 {
		ret.Config.TPDomain = consts.TPDomain
	}
	ret.Method = consts.MethodTradeCreate
	ret.Timestamp = fmt.Sprintf("%d", time.Now().Unix())
	ret.bizContent = simplejson.New()
	return ret
}

// 将Request编码成POST请求的Body
func (req *TradeCreateRequest) Encode() (string, error) {
	//加签
	req.bizContent.Set("out_order_no", req.OutOrderNo)
	req.bizContent.Set("uid", req.Uid)
	req.bizContent.Set("uid_type", req.UidType)
	req.bizContent.Set("merchant_id", req.MerchantId)
	req.bizContent.Set("total_amount", req.TotalAmount)
	req.bizContent.Set("currency", req.Currency)
	req.bizContent.Set("subject", req.Subject)
	req.bizContent.Set("body", req.Body)
	req.bizContent.Set("product_code", req.ProductCode)
	req.bizContent.Set("payment_type", req.PaymentType)
	req.bizContent.Set("trade_time", req.TradeTime)
	req.bizContent.Set("valid_time", req.ValidTime)
	req.bizContent.Set("notify_url", req.NotifyUrl)
	req.bizContent.Set("service_fee", req.ServiceFee)
	req.bizContent.Set("risk_info", req.RiskInfo)

	// Json encode
	bizContentBytes, err := req.bizContent.Encode()
	if err != nil {
		util.Debug("TradeCreateRequest Encode bizContent.Encode err: %s, bizContent %s\n", err, *req.bizContent)
		return "", util.Wrap(err, "TradeCreateRequest Encode failed when [bizContent.Encode()]")
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

// 生成此次请求logid
func (req *TradeCreateRequest) GetLogId() string {
	return fmt.Sprintf("%s_%s_%s_%s", req.AppId, req.MerchantId, req.OutOrderNo, req.Timestamp)
}

// 获取请求url地址
func (req *TradeCreateRequest) GetUrl() string {
	return req.Config.TPDomain + "/" + req.Path
}

// 提供该接口，方便业务方设置可选参数，比如product_code、payment_type等
func (req *TradeCreateRequest) SetBizContentKV(key string, val interface{}) {
	req.bizContent.Set(key, val)
}

// 预下单接口响应
type TradeCreateResponse struct {
	Data    *simplejson.Json
	TradeNo string `json:"trade_no"`
	URL     string `json:"url"`
	req     *TradeCreateRequest
}

// 初始化预下单响应
func NewTradeCreateResponse(req *TradeCreateRequest) *TradeCreateResponse {
	ret := new(TradeCreateResponse)
	ret.Data = simplejson.New()
	ret.req = req
	return ret
}

// 返回拉起小程序收银台的参数, json字符串
func (resp *TradeCreateResponse) GetCashdeskAppletParams() (string, error) {
	returnMap := make(map[string]string)
	switch resp.req.AppletVersion {
	case "1.0":
		returnJson, err := resp.getAppletParams1_0()
		if err != nil {
			return "", util.Wrap(err, "GetCashdeskAppletParams failed when[getAppletParams1_0()]")
		}
		returnMap["1.0"] = returnJson
	case "2.0":
		returnJson, err := resp.getAppletParams2_0()
		if err != nil {
			return "", util.Wrap(err, "GetCashdeskAppletParams failed when[getAppletParams2_0()]")
		}
		returnMap["2.0"] = returnJson
	case "2.0+":
		returnJson1_0, err := resp.getAppletParams1_0()
		if err != nil {
			return "", util.Wrap(err, "GetCashdeskAppletParams failed when[getAppletParams1_0()]")
		}
		returnJson2_0, err := resp.getAppletParams2_0()
		if err != nil {
			return "", util.Wrap(err, "GetCashdeskAppletParams failed when[getAppletParams2_0()]")
		}
		returnMap["1.0"] = returnJson1_0
		returnMap["2.0"] = returnJson2_0

	case "3.0":
		// 2020-02-07：新版优化之后仅需返回支付要素，不用返回版本
		returnJson, err := resp.getAppletParams2_0()
		if err != nil {
			return "", util.Wrap(err, "GetCashdeskAppletParams failed when[getAppletParams2_0()]")
		}
		return returnJson, nil

	default:
		return "", fmt.Errorf(util.ErrorFormat, "AppletVerion", "AppletVersion can only be 1.0, 2.0 or 2.0+")
	}
	returnJson, err := util.JsonMarshal(returnMap)
	if err != nil {
		return "", util.Wrap(err, "GetCashdeskAppletParams failed when[JsonMarshal()]")
	}
	return returnJson, nil
}

// 返回拉起二维码收银台的参数，URL
func (resp *TradeCreateResponse) GetCashdeskQRParams() (string, error) {
	return resp.URL, nil
}

// 小程序1.0参数
func (resp *TradeCreateResponse) getAppletParams1_0() (string, error) {
	appletParams := make(map[string]interface{})

	appletParams["app_id"] = resp.req.AppId
	appletParams["sign_type"] = resp.req.SignType
	appletParams["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())

	// 2019/08/06
	// 现在不需要从交易获取trade_no了，可以直接用out_order_no代替trade_no
	appletParams["trade_no"] = resp.req.OutOrderNo

	appletParams["merchant_id"] = resp.req.MerchantId
	appletParams["uid"] = resp.req.Uid
	appletParams["total_amount"] = resp.req.TotalAmount

	paramString, err := util.JsonMarshal(map[string]string{"url": resp.req.Params})
	if err != nil {
		return "", util.Wrap(err, "getAppletParams1_0 failed when [JsonMarshal()]")
	}
	appletParams["params"] = paramString

	appletParams["sign"] = util.BuildMd5WithSalt(appletParams, resp.req.AppSecret)

	appletParams["method"] = consts.MethodTradeConfirm // 方法要改为请求confirm
	appletParams["pay_type"] = resp.req.PaymentType1_0 // pay_type
	appletParams["pay_channel"] = resp.req.PayChannel
	appletParams["risk_info"] = resp.req.RiskInfo
	//if resp.req.ReturnUrl != "" {
	//	appletParams["return_url"] = resp.req.ReturnUrl
	//}
	//if resp.req.ShowURL != "" {
	//	appletParams["show_url"] = resp.req.ShowURL
	//}

	returnParams, err := util.JsonMarshal(appletParams)
	if err != nil {
		return "", util.Wrap(err, "getAppletParams1_0 failed when [JsonMarshal()]")
	}

	return returnParams, nil
}

// 小程序2.0参数
func (resp *TradeCreateResponse) getAppletParams2_0() (string, error) {
	cashDeskParams := make(map[string]interface{})

	cashDeskParams["app_id"] = resp.req.AppId
	cashDeskParams["sign_type"] = resp.req.SignType
	cashDeskParams["merchant_id"] = resp.req.MerchantId
	if resp.req.Uid != "" {
		cashDeskParams["uid"] = resp.req.Uid
	}
	if resp.req.OutOrderNo != "" {
		cashDeskParams["out_order_no"] = resp.req.OutOrderNo
	}
	cashDeskParams["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
	cashDeskParams["total_amount"] = strconv.Itoa(resp.req.TotalAmount)
	if resp.req.NotifyUrl != "" {
		cashDeskParams["notify_url"] = resp.req.NotifyUrl
	}
	if resp.req.TradeType != "" {
		cashDeskParams["trade_type"] = resp.req.TradeType
	}
	if resp.req.ProductCode != "" {
		cashDeskParams["product_code"] = resp.req.ProductCode
	}
	if resp.req.PaymentType != "" {
		cashDeskParams["payment_type"] = resp.req.PaymentType
	}
	if resp.req.Subject != "" {
		cashDeskParams["subject"] = resp.req.Subject
	}
	if resp.req.Body != "" {
		cashDeskParams["body"] = resp.req.Body
	}
	if resp.req.TradeTime != "" {
		cashDeskParams["trade_time"] = resp.req.TradeTime
	}
	if resp.req.ValidTime != "" {
		cashDeskParams["valid_time"] = resp.req.ValidTime
	}
	if resp.req.Currency != "" {
		cashDeskParams["currency"] = resp.req.Currency
	}
	if resp.req.Version != "" {
		cashDeskParams["version"] = resp.req.Version
	}
	if resp.req.AlipayUrl != "" {
		cashDeskParams["alipay_url"] = resp.req.AlipayUrl
	}
	if resp.req.WxUrl != "" {
		cashDeskParams["wx_url"] = resp.req.WxUrl
	}
	if resp.req.WxType != "" {
		cashDeskParams["wx_type"] = resp.req.WxType
	}
	if resp.req.LimitPay != "" {
		cashDeskParams["limit_pay"] = resp.req.LimitPay
	}
	cashDeskParams["sign"] = util.BuildMd5WithSalt(cashDeskParams, resp.req.AppSecret)
	if resp.req.RiskInfo != "" {
		cashDeskParams["risk_info"] = resp.req.RiskInfo
	}

	returnParams, err := util.JsonMarshal(cashDeskParams)
	if err != nil {
		return "", util.Wrap(err, "getAppletParams2_0 failed when [JsonMarshal()]")
	}

	return returnParams, nil
}

// 将响应json数据反序列化为对应接口
func (resp *TradeCreateResponse) Decode() error {
	var respBytes []byte
	var err error
	// 走网关的接口拿到的参数在response里
	// 二维码接口拿到的参数在data里
	switch resp.Data.Get("response").Interface() {
	case nil:
		respBytes, err = resp.Data.Get("data").Encode()
	default:
		respBytes, err = resp.Data.Get("response").Encode()
	}
	if err != nil {
		return err
	}
	if err := json.Unmarshal(respBytes, resp); err != nil {
		return err
	}
	return nil
}

// 设置原始响应
func (resp *TradeCreateResponse) SetData(data *simplejson.Json) {
	resp.Data = data
}

// 1.0版小程序参数查验
func (req *TradeCreateRequest) checkParams1_0() error {
	if err := util.CheckAppId(req.AppId); err != nil {
		return err
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

	if err := util.CheckOutOrderNo(req.OutOrderNo); err != nil {
		return err
	}

	if err := util.CheckUid(req.Uid); err != nil {
		return err
	}

	if err := util.CheckMerchantId(req.MerchantId); err != nil {
		return err
	}

	if err := util.CheckTotalAmount(req.TotalAmount); err != nil {
		return err
	}

	if err := util.CheckCurrency(req.Currency); err != nil {
		return err
	}

	if err := util.CheckSubject(req.Subject); err != nil {
		return err
	}

	if err := util.CheckBody(req.Body); err != nil {
		return err
	}

	if err := util.CheckTradeTime(req.TradeTime); err != nil {
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

// 2.0版小程序参数查验
func (req *TradeCreateRequest) checkParams2_0() error {
	if err := util.CheckAppId(req.AppId); err != nil {
		return err
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

	if err := util.CheckOutOrderNo(req.OutOrderNo); err != nil {
		return err
	}

	if err := util.CheckUid(req.Uid); err != nil {
		return err
	}

	if err := util.CheckMerchantId(req.MerchantId); err != nil {
		return err
	}

	if err := util.CheckTotalAmount(req.TotalAmount); err != nil {
		return err
	}

	if err := util.CheckCurrency(req.Currency); err != nil {
		return err
	}

	if err := util.CheckSubject(req.Subject); err != nil {
		return err
	}

	if err := util.CheckBody(req.Body); err != nil {
		return err
	}

	if err := util.CheckTradeTime(req.TradeTime); err != nil {
		return err
	}

	if err := util.CheckNotifyUrl(req.NotifyUrl); err != nil {
		return err
	}

	if err := util.CheckRiskInfo(req.RiskInfo); err != nil {
		return err
	}

	if err := util.CheckProductCode(req.ProductCode); err != nil {
		return err
	}

	if err := util.CheckPaymentType(req.PaymentType); err != nil {
		return err
	}

	if err := util.CheckCashDeskTradeType(req.TradeType); err != nil {
		return err
	}

	if err := util.CheckValidTime(req.ValidTime); err != nil {
		return err
	}

	return nil
}
