package server

const (
	StartWord             = "Start"
	StopWord              = "Stop"
	ServerIsFull          = "Server is full!"
	TimeIsOutTemplate     = "Time is out! Balance: %d"
	SuccessTemplate       = "Success! Balance: %d"
	UnexpectedMsgTemplate = "Number or command 'Stop' expected, but given: %s"
	tcp                   = "tcp"
	couldntCreateListener = "Couldn't create listener!"
	maxConnCountDefault   = 100
	secondsToSolveDefault = 10
	WsEndpointDefault     = "/ws"
	addrDefault           = "localhost:8080"
	Instruction           = "Send 'Start' to start exercise. Send 'Stop' to stop exercise."
	AlreadyStarted        = "Already started!"
	HaventStartedYet      = "Haven't started yet!"
	WrongAnswer           = "Wrong answer!"
	Reminder              = "Send 'Start' to start exercise!"
	HasBeenStopped        = "Exercise has been stopped!"
	exerciseTemplate      = "%d + %d ="
)
