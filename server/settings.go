package server

var SettingsInstance = DefaultSettings()

type Settings struct {
	Addr           string
	WsEndpoint     string
	MaxConnCount   int
	SecondsToSolve int
}

func (settings *Settings) ChangeAddr(addr string) *Settings {
	settings.Addr = addr
	return settings
}

func (settings *Settings) ChangeWsEndpoint(wsEndpoint string) *Settings {
	settings.WsEndpoint = wsEndpoint
	return settings
}

func (settings *Settings) ChangeMaxConnCount(maxConnCount int) *Settings {
	connCount = make(chan struct{}, maxConnCount)
	settings.MaxConnCount = maxConnCount
	return settings
}

func (settings *Settings) ChangeSecondsToSolve(secondsToSolve int) *Settings {
	settings.SecondsToSolve = secondsToSolve
	return settings
}

func DefaultSettings() *Settings {
	return &Settings{
		Addr:           addrDefault,
		WsEndpoint:     WsEndpointDefault,
		MaxConnCount:   maxConnCountDefault,
		SecondsToSolve: secondsToSolveDefault,
	}
}
