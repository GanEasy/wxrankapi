// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import "testing"
import "fmt"

func Test_GetSign(t *testing.T) {
	js, err := GetSign("o7UTkjr7if4AQgcPmveQ5wJ5alsA")
	if err != nil {
		panic(err)
	}
	fmt.Println(js)
}
