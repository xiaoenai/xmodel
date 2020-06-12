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

/*
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `name` varchar(20) NOT NULL DEFAULT '' COMMENT '用户名',
  `age` int(10) NOT NULL DEFAULT '0' COMMENT '用户年龄',
  `updated_at` bigint(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `created_at` bigint(11) NOT NULL COMMENT '创建时间',
  `deleted_ts` bigint(11) NOT NULL DEFAULT '0' COMMENT '删除时间(0表示未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户信息表';

CREATE TABLE `log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `text` varchar(20) NOT NULL DEFAULT '' COMMENT 'desc',
  `updated_at` bigint(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `created_at` bigint(11) NOT NULL COMMENT '创建时间',
  `deleted_ts` bigint(11) NOT NULL DEFAULT '0' COMMENT '删除时间(0表示未删除)',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='日志表';

CREATE TABLE `device` (
  `uuid` varchar(20) NOT NULL DEFAULT '' COMMENT 'uuid',
  `updated_at` bigint(20) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `created_at` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `deleted_ts` bigint(20) NOT NULL DEFAULT '0' COMMENT '删除时间(0表示未删除)',
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='设备表';
*/

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
