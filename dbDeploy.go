package earPicking

import (
	"errors"
	"fmt"
)

// DbDeploy 初始化数据库连接信息
type DbDeploy struct {
	//mysql data source name
	Dsn       string // 数据连接访问路径
	User      string // 数据库用户名
	Password  string // 数据库密码
	IPAddress string // 数据库访问IP地址
	Port      string // 数据库访问端口号
	DataType  string // 数据库类型
	DataName  string // 数据库名称
}

// 声明异常常量
var (
	ErrDataTypeNotExist = errors.New("无数据库类型")      // ErrDataType 没有该数据库类型
	ErrDataTypeIsNull   = errors.New("数据库类型为空")     // ErrDataTypeIsNull 数据库类型为空
	ErrUserIsNull       = errors.New("用户名为空")       // ErrUserIsNull 用户名为空
	ErrPasswordIsNull   = errors.New("数据库密码为空")     // ErrPasswordIsNull 数据库密码为空
	ErrIPAddressIsNull  = errors.New("数据库访问IP地址为空") // ErrIPAddressIsNull 数据库访问IP地址为空
	ErrPortIsNull       = errors.New("数据库访问端口号为空")  // ErrPortIsNull 数据库访问端口号为空
	ErrDataNameIsNull   = errors.New("数据库名称为空")     // ErrDataNameIsNull 数据库名称为空
)

// deployDBInfo 设置数据库连接信息
func (dd *DbDeploy) deployDBInfo() (err error) {
	//  "root:123456@tcp(127.0.0.1:3306)/db_config"
	// 如果Dsn未被设置，则根据参数类型拼接Dsn
	if dd.Dsn == "" {
		err = dd.checkDbDeploy()
		if err != nil {
			return err
		}
		switch dd.DataType {
		case "mySql":
			dd.Dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dd.User, dd.Password, dd.IPAddress, dd.Port, dd.DataName)
		default:
			err = ErrDataTypeNotExist
		}
	}
	return err
}

// checkDbDeploy 检查db配置参数是否符合生成dsn条件
func (dd *DbDeploy) checkDbDeploy() error {
	if dd.DataName == "" {
		return ErrDataNameIsNull
	}
	if dd.DataType == "" {
		return ErrDataTypeIsNull
	}
	if dd.User == "" {
		return ErrUserIsNull
	}
	if dd.Password == "" {
		return ErrPasswordIsNull
	}
	if dd.IPAddress == "" {
		return ErrIPAddressIsNull
	}
	if dd.Port == "" {
		return ErrPortIsNull
	}
	return nil
}
