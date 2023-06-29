package proxy

// Run start proxy server
func Run() {
	// init casters from config
	initCasters()

	// init ntrip server
	initNtripServer()
}
