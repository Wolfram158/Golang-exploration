package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader  = websocket.Upgrader{}
	connCount = make(chan struct{}, 1)
	mux       = http.NewServeMux()
	total     = 0
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case connCount <- struct{}{}:
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		go handleWsConnection(conn)
	default:
		http.Error(w, "Server is full", 400)
		return
	}
}

func handleWsConnection(conn *websocket.Conn) {
	defer func() {
		<-connCount
	}()
	defer conn.Close()
	conn.WriteMessage(1, []byte("Send 'Start' to start exercise. Send 'Stop' to stop exercise."))
	var answer = make(chan int)
	var close = make(chan struct{})
	var isStarted = false
	var start = make(chan struct{})
	var stop = make(chan struct{})
	go func() {
		for {
			_, p, err := conn.ReadMessage()
			var str = string(p)
			if err == nil {
				switch str {
				case "Start":
					if !isStarted {
						isStarted = true
						start <- struct{}{}
					} else {
						conn.WriteMessage(1, []byte("Already started!"))
					}

				case "Stop":
					if isStarted {
						isStarted = false
						stop <- struct{}{}
					} else {
						conn.WriteMessage(1, []byte("Haven't started yet!"))
					}
				default:
					num, err := strconv.Atoi(str)
					if err != nil {
						if isStarted {
							conn.WriteMessage(1, fmt.Appendf(nil, "Number expected, but given: %s", str))
						} else {
							conn.WriteMessage(1, []byte("Send 'Start' to start exercise!"))
						}
					} else {
						answer <- num
					}
				}
			} else {
				close <- struct{}{}
				return
			}
		}
	}()
loop1:
	for {
		select {
		case <-close:
			break loop1

		case <-start:
		loop2:
			for {
				n1 := rand.Intn(100)
				n2 := rand.Intn(100)
				conn.WriteMessage(1, fmt.Appendf(nil, "%d + %d = ", n1, n2))
				var ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(10))
				defer cancel()
			loop3:
				for {
					select {
					case <-ctx.Done():
						total--
						conn.WriteMessage(1, fmt.Appendf(nil, "Time is out! Balance: %d", total))
						break loop3
					case ans := <-answer:
						if ans == n1+n2 {
							total++
							conn.WriteMessage(1, fmt.Appendf(nil, "Success! Balance: %d", total))
							break loop3
						} else {
							conn.WriteMessage(1, []byte("Wrong answer!"))
						}
					case <-stop:
						conn.WriteMessage(1, []byte("Exercise has been stopped!"))
						break loop2
					}
				}
			}
		}
	}
}

func main() {
	rand.New(rand.NewSource(time.Now().Unix()))
	mux.HandleFunc("/ws", wsHandler)
	go func() {
		err := http.ListenAndServe(":8080", mux)
		if err != nil {
			log.Fatal(err)
		}
	}()
	select {}
}
