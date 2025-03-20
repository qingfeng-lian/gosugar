package sugar

import (
	"fmt"
	"testing"
)

func TestExponentialBackoff(t *testing.T) {
	// 0值
	d, err := ExponentialBackoff(1, 60, 0, false)
	fmt.Println(d, err)

	for i := 0; i < 15; i++ {
		d, err := ExponentialBackoff(i, 60, 3600*3, true)
		if err != nil {
			// t.Fatal("出现了错误", err, d)
			break
		}
		msg := fmt.Sprintf("第%d次，间隔时间：%s", i+1, d.String())
		println(msg)
	}
}
