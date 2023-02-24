module go.kuoruan.net/log

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/stretchr/testify v1.8.0
	go.uber.org/zap v1.23.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

retract (
	v0.3.0
	v0.2.0
	v0.1.0
)
