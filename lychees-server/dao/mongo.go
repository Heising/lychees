package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"lychees-server/configs"
	"lychees-server/logs"
	"lychees-server/models"
	"strconv"
	"time"
)

var mongoDB *mongo.Client
var collection *mongo.Collection

// 设置客户端链接配置
func initMongo() {
	clientOptions := options.Client().ApplyURI(
		"mongodb://" +
			configs.Config.Mongo.Host + ":" +
			strconv.Itoa(configs.Config.Mongo.Port))
	var err error
	// 链接MongoDB

	mongoDB, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logs.Logger.Fatal(err)
	} else {
		logs.Logger.Info("MongoDB连接成功")
	}
	// 检查链接
	err = mongoDB.Ping(context.Background(), nil)
	if err != nil {
		logs.Logger.Fatal(err)

	} else {
		logs.Logger.Info("MongoDB检测成功！")
	}

	collection = mongoDB.Database(configs.Config.Mongo.DBName).Collection("bookmarks")
}

// 给一个省内存？防止多次初始化
var ArrayBookmarks = [][]models.Item{}

// 初始化文档
func DocumentsInit(user *models.RegisterUser) (*mongo.InsertOneResult, error) {
	unixMilli := time.Now().UnixMilli()

	BookMark := models.Bookmarks{
		ID:       user.ID,
		UpdateAt: unixMilli,
		PersonalInfo: models.PersonalInfo{
			Nickname:     user.Nickname,
			Birthday:     user.Birthday,
			IconfontLink: models.IconfontLink{IconfontLink: user.IconfontLink.IconfontLink},
		},

		ArrayBookmarks: ArrayBookmarks,
	}

	insertOneResult, err := collection.InsertOne(context.Background(), &BookMark)
	if err != nil {
		logs.Logger.Fatal(err)
		return nil, err
	}
	logs.Logger.Info("文档添加成功")
	UpdateBookmarkStatus(fmt.Sprint(user.ID), unixMilli)

	return insertOneResult, err

}

// 添加书签
func PushSubDocument(documentID uint, index string, newBookMarkData models.Item) (updateAt int64, err error) {
	updateAt = time.Now().UnixMilli()

	// 更新BookMarks结构体中的数组元素
	filter := bson.M{"_id": documentID} // 匹配条件，假设_id为123
	//update := bson.M{
	//	"$push": bson.M{"arrayBookmarks." + index: newBookMarkData},
	//	"$pull": bson.M{"arrayBookmarks": nil},
	//	"$set":  bson.M{"updateAt": updateAt},
	//} // 使用$操作符更新特定数组元素的字段值
	//_, err = collection.UpdateOne(context.Background(), filter, bson.M{"$pull": bson.M{"ExternalIconfont": nil}})
	updatePush := bson.M{
		"$push": bson.M{"arrayBookmarks." + index: newBookMarkData},
	} // 使用$操作符更新特定数组元素的字段值
	updatePull := bson.M{
		"$pull": bson.M{"arrayBookmarks": nil},
	}
	updateSet := bson.M{
		"$set": bson.M{"updateAt": updateAt},
	}
	//_, err = collection.UpdateOne(context.Background(), filter, update)
	// 创建多个UpdateOneModel实例，每个实例代表一个更新操作
	BulkModels := []mongo.WriteModel{
		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updatePull),

		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updatePush),
		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updatePull),
		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updateSet),
		// 其他操作...
	}

	// 执行BulkWrite操作
	_, err = collection.BulkWrite(context.TODO(), BulkModels)
	if err != nil {
		logs.Logger.Error("添加失败:", err)
		return 0, err
	}

	logs.Logger.Info("更新成功")
	UpdateBookmarkStatus(strconv.Itoa(int(documentID)), updateAt)

	return updateAt, err
}

// 更新子文档中的特定字段
func UpdateSubDocument(documentID uint, index string, target string, newBookMarkData models.Item) (updateAt int64, err error) {
	updateAt = time.Now().UnixMilli()

	// 更新BookMarks结构体中的数组元素
	filter := bson.M{"_id": documentID} // 匹配条件，假设_id为123
	update := bson.M{
		"$set": bson.M{
			"updateAt":                               updateAt,
			"arrayBookmarks." + index + "." + target: newBookMarkData}}
	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		logs.Logger.Info("更新失败:", err)
		return
	}

	logs.Logger.Info("更新成功")
	UpdateBookmarkStatus(strconv.Itoa(int(documentID)), updateAt)
	return updateAt, err
}

