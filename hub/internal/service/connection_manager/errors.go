package connection_manager

import "errors"

var ErrConnectionClose = errors.New("connection close")

var ErrNotFoundConn = errors.New("agent connection not found")
