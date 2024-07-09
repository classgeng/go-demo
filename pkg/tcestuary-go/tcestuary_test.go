package tcestuary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMysqlConfig(t *testing.T) {

	// 测试期间, 临时设置配置路径
	SetConfigDirectory("./_example")

	t.Run("key-exist", func(t *testing.T) {
		db, err := GetMysqlConfig("ocloud_api3.api_sync")
		assert.NoError(t, err)
		assert.NotNil(t, db)
		t.Logf("%v\n", *db)
	})

	t.Run("key-not-exist", func(t *testing.T) {
		db, err := GetMysqlConfig("ocloud_api3.not_exist")
		assert.Equal(t, ErrNotFound, err)
		assert.Nil(t, db)
	})

	t.Run("key-format-error", func(t *testing.T) {
		db, err := GetMysqlConfig("ocloud_api3")
		assert.Equal(t, ErrKeyFormat, err)
		assert.Nil(t, db)
	})

	t.Run("scope-error", func(t *testing.T) {
		db, err := GetMysqlConfig("dbsql_tcenter_CCDB4.CCDB4")
		assert.Equal(t, ErrUsageInvalid, err)
		assert.Nil(t, db)
		db, err = GetMysqlConfig("dbsql_yje_yujie_data.not_exist")
		assert.Equal(t, ErrUsageInvalid, err)
		assert.Nil(t, db)
	})
}

func TestGetMysqlConfigAllRegion(t *testing.T) {

	// 测试期间, 临时设置配置路径
	SetConfigDirectory("./_example")

	t.Run("key-exist", func(t *testing.T) {
		dbs, err := GetMysqlConfigAllRegion("dbsql_tcenter_CCDB4.CCDB4")
		assert.NoError(t, err)
		assert.NotNil(t, dbs)
		for _, db := range dbs {
			assert.NotEqual(t, 0, db.RegionID)
			t.Logf("%+v\n", *db)
		}
	})

	t.Run("key-not-exist", func(t *testing.T) {
		db, err := GetMysqlConfigAllRegion("not_exist.not_exist")
		assert.Equal(t, ErrNotFound, err)
		assert.Nil(t, db)
	})

	t.Run("scope-error", func(t *testing.T) {
		db, err := GetMysqlConfigAllRegion("ocloud_api3.not_exist")
		assert.Equal(t, ErrUsageInvalid, err)
		assert.Nil(t, db)
		db, err = GetMysqlConfigAllRegion("dbsql_yje_yujie_data.not_exist")
		assert.Equal(t, ErrUsageInvalid, err)
		assert.Nil(t, db)
	})

}
func TestGetMysqlConfigAllZone(t *testing.T) {

	// 测试期间, 临时设置配置路径
	SetConfigDirectory("./_example")

	t.Run("key-exist", func(t *testing.T) {
		dbs, err := GetMysqlConfigAllZone("dbsql_yje_yujie_data.yujie_data")
		assert.NoError(t, err)
		assert.NotNil(t, dbs)
		for _, db := range dbs {
			assert.NotEqual(t, 0, db.RegionID)
			assert.NotEqual(t, 0, db.ZoneID)
			t.Logf("%+v\n", *db)
		}
	})

	t.Run("key-not-exist", func(t *testing.T) {
		db, err := GetMysqlConfigAllZone("not_exist.not_exist")
		assert.Equal(t, ErrNotFound, err)
		assert.Nil(t, db)
	})

	t.Run("scope-error", func(t *testing.T) {
		db, err := GetMysqlConfigAllZone("ocloud_api3.api_sync")
		assert.Equal(t, ErrUsageInvalid, err)
		assert.Nil(t, db)
		db, err = GetMysqlConfigAllZone("dbsql_tcenter_CCDB4.CCDB4")
		assert.Equal(t, ErrUsageInvalid, err)
		assert.Nil(t, db)
	})

}

func TestLoadConfigOnce(t *testing.T) {
	// 测试期间, 临时设置配置路径
	SetConfigDirectory("./_example")

	// 多次触发写 load 链路
	var err error
	_, err = GetMysqlConfig("ocloud_api3.api_sync")
	assert.NoError(t, err)
	_, err = GetMysqlConfig("ocloud_api3.api_sync")
	assert.NoError(t, err)

	// 仅在启动时写一次版本文件
	assert.Equal(t, int32(1), std.loadCounter)
}

func TestWriteVersionOnce(t *testing.T) {
	// 测试期间, 临时设置配置路径
	SetConfigDirectory("./_example")

	// 多次触发写 Version 链路
	GetMysqlConfig("ocloud_api3.api_sync")
	GetMysqlConfig("ocloud_api3.api_sync")

	// 仅在启动时写一次版本文件
	assert.Equal(t, int32(1), std.writeCounter)
}

func TestDebug(t *testing.T) {

	// 测试期间, 临时设置配置路径
	SetConfigDirectory("./_example")

	// 触发配置加载
	GetMysqlConfig("ocloud_api3.api_sync")

	// 测试期间, 输出 Debug 信息
	Debug()
}

func TestGetRegion(t *testing.T) {

	// 测试期间, 临时设置配置路径
	SetConfigDirectory("./_example")

	t.Run("key-exist", func(t *testing.T) {
		r, err := GetRegion(50000005)
		assert.NoError(t, err)
		assert.NotNil(t, r)
		assert.NotEqual(t, 0, r.RegionID)
		assert.NotEqual(t, "", r.RegionName)
	})
	t.Run("key-not-exist", func(t *testing.T) {
		r, err := GetRegion(50000000)
		assert.Equal(t, ErrNotFound, err)
		assert.Nil(t, r)
	})
}

func TestGetZone(t *testing.T) {

	// 测试期间, 临时设置配置路径
	SetConfigDirectory("./_example")

	t.Run("key-exist", func(t *testing.T) {
		r, err := GetZone(50000005, 50050002)
		assert.NoError(t, err)
		assert.NotNil(t, r)
		assert.NotEqual(t, 0, r.RegionID)
		assert.NotEqual(t, "", r.ZoneName)
	})
	t.Run("key-not-exist", func(t *testing.T) {
		r, err := GetZone(50000005, 50050000)
		assert.Equal(t, ErrNotFound, err)
		assert.Nil(t, r)
	})
}
