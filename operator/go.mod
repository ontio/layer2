module github.com/ontio/layer2/operator

go 1.14

require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/ontio/ontology v1.10.0
	github.com/ontio/ontology-go-sdk v0.0.0
	github.com/urfave/cli v1.22.4
)

replace (
	github.com/ontio/ontology => github.com/blockchain-develop/ontology v1.11.1-0.20200728102543-adcdf8013053
	github.com/ontio/ontology-go-sdk => github.com/blockchain-develop/ontology-go-sdk v1.11.5-0.20200729075224-2da2e1bbdfbd
)
