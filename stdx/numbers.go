// Copyright (C) 2025 T-Force I/O
//
// TF GoLib is licensed under the MIT license.
// You should receive a copy of MIT along with this software.
// If not, see <https://opensource.org/license/mit>

package stdx

import (
	"math/big"
)

// A Uint256 is wrapped *big.Int to intergrate with system that use large number and provide fluent apis for type casting.
// Available since v0.1.0
type Uint256 struct {
	value *big.Int
}

// Cast Uint256 to *big.Int.
// Available since v0.1.0
func (s *Uint256) BigInt() *big.Int {
	return s.value
}

// Cast Uint256 to uint64. This may cause data loss when the value is larger than 18446744073709551615.
// Available since v0.1.0
func (s *Uint256) Uint64() uint64 {
	return s.value.Uint64()
}

// Return underlying value of Bytes as *big.Int.
// Available since v0.1.0
func (s *Uint256) Value() *big.Int {
	return s.value
}
