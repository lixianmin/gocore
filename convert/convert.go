/********************************************************************
created:    2018-09-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

package convert

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"
	"unsafe"
)

func ToString(v interface{}) string {
	switch v := v.(type) {
	case []byte:
		return *(*string)(unsafe.Pointer(&v))
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case string:
		return v
	default:
		return fmt.Sprintf("convert.ToString(): unexpected type for ToString, got type %T", v)
	}
}

var _toJson = json.Marshal
var _fromJson = json.Unmarshal

func InitJson(toJson func(v interface{}) ([]byte, error), fromJson func(data []byte, v interface{}) error) {
	if toJson == nil {
		panic("toJson is nil")
	}

	if fromJson == nil {
		panic("fromJson is nil.")
	}

	_toJson = toJson
	_fromJson = fromJson
}

func ToJson(v interface{}) ([]byte, error) {
	return _toJson(v)
}

func FromJson(data []byte, v interface{}) error {
	return _fromJson(data, v)
}

func ToHuman(num uint64) string {
	if num >= 1073741824 {
		var v = float64(num) / 1073741824
		return fmt.Sprintf("%.2fG", v)
	} else if num >= 1048576 {
		var v = float64(num) / 1048576
		return fmt.Sprintf("%.2fM", v)
	} else if num >= 1024 {
		var v = float64(num) / 1024
		return fmt.Sprintf("%.2fK", v)
	} else {
		return fmt.Sprintf("%dB", num)
	}
}

//////////////////////////////////////////////////////////////////

var zWriterPool = sync.Pool{New: func() interface{} {
	return zlib.NewWriter(nil)
}}

func ToZLib(input []byte, buffer *bytes.Buffer) error {
	if input == nil || buffer == nil {
		return nil
	}

	var writer = zWriterPool.Get().(*zlib.Writer)
	writer.Reset(buffer)

	_, err := writer.Write(input)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	zWriterPool.Put(writer)
	return nil
}

func FromZLib(input []byte) ([]byte, error) {
	if input == nil {
		return nil, nil
	}

	buffer := bytes.NewReader(input)
	var reader, err = zlib.NewReader(buffer)
	if err != nil {
		return nil, err
	}

	output, err := ioutil.ReadAll(reader)
	reader.Close()

	if err != nil {
		return output, err
	}

	return output, nil
}

func ToJsonZlib(v interface{}, buffer *bytes.Buffer) error {
	var jsonData, err = _toJson(v)
	if nil != err {
		return err
	}

	return ToZLib(jsonData, buffer)
}