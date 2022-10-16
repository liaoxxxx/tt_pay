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

// 退款查询接口
func RefundQuery(ctx context.Context, req *RefundQueryRequest) (*RefundQueryResponse, error) {
	if err := req.checkParams(); err != nil {
		return nil, err
	}
	resp := NewRefundQueryResponse()
	err := Execute(ctx, req.TPClientTimeoutMs, req, resp)
	// 当出现请求失败错误时，不封装
	if _, ok := err.(*util.Error); ok {
		return nil, err
	}
	if err != nil {
		return nil, util.Wrap(err, "RefundQuery failed when [Execute()]")
	}
	return resp, nil
}

// 退款查询Request
type RefundQueryRequest struct {
	config.Config
	Method      string
	Format      string
	Charset     string
	Uid         string
	SignType    string
	Version     string
	Timestamp   string
	OutRefundNo string
	RefundNo    string
	path        string
	bizContent  *simplejson.Json
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
func NewRefundQueryRequest(config config.Config) *RefundQueryRequest {
	ret := new(RefundQueryRequest)
	ret.Config = config
	ret.Version = "1.0"
	ret.SignType = "MD5"
	ret.Format = "JSON"
	ret.Charset = "utf-8"
	ret.path = consts.TPPath
	if len(ret.Config.TPDomain) == 0 {
		ret.Config.TPDomain = consts.TPDomain
	}
	ret.Method = consts.MethodRefundQuery
	ret.Timestamp = fmt.Sprintf("%d", time.Now().Unix())
	ret.bizContent = simplejson.New()
	return ret
}

// 将Request编码成POST请求的Body
func (req *RefundQueryRequest) Encode() (string, error) {
	// 加签
	req.bizContent.Set("out_refund_no", req.OutRefundNo)
	req.bizContent.Set("refund_no", req.RefundNo)
	req.bizContent.Set("merchant_id", req.Config.MerchantId)
	req.bizContent.Set("uid", req.Uid)

	bizContentBytes, err := req.bizContent.Encode()
	if err != nil {
		util.Debug("RefundQueryRequest Encode bizContent.Encode err: %s, bizContent %s\n", err, *req.bizContent)
		return "", util.Wrap(err, "RefundQueryRequest Encode failed when [bizContent.Encode()]")
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
	//序列化
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
// refund_no 和 out_refund_no 哪个不空用哪个，都不空优先用out_refund_no
func (req *RefundQueryRequest) GetLogId() string {
	id := ""
	if len(req.RefundNo) != 0 {
		id = req.RefundNo
	}
	if len(req.OutRefundNo) != 0 {
		id = req.OutRefundNo
	}
	return fmt.Sprintf("%s_%s_%s_%s", req.Config.AppId, req.Config.MerchantId, id, req.Timestamp)
}

// 获取请求url地址
func (req *RefundQueryRequest) GetUrl() string {
	return req.Config.TPDomain + "/" + req.path
}

// 提供该接口，方便业务方设置可选参数，比如product_code、payment_type等
func (req *RefundQueryRequest) SetBizContentKV(key string, val interface{}) {
	req.bizContent.Set(key, val)
}

// 退款查询接口响应
type RefundQueryResponse struct {
	Data         *simplejson.Json
	OutRefundNo  string `json:"out_refund_no"`
	RefundNo     string `json:"refund_no"`
	TradeNo      string `json:"trade_no"`
	RefundAmount string `json:"refund_amount"`
	RefundStatus string `json:"refund_status"`
	ChannelExt   string `json:"channel_ext"`
}

// 初始化退款查询响应
func NewRefundQueryResponse() *RefundQueryResponse {
	ret := new(RefundQueryResponse)
	ret.Data = simplejson.New()
	return ret
}

// 将响应json数据反序列化为对应接口
func (resp *RefundQueryResponse) Decode() error {
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
func (resp *RefundQueryResponse) SetData(data *simplejson.Json) {
	resp.Data = data
}

// 目前只查验大写字母开头的参数(用户必传参数)
func (req *RefundQueryRequest) checkParams() error {
	if req.Method != consts.MethodRefundQuery {
		return fmt.Errorf(util.ErrorFormat, "Method", "must be tp.refund.query")
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
	if req.OutRefundNo == "" && req.RefundNo == "" {
		return errors.New("OurRefundNo and RefundNo can't both be blank")
	}

	if req.OutRefundNo != "" {
		if err := util.CheckOutRefundNo(req.OutRefundNo); err != nil {
			return err
		}
	}

	if req.RefundNo != "" {
		if err := util.CheckRefundNo(req.RefundNo); err != nil {
			return err
		}
	}

	return nil
}
