module github.com/masibw/goone

go 1.16

require (
	github.com/gostaticanalysis/analysisutil v0.6.1
	// idk why, but we need this indirect dependency to pass the test github.com/masibw/goone_test v0.0.0-20210112093021-7d2e0b363db0
	github.com/masibw/goone_test v0.0.0-20210112093021-7d2e0b363db0 // indirect
	golang.org/x/tools v0.1.11
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/yaml.v2 v2.4.0
)
