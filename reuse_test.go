package reuse

import (
	ctx "context"
	"net"
	"sync"
	"testing"
	"time"
)

func TestReuseServerPort(t *testing.T) {
	network := "tcp"
	addr := ":9999"
	msg := "Hello World"
	wg := sync.WaitGroup{}
	listeners := make([]net.Listener, 2)
	for i := 0; i < 2; i++ {
		if listener, err := ListenConfig.Listen(ctx.Background(), network, addr); err != nil {
			t.Errorf("listen %s fail: %s", addr, err)
		} else {
			listeners[i] = listener
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					conn, err := listener.Accept()
					if err != nil {
						break
					}
					wg.Add(1)
					go func(conn net.Conn) {
						defer wg.Done()
						defer conn.Close()
						buf := make([]byte, 1024)
						n, err := conn.Read(buf)
						if err != nil {
							t.Error(err)
						}
						_, err = conn.Write(buf[:n])
						if err != nil {
							t.Error(err)
						}
					}(conn)
				}
			}()
		}
	}
	time.Sleep(time.Millisecond * 100)
	localPort := 8888
	d := net.Dialer{LocalAddr: &net.TCPAddr{Port: localPort}, Control: Control}
	conn, err := d.Dial(network, addr)
	if err != nil {
		t.Error("dial failed:", err)
		return
	}
	if _, err := conn.Write([]byte(msg)); err != nil {
		t.Error(err)
		return
	}
	buf := make([]byte, 1024)
	if n, err := conn.Read(buf); err != nil {
		t.Error(err)
		return
	} else if n != len(msg) {
		t.Errorf("%d %d", n, len(msg))
	}
	conn.Close()
	for i := 0; i < 2; i++ {
		listeners[i].Close()
	}
	wg.Wait()
}

func TestReuseClientPort(t *testing.T) {
	network := "tcp"
	addr1 := ":9997"
	addr2 := ":9998"
	msg := "Hello World"
	wg := sync.WaitGroup{}
	listener1, err := ListenConfig.Listen(ctx.Background(), network, addr1)
	if err != nil {
		t.Errorf("listen %s fail: %s", addr1, err)
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				conn, err := listener1.Accept()
				if err != nil {
					break
				}
				wg.Add(1)
				go func(conn net.Conn) {
					defer wg.Done()
					defer conn.Close()
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err != nil {
						t.Error(err)
					}
					_, err = conn.Write(buf[:n])
					if err != nil {
						t.Error(err)
					}
				}(conn)
			}
		}()
	}
	listener2, err := ListenConfig.Listen(ctx.Background(), network, addr2)
	if err != nil {
		t.Errorf("listen %s fail: %s", addr2, err)
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				conn, err := listener2.Accept()
				if err != nil {
					break
				}
				wg.Add(1)
				go func(conn net.Conn) {
					defer wg.Done()
					defer conn.Close()
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err != nil {
						t.Error(err)
					}
					_, err = conn.Write(buf[:n])
					if err != nil {
						t.Error(err)
					}
				}(conn)
			}
		}()
	}
	time.Sleep(time.Millisecond * 100)
	localPort := 8888
	d := net.Dialer{LocalAddr: &net.TCPAddr{Port: localPort}, Control: Control}
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := d.Dial(network, addr1)
		time.Sleep(time.Millisecond * 100)
		if err != nil {
			t.Error("dial failed:", err)
			return
		}
		if _, err := conn.Write([]byte(msg)); err != nil {
			t.Error(err)
			return
		}
		buf := make([]byte, 1024)
		if n, err := conn.Read(buf); err != nil {
			t.Error(err)
			return
		} else if n != len(msg) {
			t.Errorf("%d %d", n, len(msg))
		}
		time.Sleep(time.Millisecond * 500)
		conn.Close()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := d.Dial(network, addr2)
		time.Sleep(time.Millisecond * 100)
		if err != nil {
			t.Error("dial failed:", err)
			return
		}
		if _, err := conn.Write([]byte(msg)); err != nil {
			t.Error(err)
			return
		}
		buf := make([]byte, 1024)
		if n, err := conn.Read(buf); err != nil {
			t.Error(err)
			return
		} else if n != len(msg) {
			t.Errorf("%d %d", n, len(msg))
		}
		time.Sleep(time.Millisecond * 500)
		conn.Close()
	}()
	time.Sleep(time.Second)
	listener1.Close()
	listener2.Close()
	wg.Wait()
}
