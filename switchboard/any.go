package switchboard

import "gitlab.com/xx_network/primitives/id"

// ID to respond to any message type
const AnyType = int32(0)

//ID to respond to any user
func AnyUser() *id.ID {
	return &id.ZeroUser
}
