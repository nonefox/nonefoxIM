package tools

import (
	"crypto/tls"
	"fmt"
	"github.com/google/uuid"
	"github.com/jordan-wright/email"
	"im/define"
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	"crypto/md5"
	"github.com/golang-jwt/jwt/v4"
)

// UserClaims 用户的声明结构（我们会对他进行签名生成用户的token信息）
type UserClaims struct {
	Identity string `json:"identity"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

//这是我们的key，我们会用她来签名
var myKey = []byte("im")

// GetMd5 把密码进行盐值加密处理
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// GenerateToken 通过用户的identity和email生成token
func GenerateToken(identity, email string) (string, error) {
	//定义一个需要签名的用户声明信息
	userClaim := UserClaims{
		Identity:         identity,
		Email:            email,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	//设置签名的加密算法，并且生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	//加上我们自己的关键字生成一个自己的JWTToken
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

// AnalyseToken 解析token并且把解析出来的用户声明返回出去
func AnalyseToken(tokenString string) (*UserClaims, error) {
	//定义一个用户声明，用来返回token中解析出来的用户声明信息
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(*jwt.Token) (interface{}, error) {
		//这个方法属于引用的包中的方法，用来得到我们的key，所以这里我们直接使用匿名方法
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims.Valid {
		return nil, fmt.Errorf("解析TOken出错：%v", err)
	}
	return userClaim, nil
}

// SendCode 通过邮箱发送验证码
func SendCode(userEmail string, code string) error {
	e := email.NewEmail() //new一个email对象
	//配置email的基本信息
	e.From = "Get <getUserEmail@163.com>" //用自己的邮箱
	e.To = []string{userEmail}
	e.Subject = "验证码已发送，请查收！"
	e.HTML = []byte("您的验证码：<b>" + code + "</b>")
	//设置163邮箱的中转协议，与邮箱密码（自己的邮箱密码）
	return e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "getUserEmail@163.com", define.MailPassword, "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
}

// GenerateCode 生成用来发送给用户邮箱中的验证码
func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	res := ""
	for i := 0; i < 6; i++ {
		res += strconv.Itoa(rand.Intn(10))
	}
	return res
}

// GetUUID 获取UUID作为我们用户的id使用
func GetUUID() string {
	u := uuid.New().String()
	if u == "" {
		log.Printf("生成uuid失败")
		return ""
	}
	return u
}