// 移动子文档中的特定字段
func MoveSubDocument(documentID uint, index string, target string, newBookMarkData models.MoveItem) (updateAt int64, err error) {
	updateAt = time.Now().UnixMilli()

	// 更新BookMarks结构体中的数组元素
	// 更新BookMarks结构体中的数组元素
	filter := bson.M{"_id": documentID} // 匹配条件，假设_id为123
	updateUnset := bson.M{
		"$unset": bson.M{"arrayBookmarks." + index + "." + target: nil},
	} // 使用$操作符更新特定数组元素的字段值

	updatePull := bson.M{
		"$pull": bson.M{

			"arrayBookmarks." + index: nil,
			//"arrayBookmarks." + fmt.Sprint(newBookMarkData.NewRowIndex): nil,
		},
	}
	updateMove := bson.M{
		"$push": bson.M{
			"arrayBookmarks." + fmt.Sprint(newBookMarkData.NewRowIndex): bson.M{"$each": []models.Item{newBookMarkData.Item}, "$position": newBookMarkData.NewColIndex},
		},
	}

	updatePullArray := bson.M{
		"$pull": bson.M{
			"arrayBookmarks": nil,

			//"arrayBookmarks." + fmt.Sprint(newBookMarkData.NewRowIndex): nil,
		},
	}
	updateSet := bson.M{
		"$set": bson.M{"updateAt": updateAt},
	}
	// 创建多个UpdateOneModel实例，每个实例代表一个更新操作
	BulkModels := []mongo.WriteModel{
		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updateUnset),

		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updatePull),

		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updateMove),

		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updatePullArray),

		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updateSet),
		// 其他操作...
	}
	// 执行BulkWrite操作
	_, err = collection.BulkWrite(context.Background(), BulkModels)
	if err != nil {
		logs.Logger.Info("更新失败:", err)
		return
	}

	logs.Logger.Info("更新成功")
	UpdateBookmarkStatus(strconv.Itoa(int(documentID)), updateAt)
	return updateAt, err
}

// 移除子文档中的特定字段
func RemoveSubDocument(documentID uint, rowIndex string, colIndex string) (updateAt int64, err error) {
	updateAt = time.Now().UnixMilli()

	// 更新BookMarks结构体中的数组元素
	filter := bson.M{"_id": documentID} // 匹配条件，假设_id为123
	updateUnset := bson.M{
		"$unset": bson.M{"arrayBookmarks." + rowIndex + "." + colIndex: nil},
	} // 使用$操作符更新特定数组元素的字段值
	updatePull := bson.M{
		"$pull": bson.M{"arrayBookmarks." + rowIndex: nil},
	}
	updateSet := bson.M{
		"$set": bson.M{"updateAt": updateAt},
	}
	// 创建多个UpdateOneModel实例，每个实例代表一个更新操作
	BulkModels := []mongo.WriteModel{
		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updateUnset),
		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updatePull),
		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updateSet),
		// 其他操作...
	}

	// 执行BulkWrite操作
	_, err = collection.BulkWrite(context.Background(), BulkModels)

	if err != nil {
		logs.Logger.Errorf("移除失败:%s", err)
		return
	}

	logs.Logger.Info("移除成功")
	UpdateBookmarkStatus(strconv.Itoa(int(documentID)), updateAt)

	return updateAt, err
}

// 更新Iconfont
func UpdateIconfont(user *models.User, newIconfontLink *models.IconfontLink) (updateAt int64, err error) {
	updateAt = time.Now().UnixMilli()
	// 根据条件更新
	// 更新BookMarks结构体中的数组元素
	filter := bson.M{"_id": user.ID} // 匹配条件，假设_id为123
	update := bson.M{"$set": bson.M{
		"updateAt":                  updateAt,
		"personalinfo.iconfontlink": newIconfontLink}} // 使用$操作符更新特定数组元素的字段值

	_, err = collection.UpdateOne(context.Background(), filter, update)

	//if err != nil {
	//	logs.Logger.Fatal(err)
	//	return
	//}
	// 移除数组中索引为null的元素
	//_, err = collection.UpdateOne(context.Background(), filter, bson.M{"$pull": bson.M{"ExternalIconfont": nil}})

	if err != nil {
		logs.Logger.Fatal("更新Iconfont失败:", err)
		return
	}

	logs.Logger.Info("更新Iconfont成功")
	UpdateBookmarkStatus(strconv.Itoa(int(user.ID)), updateAt)

	return updateAt, err

}

