package util

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"regexp"
)

// 检查各种URL统一使用urlRegExp

const (
	// 第一个占位符为变量名占位，第二个占位符为错误信息占位
	ErrorFormat = "invalid param: %s %s"

	// TODO: 验证参数长度限制
	RegexpAppid       = "^[0-9a-zA-Z-_]{1,32}$"
	RegexpMerchantid  = "^[0-9a-zA-Z-_]{1,32}$"
	RegexpUid         = "^[0-9a-zA-Z-_]{1,32}$"
	RegexpSigntype    = "^[0-9a-zA-Z]{1,10}$"
	RegexpVersion     = "^[0-9].[0-9]$"
	RegexpTimestamp   = "^[0-9]{1,19}$"
	RegexpOutRefundno = "^[a-zA-Z0-9-_]{1,32}$"
	RegexpRefundno    = "^[a-zA-Z0-9-_]{1,64}$"
	RegexpUrl         = `^(.*):(.*)$|(https|http):\/\/[-A-Za-z0-9+&@#\/%?=~_|!:,.;]+[-A-Za-z0-9+&@#\/%=~_|]`
	RegexpJson        = `^{(".*":.*,)*".*":.*}$`
	RegexpOutorderno  = `^[0-9a-zA-Z-_]{1,32}$`
	RegexpTradeno     = `^[a-zA-Z0-9-_]{1,64}$`
	RegexpTotalamount = `^[1-9][0-9]*$`

	MsgUrl      = "can only be one of the tree forms: rpc: [abc:abc]; http: [http://www.bytedance.com]; https: [https://www.bytedance.com]"
	MsgId       = "can only contains digits, lettes and special characters including '-' and '_'"
	MsgJson     = "must be a valid json string"
	MsgNumber   = "can only contains digits"
	MsgRequired = "is required"
	MsgVersion  = "must be the form of: a.b, eg. 1.0"
	MsgInteger  = "must be a positive number"
	MsgNonnil   = "must be non-nil"
)

var (
	appIdRegexp       *regexp.Regexp
	merchantIdRegexp  *regexp.Regexp
	uidRegExp         *regexp.Regexp
	signTypeRegExp    *regexp.Regexp
	versionRegExp     *regexp.Regexp
	timeStampRegExp   *regexp.Regexp
	outRefundNoRegExp *regexp.Regexp
	refundNoRegExp    *regexp.Regexp
	urlRegExp         *regexp.Regexp
	jsonRegExp        *regexp.Regexp
	outOrderNoRegExp  *regexp.Regexp
	tradeNoRegExp     *regexp.Regexp
	totalAmountRegExp *regexp.Regexp
)

func init() {
	appIdRegexp = regexp.MustCompile(RegexpAppid)
	merchantIdRegexp = regexp.MustCompile(RegexpMerchantid)
	uidRegExp = regexp.MustCompile(RegexpUid)
	signTypeRegExp = regexp.MustCompile(RegexpSigntype)
	versionRegExp = regexp.MustCompile(RegexpVersion)
	timeStampRegExp = regexp.MustCompile(RegexpTimestamp)
	outRefundNoRegExp = regexp.MustCompile(RegexpOutRefundno)
	refundNoRegExp = regexp.MustCompile(RegexpRefundno)
	urlRegExp = regexp.MustCompile(RegexpUrl)
	jsonRegExp = regexp.MustCompile(RegexpJson)
	outOrderNoRegExp = regexp.MustCompile(RegexpOutorderno)
	tradeNoRegExp = regexp.MustCompile(RegexpTradeno)
	totalAmountRegExp = regexp.MustCompile(RegexpTotalamount)
}

// 检查小程序版本，枚举值："1.0", "2.0", "2.0+"
func CheckAppletVersion(version string) error {
	if version != "1.0" && version != "2.0" && version != "2.0+" {
		return fmt.Errorf(ErrorFormat, "AppletVersion", "AppletVersion can only be '1.0', '2.0' or '2.0+'")
	}
	return nil
}

