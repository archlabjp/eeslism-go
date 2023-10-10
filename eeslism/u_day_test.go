package eeslism

import (
	"testing"

	"gotest.tools/assert"
)

func Test_FNNday(t *testing.T) {
	assert.Equal(t, 1, FNNday(1, 1))     // 1月1日: 1日目
	assert.Equal(t, 32, FNNday(2, 1))    // 2月1日: 32日目
	assert.Equal(t, 365, FNNday(12, 31)) // 12月31日: 365日目
}
