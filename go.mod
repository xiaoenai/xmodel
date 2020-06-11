module github.com/xiaoenai/xmodel

go 1.13

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

require (
	github.com/coreos/etcd v3.3.17+incompatible
	github.com/go-redis/redis/v7 v7.0.0-beta.4
	github.com/go-sql-driver/mysql v1.4.1
	github.com/henrylee2cn/cfgo v0.0.0-20180417024816-e6c3cc325b21
	github.com/henrylee2cn/erpc/v6 v6.3.1
	github.com/henrylee2cn/goutil v0.0.0-20200609142700-8e4679f1c13f
	github.com/urfave/cli v1.22.1
	github.com/xiaoenai/tp-micro/v6 v6.1.2
	golang.org/x/crypto v0.0.0-20191112222119-e1110fd1c708
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)
