package dao

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lychees-server/configs"
	"lychees-server/logs"
	"lychees-server/models"
	"strconv"
)

var postgresqlDB *gorm.DB

func initPostgresql() {
	dsn := "host=" + configs.Config.PostgreSQL.Host +
		" user=" + configs.Config.PostgreSQL.User +
		" password=" + configs.Config.PostgreSQL.Password +
		" dbname=" + configs.Config.PostgreSQL.DBName +
		" port=" + strconv.Itoa(configs.Config.PostgreSQL.Port) +
		" sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	postgresqlDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true, // 缓存预编译语句
	})
	if err != nil {
		logs.Logger.Fatal(err)
	}
	// 不会校验账号密码是否正确

	//自动迁移模型
	err = postgresqlDB.AutoMigrate(&models.User{})

	//调试所有sql
	//postgresqlDB = postgresqlDB.Debug()

	if err != nil {
		logs.Logger.Fatal(err)
	}
	// 设置默认DB对象
	// 检查数据库连接是否成功
	if postgresqlDB != nil {
		logs.Logger.Info("Postgresql数据库迁移成功")
	} else {
		logs.Logger.Fatal("Postgresql数据库连接失败")
	}

}

// 添加一个用户
func AddUser(user *models.User) error {
	err := postgresqlDB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

// 根据id查找存在的用户
func FindUser(user *models.User) (existingUser *models.User, err error) {
	// 初始化
	existingUser = &models.User{}

	err = postgresqlDB.Where("id = ?", user.ID).First(existingUser).Error
	return existingUser, err
}
func FindId(userId uint) (existingUser *models.User, err error) {
	// 初始化
	existingUser = &models.User{}

	err = postgresqlDB.Where("id = ?", userId).First(existingUser).Error
	return existingUser, err
}

// 重置密码
func PasswordReset(existingUser *models.User, newPassword string) (err error) {
	err = postgresqlDB.Model(existingUser).Update("encrypted", newPassword).Error
	return err
}

// 根据email查找存在的用户
func FindEmail(email *string) (existingUser *models.User, err error) {
	// 初始化
	existingUser = &models.User{}

	err = postgresqlDB.Where("email = ?", email).First(existingUser).Error

	return existingUser, err
}

// 根据email查找存在的用户 返回安全的字段
func FindFilterUser(user *models.User) (existingUser *models.ResponseUser, err error) {
	//var existingUser models.User
	// 在查询时，GORM 会自动选择 `id `, `name` 字段
	existingUser = &models.ResponseUser{}

	err = postgresqlDB.Model(&models.User{}).Where("email = ?", user.Email).First(existingUser).Error

	return existingUser, err
}

// 更新邮箱
func EmailUpdate(user *models.EmailUpdater) (err error) {

	return postgresqlDB.Model(&models.User{}).Where("email = ?", user.Email).Update("email", user.NewEmail).Error

}
