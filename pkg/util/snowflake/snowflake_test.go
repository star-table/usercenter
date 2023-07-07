package snowflake

import (
	"fmt"
	"github.com/bmizerany/assert"
	"testing"
)

func TestWorker_GetId(t *testing.T) {
	worker, err := NewNode(2)
	assert.Equal(t, err, nil)
	fmt.Println(worker.Generate())
}
