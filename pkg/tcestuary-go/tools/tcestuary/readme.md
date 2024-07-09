#### 工具背景说明:

1. 镜像内的 shell 脚本需要访问数据库(日常数据备份/清理工作);
2. 部署过程需要向数据库中写入软件版本信息;

#### 工具使用说明:

##### getMysqlConfig 

单行输出: host ip port user password

```
./tcestuary getMysqlConfig --configDirectory="." ocloud_api3.api_sync
```

##### getMysqlConfigAllRegion

多行输出: host ip port user password regionid

```
./tcestuary getMysqlConfigAllRegion --configDirectory="." dbsql_tcenter_CCDB4.CCDB4
```

##### getMysqlConfigAllZone

多行输出: host ip port user password regionid zoneid

```
./tcestuary getMysqlConfigAllZone --configDirectory="." dbsql_yje_yujie_data.yujie_data
```

#### 修改配置文件路径
```
--configDirectory="."   // 注意: 此处是 “路径”
```

#### 错误码处理 (自动化工具)
1. 0 表示成功, 向 stdout 写出信息;
2. 其它返回码 表示错误, 向 stderr 写出信息;

#### 命令行处理参考

stderr 重定向到 stdout, 用于日志输出
```
mysql=`./tcestuary getMysqlConfig --configDirectory="." ocloud_api3.api_sync 2>&1`
if [ $? != 0 ]; then
    echo "tcestuary error $mysql"
    exit 100
fi

read -r host ip port user passwd <<< $mysql
echo $host $ip $port $user $passwd
```
