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

func main() {


	var dbWorker earPicking.DbWorker


	//ret := Activity{}
	dbWorker.User = "root"
	dbWorker.Password = "123456"
	dbWorker.IPAddress = "127.0.0.1"
	dbWorker.DataName = "db_config"
	dbWorker.Port = "3306"
	dbWorker.DataType = "mySql"
	//dbWorker.Dsn = "root:123456@tcp(127.0.0.1:3306)/db_config"

	testSelect(dbWorker)
	testWhereSelect(dbWorker)
	testGroupBy(dbWorker)
	testOrderByList(dbWorker)


	//err := dbWorker.QueryData("SELECT * FROM data_source WHERE id = ?", 1).Unique(&ret)


	// Activity结构体数组
	//retList := []Activity{}

	//err := dbWorker.QueryData("SELECT * FROM data_source").List(&retList)
	//fmt.Printf("%+v\n", err)
	//fmt.Printf("%+v\n", retList)

	//code := dbWorker.InsertData("INSERT INTO data_source (id, dataName, jdbcUrl, driverClass, user, password, writeOrRead, createTime, remake, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",4,"db_config","jdbc:mysql://127.0.0.1:11306/db_config?useUnicode=true&characterEncoding=utf8","com.mysql.jdbc.Driver", "root", "123456", 1, 0, "", 1)
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
	retList := []Activity{}
	dbWorker.SetTableName("data_source").SelectAll(&retList)
	fmt.Printf("testSelect %+v\n", retList)
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
	retList := []Activity{}
	dbWorker.SetTableName("data_source").OrderBy("id", earPicking.SQL_OB_DESC).SelectAll(&retList)
	fmt.Printf("testOrderBy %+v\n", retList)
}

//
//type Addition struct {}
//
//func (Addition) Apply(lval, rval int) int{
//	return lval + rval
//}
