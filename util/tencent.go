package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"strconv"
)

type AccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      int    `json:"errmsg"`
}

type OpenIdResp struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     int    `json:"errmsg"`
}

// Func: Get Mini App Access Token
func GetAcceessToken(appId string, appSecret string) (string, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var arg = &fasthttp.Args{}
	arg.Set("grant_type", "client_credential")
	arg.Set("appid", appId)
	arg.Set("secret", appSecret)

	b, _, err := SendKVReq(arg, "GET", "https://api.weixin.qq.com/cgi-bin/token", nil)
	if err != nil {
		Logger(ERROR_LEVEL, "Tencent", "Send Get Access Token Request Err:"+err.Error())
		return "", err
	}

	atr := AccessTokenResp{}
	err = json.Unmarshal(b, &atr)
	if err != nil {
		Logger(ERROR_LEVEL, "Tencent", "Decode Resp Err:"+err.Error())
		return "", err
	}

	if atr.ErrCode == 0 {
		return atr.AccessToken, nil
	} else {
		Logger(ERROR_LEVEL, "Tencent", "Tencent Get AccessToken Resp Err:"+strconv.Itoa(atr.ErrCode))
		return "", errors.New("Tencent Get AccessToken Resp Err:" + strconv.Itoa(atr.ErrCode))
	}
}

// Func: Get OpenId From Tencent
func GetOpenId(appId string, appSecret string, code string, encrypt string, iv string) (string, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var arg = &fasthttp.Args{}
	arg.Set("appid", appId)
	arg.Set("secret", appSecret)
	arg.Set("js_code", code)
	arg.Set("grant_type", "authorization_code")

	b, _, err := SendKVReq(arg, "GET", "https://api.weixin.qq.com/sns/jscode2session", nil)
	if err != nil {
		Logger(ERROR_LEVEL, "Tencent", "Send Get OpenId Request Err:"+err.Error())
		return "", err
	}

	oir := OpenIdResp{}
	err = json.Unmarshal(b, &oir)
	if oir.ErrCode == 0 {
		if oir.UnionId != "" {
			return oir.UnionId, nil
		} else {
			return DecryptUnionId(encrypt, iv, oir.SessionKey)
		}
	} else {
		Logger(ERROR_LEVEL, "Tencent", "Tencent Get OpenId Resp Err:"+strconv.Itoa(oir.ErrCode))
		return "", errors.New("Tencent Get OpenId Resp Err:" + strconv.Itoa(oir.ErrCode))
	}
}

// Func: Decrypt Union Id
func DecryptUnionId(encrypt string, iv string, key string) (string, error) {
	d_key, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", errors.New("Decode Key Error ")
	}

	d_iv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return "", errors.New("Decode IV Error ")
	}

	d_encrypt, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", errors.New("Decode Encrypt Data Error ")
	}

	block, err := aes.NewCipher(d_key)
	if err != nil {
		return "", err
	}

	decrypted := make([]byte, len(d_encrypt))

	cbc_decrypt := cipher.NewCBCDecrypter(block, d_iv)
	cbc_decrypt.CryptBlocks(decrypted, d_encrypt)

	decrypted = PKCS7UnPadding(decrypted, block.BlockSize())

	return string(decrypted), nil
}

// Func: PKCS UnPadding
func PKCS7UnPadding(plantText []byte, blockSize int) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

// Func: Get OpenId From Tencent
func WeChatLogin(appId, appSecret, code string) (res OpenIdResp, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var arg = &fasthttp.Args{}
	arg.Set("appid", appId)
	arg.Set("secret", appSecret)
	arg.Set("js_code", code)
	arg.Set("grant_type", "authorization_code")
	Logger(INFO_LEVEL, "Tencent", "code:"+code)
	b, _, err := SendKVReq(arg, "GET", "https://api.weixin.qq.com/sns/jscode2session", nil)
	if err != nil {
		Logger(ERROR_LEVEL, "Tencent", "Send Get OpenId Request Err:"+err.Error())
		return res, err
	}

	err = json.Unmarshal(b, &res)
	if res.ErrCode != 0 {
		Logger(ERROR_LEVEL, "Tencent", "Tencent Get OpenId Resp Err:"+strconv.Itoa(res.ErrCode))
		return res, errors.New("Tencent Get OpenId Resp Err:" + strconv.Itoa(res.ErrCode))
	}
	return res, err
}
