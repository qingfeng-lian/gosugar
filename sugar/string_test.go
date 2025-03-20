package sugar

import (
	"fmt"
	"testing"
)

func TestStringSub(t *testing.T) {
	text := "聊天内容 ai 喂，您好ai 家长您好，我是希望学的老师，我们推出了19元最高34课时的重难点提升课，帮助孩子快速提分。您家孩子现在是几年级呢？user 喂。ai 咱家孩子今年几年级了呀？sys [坐席已弹屏]sys [API调用-跟踪]sys [API调用-转人工]sys [坐席已介入]sys [主叫挂机-坐席]录音连接 https://robot-audio.oss-cn-beijing.aliyuncs.com/S/2025/01/14/974/974-1170-S48919842633-3729835432620327912-20250114153258.mp3"
	//for i := 1; i <= 30; i++ {
	//	fmt.Println(StringSub(text, i))
	//}
	fmt.Println(StringSub(text, 500))
}
