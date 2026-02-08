package server

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

const suffix = " ="

var urlStrTemplate = "ws://%s%s"

func assertNil(t *testing.T, some any, label int) {
	if err, ok := some.(error); ok {
		t.Fatalf("%v", err)
	}
	if some != nil {
		t.Fatalf("Something expected to be nil, but currently is not nil! Label: %d", label)
	}
}

func assertStringsEqual(t *testing.T, str1 string, str2 string) {
	if str1 != str2 {
		t.Fatalf("Strings %s and %s expected to be equals", str1, str2)
	}
}

func assertStringContains(t *testing.T, str string, substr string) {
	if !strings.Contains(str, substr) {
		t.Fatalf("String %s expected to contain %s, but it does not!", str, substr)
	}
}

func calculateAnswer(input string) (int, error) {
	numbers := strings.Split(input, " + ")
	n1, err := strconv.Atoi(numbers[0])
	if err != nil {
		return 0, err
	}
	before, found := strings.CutSuffix(numbers[1], suffix)
	if !found {
		return 0, fmt.Errorf("Expected suffix %s not found", suffix)
	}
	n2, err := strconv.Atoi(before)
	if err != nil {
		return 0, err
	}
	return n1 + n2, nil
}

func TestStartExerciseStop(t *testing.T) {
	LaunchServer()
	conn, _, err := websocket.DefaultDialer.Dial(
		fmt.Sprintf(urlStrTemplate, SettingsInstance.Addr, SettingsInstance.WsEndpoint), nil,
	)
	assertNil(t, err, 1)
	_, bytes, err := conn.ReadMessage()
	assertStringsEqual(t, Instruction, string(bytes))
	err = conn.WriteMessage(websocket.TextMessage, []byte(StartWord))
	assertNil(t, err, 2)
	for i := 1; i <= 1000; i++ {
		_, bytes, err := conn.ReadMessage()
		assertNil(t, err, 3)
		answer, err := calculateAnswer(string(bytes))
		assertNil(t, err, 4)
		err = conn.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(answer)))
		assertNil(t, err, 5)
		_, bytes, err = conn.ReadMessage()
		assertNil(t, err, 6)
		assertStringContains(t, string(bytes), strconv.Itoa(i))
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(StopWord))
	assertNil(t, err, 7)
	_, bytes, err = conn.ReadMessage()
	_, bytes, err = conn.ReadMessage()
	assertNil(t, err, 8)
	assertStringsEqual(t, string(bytes), HasBeenStopped)
}
