package main

import (
	"ch2/earPicking"
	"fmt"
)

// Activity测试结构体
type Activity struct {
	ID   int64  `col:"id" json:"id"`
	Name string `col:"dataName" json:"dataName"`
}

func main() {

	// dbWorker := util.DbWorker{
	// 	Dsn: "root:123456@tcp(127.0.0.1:3306)/db_config",
	// }
	var dbWorker earPicking.DbWorker
	var ret = Activity{}
	dbWorker.User = "root"
	dbWorker.Password = "123456"
	dbWorker.IPAddress = "127.0.0.1"
	dbWorker.DataName = "db_config"
	dbWorker.Port = "3306"
	dbWorker.DataType = "mySql"
	//dbWorker.Dsn = "root:123456@tcp(127.0.0.1:3306)/db_config"
	var paramss = make([]interface{}, 0)
	paramss = append(paramss, 1)
	// paramss = append(paramss, "mui")
	err := dbWorker.QueryData("SELECT * FROM data_source WHERE id = ?", 1).Unique(&ret)
	fmt.Printf("%+v\n", err)
	fmt.Printf("%+v\n", ret)

	// Activity结构体数组
	retList := []Activity{}
	err = dbWorker.QueryData("SELECT * FROM data_source").List(&retList)
	fmt.Printf("%+v\n", err)
	fmt.Printf("%+v\n", retList)
	//dbWorker.QueryData("SELECT * FROM data_source WHERE id = ? and dataName = ?", nil)
}
