package modules

import (
	"github.com/templwind/sass-starter/modules/app"
	"github.com/templwind/sass-starter/modules/www"
)

// Module registry for all modules
// This	is where you register all the modules used in your application
var registry = map[string]Module{
	"app": app.Module(),
	"www": www.Module(),
}
