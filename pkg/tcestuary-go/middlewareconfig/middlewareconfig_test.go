package middlewareconfig

import (
	"fmt"
	"testing"
)

func TestMysqlGetConfig(t *testing.T) {
	mysqlConfig := &MysqlConfig{}
	mysqlConfig.GetConfig("dbsql_billing")
	fmt.Printf("%v", mysqlConfig)

}

func TestRedisGetConfig(t *testing.T) {
	redisConfig := &RedisConfig{}
	redisConfig.GetConfig("ckv_cas")
	fmt.Printf("%v", redisConfig)

}

func TestKafkaGetConfig(t *testing.T) {
	kafkaConfig := &KafkaConfig{}
	kafkaConfig.GetConfig("wtag")
	fmt.Printf("%v", kafkaConfig)

}

func TestCMQGetConfig(t *testing.T) {
	cmqConfig := &CMQConfig{}
	cmqConfig.GetConfig("mq_waccount")
	fmt.Printf("%v", cmqConfig)
}

func TestCSPGetConfig(t *testing.T) {
	cspConfig := &CSPConfig{}
	cspConfig.GetConfig("yehe_file")
	fmt.Printf("%v", cspConfig)
}

func TestMongodbConfig(t *testing.T) {
	mongodbConfig := &MongodbConfig{}
	mongodbConfig.GetConfig("bsp_document")
	fmt.Printf("%v", mongodbConfig)
}

func TestESConfig(t *testing.T) {
	esConfig := &ESConfig{}
	esConfig.GetConfig("es_audit")
	fmt.Printf("%v", esConfig)
}

func TestZKConfig(t *testing.T) {
	zkConfig := &ZKConfig{}
	zkConfig.GetConfig("yunapi3_zk")
	fmt.Printf("%v", zkConfig)
}

func TestHdfsConfig(t *testing.T) {
	hdfsConfig := &HdfsConfig{}
	hdfsConfig.GetConfig("file_server")
	fmt.Printf("%v", hdfsConfig)
}
