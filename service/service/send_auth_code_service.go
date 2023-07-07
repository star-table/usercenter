package service

import (
	"fmt"
	"strings"

	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/errs"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/core/nacos"
	"github.com/star-table/usercenter/pkg/util/json"
	"github.com/star-table/usercenter/pkg/util/rand"
	"github.com/star-table/usercenter/pkg/util/temp"
	"github.com/star-table/usercenter/service/domain"
	"github.com/star-table/usercenter/service/model/req"
	"github.com/star-table/usercenter/service/model/resp"
	"github.com/star-table/usercenter/service/model/vo"
	"github.com/star-table/usercenter/service/model/vo/msgvo"
	"upper.io/db.v3"
)

const defaultAuthCode = "000000"

func SendSMSLoginCode(phoneNumber string) (bool, errs.SystemErrorInfo) {
	return SendAuthCode(vo.SendAuthCodeReq{
		AuthType:    consts.AuthCodeTypeLogin,
		AddressType: consts.ContactAddressTypeMobile,
		Address:     phoneNumber,
	})
}

// SendAuthCode 发送验证码
func SendAuthCode(input vo.SendAuthCodeReq) (bool, errs.SystemErrorInfo) {
	addressType := input.AddressType
	authType := input.AuthType
	contactAddress := input.Address

	if addressType != consts.ContactAddressTypeMobile && addressType != consts.ContactAddressTypeEmail {
		return false, errs.NotSupportedContactAddressType
	}

	err := domain.CheckSMSLoginCodeFreezeTime(authType, addressType, contactAddress)
	if err != nil {
		return false, err
	}

	//如果不是注册，登录，绑定，该账号必须存在
	if authType != consts.AuthCodeTypeRegister && authType != consts.AuthCodeTypeLogin && authType != consts.AuthCodeTypeBind {
		_, dbErr := domain.GetUserByLoginNameOrMobileOrEmail("", contactAddress, contactAddress)
		if dbErr != nil {
			if dbErr == db.ErrNoMoreRows {
				return false, errs.NotBindAccountError
			}
			logger.Error(dbErr)
			return false, errs.MysqlOperateError
		}
	}

	//如果是注册或者绑定，该账号必须不存在
	if authType == consts.AuthCodeTypeRegister || authType == consts.AuthCodeTypeBind {
		_, dbErr := domain.GetUserByLoginNameOrMobileOrEmail("", contactAddress, contactAddress)
		if dbErr != nil {
			if dbErr != db.ErrNoMoreRows {
				logger.Error(dbErr)
				return false, errs.MysqlOperateError
			}
		} else {
			// 账号存在
			if addressType == consts.ContactAddressTypeMobile {
				if authType == consts.AuthCodeTypeRegister {
					return false, errs.MobileAlreadyRegister
				} else {
					return false, errs.MobileAlreadyBind
				}
			} else {
				return false, errs.EmailAlreadyRegister
			}
		}

	}

	authCode := defaultAuthCode
	if !IsInWhiteList(contactAddress) {
		authCode = rand.RandomVerifyCode(6)
		//异步发送
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Error(r)
					fmt.Printf("捕获到的错误：%s\n", r)
				}
			}()
			switch addressType {
			case consts.ContactAddressTypeMobile:
				// 如果是 `+86-` 前缀的，则去除该前缀
				sendSmsAddress := strings.Replace(contactAddress, "+86-", "", 1)
				err = sendSmsAuthCode(authType, sendSmsAddress, authCode)
				if err != nil {
					logger.Error(err)
				}
			case consts.ContactAddressTypeEmail:
				err = sendMailAuthCode(authType, contactAddress, authCode)
				if err != nil {
					logger.Error(err)
				}
			}
		}()
	}
	setFreezeErr := domain.SetSMSLoginCodeFreezeTime(authType, addressType, contactAddress, 1)
	if setFreezeErr != nil {
		//这里不要影响主流程
		logger.Error(setFreezeErr)
	}
	setLoginCode := domain.SetSMSLoginCode(authType, addressType, contactAddress, authCode)
	if setLoginCode != nil {
		//这里不要影响主流程
		logger.Error(setLoginCode)
	}
	return true, nil
}