// CheckAppId 检查AppId是否有效
func CheckAppId(appId string) error {
	isMatch := appIdRegexp.MatchString(appId)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "AppId", MsgId)
	}
	return nil
}

func CheckMerchantId(merchantId string) error {
	isMatch := merchantIdRegexp.MatchString(merchantId)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "MerchantId", MsgId)
	}
	return nil
}

func CheckAppSecret(appSecret string) error {
	if len(appSecret) == 0 {
		return fmt.Errorf(ErrorFormat, "AppSecret", MsgRequired)
	}
	return nil
}

func CheckUid(uid string) error {
	isMatch := uidRegExp.MatchString(uid)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "Uid", MsgId)
	}
	return nil
}

func CheckBizContent(bizContent *simplejson.Json) error {
	if bizContent == nil {
		return fmt.Errorf(ErrorFormat, "bizContent", MsgNonnil)

	}
	return nil
}

func CheckSignType(signType string) error {
	if signType != "MD5" {
		return fmt.Errorf(ErrorFormat, "SignType", "Only MD5 is supported in current version")
	}
	return nil
}

func CheckFormat(format string) error {
	if format != "JSON" {
		return fmt.Errorf(ErrorFormat, "Format", "Only JSON is supported in current verison")
	}
	return nil
}

func CheckCharset(charset string) error {
	if charset != "utf-8" {
		return fmt.Errorf(ErrorFormat, "Charset", "Only utf-8 is supported in current verison")
	}
	return nil
}

func CheckVersion(version string) error {
	isMatch := versionRegExp.MatchString(version)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "Version", MsgVersion)
	}
	return nil
}

func CheckTimeStamp(timestamp string) error {
	isMatch := timeStampRegExp.MatchString(timestamp)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "Timestamp", MsgNumber)
	}
	return nil
}

func CheckTradeTime(tradeTime string) error {
	isMatch := timeStampRegExp.MatchString(tradeTime)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "TradeTime", MsgNumber)
	}
	return nil
}

func CheckValidTime(validTime string) error {
	isMatch := timeStampRegExp.MatchString(validTime)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "ValidTime", MsgNumber)
	}
	return nil
}

func CheckOutRefundNo(outRefundNo string) error {
	isMatch := outRefundNoRegExp.MatchString(outRefundNo)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "OutRefundNo", MsgId)
	}
	return nil
}

func CheckRefundNo(refundNo string) error {
	isMatch := refundNoRegExp.MatchString(refundNo)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "RefundNo", MsgId)
	}
	return nil
}

func CheckNotifyUrl(url string) error {
	isMatch := urlRegExp.MatchString(url)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "NotifyUrl", MsgUrl)
	}
	return nil
}

func CheckReturnUrl(returnUrl string) error {
	isMatch := urlRegExp.MatchString(returnUrl)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "ReturnUrl", MsgUrl)
	}
	return nil
}

func CheckRefundAmount(num int) error {
	if num <= 0 {
		return fmt.Errorf(ErrorFormat, "RefundAmount", MsgInteger)
	}
	return nil
}

func CheckRiskInfo(riskInfo string) error {
	isMatch := jsonRegExp.MatchString(riskInfo)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "RiskInfo", MsgJson)
	}
	return nil
}

func CheckExtParam(extParam string) error {
	isMatch := jsonRegExp.MatchString(extParam)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "ExtParam", MsgJson)
	}
	return nil
}

func CheckExt(ext string) error {
	isMatch := jsonRegExp.MatchString(ext)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "Ext", MsgJson)
	}
	return nil
}

func CheckSettlementExt(settlementExt string) error {
	isMatch := jsonRegExp.MatchString(settlementExt)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "SettlementExt", MsgJson)
	}
	return nil
}

func CheckParamsForApplet(params string) error {
	isMatch := jsonRegExp.MatchString(params)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "Params", MsgJson)
	}
	return nil
}

