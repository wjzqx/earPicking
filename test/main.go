package main

import (
	"ch2/earPicking"
	"fmt"
)

// Activity测试结构体
type Activity struct  {
	ID   int64  `col:"id" json:"id"`
	Name string `col:"dataName" json:"dataName"`
}

type Source struct {
	ID          int64  `col:"id" json:"id"`
	Name        string `col:"dataName" json:"dataName"`
	JdbcUrl     string `col:"jdbcUrl" json:"jdbcUrl"`
	DriverClass string `col:"driverClass" json:"driverClass"`
	User        string `col:"user" json:"user"`
	Password    string `col:"password" json:"password"`
	WriteOrRead int8   `col:"writeOrRead" json:"writeOrRead"`
	CreateTime  int64  `col:"createTime" json:"createTime"`
	Remake      string `col:"remake" json:"remake"`
	Status      int8   `col:"status" json:"status"`
}

func main() {


	var dbWorker earPicking.DbWorker


	//ret := Activity{}

	// 设置数据库连接信息
	dbWorker.User = "root"
	dbWorker.Password = "123456"
	dbWorker.IPAddress = "127.0.0.1"
	dbWorker.DataName = "db_config"
	dbWorker.Port = "3306"
	dbWorker.DataType = "mySql"
	//dbWorker.Dsn = "root:123456@tcp(127.0.0.1:3306)/db_config"

	// 查询数据
	//testSelect(dbWorker)

	// 查询数返回为MAP
	//testSelectToMap(dbWorker)

	// 条件查询
	//testWhereSelect(dbWorker)
	// 分组查询
	//testGroupBy(dbWorker)
	// 排序查询
	//testOrderByList(dbWorker)

	// 新增数据
	//testInster(dbWorker)
	// 删除数据
	//testDel(dbWorker)
	// 修改数据
	//testUpdate(dbWorker)


	m := dbWorker.QueryData("SELECT * FROM data_source WHERE id = ?", 1).ForMap()
	fmt.Printf("%+v\n", m)

	// Activity结构体数组
	//retList := []Activity{}

	//err := dbWorker.QueryData("SELECT * FROM data_source").List(&retList)
	//fmt.Printf("%+v\n", err)
	//fmt.Printf("%+v\n", retList)

	//code := dbWorker.ExecDate("INSERT INTO data_source (id, dataName, jdbcUrl, driverClass, user, password, writeOrRead, createTime, remake, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",5,"db_config111","jdbc:mysql://127.0.0.1:11306/db_config?useUnicode=true&characterEncoding=utf8","com.mysql.jdbc.Driver", "root", "123456", 1, 0, "", 1)
	//fmt.Printf("%+v\n", code)

	//code := dbWorker.DeleteData("DELETE FROM data_source WHERE id=?", 4)
	//fmt.Printf("%+v\n", code)
	//dbWorker.QueryData("SELECT * FROM data_source WHERE id = ? and dataName = ?", nil)



	// type db earPicking.DbWorker
	//util.RegistryType((*db)(nil))
	//
	//structName := "db"
	//
	//s, ok := util.NewStruct(structName)
	//if !ok {
	//	return
	//}
	//t, ok := s.(db)
	//if !ok {
	//	return\

	//}
	//fmt.Println(t, reflect.TypeOf(t))
	//t.User = "root"
	//t.Password = "123456"
	//t.IPAddress = "127.0.0.1"
	//t.DataName = "db_config"
	//t.Port = "3306"
	//t.DataType = "mySql"
	//fmt.Println(s, reflect.TypeOf(s))

	//operation := util.Operation{Addition{}}
	//
	//res := operation.Operate(1,2)
	//fmt.Println(res)

}


func testSelect(dbWorker earPicking.DbWorker){
	var retList []Activity
	dbWorker.SetTableName("data_source").OrderBy("id", earPicking.SQL_OB_ASC).Limit("1,3").SelectAll(&retList)
	fmt.Printf("testSelect %+v\n", retList)
}

func testSelectToMap(dbWorker earPicking.DbWorker){
	myMap, _ := dbWorker.SetTableName("data_source").Where("id = 1").ToMap()
	fmt.Printf("testSelectToMap %+v\n", myMap)
}

func testWhereSelect(dbWorker earPicking.DbWorker){
	ret := Activity{}
	dbWorker.SetTableName("data_source").Where("id = 1").Select(&ret)
	fmt.Printf("WhereSelect %+v\n", ret)
}

func testGroupBy(dbWorker earPicking.DbWorker){
	ret := Activity{}
	dbWorker.SetTableName("data_source").GroupBy("remake").Select(&ret)
	fmt.Printf("testGroupBy %+v\n", ret)
}

func testOrderByList(dbWorker earPicking.DbWorker){
	var retList  []Activity
	dbWorker.SetTableName("data_source").OrderBy("id", earPicking.SQL_OB_DESC).Where("remake = 1").SelectAll(&retList)
	fmt.Printf("testOrderBy %+v\n", retList)
}

func testInster(dbWorker earPicking.DbWorker){
	var a Source
	a.ID = 4
	a.Name = "DB_CONFIG"
	a.JdbcUrl = "jdbc:mysql://127.0.0.1:11306/db_config?useUnicode=true&characterEncoding=utf8"
	a.DriverClass = "com.mysql.jdbc.Driver"
	a.User = "root"
	a.WriteOrRead= 1


	code := dbWorker.SetTableName("data_source").InsertData(&a)
	fmt.Printf("%+v\n", code)
}

func testDel(dbWorker earPicking.DbWorker){
	code := dbWorker.SetTableName("data_source").Where("id=4").DeleteData()
	fmt.Printf("%+v\n", code)
}

func testUpdate(dbWorker earPicking.DbWorker){
	var a Source
	a.Name = "DB_CONFIG123"

	a.User = "admin"
	a.Password ="654321"
	a.WriteOrRead= 2
	code := dbWorker.SetTableName("data_source").Where("id=4").ModifyData(&a)
	fmt.Printf("%+v\n", code)
}

//
//type Addition struct {}
//
//func (Addition) Apply(lval, rval int) int{
//	return lval + rval
//}