func sendSmsAuthCode(authType int, mobile string, authCode string) errs.SystemErrorInfo {
	switch authType {
	case consts.AuthCodeTypeLogin:
		return SendSMSPost(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeLoginAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	case consts.AuthCodeTypeRegister:
		return SendSMSPost(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeRegisterAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	case consts.AuthCodeTypeResetPwd:
		return SendSMSPost(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeResetPwdAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	case consts.AuthCodeTypeRetrievePwd:
		return SendSMSPost(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeRetrievePwdAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	case consts.AuthCodeTypeBind:
		return SendSMSPost(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeBindAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	case consts.AuthCodeTypeUnBind:
		return SendSMSPost(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeUnBindAuthCode, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	case consts.AuthCodeTypeChangeSuperAdmin:
		return SendSMSPost(mobile, consts.SMSSignNameBeiJiXing, consts.SMSTemplateCodeChangeSuperAdmin, map[string]string{
			consts.SMSParamsNameCode: authCode,
		})
	}
	return errs.NotSupportedAuthCodeType
}

func SendSMSPost(mobile string, signName string, templateCode string, params map[string]string) errs.SystemErrorInfo {
	_, _, err := nacos.DoPost("msgsvc", "api/msgsvc/sendSMS", map[string]interface{}{}, json.ToJsonIgnoreError(msgvo.SendSMSReqVoReqData{
		Mobile:       mobile,
		Params:       params,
		SignName:     signName,
		TemplateCode: templateCode,
	}))

	if err != nil {
		logger.Error(err)
		return errs.BuildSystemErrorInfo(errs.SystemError, err)
	}

	return nil
}

func sendMailAuthCode(authType int, email string, authCode string) errs.SystemErrorInfo {
	emails := []string{email}

	switch authType {
	case consts.AuthCodeTypeLogin:
		return SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeLogin, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode:   authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionLogin,
		}))
	case consts.AuthCodeTypeRegister:
		return SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeRegister, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode:   authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionRegister,
		}))
	case consts.AuthCodeTypeResetPwd:
		return SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeResetPwd, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode:   authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionResetPwd,
		}))
	case consts.AuthCodeTypeRetrievePwd:
		return SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeRetrievePwd, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode:   authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionRetrievePwd,
		}))
	case consts.AuthCodeTypeBind:
		return SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeBind, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode:   authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionBind,
		}))
	case consts.AuthCodeTypeUnBind:
		return SendMailRelaxed(emails, consts.MailTemplateSubjectAuthCodeUnBind, temp.RenderIgnoreError(consts.MailTemplateContentAuthCode, map[string]string{
			consts.SMSParamsNameCode:   authCode,
			consts.SMSParamsNameAction: consts.SMSAuthCodeActionUnBind,
		}))
	}
	return errs.NotSupportedAuthCodeType
}

func SendMailRelaxed(emails []string, subject string, content string) errs.SystemErrorInfo {
	_, _, err := nacos.DoPost("msgsvc", "api/msgsvc/sendMail", map[string]interface{}{}, json.ToJsonIgnoreError(msgvo.SendMailReqData{
		Emails:  emails,
		Subject: subject,
		Content: content,
	}))

	if err != nil {
		logger.Error(err)
		return errs.BuildSystemErrorInfo(errs.SystemError, err)
	}

	return nil
}

func IsInWhiteList(phoneNumber string) bool {
	whiteList, err := domain.GetPhoneNumberWhiteList()
	if err != nil {
		logger.Error(err)
		return false
	}
	for _, v := range whiteList {
		if v == phoneNumber {
			return true
		}
	}

	return false
}

func VerifyCaptcha(captchaID, captchaPassword *string, phoneNumber string) errs.SystemErrorInfo {
	if IsInWhiteList(phoneNumber) {
		return nil
	}
	if captchaID == nil || captchaPassword == nil {
		return errs.CaptchaError
	}

	res, err := domain.GetPwdLoginCode(*captchaID)
	if err != nil {
		logger.Error(err)
		return err
	}

	clearErr := domain.ClearPwdLoginCode(*captchaID)
	if clearErr != nil {
		logger.Error(clearErr)
		return clearErr
	}

	if res != *captchaPassword {
		return errs.CaptchaError
	}

	return nil
}

// AuthSmsCode 校验短信验证码是否正确，并返回临时 token
func AuthSmsCode(reqParam req.AuthSmsCodeReq) (*resp.AuthSmsCodeResp, errs.SystemErrorInfo) {
	if !domain.CheckAuthTypeValid(reqParam.AuthType) {
		return nil, errs.NotSupportedAuthCodeType
	}
	// 校验验证码是否正确
	if reqParam.CaptchaPassword == "" {
		return nil, errs.SMSLoginCodeNotMatch
	}
	err := domain.AuthCodeVerify(reqParam.AuthType, consts.ContactAddressTypeMobile, reqParam.PhoneNumber, reqParam.CaptchaPassword)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	// 生成临时 token
	tmpToken := rand.RandomVerifyCode(18)
	// 临时存储 token，校验时，提取 token 校验。
	err = domain.SetSMSLoginCode(reqParam.AuthType, consts.ContactAddressTypeMobile, reqParam.PhoneNumber, tmpToken)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	resp := &resp.AuthSmsCodeResp{
		Token:    tmpToken,
		AuthType: reqParam.AuthType,
	}

	return resp, nil
}
