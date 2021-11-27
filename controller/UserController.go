package controller

import (
	"log"
	"net/http"
	"tricyzhou.com/ginessential/dto"
	"tricyzhou.com/ginessential/response"

	"golang.org/x/crypto/bcrypt"

	"tricyzhou.com/ginessential/util"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"tricyzhou.com/ginessential/common"
	"tricyzhou.com/ginessential/model"
)

func Register(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser = model.User{}
	//json.NewDecoder(ctx.Request.Body).Decode(&requestUser)
	ctx.Bind(&requestUser)
	name := requestUser.Name
	tel := requestUser.Telephone
	passwd := requestUser.Password

	// 手机号11位
	if len(tel) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位1")
		return
	}

	// 密码不少于6位
	if len(passwd) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	// 没传name随机生成
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	// 电话号码是否存在
	if isTelephoneExist(DB, tel) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	// 新建用户
	newUser := model.User{
		Name:      name,
		Telephone: tel,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)

	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}

	response.Success(ctx, gin.H{"token": token}, "注册成功")
}

func Login(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser = model.User{}
	//json.NewDecoder(ctx.Request.Body).Decode(&requestUser)
	ctx.Bind(&requestUser)
	tel := requestUser.Telephone
	passwd := requestUser.Password

	// 手机号11位
	if len(tel) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位2")
		return
	}

	// 密码不少于6位
	if len(passwd) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	// 判断手机号是否存在
	var user model.User
	DB.Where("telephone = ?", tel).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}
	// 判断密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwd)); err != nil {
		response.Fail(ctx, nil, "密码错误")
		return
	}
	// 发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}
	response.Success(ctx, gin.H{"token": token}, "登陆成功")

}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}

func isTelephoneExist(db *gorm.DB, tel string) bool {
	var user model.User
	db.Where("telephone = ?", tel).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
