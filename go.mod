module github.com/xiaoenai/xmodel

go 1.13

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

require (
	github.com/coreos/etcd v3.3.17+incompatible
	github.com/coreos/go-systemd v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/go-redis/redis/v7 v7.0.0-beta.4
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/henrylee2cn/cfgo v0.0.0-20180417024816-e6c3cc325b21
	github.com/henrylee2cn/goutil v0.0.0-20200609142700-8e4679f1c13f
	github.com/urfave/cli v1.22.1
	github.com/xiaoenai/glog v0.0.0-20200611142840-66249c007189
	go.etcd.io/bbolt v1.3.3 // indirect
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/crypto v0.0.0-20191112222119-e1110fd1c708
	google.golang.org/genproto v0.0.0-20200610212329-df9b449b0ff2 // indirect
	google.golang.org/grpc v1.29.1 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)