func CheckLimitPay(limitPay string) error {
	isMatch := jsonRegExp.MatchString(limitPay)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "LimitPay", MsgJson)
	}
	return nil
}

func CheckCashdeskExts(cashdeskExts string) error {
	isMatch := jsonRegExp.MatchString(cashdeskExts)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "CashdeskExts", MsgJson)
	}
	return nil
}

func CheckOutOrderNo(outOrderNo string) error {
	isMatch := outOrderNoRegExp.MatchString(outOrderNo)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "OutOrderNo", MsgId)
	}
	return nil
}

func CheckTradeNo(tradeNo string) error {
	isMatch := tradeNoRegExp.MatchString(tradeNo)
	if !isMatch {
		return fmt.Errorf(ErrorFormat, "TradeNo", MsgId)
	}
	return nil
}

func CheckTotalAmount(totalAmount int) error {
	if totalAmount <= 0 {
		return fmt.Errorf(ErrorFormat, "TotalAmount", MsgInteger)
	}
	return nil
}

func CheckUidType(uidType string) error {
	if len(uidType) == 0 {
		return errors.New("invalid param uidType")
	}
	return nil
}

func CheckCurrency(currency string) error {
	if currency == "" {
		return fmt.Errorf(ErrorFormat, "Currency", MsgRequired)
	}
	return nil
}

func CheckSubject(subject string) error {
	if subject == "" {
		return fmt.Errorf(ErrorFormat, "Subject", MsgRequired)
	}
	return nil
}

func CheckBody(body string) error {
	if body == "" {
		return fmt.Errorf(ErrorFormat, "Body", MsgRequired)
	}
	return nil
}

func CheckProductCode(productCode string) error {
	if len(productCode) == 0 {
		return errors.New("invalid param: ProductCode")
	}
	return nil
}

func CheckPaymentType(paymentType string) error {
	if len(paymentType) == 0 {
		return errors.New("invalid param： PaymentType")
	}
	return nil
}

func CheckServiceFee(serviceFee string) error {
	if len(serviceFee) == 0 {
		return errors.New("invalid param: ServiceFee")
	}
	return nil
}

func CheckSettlementProductCode(settlementProductCode string) error {
	if len(settlementProductCode) == 0 {
		return errors.New("invalid param SettlementProductCode")
	}
	return nil
}

func CheckSellerMerchantId(sellerMerchantId string) error {
	if len(sellerMerchantId) == 0 {
		return errors.New("invalid param SellerMerchantId")
	}
	return nil
}

func CheckRoyaltyParameters(royaltyParameters string) error {
	if len(royaltyParameters) == 0 {
		return errors.New("invalid param: RoyaltyParameters")
	}
	return nil
}

func CheckTransCode(transCode string) error {
	if len(transCode) == 0 {
		return errors.New("invalid param: TransCode")
	}
	return nil
}

func CheckCashDeskTradeType(tradeType string) error {
	if len(tradeType) == 0 {
		return errors.New("invalid param: CashDeskTradeType")
	}
	return nil
}

func CheckPayChannel(payChannel string) error {
	if len(payChannel) == 0 {
		return errors.New("invalid param: PayChannel")
	}
	return nil
}

func CheckPayType(payType string) error {
	if len(payType) == 0 {
		return errors.New("invalid param: PayType")
	}
	return nil
}

func CheckOutTradeNo(outTradeNo string) error {
	if len(outTradeNo) == 0 {
		return errors.New("invalid patam: OutTradeNo")
	}
	return nil
}

func CheckTradeName(tradeName string) error {
	if len(tradeName) == 0 {
		return errors.New("invalid patam: TradeName")
	}
	return nil
}

func CheckTradeDesc(tradeDesc string) error {
	if len(tradeDesc) == 0 {
		return errors.New("invalid patam: TradeDesc")
	}
	return nil
}
