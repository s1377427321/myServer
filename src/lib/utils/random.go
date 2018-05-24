package utils

import (
	"sync"
	"math/rand"
)

var o *rand.Rand = rand.New(rand.NewSource(TimestampNano()))

var random_mux_ sync.Mutex

// 区间
func RandInt32Section(min, max int32) (r int32) {
	random_mux_.Lock()
	defer random_mux_.Unlock()
	return (o.Int31n(max-min) + min)
}