package earPicking

import (
	"database/sql"
	"errors"

	"reflect"
	"strconv"
)

type (

	// queryRes 查询结果
	queryRes struct {
		isRun bool               // 是否查询过数据
		data []map[string]string // 存储数据库查询数据
		err  error               // 异常
	}

	// tableInfo 表信息结构体
	tableInfo struct {
		tabName       string // 表名
		primaryKey    string // 主键标识
		limit         int    // 分页起始页
		offset        int    // 分页结束页
		sqlTemp       string // sql临时存储
		whereTemp     string // sql条件临时存储
		orderByTemp   string // 排序存储
		groupByTemp   string // 分组聚合
		limitTemp     string
	}
)

var typeRegistry = make(map[string]reflect.Type)

/**
 * formatRes 处理查询响应数据, map对象映射（私有方法）
 * @param rows 行是查询的结果
 * @return error 错误信息，正常时，该值为nil
 */
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
	// jsonStr, err := util.ToJsonStr(recordList)
	// fmt.Println(jsonStr)
	// checkErr(err)
	rows.Close()

	return nil

}

/**
 * ForMap map对象映射（公共方法）
 * @return map[string]string 返回map集合
 */
func (qr *queryRes) ForMap() (map[string]string){


	if len(qr.data) > 0 {
		return qr.data[0]
	}
	return nil
}

/**
 * ToObject 实体对象映射（公共方法）
 * @param v reflect.Value 填充数据的对象
 * @return error 错误信息，正常时，该值为nil
 */
func (qr *queryRes) ToObject(in interface{}) error {
	if len(qr.data) > 0 {
		return qr.mapping(qr.data[0], reflect.ValueOf(in))
	}
	return nil
}

/**
 * toObjForList 集合实体对象映射（公共方法）
 * @param v reflect.Value 填充数据的对象
 * @return error 错误信息，正常时，该值为nil
 */
func (qr *queryRes) ToObjForList(in interface{}) error {

	if qr.err != nil {
		return qr.err
	}

	length := len(qr.data)


	if length > 0 {
		v := reflect.ValueOf(in)
		qr.toObjForList(v)
	}

	return nil
}

/**
 * toObject 实体对象映射（私有方法）
 * @param v reflect.Value 填充数据的对象
 * @return error 错误信息，正常时，该值为nil
 */
func (qr *queryRes) toObject(v reflect.Value)error{
	if len(qr.data) > 0 {
		return qr.mapping(qr.data[0], v)
	}
	return nil
}

/**
 * toObjForList 集合实体对象映射（私有方法）
 * @param v reflect.Value 填充数据的对象
 * @return error 错误信息，正常时，该值为nil
 */
func (qr *queryRes) toObjForList(v reflect.Value) error{

	// 1.reflect.ValueOf->获取接口保管的具体值(实例化)
	// 2.Indirect->获取该值的指针值
	sliceValue := reflect.Indirect(v)

	// 如果该对象类型不是切片类型，则返回类型错误
	if sliceValue.Kind() != reflect.Slice {
		return ErrorNeedPointerToSlice
	}

	// 获取切片集合中的类型
	sliceElementType := sliceValue.Type().Elem()

	for _, results := range qr.data {
		// 根据切片集合中的类型，创建新的实体
		newValue := reflect.New(sliceElementType)
		// 映射数据
		err := qr.mapping(results, newValue)
		if err != nil {
			return err
		}

		//fmt.Printf("sliceValue ： %+v\n", sliceValue)

		// 获取映射后的对象实例的指针值
		rTmep := reflect.Indirect(reflect.ValueOf(newValue.Interface()))

		//fmt.Printf("rTmep ： %+v\n", rTmep)
		// 将映射后的对象实例的指针值复制到源对象指针中
		sliceValue.Set(reflect.Append(sliceValue, rTmep))

		//fmt.Printf("sliceValue111 ： %+v\n", sliceValue)
	}


	return nil
}



/**
 * 处理结构体于查询数据之间的映射关系
 * @param m map[string]string 数据源对象，存储数据源
 * @param v reflect.Value 数据填充对象，用来接收数据
 * @return error 错误信息，正常时，该值为nil
 */
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
			// 判读填充的数据类型，将值转换成对应的数据类型
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

func OjbToString(kind reflect.Kind, v interface{})(valStr string, err error){

	meta := reflect.ValueOf(v)
	switch kind {
	case reflect.String:
		str := "\"" +meta.String()+  "\""
		return str, nil
	case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
		str := strconv.FormatInt(meta.Int(), 10)
		return str,nil
	case reflect.Float32, reflect.Float64:
		str := strconv.FormatFloat(meta.Float(), 'f', -1, 32)
		return str, nil
	case reflect.Bool:
		str := strconv.FormatBool(meta.Bool())
		return str, nil
	default:
		return "", errors.New("没有对应的数据类型")
	}

}

