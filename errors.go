package socks

import "errors"

var (
	// ErrServerClosed is returned by the Server's Serve, ServeTLS, ListenAndServe,
	// and ListenAndServeTLS methods after a call to Shutdown or Close.
	ErrServerClosed = errors.New("http: Server closed")
)
