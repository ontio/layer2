module github.com/ontio/layer2/server

go 1.14

require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/ontio/ontology-go-sdk v0.0.0
	github.com/urfave/cli v1.22.4
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0
)

replace github.com/ontio/ontology-go-sdk => github.com/blockchain-develop/ontology-go-sdk v1.11.5-0.20200729075224-2da2e1bbdfbd
