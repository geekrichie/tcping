package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"syscall"
	"time"
)

var host string
var port int
var count int
func init() {
	flag.StringVar(&host,"host", "localhost", " hostname for your testing server")
	flag.IntVar(&port, "port", 80, "port number")
	flag.IntVar(&count, "count", 4, "tcp ping count")
}

func main() {
	flag.Parse()

	addr, err := net.LookupIP(host)
	if err != nil {
		fmt.Printf(" gethostname(%s) error : %s\n", host,err)
		return
	}
	var (
		fail = 0
		success = 0
		maxtime float64= -1 << 63
		mintime float64= 1<<63 - 1
		avgtime = 0.0
	)
	for cnt := count;cnt >0;cnt = cnt-1 {
		timecost,ip,err := connect(addr)
		if err != nil {
			fmt.Printf("连接%s失败:%s\n", host, err.Error())
			fail++
			continue
		}
		success++
		fmt.Printf("来自%s的回复：时间=%dms\n", ip, timecost)
		if float64(timecost) > maxtime {
			maxtime = float64(timecost)
		}
		if float64(timecost) < mintime {
			mintime = float64(timecost)
		}
		avgtime += float64(timecost)
	}
	fmt.Printf("--- %s:%d ping statistics ---\n", host, port);
	fmt.Printf("%d responses, %d ok, %3.2f%% failed\n", count, success, 100.0*float64(fail)/float64(count))
	if success   > 0 {
		fmt.Printf("round-trip min/avg/max = %.1f/%.1f/%.1f ms\n", mintime, float64(avgtime)/float64(success), maxtime);
	}
}

func connect(addrs []net.IP) (time.Duration, string, error){
	var ipAddr string
	for _,addr := range addrs {
		//创建socket文件描述符
		syscall.ForkLock.RLock()
		fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM,0)
		syscall.ForkLock.RUnlock()
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err  = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
			fmt.Println("setopt")
			syscall.Close(fd)
			continue
		}
		sa := &syscall.SockaddrInet4{
			Port: port,
		}

		copy(sa.Addr[:], addr[len(addr)-4:])
		var start = time.Now()
		err = handleConnect(fd, sa)
		if err != nil {
			syscall.Close(fd)
			continue
		}
		syscall.Close(fd)
		last := time.Since(start)
		ipAddr = addr.String()
		return last*time.Nanosecond/time.Millisecond, ipAddr, nil
	}
	return 0, "",errors.New("connect time out")

}

func handleConnect(fd syscall.Handle, sa syscall.Sockaddr) error{
	timePeriod := time.Duration(5)* time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timePeriod)
	connectChan := make(chan error,1)
	go func (){
		err := syscall.Connect(fd, sa)
		connectChan <- err
	}()
	defer cancel()
	select {
	case err := <- connectChan:
		return err
	case <- ctx.Done():
		return ctx.Err()

	}
}
