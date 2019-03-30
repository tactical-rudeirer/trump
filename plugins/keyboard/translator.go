package keyboard

import (
	"trump/pkg/middleware"
	"os/exec"
	"strings"
	"log"
	"regexp"
	"strconv"
	"trump/pkg/capture/dev-input"
	"fmt"
)

var Translator = middleware.Plugin{
	Process:  ProcessMsg,
	Priority: 0,
	Init:     InitTranslator,
	Shutdown: func() {},
}

func InitTranslator() {
	keymapRaw, err := exec.Command("dumpkeys").Output()
	if err != nil {
		log.Fatalf("failed to load keymap: %v", err)
	}
	keymapString := string(keymapRaw)
	parseKeymap(keymapString)
}

var KEYCODE = regexp.MustCompile("\\s*(?P<modifiers>((shift|altgr|control|alt|shiftl|shiftr|ctrll|ctrlr|capsshift)\\s+)*)keycode\\s+(?P<keycode>[\\d]+)\\s+=\\s+\\+?(?P<keysym>\\w+)\\s*")

func mapSubexpNames(m, n []string) map[string]string {
	m, n = m[1:], n[1:]
	r := make(map[string]string, len(m))
	for i, _ := range n {
		r[n[i]] = m[i]
	}
	return r
}

type key struct {
	modifiers string
	code      uint16
}

var keymap map[key]string

func parseKeymap(str string) {
	keymap = map[key]string{}
	for _, l := range strings.Split(str, "\n") {
		if KEYCODE.MatchString(l) {
			m := mapSubexpNames(KEYCODE.FindStringSubmatch(l), KEYCODE.SubexpNames())
			keycode, err := strconv.ParseInt(m["keycode"], 10, 16)
			if err != nil {
				log.Printf("failed to parse keycode %s: %v", m["keycode"], err)
			}
			keymap[key{modifiers: strings.Join(strings.Fields(m["modifiers"]),","), code: uint16(keycode)}] = m["keysym"]
		}
	}
}

var SHIFT = false
var ALTGR = false
var CONTROL = false
var ALT = false
var SHIFTL = false
var SHIFTR = false
var CTRLL = false
var CTRLR = false
var CAPSSHIFT = false

func ProcessMsg(arg interface{}) interface{} {
	k := arg.(dev_input.DevInputData)
	if k.Value > 2 || k.Value < 0 {
		return nil
	}
	switch k.Code {
	case 29:
		CONTROL = k.Value > 0
		CTRLL = k.Value > 0

	case 42:
		SHIFT = k.Value > 0
		SHIFTL = k.Value > 0

	case 54:
		SHIFT = k.Value > 0
		SHIFTR = k.Value > 0

	case 56:
		ALT = k.Value > 0

	case 58:
		CAPSSHIFT = k.Value > 0

	case 97:
		CONTROL = k.Value > 0
		CTRLR = k.Value > 0

	case 100:
		ALTGR = k.Value > 0
	}
	m := ""
	if SHIFT {
		m += "shift,"
	}
	if ALTGR {
		m += "altgr,"
	}
	if CONTROL {
		m += "control,"
	}
	if ALT {
		m += "alt,"
	}
	if SHIFTL {
		m += "shiftl,"
	}
	if SHIFTR {
		m += "shiftr,"
	}
	if CTRLL {
		m += "ctrll,"
	}
	if CTRLR {
		m += "ctrlr,"
	}
	if CAPSSHIFT {
		m += "capsshift,"
	}
	if len(m) > 0 {
		m = m[:len(m) - 1]
	}
	if k.Value > 0 {
		fmt.Println(keymap[key{modifiers: m, code: k.Code}])
	}
	return arg
}
