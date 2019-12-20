package earPicking

import (
	"errors"
	"fmt"
	"strings"
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
	ErrDataTypeNotExist     = errors.New("无数据库类型")      // ErrDataType 没有该数据库类型
	ErrDataTypeIsNull       = errors.New("数据库类型为空")     // ErrDataTypeIsNull 数据库类型为空
	ErrUserIsNull           = errors.New("用户名为空")       // ErrUserIsNull 用户名为空
	ErrPasswordIsNull       = errors.New("数据库密码为空")     // ErrPasswordIsNull 数据库密码为空
	ErrIPAddressIsNull      = errors.New("数据库访问IP地址为空") // ErrIPAddressIsNull 数据库访问IP地址为空
	ErrPortIsNull           = errors.New("数据库访问端口号为空")  // ErrPortIsNull 数据库访问端口号为空
	ErrDataNameIsNull       = errors.New("数据库名称为空")     // ErrDataNameIsNull 数据库名称为空
	ErrorNeedPointerToSlice = errors.New("需要指向切片的指针")   // ErrorNeedPointerToSlice 需要指向切片的指针
	ErrDateType             = errors.New("数据类型不正确")
	ErrTableNameIsNull      = errors.New("表名为空，请设置表名")
)

// 声明sql语句常量
var (
	SQL_SELECT  = "SELECT _cols_ from _tableName_ _WHERE_ _GROUPBY_ _ORDERBY_ "
	SQL_WHERE   = "WHERE _colContent_"
	SQL_INSERT  = "INSERT INTO _tableName_ (_colName_) VALUES (_colValue_)"
	SQL_UPDATE  = "UPDATE _tableName_ SET _colContent_"
	SQL_CONTENT = "_colName_ = _colValue_"
	SQL_GROUPBY = "GROUP BY _colContent_"
	SQL_ORDERBY = "ORDER BY _colContent_"
	SQL_OB_ASC = "ASC"
	SQL_OB_DESC = "DESC"

	STR_SELECT = "select"
	STR_INSERT = "insert"
	STR_UPDATE = "update"

	STR_WHERE   = "_WHERE_"
	STR_GROUPBY = "_GROUPBY_"
	STR_ORDERBY = "_ORDERBY_"


	STR_COLS      = "_cols_"
	STR_COLNAME   = "_colName_"
	STR_CONTENT   = "_colContent_"
	STR_COLVALUE  = "_colValue_"
	STR_TABLENAME = "_tableName_"

)

// deployDBInfo 设置数据库连接信息
func (dd *DbDeploy) deployDBInfo() (err error) {
	//  "root:123456@tcp(127.0.0.1:3306)/db_config"
	// 如果Dsn未被设置，则根据参数类型拼接Dsn
	if dd.Dsn == "" {
		err = checkDbDeploy(dd)
		if err != nil {
			return err
		}
		// 统一转换成小写，根据数据库类型设置dsn
		switch strings.ToLower(dd.DataType) {
		case "mysql":
			dd.Dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dd.User, dd.Password, dd.IPAddress, dd.Port, dd.DataName)
		default:
			err = ErrDataTypeNotExist
		}
	}
	return err
}

// checkDbDeploy 检查db配置参数是否符合生成dsn条件
func  checkDbDeploy(dd *DbDeploy) error {
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
