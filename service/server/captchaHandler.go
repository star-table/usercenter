package server

import (
	"strconv"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/star-table/usercenter/core/consts"
	"github.com/star-table/usercenter/core/logger"
	"github.com/star-table/usercenter/service/domain"
)

type CaptchaResp struct {
	Code    int32           `json:"code"`
	Message string          `json:"message"`
	Data    CaptchaResponse `json:"data"`
}

type CaptchaResponse struct {
	CaptchaId string `json:"captchaId"` //验证码Id
	ImageUrl  string `json:"imageUrl"`  //验证码图片url
}

type RedisCache struct {
}

func (s *RedisCache) Set(id string, digits []byte) {
	pwd := ""
	for _, b := range digits {
		pwd += strconv.Itoa(int(b))
	}
	err := domain.SetPwdLoginCode(id, pwd)
	if err != nil {
		logger.Error(err)
	}
}

func (s *RedisCache) Get(id string, clear bool) (digits []byte) {
	resp, err := domain.GetPwdLoginCode(id)
	if err != nil {
		logger.Error(err)
	}
	res := []byte{}
	for _, i2 := range resp {
		v, _ := strconv.Atoi(string(i2))
		res = append(res, byte(v))
	}
	return res
}

func CaptchaGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//length := captcha.DefaultLen
		length := 4
		captchaId := captcha.NewLen(length)
		var captcha CaptchaResponse
		captcha.CaptchaId = captchaId
		captcha.ImageUrl = "/captcha/" + captchaId + ".png"
		c.JSON(200, CaptchaResp{Data: captcha, Code: 0, Message: "success"})
	}
}

func CaptchaShowHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		width := captcha.StdWidth
		widthParam := c.Query("w")
		if widthParam != consts.BlankString {
			mid, err := strconv.Atoi(widthParam)
			if err != nil {
				logger.Error(err)
			} else {
				width = mid
			}
		}

		height := captcha.StdHeight
		heightParam := c.Query("h")
		if heightParam != consts.BlankString {
			mid, err := strconv.Atoi(heightParam)
			if err != nil {
				logger.Error(err)
			} else {
				height = mid
			}
		}
		if width > 2000 || height > 2000 {
			c.JSON(500, "You are a bad boy!")
			return
		}
		captcha.Server(width, height).ServeHTTP(c.Writer, c.Request)
	}
}
