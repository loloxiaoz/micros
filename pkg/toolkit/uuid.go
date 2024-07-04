package toolkit

import (
	"time"
)

const ID_OFFSET = 1251414000
const SEQ_BIT_LENGTH = 18
const INSTANCE_SEQ_BIT_LENGTH = 6 + 18
const SEQ_BIT_VALUE = 1 << 18

var seq int64 = 0

// GenID 生成ID
func GenID(instance int) uint {
	uuid := time.Now().Unix()
	uuid -= ID_OFFSET
	uuid <<= INSTANCE_SEQ_BIT_LENGTH
	uuid += int64(instance << SEQ_BIT_LENGTH)
	seq++
	uuid += (seq % SEQ_BIT_VALUE)
	return uint(uuid)
}
