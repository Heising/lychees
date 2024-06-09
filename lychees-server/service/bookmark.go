package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"lychees-server/dao"
	"lychees-server/logs"
	"lychees-server/models"
	"lychees-server/utils"
	"net/http"
)

func UpdateIconfont(ctx *gin.Context) {
	var newIconfont models.IconfontLink
	var err error
	if err = ctx.ShouldBindJSON(&newIconfont); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "字段错误"})
		return
	}

	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	var existingUser = &models.User{ID: userID}

	updateAt, err := dao.UpdateIconfont(existingUser, &newIconfont)
	if err != nil {

		utils.ServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"updateAt": updateAt})
}
func UpdatePersonalInfo(ctx *gin.Context) {
	var newPersonalInfo models.PersonalInfo
	var err error
	if err = ctx.ShouldBindJSON(&newPersonalInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "字段错误"})
		return
	}

	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.ServerError(ctx)
		return
	}

	var existingUser = &models.User{ID: userID}

	updateAt, err := dao.UpdatePersonalInfo(existingUser, &newPersonalInfo)
	if err != nil {

		utils.ServerError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"updateAt": updateAt})
}

// 查找书签
func FindDocuments(ctx *gin.Context) {

	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.ServerError(ctx)
	}
	updateAt := ctx.Param("updateAt")
	if !utils.CheckDigit(updateAt) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}
	if updateAt != "0" && dao.CheckBookmarkStatus(userID, updateAt) {
		ctx.Status(http.StatusNoContent)
		return
	}
	find, err := dao.FindDocuments(userID)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "未找到书签数据",
				"data":  find,
			})
			return
		}
		logs.Logger.Fatal(err)
		return
	}
	//if err != nil {
	//	utils.Warning(ctx)
	//	return
	//}

	ctx.JSON(http.StatusOK, gin.H{
		"data": find,
	})

}

// 添加书签
func PushSubDocument(ctx *gin.Context) {
	index := ctx.Param("index") // 使用Param方法获取路径中的参数
	if !utils.CheckDigit(index) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "非法参数！",
		})
		return
	}
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.ServerError(ctx)
	}

	var newBookMarkData models.Item
	if err := ctx.ShouldBindJSON(&newBookMarkData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "字段错误", "data": newBookMarkData})
		logs.Logger.Info(err)
		return
	}

	if newBookMarkData.Title == "" || !newBookMarkData.IsSvg && newBookMarkData.Icon == "" {
		logs.Logger.Info("尝试抓取网站信息")
		utils.GetFaviconTitle(&newBookMarkData)
	} else {
		logs.Logger.Info("不需要抓取网站信息")
	}

	updateAt, err := dao.PushSubDocument(userID, index, newBookMarkData)
	if err != nil {
		logs.Logger.Info(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "添加失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"updateAt": updateAt,
		"data":     newBookMarkData,
	})

}

// 更新文档嵌套的书签
func UpdateSubDocument(ctx *gin.Context) {
	index := ctx.Param("index")   // 使用Param方法获取路径中的参数
	target := ctx.Param("target") // 使用Param方法获取路径中的参数
	if !utils.CheckDigit(index) && !utils.CheckDigit(target) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "非法参数！",
		})
		return
	}
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.ServerError(ctx)
	}
	var newBookMarkData models.Item
	if err := ctx.ShouldBindJSON(&newBookMarkData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "字段错误"})
		return
	}
	if newBookMarkData.Title == "" || !newBookMarkData.IsSvg && newBookMarkData.Icon == "" {
		logs.Logger.Info("尝试抓取网站信息")
		utils.GetFaviconTitle(&newBookMarkData)
	} else {
		logs.Logger.Info("不需要抓取网站信息")
	}

	updateAt, err := dao.UpdateSubDocument(userID, index, target, newBookMarkData)
	if err != nil {
		logs.Logger.Info(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "添加失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"updateAt": updateAt,
		"data":     newBookMarkData,
	})

}

// 移动文档嵌套的书签
func MoveSubDocument(ctx *gin.Context) {
	index := ctx.Param("index")   // 使用Param方法获取路径中的参数
	target := ctx.Param("target") // 使用Param方法获取路径中的参数
	if !utils.CheckDigit(index) && !utils.CheckDigit(target) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "非法参数！",
		})
		return
	}
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.ServerError(ctx)
	}
	var newBookMarkData models.MoveItem
	if err := ctx.ShouldBindJSON(&newBookMarkData); err != nil {
		logs.Logger.Info(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "字段错误"})
		return
	}

	if newBookMarkData.Item.Title == "" || !newBookMarkData.Item.IsSvg && newBookMarkData.Item.Icon == "" {
		logs.Logger.Info("尝试抓取网站信息")
		utils.GetFaviconTitle(&newBookMarkData.Item)
	} else {
		logs.Logger.Info("不需要抓取网站信息")
	}
	updateAt, err := dao.MoveSubDocument(userID, index, target, newBookMarkData)
	if err != nil {
		logs.Logger.Info(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "添加失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"updateAt": updateAt,
		"data":     newBookMarkData.Item,
	})

}

// 移除文档嵌套的书签
func RemoveSubDocument(ctx *gin.Context) {
	rowIndex := ctx.Param("rowIndex") // 使用Param方法获取路径中的参数
	colIndex := ctx.Param("colIndex") // 使用Param方法获取路径中的参数
	if !utils.CheckDigit(rowIndex) && !utils.CheckDigit(colIndex) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "非法参数！",
		})
		return
	}

	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.ServerError(ctx)
	}
	updateAt, err := dao.RemoveSubDocument(userID, rowIndex, colIndex)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "移除失败！",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"updateAt": updateAt,
	})

}

// 清除页面
func ClearEmptyArray(ctx *gin.Context) {
	rowIndexs := make([]uint, 1)

	if err := ctx.ShouldBindJSON(&rowIndexs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "字段错误", "data": rowIndexs})
		logs.Logger.Info(err)
		return
	}

	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.ServerError(ctx)
	}
	updateAt, err := dao.ClearEmptyArray(userID, rowIndexs)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "移除失败！",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"updateAt": updateAt,
	})
}

// 交换书签
func SwapSubDocument(ctx *gin.Context) {
	index := ctx.Param("index") // 使用Param方法获取路径中的参数

	if !utils.CheckDigit(index) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "非法参数！",
		})
		return
	}
	userID, err := utils.GetUserID(ctx)
	if err != nil {
		utils.ServerError(ctx)
	}

	swapBookMarkData := make(map[uint]models.Item)

	if err := ctx.ShouldBindJSON(&swapBookMarkData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "字段错误", "data": swapBookMarkData})
		logs.Logger.Info(err)
		return
	}

	updateAt, err := dao.SwapSubDocument(userID, index, swapBookMarkData)
	if err != nil {
		logs.Logger.Info(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "交换失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"updateAt": updateAt,
	})

}
