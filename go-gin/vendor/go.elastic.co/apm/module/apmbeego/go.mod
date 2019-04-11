module go.elastic.co/apm/module/apmbeego

require (
	github.com/astaxie/beego v1.11.1
	github.com/stretchr/testify v1.2.2
	go.elastic.co/apm v1.3.0
	go.elastic.co/apm/module/apmhttp v1.3.0
	go.elastic.co/apm/module/apmsql v1.3.0
)

replace go.elastic.co/apm => ../..

replace go.elastic.co/apm/module/apmhttp => ../apmhttp

replace go.elastic.co/apm/module/apmsql => ../apmsql
