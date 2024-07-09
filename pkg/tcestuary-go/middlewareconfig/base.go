package middlewareconfig

import (
	"encoding/json"
	"fmt"
	"git.code.oa.com/tce-config/tcestuary-go/v4"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

type MiddleWareConfig interface {
	GetConfig(name string) error
}

const encryptedTagName = "encrypted"

var (
	storageSecurity tcestuary.StorageSecurity
	config          *Config
)

func init() {
	config = &Config{}
	configPath := os.Getenv("CONFIG_FILE_PATH")
	if configPath != "" {
		tcestuary.SetConfigDirectory(configPath)
	}
	b, err := ioutil.ReadFile(tcestuary.GetConfigDirectory())
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, config)
	if err != nil {
		panic(err)
	}

	storageSecurity, err = tcestuary.NewPasswdSecret()
	if err != nil {
		panic(err)
	}

}

type Config struct {
	Mysql   map[string]*json.RawMessage `json:"mysql"`
	ES      map[string]*json.RawMessage `json:"es"`
	Redis   map[string]*json.RawMessage `json:"redis"`
	ZK      map[string]*json.RawMessage `json:"zk"`
	Mongodb map[string]*json.RawMessage `json:"mongodb"`
	CSP     map[string]*json.RawMessage `json:"csp"`
	Kafka   map[string]*json.RawMessage `json:"kafka"`
	CMQ     map[string]*json.RawMessage `json:"cmq"`
	Hdfs    map[string]*json.RawMessage `json:"hdfs"`
}

func (c *Config) GetMiddleWareFiled(name, filed string) (*json.RawMessage, error) {
	t := reflect.TypeOf(c)
	v := reflect.ValueOf(c)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("config is not a struct")
	}
	filedNum := t.NumField()
	for i := 0; i < filedNum; i++ {
		if strings.ToLower(t.Field(i).Name) == strings.ToLower(name) {
			val := v.FieldByName(t.Field(i).Name)
			if message, ok := val.Interface().(map[string]*json.RawMessage); ok {
				return message[filed], nil
			}
		}
	}
	return nil, nil

}

func unmarshallConfig(name, filed string, inStructPtr interface{}) error {
	val, err := config.GetMiddleWareFiled(name, filed)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(*val, inStructPtr); err != nil {
		return err
	}
	return structByReflect(inStructPtr)
}

func structByReflect(inStructPtr interface{}) error {
	rType := reflect.TypeOf(inStructPtr)
	rVal := reflect.ValueOf(inStructPtr)

	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
		rVal = rVal.Elem()
	} else {
		return fmt.Errorf("config must be ptr to struct")
	}
	//
	for i := 0; i < rType.NumField(); i++ {
		t := rType.Field(i)
		f := rVal.Field(i)

		// 基础类型
		if t.Type.Kind() == reflect.String {
			encrypted := t.Tag.Get(encryptedTagName)
			if encrypted == "true" {
				//decrypt string
				res, err := storageSecurity.Decrypt(f.String())
				if err != nil {
					return err
				}
				f.Set(reflect.ValueOf(res))
			}

		} else if t.Type.Kind() == reflect.Slice {
			ele := t.Type.Elem()
			if ele.Kind() == reflect.Ptr {
				ele = ele.Elem()
			}
			if ele.Kind() == reflect.String {
				encrypted := t.Tag.Get(encryptedTagName)
				if encrypted == "true" {
					for j := 0; j < f.Len(); j++ {
						//decrypt string
						val := f.Index(j)
						res, err := storageSecurity.Decrypt(val.String())
						if err != nil {
							return err
						}
						val.Set(reflect.ValueOf(res))
					}
				}
			} else if ele.Kind() == reflect.Struct {
				for j := 0; j < f.Len(); j++ {
					if err := structByReflect(f.Index(j)); err != nil {
						return err
					}
				}
			}
		} else if t.Type.Kind() == reflect.Struct {
			if err := structByReflect(f); err != nil {
				return err
			}
		}
	}
	return nil
}
