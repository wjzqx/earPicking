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
)

/**
 * QueryData 查询公共方法
 * @param sql 执行sql
 * @return DbWorker
 */
func (dw *DbWorker) QueryData(sql string, args ...interface{}) *DbWorker {
	db := openDb(dw)
	rows, err := db.Query(sql, args...)
	checkErr(err)
	defer db.Close()

	var err1 = dw.formatRes(rows)
	checkErr(err1)
	return dw
}

/**
  * 设置查询的表名 SetTableName
  * @param tableName 表名
  * @return DbWorker
  */
func (dw *DbWorker) SetTableName(tableName string) *DbWorker{
	dw.tableInfo.tabName = tableName
	return dw
}

/**
  * 设置查询条件 Where
  * @param s 条件参数
  * @return DbWorker
  */
func (dw *DbWorker) Where(s string) *DbWorker{
	dw.whereTemp = s
	return dw
}

/**
  * 设置分组聚合字段 GroupBy
  * @param s 条件参数
  * @return DbWorker
  */
func (dw *DbWorker)GroupBy(s string) *DbWorker{
	dw.groupByTemp = s
	return dw
}

/**
 * 设置排序字段 OrderBy
 * @param col 表字段
 * @param s   条件参数
 * @return DbWorker
 */
func (dw *DbWorker)OrderBy(col string, s string) *DbWorker{
	if s == SQL_OB_DESC {
		dw.orderByTemp = col + " " + s
	}else {
		dw.orderByTemp = col + " " + SQL_OB_ASC
	}
	return dw
}


/**
 * 设置分页字段 Limit
 * @param s   条件参数
 * @return DbWorker
 */
func (dw *DbWorker)Limit(s string) *DbWorker{
	dw.limitTemp = s
	return dw
}

/**
 * 查询语句 Select
 * @param in 条件参数
 * @return DbWorker
 */
func (dw *DbWorker) Select(in interface{}) (err error){
	v := reflect.ValueOf(in)
	cols, _ , _:= formatCols(v, STR_SELECT)

	dw.sqlTemp, err = dw.selectSql(cols)

	checkErr(err)
	if err != nil{
		return err
	}

	fmt.Printf("%+v\n", dw.sqlTemp)
	dw.QueryData(dw.sqlTemp).toObject(v)

	return nil
}

/**
 * 查询一条数据，返回map对象 QueryToMap
 * @return m    数据源
 *	       err  nil
 */
func (dw *DbWorker) QueryToMap()(m map[string]string, err error){
	dw.sqlTemp, err = dw.selectSql("*")
	checkErr(err)
	if err != nil{
		return nil, err
	}
	fmt.Printf("%+v\n", dw.sqlTemp)
	return dw.QueryData(dw.sqlTemp).ForMap(), nil
}

/**
 * 查询所有数据
 * @param in 转换数据对象
 * @return err nil
 */
func (dw *DbWorker) SelectAll(in interface{}) (err error){
	v := reflect.ValueOf(in)
	cols := formatColsList(v, STR_SELECT)
	dw.sqlTemp, err = dw.selectSql(cols)

	if err != nil{
		return err
	}

	fmt.Printf("%+v\n", dw.sqlTemp)
	dw.QueryData(dw.sqlTemp).toObjForList(v)

	return nil
}

/**
 * InsertData 新增数据公共方法
 * @param in 参数对象
 * @return code 0，执行失败，1,执行成功
 */
func (dw *DbWorker) InsertData(in interface{}) int {
	v := reflect.ValueOf(in)
	var err error
	cols, vals, _ := formatCols(v, STR_INSERT)
	dw.sqlTemp, err = dw.insertSql(cols,vals)
	checkErr(err)
	if err != nil{
		return 0
	}

	fmt.Printf("sql: %+v\n", dw.sqlTemp)
	code, err := dw.dbExec(dw.sqlTemp)

	checkErr(err)
	return code
}

/**
 * 执行增，删，改的sql语句
 * @param sql 数据库语句
 * @param args 参数
 * @return code 执行后成功行数
 */
func (dw *DbWorker) ExecDate(sql string, args ...interface{}) int {
	code, err := dw.dbExec(sql, args...)
	checkErr(err)
	return code
}

/**
 * ModifyData 修改数据公共方法
 * @param in 参数对象
 * @return code 0，执行失败，1,执行成功
 */
func (dw *DbWorker) ModifyData(in interface{}) int {
	code := 0

	v := reflect.ValueOf(in)
	var err error
	_, _, content := formatCols(v, STR_UPDATE)
	dw.sqlTemp, err = dw.updateSql(content)
	checkErr(err)
	if err != nil{
		return 0
	}
	fmt.Printf("sql: %+v\n", dw.sqlTemp)


	//	db := openDb(dw.Dsn)
	code, err = dw.dbExec(dw.sqlTemp)
	//checkErr(err)
	return code
}

/**
 * DeleteData 删除数据公共方法
 * @return code 0，执行失败，1,执行成功
 */
