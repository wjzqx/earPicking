package earPicking

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" //..
)

// DbWorker ...
type DbWorker struct {
	DbDeploy // 连接库连接信息
	queryRes // 查询结果存储对象
}

// ParamMap 查询条件参数
type ParamMap struct {
	Key   string      // 参数名称
	Value interface{} // 参数值
}

// checkErr 错误检测
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// ParamList 查询条件参数集合
type ParamList []ParamMap

// openDb 打开数据库连接
func openDb(dbWorker *DbWorker) *sql.DB {
	//	dbWorker..deployDBInfo(&dbWorker.DbDeploy)
	err := dbWorker.DbDeploy.deployDBInfo()
	checkErr(err)
	//fmt.Printf("%+v\n", dbWorker)
	db, err := sql.Open("mysql", dbWorker.Dsn)
	checkErr(err)
	return db
}

// QueryData 查询公共方法
// @param sql 执行sql
func (dw *DbWorker) QueryData(sql string, args ...interface{}) *DbWorker {
	db := openDb(dw)
	rows, err := db.Query(sql, args...)
	checkErr(err)
	defer db.Close()

	var err1 = dw.formatRes(rows)
	checkErr(err1)
	return dw
}

// InsertData 新增数据公共方法
// @param sql 执行sql
func (dw *DbWorker) InsertData(sql string) int {
	db := openDb(dw)
	rows, err := db.Query(sql)
	fmt.Print(rows)
	checkErr(err)
	defer db.Close()

	return 0
}

// ModifyData 修改数据公共方法
// @param sql 执行sql
func (dw *DbWorker) ModifyData(sql string) int {
	//	db := openDb(dw.Dsn)
	return 0
}

// DeleteData 删除数据公共方法
// @param sql 执行sql
func (dw *DbWorker) DeleteData(sql string) int {
	return 0
}
