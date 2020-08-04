module gitlab.com/xx_network/primitives

require (
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/primitives v0.0.0-20200731184040-494269b53b4d
)

replace (
	gitlab.com/xx_network/collections/ring => gitlab.com/xx_network/collections/ring.git v0.0.1
)


go 1.13
