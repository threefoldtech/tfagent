package pkg

// Agent is the main tft agent entity
type Agent struct { // 	network ConnectionManager
}

// SetupAdminControl for the agent, allowing it to be controlled over an http
// interface. WARNING: this interface has full controll over the agent, including
// wallets and their secrets
func (a *Agent) SetupAdminControl() {
	// r := mux.NewRouter()
	// TODO
}
