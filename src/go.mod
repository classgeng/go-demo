module godemo

go 1.13

replace git.code.oa.com/tce-config/tcestuary-go/v4 => ./../pkg/tcestuary-go

require (
	git.code.oa.com/tce-config/tcestuary-go/v4 v4.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
