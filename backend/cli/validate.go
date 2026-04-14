package cli

import "errors"

var ErrNoDeckOpen = errors.New("no deck opened, please open deck first")
var ErrNoLogin = errors.New("not login, please login first")
