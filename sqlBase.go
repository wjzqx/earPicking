package earPicking

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // ..
	"reflect"
	"strings"
)

type (
	// DbWorker ...
	DbWorker struct {
		DbDeploy  // 连接库连接信息
		queryRes  // 查询结果存储对象
		tableInfo // 表信息结构体
	}

	// ParamMap 查询条件参数
	ParamMap struct {
		Key   string      // 参数名称
		Value interface{} // 参数值
	}
)

// checkErr 错误检测
func checkErr(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
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

// SetTableName 设置查询的表名
func (dw *DbWorker) SetTableName(tableName string) *DbWorker{
	dw.tableInfo.tabName = tableName
	return dw
}

// 设置查询条件
func (dw *DbWorker) Where(s string) *DbWorker{
	dw.whereTemp = s
	return dw
}

// 设置分组聚合字段
func (dw *DbWorker)GroupBy(s string) *DbWorker{
	dw.groupByTemp = s
	return dw
}

// 设置排序字段
func (dw *DbWorker)OrderBy(col string, s string) *DbWorker{
	if s == SQL_OB_DESC {
		dw.orderByTemp = col + " " + s
	}else {
		dw.orderByTemp = col + " " + SQL_OB_ASC
	}
	return dw
}

// 查询语句
func (dw *DbWorker) Select(in interface{}) (err error){
	v := reflect.ValueOf(in)
	cols := formatCols(v)

	dw.sqlTemp, err = dw.selectSql(cols)

	if err != nil{
		return err
	}

	fmt.Printf("%+v\n", dw.sqlTemp)

	dw.QueryData(dw.sqlTemp).unqiue(v)


	return nil
}

/**
 * 查询所有数据
 */
func (dw *DbWorker) SelectAll(in interface{}) (err error){

	v := reflect.ValueOf(in)
	cols := formatColsList(v)

	dw.sqlTemp, err = dw.selectSql(cols)

	if err != nil{
		return err
	}

	fmt.Printf("%+v\n", dw.sqlTemp)

	dw.QueryData(dw.sqlTemp).list(v)


	return nil
}

func formatCols(v reflect.Value) string{

	var tag string

	t := v.Type()
	val := v.Elem()
	typ := t.Elem()

	// 检查结构体的数据类型
	if !val.IsValid() {
		//return ErrDateType
	}

	for i := 0; i < val.NumField(); i++ {
		// 获取字段注解
		tag =  tag + typ.Field(i).Tag.Get("col") + ","
	}
	tag = strings.TrimRight(tag, ",")

	return tag
}

func formatColsList (v reflect.Value)  string{
	sliceValue := reflect.Indirect(v)
	// 获取切片集合中的类型
	sliceElementType := sliceValue.Type().Elem()

	newValue := reflect.New(sliceElementType)
	res := formatCols(newValue)
	return res
}

func (ti *tableInfo) selectSql(cols string) (string,error){
	var res = SQL_SELECT

	// 表名为必填项，没有设置表名则提示错误。
	if ti.tabName == ""{
		return "", ErrTableNameIsNull
	}

	// 组装select语句
	res = strings.Replace(res, STR_COLS, cols, -1)
	res = strings.Replace(res, STR_TABLENAME, ti.tabName, -1)

	// 组装查询条件
	if ti.whereTemp != ""{
		whereSql := strings.Replace(SQL_WHERE, STR_CONTENT, ti.whereTemp, -1)
		res = strings.Replace(res, STR_WHERE, whereSql, -1)
	}else {
		res = strings.Replace(res, STR_WHERE + " ", "", -1)
	}

	// 组装分组聚合条件
	if ti.groupByTemp != ""{
		groupBySql := strings.Replace(SQL_GROUPBY, STR_CONTENT, ti.groupByTemp, -1)
		res  = strings.Replace(res, STR_GROUPBY, groupBySql, -1)
		//ti.sqlTemp = ti.sqlTemp + " " + strings.Replace(SQL_GROUPBY, STR_CONTENT, ti.groupByTemp, -1)
	}else {
		res = strings.Replace(res, STR_GROUPBY + " ", "", -1)
	}

	// 组装排序字段条件
	if ti.orderByTemp != "" {
		orderBySql := strings.Replace(SQL_ORDERBY, STR_CONTENT, ti.orderByTemp, -1)
		res  = strings.Replace(res, STR_ORDERBY, orderBySql, -1)
	}else {
		res = strings.Replace(res, STR_ORDERBY , "", -1)
	}

	return res, nil
}


// InsertData 新增数据公共方法
// @param sql 执行sql
func (dw *DbWorker) dbExec(sql string, args ...interface{}) (int, error) {
	db := openDb(dw)
	res, err := db.Exec(sql, args...)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	defer db.Close()
	return int(id), nil
}

// InsertData 新增数据公共方法
// @param sql 执行sql
func (dw *DbWorker) InsertData(sql string,  args ...interface{}) int {
	code, err := dw.dbExec(sql, args...)
	checkErr(err)
	return code
}

// ModifyData 修改数据公共方法
// @param sql 执行sql
func (dw *DbWorker) ModifyData(sql string, args ...interface{}) int {
	//	db := openDb(dw.Dsn)
	code, err := dw.dbExec(sql, args...)
	checkErr(err)
	return code
}

// DeleteData 删除数据公共方法
// @param sql 执行sql
func (dw *DbWorker) DeleteData(sql string, args ...interface{}) int {
	code, err := dw.dbExec(sql, args...)
	checkErr(err)
	return code
}



