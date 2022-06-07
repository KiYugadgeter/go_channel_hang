package investment

import (
	"encoding/json"
	"fmt"
	"investment/calc"
	"log"
	"net"
	"strconv"
	"sync"
)

type Servable interface {
	int | calc.Candle
}

func ByteLength(b []byte) ([]byte, int) {
	buff := make([]byte, 8)
	b_len := len(b)
	b_len64 := uint64(b_len)
	for i := 0; i < 8; i++ {
		buff[i] = byte((b_len64 >> ((7 - i) * 8))) & byte(0xff)
	}
	buff = append(buff, b...)
	return buff, b_len
}

type CandleServer[T Servable] struct {
	Conns *sync.Map
	C     chan T
}

func (t *CandleServer[T]) Run(addr string, port int) {
	defer func() {
		t.Conns.Range(func(k, v any) bool {
			c := k.(*net.Conn)
			(*c).Close()
			return true
		})
	}()
	go func() {
		var count int = 0
		for {
			select {
			case v := <-t.C:
				count += 1
				fmt.Println("candle_server", "case", count)

				b, err := json.Marshal(v)
				if err != nil {
					fmt.Println("Errhang ", err)
				}
				result_buff, _ := ByteLength(b)
				t.Conns.Range(func(k, v any) bool {
					c := k.(*net.Conn)
					_, err := (*c).Write(result_buff)
					if err != nil {
						(*c).Close()
						t.Conns.Delete(k)
					}
					return true
				})
				fmt.Println("candle_server", "after range")
			}
		}
	}()
	ln, err := net.Listen("tcp", addr+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("Listen Error", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Accept Error", err)
		}
		t.Conns.Store(&conn, struct{}{})
	}

}

func NewCandleServer[T Servable](c chan T) *CandleServer[T] {
	return &CandleServer[T]{C: c, Conns: &sync.Map{}}
}
