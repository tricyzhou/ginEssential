package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"tricyzhou.com/ginessential/model"
	"tricyzhou.com/ginessential/repository"
	"tricyzhou.com/ginessential/response"
	"tricyzhou.com/ginessential/vo"
)

type ICategoryController interface {
	RestController
}

type CategoryController struct {
	Repository repository.CategoryRepository
}

func NewCategoryController() ICategoryController{
	repository := repository.NewCategoryRepository()
	repository.DB.AutoMigrate(model.Category{})
	return CategoryController{Repository: repository}
}

func (c CategoryController) Create(ctx *gin.Context) {
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err!=nil{
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}

	category, err := c.Repository.Create(requestCategory.Name)
	if err!=nil{
		//response.Fail(ctx, nil, "创建失败")
		panic(err)
		return
	}
	response.Success(ctx, gin.H{"category": category}, "")
}

func (c CategoryController) Update(ctx *gin.Context) {
	// 绑定body中参数
	var requestCategory vo.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err!=nil{
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}
	// 获取path中参数
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	updateCategory, err := c.Repository.SelectById(categoryId)
	if err!=nil {
		response.Fail(ctx, nil, "分类不存在")
		return
	}
	// 更新分类
	// map
	// struct
	// name value
	category, err := c.Repository.Update(*updateCategory, requestCategory.Name)
	if err!=nil{
		panic(err)
	}
	response.Success(ctx, gin.H{"category": category}, "修改成功")
}

func (c CategoryController) Show(ctx *gin.Context) {
	// 获取path中参数
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	category, err := c.Repository.SelectById(categoryId)
	if err!=nil {
		response.Fail(ctx, nil, "分类不存在")
		return
	}
	response.Success(ctx, gin.H{"category": category}, "")
}

func (c CategoryController) Delete(ctx *gin.Context) {
	// 获取path中参数
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	if err:=c.Repository.DeleteById(categoryId); err!=nil {
		response.Fail(ctx, nil, "删除失败请重试")
		return
	}
	response.Success(ctx, nil, "")
}

