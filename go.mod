module github.com/xiaoenai/xmodel

go 1.13

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

require (
	github.com/coreos/etcd v2.3.8+incompatible
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-redis/redis/v7 v7.0.0-beta.4
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/henrylee2cn/cfgo v0.0.0-20180417024816-e6c3cc325b21
	github.com/henrylee2cn/goutil v0.0.0-20191029125303-21920e347847
	github.com/kr/pretty v0.1.0 // indirect
	github.com/urfave/cli v1.22.9
	github.com/xiaoenai/glog v0.0.0-20200611142840-66249c007189
	golang.org/x/crypto v0.0.0-20190510104115-cbcb75029529
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)
