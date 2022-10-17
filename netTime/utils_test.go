///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package netTime

// mockTimeSource is a utility for testing. This gives you an object
// adhering to the TimeSource interface. It will return via NowMs
// the return time set in the structure.
type mockTimeSource struct {
	returnTime int64
}

// NowMs will return the returnTime set in the structure.
func (m *mockTimeSource) NowMs() int64 {
	return m.returnTime
}
