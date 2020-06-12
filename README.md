# xmodel
Golang model工具集/脚手架(MySql/Redis/Etcd/MongoDB)

## Install

```
go get -u -f -d github.com/xiaoenai/xmodel/...
cd $GOPATH/src/github.com/xiaoenai/xmodel/cmd/xmodel
go install
```

## Command

```
NAME:
   XModel - a deployment tools of xmodel frameware

USAGE:
   xmodel [global options] command [command options] [arguments...]

VERSION:
   v1.0.0

AUTHOR:
   xiaoenai

COMMANDS:
   gen      Generate a xmodel code
   tpl      Add mysql model struct code to project template
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## xmodel gen

- 执行`xmoel gen`生成模板文件
- 模板文件配置在__model__tpl__.go中，根据数据库表在模板中编写对应的结构体，将结构体名添加到`__MYSQL_MODEL__`或`__MONGO_MODEL__`中执行`xmoel gen`命令生成

```
.
├── __model__gen__.lock
├── __model__tpl__.go
├── args
│   ├── const.gen.go
│   └── type.gen.go // 类型
├── go.mod
└── model
    ├── init.go // DB初始化
    ├── mongo_meta.gen.go // MongoDB表  
    ├── mysql_device.gen.go // MySql表
    ├── mysql_log.gen.go
    └── mysql_user.gen.go

2 directories, 10 files
```

## xmodel tpl

- 用于直接连接数据库生成表结构

```
xmodel tpl -host 127.0.0.1 -port 3306 -username qas -password 123456 -db test -table test -ssh_user linux_user -ssh_host 127.0.0.1 -ssh_port 22
```

## __model__tpl__.go

```
// Command xmodel is the xmodel tools.
// The framework reference: https://github.com/xiaoenai/xmodel
package __TPL__

// __MYSQL_MODEL__ create mysql model
type __MYSQL_MODEL__ struct {
	User
	Log
	Device
}

// __MONGO_MODEL__ create mongodb model
type __MONGO_MODEL__ struct {
	Meta
}

// User user info
type User struct {
	Id   int64  `key:"pri"`
	Name string `key:"uni"`
	Age  int32
}

type Log struct {
	Text string
}

type Device struct {
	UUID string `key:"pri"`
}

type Meta struct {
	Hobby []string
	Tags  []string
}
```

## Example

```
package example

import (
	"context"
	"testing"
	"time"

	"github.com/xiaoenai/glog"
	"github.com/xiaoenai/xmodel/example/model"
	"github.com/xiaoenai/xmodel/mongo"
	"github.com/xiaoenai/xmodel/mysql"
	"github.com/xiaoenai/xmodel/redis"
	"github.com/xiaoenai/xmodel/sqlx"
)

func TestXmodel(t *testing.T) {
	// mysql
	mysqlConfig := &mysql.Config{
		Database:     "xmodel",
		Username:     "root",
		Password:     "",
		Host:         "127.0.0.1",
		Port:         3306,
		MaxIdleConns: 50,
		MaxOpenConns: 50,
		NoCache:      false,
	}

	// mongodb
	mongodbConfig := &mongo.Config{
		Addrs:     []string{"127.0.0.1:27017"},
		Timeout:   10,
		PoolLimit: 256,
		Username:  "root",
		Password:  "",
		Database:  "test",
	}

	// redis
	redisConfig := &redis.Config{
		DeployType: redis.TypeSingle,
		ForSingle:  redis.SingleConfig{Addr: "127.0.0.1:6379"},
	}
	if err := model.Init(mysqlConfig, mongodbConfig, redisConfig, time.Duration(24)*time.Hour); err != nil {
		panic(err)
	}

	// insert
	if _, err := model.InsertUser(&model.User{
		Name: "xmodel",
	}); err != nil {
		glog.Errorf("TestXmodel: insert uer err-> %v", err)
		return
	}

	// transaction
	if err := model.GetMysqlDB().TransactCallbackInSession(func(ctx context.Context, tx *sqlx.Tx) error {
		// insert
		if _, err := model.InsertUser(&model.User{
			Name: "xmodel",
		}); err != nil {
			glog.Errorf("TestXmodel: insert uer err-> %v", err)
			return err
		}
		return nil
	}); err != nil {
		glog.Errorf("TestXmodel: insert user transaction err-> %v", err)
	}

	// select
	user, exists, err := model.GetUserByPrimary(1)
	if err != nil {
		glog.Errorf("TestXmodel: GetUserByPrimary err-> %v", err)
		return
	}
	if !exists {
		glog.Errorf("TestXmodel: user not exists")
		return
	}
	glog.Infof("TestXmodel: user.name-> %s", user.Name)
}
```