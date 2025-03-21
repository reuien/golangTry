package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"time"
)

var (
	timeout      int64
	size         int
	count        int
	typ          uint8 = 8
	code         uint8 = 0
	sendCount    int
	successCount int
	failCount    int
	minTx        int64 = math.MaxInt64
	maxTx        int64 = 0
	totalTs      int64
)

type ICMP struct {
	Type        uint8
	Code        uint8
	CheckSum    uint16
	ID          uint16
	SequenceNum uint16
} //顺序不能改变

func main() {
	getCommandArgs()
	desIp := os.Args[len(os.Args)-1]
	conn, err := net.DialTimeout("ip:icmp", desIp, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	fmt.Printf("正在Ping %s [%s] 具有 %d 字节的数据: \n", desIp, conn.RemoteAddr(), size)
	for i := 0; i < count; i++ {
		sendCount++
		t1 := time.Now()
		icmp := &ICMP{
			Type:        typ,
			Code:        code,
			CheckSum:    0,
			ID:          1,
			SequenceNum: 1,
		}

		data := make([]byte, size)
		var buffer bytes.Buffer
		binary.Write(&buffer, binary.BigEndian, icmp)
		buffer.Write(data)
		data = buffer.Bytes()

		checkSum := checkSum(data)
		//填充到高八位的值
		data[2] = byte(checkSum >> 8)
		data[3] = byte(checkSum)
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
		n, err := conn.Write(data)
		if err != nil {
			failCount++
			log.Println(err)
			continue
		}
		buf := make([]byte, 65535)
		n, err = conn.Read(buf)
		if err != nil {
			failCount++
			log.Println(err)
			continue
		}
		successCount++
		ts := time.Since(t1).Milliseconds()
		if minTx > ts {
			minTx = ts
		}
		if maxTx < ts {
			maxTx = ts
		}
		totalTs += ts
		fmt.Printf("来自 %d.%d.%d.%d 的回复: 字节 = %d 时间 = %d ms TTL = %d\n", buf[12], buf[13], buf[14], buf[15], n-28, ts, buf[8])
		time.Sleep(time.Second)
	}
	fmt.Printf("%s 的 Ping 统计信息：\n   数据包：已发送 =%d ,已接收 = %d  丢失 %d, (%.2f)丢失，\n往返行程的估计时间(以毫秒为单位):\n     最短 = %dms,最长 = %dms,平均时间 = %dms", conn.RemoteAddr(), sendCount, successCount, failCount, float64(failCount)/float64(sendCount), minTx, maxTx, totalTs/int64(sendCount))
}

func getCommandArgs() {
	flag.Int64Var(&timeout, "w", 1000, "请求超时时长")
	flag.IntVar(&size, "l", 32, "请求发送缓冲区大小，单位字节")
	flag.IntVar(&count, "n", 4, "发送请求参数")
	flag.Parse()
}

// 签名,这里实现的是一个校验和的算法
func checkSum(data []byte) uint16 {
	length := len(data)
	index := 0
	var sum uint32
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		length -= 2
		index += 2
	}

	if length != 0 {
		sum += uint32(data[index])
	}
	hi16 := sum >> 16
	for hi16 != 0 {
		sum = hi16 + uint32(uint16(sum))
		hi16 = sum >> 16
	}
	return uint16(^sum)
}
