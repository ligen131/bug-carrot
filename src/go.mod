module bug-carrot

go 1.17

require (
	github.com/labstack/echo/v4 v4.7.2
	github.com/sirupsen/logrus v1.8.1
	github.com/togatoga/goforces v0.0.0-20200804081705-45bb4957d135
	github.com/yanyiwu/gojieba v1.1.2
	go.mongodb.org/mongo-driver v1.9.0
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require github.com/google/go-querystring v1.1.0 // indirect

replace github.com/yanyiwu/gojieba v1.1.2 => github.com/ttys3/gojieba v1.1.3
