// khan
// https://github.com/jpholanda/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package util

import "time"

// NowMilli returns now in milliseconds since epoch
func NowMilli() int64 {
	return time.Now().UnixNano() / 1000000
}
