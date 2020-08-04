module gitlab.com/elixxir/primitives

go 1.13

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/xx_network/primitives v0.0.0-20200803231956-9b192c57ea7c
)

replace (
	gitlab.com/xx_network/collections/ring => gitlab.com/xx_network/collections/ring.git v0.0.1
)