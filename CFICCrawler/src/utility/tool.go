package utility

import (
	"reflect"
	"fmt"
	"os"
	"github.com/axgle/mahonia"
	"bufio"
	"path"
	"strconv"
	"time"
)

var logger = GetLogger()

// To check if the list contains the elem
func Contains(list interface{}, elem interface{}) bool {
	value := reflect.ValueOf(list)
	if value.Kind() != reflect.Slice {
		logger.Errorf("Input type is not an array or slice type: %v, kind:%s", value, value.Kind())

		return false
	}

	for i:=0; i<value.Len();i++ {
		if value.Index(i).Interface() == elem {
			return true
		}
	}
	return false
}

// To get map keys
func Keys(i interface{}) (keys []string) {
	v := reflect.ValueOf(i)

	if v.Kind() != reflect.Map {
		fmt.Fprintf(os.Stderr, "Input type is not a map type: %v", v)
		return nil
	}

	for _,key := range v.MapKeys() {
		keys = append(keys, key.Interface().(string))
	}

	return keys
}

// To get map values
func Values(i interface{}) ([]interface{}) {
	var result []interface{}
	v := reflect.ValueOf(i)

	if v.Kind() != reflect.Map {
		fmt.Fprintf(os.Stderr, "Input type is not a map type: %v", v)
		return nil
	}

	for _,key := range (v.MapKeys()) {
		result = append(result, v.MapIndex(key).Interface())
	}

	return result
}

// Write one line to file.
func WriteToFile(file string, line string) error{
	os.MkdirAll(path.Dir(file), 0777)

	f, err:= os.OpenFile(file, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARN: Open file %s failed, %s\n", file, err)
		return err
	}
	defer f.Close()

	encoder := mahonia.NewEncoder("gbk")
	writer := bufio.NewWriter(encoder.NewWriter(f))
	writer.WriteString(line + "\n")
	writer.Flush()
	//io.Copy(writer, strings.NewReader(line))
	return nil
}


func IsFileExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func String2Int64(value string) int64 {
	if v,err := strconv.ParseInt(value, 0, 8); err != nil {
		return 0
	} else {
		return v
	}
}

func String2Folat32(value string) float32{
	if v,err := strconv.ParseFloat(value, 32); err != nil {
		return 0.0
	} else {
		return float32(v)
	}

}


func String2Folat64(value string) float64{
	if v,err := strconv.ParseFloat(value, 64); err != nil {
		return 0.0
	} else {
		return float64(v)
	}

}

func String2Date(value string) time.Time {
	result, _ := time.Parse("2006-01-02", value)

	return result
}