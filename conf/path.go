package conf

import (
	"path/filepath"
)

var (
	ROOT_PATH       string
	CONTROLLER_PATH string
	CONF_PATH       string
	TPL_PATH        string
	STATIC_PATH     string
	IMG_PATH        string
	JS_PATH         string
	CSS_PATH        string
)

func init() {
	ROOT_PATH, _ = filepath.Abs("./")
	CONTROLLER_PATH = ROOT_PATH + "/app/admin/controller"
	CONF_PATH = ROOT_PATH + "/conf"
	TPL_PATH = ROOT_PATH + "/core/lazy/tpl"
	STATIC_PATH = ROOT_PATH + "/static"
	IMG_PATH = STATIC_PATH + "/img"
	JS_PATH = STATIC_PATH + "/js"
	CSS_PATH = STATIC_PATH + "/css"

}
