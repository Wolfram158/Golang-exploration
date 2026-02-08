package server

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader              = websocket.Upgrader{}
	connCount             = make(chan struct{}, maxConnCountDefault)
	mux                   = http.NewServeMux()
	instructionBytes      = []byte(Instruction)
	alreadyStartedBytes   = []byte(AlreadyStarted)
	haventStartedYetBytes = []byte(HaventStartedYet)
	wrongAnswerBytes      = []byte(WrongAnswer)
	reminderBytes         = []byte(Reminder)
	hasBeenStoppedBytes   = []byte(HasBeenStopped)
)

type Session struct {
	answer    chan int
	close     chan struct{}
	start     chan struct{}
	stop      chan struct{}
	isStarted bool
	total     int
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case connCount <- struct{}{}:
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		go handleWsConnection(conn)
	default:
		http.Error(w, ServerIsFull, http.StatusBadRequest)
		return
	}
}

func handleWsConnection(
	conn *websocket.Conn,
) {
	defer func() {
		<-connCount
	}()
	defer conn.Close()
	conn.WriteMessage(websocket.TextMessage, instructionBytes)
	var ctrl = Session{
		answer:    make(chan int),
		close:     make(chan struct{}),
		start:     make(chan struct{}),
		stop:      make(chan struct{}),
		isStarted: false,
		total:     0,
	}
	go ctrl.dispatchMsg(conn)
	ctrl.handleMsg(conn)
}

func (ctrl *Session) dispatchMsg(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		var str = string(p)
		if err != nil {
			ctrl.close <- struct{}{}
			return
		}
		switch str {
		case StartWord:
			if !ctrl.isStarted {
				ctrl.isStarted = true
				ctrl.start <- struct{}{}
			} else {
				conn.WriteMessage(websocket.TextMessage, alreadyStartedBytes)
			}

		case StopWord:
			if ctrl.isStarted {
				ctrl.isStarted = false
				ctrl.stop <- struct{}{}
			} else {
				conn.WriteMessage(websocket.TextMessage, haventStartedYetBytes)
			}
		default:
			num, err := strconv.Atoi(str)
			if err != nil {
				if ctrl.isStarted {
					conn.WriteMessage(websocket.TextMessage, fmt.Appendf(nil, UnexpectedMsgTemplate, str))
				} else {
					conn.WriteMessage(websocket.TextMessage, reminderBytes)
				}
			} else {
				ctrl.answer <- num
			}
		}
	}
}

func (ctrl *Session) handleMsg(
	conn *websocket.Conn,
) {
loop1:
	for {
		select {
		case <-ctrl.close:
			break loop1

		case <-ctrl.start:
		loop2:
			for {
				n1 := rand.Intn(100)
				n2 := rand.Intn(100)
				conn.WriteMessage(websocket.TextMessage, fmt.Appendf(nil, exerciseTemplate, n1, n2))
				var ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(SettingsInstance.SecondsToSolve))
			loop3:
				for {
					select {
					case <-ctx.Done():
						ctrl.total--
						conn.WriteMessage(websocket.TextMessage, fmt.Appendf(nil, TimeIsOutTemplate, ctrl.total))
						break loop3
					case ans := <-ctrl.answer:
						if ans == n1+n2 {
							ctrl.total++
							conn.WriteMessage(websocket.TextMessage, fmt.Appendf(nil, SuccessTemplate, ctrl.total))
							break loop3
						} else {
							conn.WriteMessage(websocket.TextMessage, wrongAnswerBytes)
						}
					case <-ctrl.stop:
						conn.WriteMessage(websocket.TextMessage, hasBeenStoppedBytes)
						cancel()
						break loop2
					}
				}
			_:
				cancel()
			}
		}
	}
}

func LaunchServer() chan error {
	rand.New(rand.NewSource(time.Now().Unix()))
	mux.HandleFunc(SettingsInstance.WsEndpoint, wsHandler)
	listener, err := net.Listen(tcp, SettingsInstance.Addr)
	if err != nil {
		log.Fatalf(couldntCreateListener)
	}
	var errChan = make(chan error)
	go func() {
		if err = http.Serve(listener, mux); err != nil {
			errChan <- err
		}
	}()
	return errChan
}
