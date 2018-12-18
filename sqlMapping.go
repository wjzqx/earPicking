package earPicking

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/mySimpleCache/util"
)

// queryRes 查询结果
type queryRes struct {
	data []map[string]string // 存储数据库查询数据
	err  error               // 异常
}

// formatRes 处理查询响应数据
// @param rows 行是查询的结果
func (qr *queryRes) formatRes(rows *sql.Rows) error {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([][]byte, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	recordList := make([]map[string]string, 0)
	for rows.Next() {

		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		checkErr(err)
		record := make(map[string]string)
		for i, col := range values {

			if col != nil {
				record[columns[i]] = string(col)
			}
		}
		recordList = append(recordList, record)

	}

	qr.data = recordList
	jsonStr, err := util.ToJsonStr(recordList)
	fmt.Println(jsonStr)
	checkErr(err)
	rows.Close()

	return nil

}

// Unique ...
func (qr *queryRes) Unique(in interface{}) error {
	if len(qr.data) > 0 {
		return qr.mapping(qr.data[0], reflect.ValueOf(in))
	}
	return nil
}

// List 集合对象映射
func (qr *queryRes) List(in interface{}) error {

	if qr.err != nil {
		return qr.err
	}

	length := len(qr.data)

	if length > 0 {
		v := reflect.ValueOf(in).Elem()
		fmt.Printf("v  %+v\n", v)
		newv := reflect.MakeSlice(v.Type(), 0, length)
		fmt.Printf("v  %+v\n", v)

		v.Set(newv)
		v.SetLen(length)

		index := 0
		for i := 0; i < length; i++ {
			k := v.Type().Elem()
			newObj := reflect.New(k)
			err := qr.mapping(qr.data[i], newObj)
			if err != nil {
				return err
			}

			n := reflect.ValueOf(newObj)
			fmt.Printf("n  %+v\n", n)
			m := reflect.ValueOf(v.Index(index))
			m = n
			fmt.Printf("v.  %+v\n", m)

			//v.Index(index).Set(newObj)
			index++
		}
		v.SetLen(index)
	}
	return nil
}

// 处理结构体于查询数据之间的映射关系
func (qr *queryRes) mapping(m map[string]string, v reflect.Value) error {

	t := v.Type()
	val := v.Elem()
	typ := t.Elem()

	// 检查结构体的数据类型
	if !val.IsValid() {
		return errors.New("数据类型不正确")
	}

	for i := 0; i < val.NumField(); i++ {
		value := val.Field(i)
		kind := value.Kind()
		// 获取字段注解
		tag := typ.Field(i).Tag.Get("col")

		if len(tag) > 0 {
			meta, ok := m[tag]
			if !ok {
				continue
			}
			// 判断字段是否有读写权限
			if !value.CanSet() {
				return errors.New("结构体字段没有读写权限")
			}

			if len(meta) == 0 {
				continue
			}

			if kind == reflect.String {
				value.SetString(meta)
			} else if kind == reflect.Float32 {
				f32, err := strconv.ParseFloat(meta, 32)

				if err != nil {
					return err
				}
				value.SetFloat(f32)
			} else if kind == reflect.Float64 {
				f64, err := strconv.ParseFloat(meta, 64)

				if err != nil {
					return err
				}
				value.SetFloat(f64)
			} else if kind == reflect.Int64 {
				i64, err := strconv.ParseInt(meta, 10, 64)
				if err != nil {
					return err
				}
				value.SetInt(i64)
			} else if kind == reflect.Int {
				i, err := strconv.Atoi(meta)
				if err != nil {
					return err
				}

				value.SetInt(int64(i))
			} else if kind == reflect.Bool {
				b, err := strconv.ParseBool(meta)
				if err != nil {
					return err
				}
				value.SetBool(b)
			} else {
				return errors.New("没有对应的数据类型")
			}
		}
	}
	return nil
}
