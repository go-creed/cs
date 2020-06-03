# cs
文件存储系统

## How To Start

You must have etcd node and then execute the following command to add the App configuration

```shell
etcdctl put /cs/app/config/go.micro.cs.service.auth "{\"name\":\"go.micro.cs.service.auth\",\"address\":\"localhost:12004\",\"version\":\"latest\"}"
```

Start after

```shell
go run main.go --registry=etcd --registry_address=127.0.0.1:2379 --cc=127.0.0.1:2379
```

If the startup fails, please check whether the configuration files of redis and mysql are correct

1. Add a mysql configuration file

   ```shell
   etcdctl put /cs/app/config/mysql "{\"user\":\"root\",\"password\":\"root\",\"address\":\"localhost:3306\",\"db_name\":\"cs\",\"logMode\":true}"
   ```

2. Add a redis configuration file

   ```shell
   etcdctl put /cs/app/config/redis "{\"user\":\"root\",\"password\":\"root\",\"address\":\"localhost:6379\"}"
   ```

   