func (dw *DbWorker) DeleteData() int {

	var err error
	dw.sqlTemp, err = dw.deleteSql()
	fmt.Printf("sql: %+v\n", dw.sqlTemp)
	code, err := dw.dbExec(dw.sqlTemp)
	checkErr(err)
	return code
}

// checkErr 错误检测
func checkErr(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}
}

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

func formatCols(v reflect.Value,sqlType string)(tag string, colVal string,content string){
	t := v.Type()
	val := v.Elem()
	typ := t.Elem()

	// 检查结构体的数据类型
	if !val.IsValid() {
		//return ErrDateType
	}

	switch sqlType {
		case STR_SELECT:
			for i := 0; i < val.NumField(); i++ {
				// 获取字段注解
				tag =  tag + typ.Field(i).Tag.Get("col") + ","
			}
			tag = strings.TrimRight(tag, ",")
		case STR_INSERT, STR_UPDATE:
			for i := 0; i < val.NumField(); i++ {
				value := val.Field(i)
				var _colVal = value.Interface()
				kind := value.Kind()
				if _colVal != "" && _colVal != nil {
					//
					v, _ := OjbToString(kind, _colVal)
					if v != "0"{
						// 获取字段注解
						tag =  tag + typ.Field(i).Tag.Get("col") + ","
						colVal = colVal + v + ","
						c := typ.Field(i).Tag.Get("col") + "=" + v + ","
						content = content + c
					}

				}
			}
			tag = strings.TrimRight(tag, ",")
			colVal = strings.TrimRight(colVal, ",")
			content = strings.TrimRight(content, ",")
		}

	return tag, colVal,content
}

func formatColsList (v reflect.Value, sqlType string)  string{
	sliceValue := reflect.Indirect(v)
	// 获取切片集合中的类型
	sliceElementType := sliceValue.Type().Elem()

	newValue := reflect.New(sliceElementType)
	res,  _, _ := formatCols(newValue, sqlType)
	return res
}

/**
 * 拼装insert语句
 * @param cols 查询字段
 * @return string sql语句
 *         error
 */
func (ti *tableInfo) insertSql(cols string, val string) (string, error){
	var res = SQL_INSERT

	// 表名为必填项，没有设置表名则提示错误。
	if ti.tabName == ""{
		return "", ErrTableNameIsNull
	}

	// 组装INSERT语句
	res = strings.Replace(res, STR_COLNAME, cols, -1)
	res = strings.Replace(res, STR_COLVALUE, val, -1)
	res = strings.Replace(res, STR_TABLENAME, ti.tabName, -1)
	return res, nil
}

/**
 * 拼装delete语句
 * @param cols 查询字段
 * @return string sql语句
 *         error
 */
func (ti *tableInfo) deleteSql()(string, error){

	var res = SQL_DELETE

	// 表名为必填项，没有设置表名则提示错误。
	if ti.tabName == ""{
		return "", ErrTableNameIsNull
	}

	// 组装表名
	res = strings.Replace(res, STR_TABLENAME, ti.tabName, -1)

	// 组装查询条件
	if ti.whereTemp != ""{
		whereSql := strings.Replace(SQL_WHERE, STR_CONTENT, ti.whereTemp, -1)
		res = strings.Replace(res, STR_WHERE, whereSql, -1)
	}else {
		res = strings.Replace(res, STR_WHERE + " ", "", -1)
	}

	return res, nil
}

/**
 * 拼装update语句
 * @param cols 查询字段
 * @return string sql语句
 *         error
 */
func (ti *tableInfo) updateSql(content string)(string, error){
	var res = SQL_UPDATE
	// 表名为必填项，没有设置表名则提示错误。
	if ti.tabName == ""{
		return "", ErrTableNameIsNull
	}

	// 组装表名
	res = strings.Replace(res, STR_TABLENAME, ti.tabName, -1)
	res = strings.Replace(res, STR_CONTENT, content, -1)

	// 组装查询条件
	if ti.whereTemp != ""{
		whereSql := strings.Replace(SQL_WHERE, STR_CONTENT, ti.whereTemp, -1)
		res = strings.Replace(res, STR_WHERE, whereSql, -1)
	}else {
		res = strings.Replace(res, STR_WHERE + " ", "", -1)
	}

	return res, nil
}

/**
 * 拼装select语句
 * @param cols 查询字段
 * @return string sql语句
 *         error
 */
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

	// 设置分页条件
	if ti.limitTemp != "" {
		limitSql := strings.Replace(SQL_LIMIT, STR_CONTENT, ti.limitTemp, -1)
		res  = strings.Replace(res, STR_LIMIT, limitSql, -1)
	}else {
		res = strings.Replace(res, STR_LIMIT , "", -1)
	}

	return res, nil
}

/**
 * 执行增，删，改的业务方法
 * @param sql 数据库语句
 * @param args 参数
 * @return code 执行后成功行数
 */
func (dw *DbWorker) dbExec(sql string, args ...interface{}) (int, error) {
	db := openDb(dw)
	res, err := db.Exec(sql, args...)
	if err != nil {
		return -1, err
	}
	id, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	defer db.Close()
	return int(id), nil
}