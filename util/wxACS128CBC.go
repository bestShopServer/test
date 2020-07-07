package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// CBC 模式
//解密
/**
* rawData 原始加密数据
* key  密钥
* iv  向量
 */
func DncryptWx(rawData, key, iv string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	key_b, err_1 := base64.StdEncoding.DecodeString(key)
	iv_b, _ := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return "", err
	}
	if err_1 != nil {
		return "", err_1
	}
	dnData, err := AesCBCDncryptWx(data, key_b, iv_b)
	if err != nil {
		return "", err
	}
	return string(dnData), nil
}

// 解密
func AesCBCDncryptWx(encryptData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	blockSize := block.BlockSize()
	if len(encryptData) < blockSize {
		panic("ciphertext too short")
	}
	if len(encryptData)%blockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptData, encryptData)
	// 解填充
	encryptData, err = PKCS7UnPaddingWx(encryptData)
	return encryptData, err
}

//去除填充
func PKCS7UnPaddingWx(origData []byte) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])
	//defer Error("PKCS7UnPaddingWx 出现异常...")
	defer func() {
		if r := recover(); r != nil {
			Error("捕获到的错误：%s\n", r)
		}
	}()
	idx := length - unpadding
	if idx <= 0 {
		Error("数组越界:%v", idx)
		return nil, errors.New("数据处理失败")
	}
	return origData[:(idx)], nil
}