// 更新用户昵称和生日
func UpdatePersonalInfo(user *models.User, newPersonalInfo *models.PersonalInfo) (updateAt int64, err error) {
	updateAt = time.Now().UnixMilli()
	// 根据条件更新
	// 更新BookMarks结构体中的数组元素
	filter := bson.M{"_id": user.ID} // 匹配条件，假设_id为123
	update := bson.M{"$set": bson.M{
		"updateAt":              updateAt,
		"personalinfo.nickname": newPersonalInfo.Nickname,
		"personalinfo.birthday": newPersonalInfo.Birthday,
	}} // 使用$操作符更新特定数组元素的字段值

	_, err = collection.UpdateOne(context.Background(), filter, update)

	//if err != nil {
	//	logs.Logger.Fatal(err)
	//	return
	//}
	// 移除数组中索引为null的元素
	//_, err = collection.UpdateOne(context.Background(), filter, bson.M{"$pull": bson.M{"ExternalIconfont": nil}})

	if err != nil {
		logs.Logger.Fatal("更新用户信息失败:", err)
		return
	}

	logs.Logger.Info("更新用户信息成功")
	UpdateBookmarkStatus(strconv.Itoa(int(user.ID)), updateAt)

	return updateAt, err

}

// FindDocuments 查找书签信息
func FindDocuments(id uint) (result *models.Bookmarks, err error) {
	//包装类型
	filter := bson.M{"_id": id}
	//初始化
	result = &models.Bookmarks{}
	err = collection.FindOne(context.Background(), filter).Decode(result)

	UpdateBookmarkStatus(strconv.Itoa(int(id)), result.UpdateAt)

	return result, err
}

func ClearEmptyArray(documentID uint, rowIndex []uint) (updateAt int64, err error) {
	updateAt = time.Now().UnixMilli()
	array := bson.M{}

	// Iterate over rowIndex slice
	for _, index := range rowIndex {
		array["arrayBookmarks."+fmt.Sprint(index)] = nil
	}

	// 更新BookMarks结构体中的数组元素
	filter := bson.M{"_id": documentID} // 匹配条件，假设_id为123
	updateUnset := bson.M{
		"$unset": array,
	} // 使用$操作符更新特定数组元素的字段值
	updatePull := bson.M{
		"$pull": bson.M{"arrayBookmarks": nil},
	}
	updateSet := bson.M{
		"$set": bson.M{"updateAt": updateAt},
	}
	// 创建多个UpdateOneModel实例，每个实例代表一个更新操作
	BulkModels := []mongo.WriteModel{
		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updateUnset),
		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updatePull),
		mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updateSet),
		// 其他操作...
	}

	// 执行BulkWrite操作
	_, err = collection.BulkWrite(context.TODO(), BulkModels)

	if err != nil {
		logs.Logger.Errorf("移除失败:%s", err)
		return
	}

	logs.Logger.Info("移除成功")
	UpdateBookmarkStatus(strconv.Itoa(int(documentID)), updateAt)

	return updateAt, err
}

func SwapSubDocument(documentID uint, index string, swapBookMarkData map[uint]models.Item) (updateAt int64, err error) {
	updateAt = time.Now().UnixMilli()
	array := bson.M{}
	for i := range swapBookMarkData {
		array["arrayBookmarks."+fmt.Sprint(index, ".", i)] = swapBookMarkData[i]
	}
	array["updateAt"] = updateAt
	// 更新BookMarks结构体中的数组元素
	filter := bson.M{"_id": documentID} // 匹配条件，假设_id为123
	update := bson.M{
		"$set": array,
	}
	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		logs.Logger.Fatal("交换失败:", err)
		return
	}

	logs.Logger.Info("交换成功")
	UpdateBookmarkStatus(strconv.Itoa(int(documentID)), updateAt)
	// 移除数组中索引为null的元素
	return updateAt, err
}
