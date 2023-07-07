package lock

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestLockAndUnlock(t *testing.T) {

	go Action1()

	time.Sleep(time.Duration(10) * time.Second)

	go Action2()
}

func Action1() {
	for i := 0; i < 100; i++ {
		fmt.Printf("lock success: %d \n", i)
		Lock(strconv.Itoa(i))
		fmt.Printf("unlock success: %d \n", i)
	}
}

func Action2() {
	for i := 0; i < 100; i++ {
		Unlock(strconv.Itoa(i))
		fmt.Printf("unlock: %d \n", i)
		time.Sleep(time.Duration(1) * time.Second)
	}
}
