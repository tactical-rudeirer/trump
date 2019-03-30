package keyboard

import (
	"trump/pkg/middleware"
	"os/exec"
	"strings"
	"trump/pkg/capture/dev-input"
	"regexp"
	"strconv"
	"log"
	fmt "fmt"
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
			keymap[key{modifiers: strings.Join(strings.Fields(m["modifiers"]), ","), code: uint16(keycode)}] = m["keysym"]
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
		m = m[:len(m)-1]
	}
	if k.Value > 0 {
		//fmt.Println(keymap[key{modifiers: m, code: k.Code}])
		fmt.Print(charMap[keymap[key{modifiers: m, code: k.Code}]])
	}
	return arg
}

var charMap = map[string]string{
	"a": "a",
	"A": "A",
	"adiaeresis": "ä",
	"Adiaeresis": "Ä",
	"ae": "æ",
	"AE": "Æ",
	"Alt": "",
	"AltGr": "",
	"ampersand": "&",
	"apostrophe": "'",
	"Ascii_0": "0",
	"Ascii_1": "1",
	"Ascii_2": "2",
	"Ascii_3": "3",
	"Ascii_4": "4",
	"Ascii_5": "5",
	"Ascii_6": "6",
	"Ascii_7": "7",
	"Ascii_8": "8",
	"Ascii_9": "9",
	"asciitilde": "~",
	"asterisk": "*",
	"at": "@",
	"b": "b",
	"B": "B",
	"backslash": "\\",
	"BackSpace": "\b",
	"bar": "|",
	"Boot": "",
	"braceleft": "{",
	"braceright": "}",
	"bracketleft": "[",
	"bracketright": "]",
	"Break": "",
	"c": "c",
	"C": "C",
	"colon": ":",
	"comma": ",",
	"Compose": "",
	"Console_1": "",
	"Console_10": "",
	"Console_11": "",
	"Console_12": "",
	"Console_13": "",
	"Console_14": "",
	"Console_15": "",
	"Console_16": "",
	"Console_17": "",
	"Console_18": "",
	"Console_19": "",
	"Console_2": "",
	"Console_20": "",
	"Console_21": "",
	"Console_22": "",
	"Console_23": "",
	"Console_24": "",
	"Console_25": "",
	"Console_26": "",
	"Console_27": "",
	"Console_28": "",
	"Console_29": "",
	"Console_3": "",
	"Console_30": "",
	"Console_31": "",
	"Console_32": "",
	"Console_33": "",
	"Console_34": "",
	"Console_35": "",
	"Console_36": "",
	"Console_4": "",
	"Console_5": "",
	"Console_6": "",
	"Console_7": "",
	"Console_8": "",
	"Console_9": "",
	"Control": "",
	"Control_a": "",
	"Control_asciicircum": "",
	"Control_b": "",
	"Control_backslash": "",
	"Control_bracketright": "",
	"Control_c": "",
	"Control_d": "",
	"Control_e": "",
	"Control_f": "",
	"Control_g": "",
	"Control_k": "",
	"Control_l": "",
	"Control_m": "",
	"Control_n": "",
	"Control_o": "",
	"Control_p": "",
	"Control_q": "",
	"Control_r": "",
	"Control_s": "",
	"Control_t": "",
	"Control_u": "",
	"Control_underscore": "",
	"Control_v": "",
	"Control_w": "",
	"Control_x": "",
	"Control_y": "",
	"Control_z": "",
	"CtrlL_Lock": "",
	"d": "d",
	"D": "D",
	"dead_acute": "´",
	"dead_cedilla": "¸",
	"dead_circumflex": "^",
	"dead_diaeresis": "¨",
	"dead_grave": "`",
	"dead_tilde": "~",
	"Decr_Console": "",
	"Delete": "\b",
	"Do": "",
	"dollar": "$",
	"Down": "",
	"e": "e",
	"E": "E",
	"eight": "8",
	"equal": "=",
	"Escape": "",
	"eth": "ð",
	"ETH": "Ð",
	"exclam": "!",
	"f": "f",
	"F": "F",
	"F1": "",
	"F10": "",
	"F11": "",
	"F12": "",
	"F13": "",
	"F14": "",
	"F15": "",
	"F16": "",
	"F17": "",
	"F18": "",
	"F19": "",
	"F2": "",
	"F20": "",
	"F21": "",
	"F22": "",
	"F23": "",
	"F24": "",
	"F25": "",
	"F26": "",
	"F27": "",
	"F28": "",
	"F29": "",
	"F3": "",
	"F30": "",
	"F31": "",
	"F32": "",
	"F33": "",
	"F34": "",
	"F35": "",
	"F36": "",
	"F37": "",
	"F38": "",
	"F39": "",
	"F4": "",
	"F40": "",
	"F41": "",
	"F42": "",
	"F43": "",
	"F44": "",
	"F45": "",
	"F46": "",
	"F47": "",
	"F48": "",
	"F5": "",
	"F6": "",
	"F7": "",
	"F8": "",
	"F9": "",
	"Find": "",
	"five": "5",
	"four": "4",
	"g": "g",
	"G": "G",
	"greater": ">",
	"h": "h",
	"H": "H",
	"Help": "",
	"Hex_0": "0",
	"Hex_1": "1",
	"Hex_2": "2",
	"Hex_3": "3",
	"Hex_4": "4",
	"Hex_5": "5",
	"Hex_6": "6",
	"Hex_7": "7",
	"Hex_8": "8",
	"Hex_9": "9",
	"Hex_A": "A",
	"Hex_B": "B",
	"Hex_C": "C",
	"Hex_D": "D",
	"Hex_E": "E",
	"Hex_F": "F",
	"i": "i",
	"I": "I",
	"Incr_Console": "",
	"Insert": "",
	"j": "j",
	"J": "J",
	"k": "k",
	"K": "K",
	"KeyboardSignal": "",
	"KP_0": "0",
	"KP_1": "1",
	"KP_2": "2",
	"KP_3": "3",
	"KP_4": "4",
	"KP_5": "5",
	"KP_6": "6",
	"KP_7": "7",
	"KP_8": "8",
	"KP_9": "9",
	"KP_Add": "+",
	"KP_Comma": ",",
	"KP_Divide": "/",
	"KP_Enter": "\n",
	"KP_MinPlus": "∓",
	"KP_Multiply": "*",
	"KP_Period": ".",
	"KP_Subtract": "-",
	"l": "l",
	"L": "L",
	"Last_Console": "",
	"Left": "",
	"less": "<",
	"Linefeed": "\n",
	"m": "m",
	"M": "M",
	"Macro": "",
	"masculine": "º",
	"Meta_a": "",
	"Meta_A": "",
	"Meta_ampersand": "",
	"Meta_apostrophe": "",
	"Meta_asciitilde": "",
	"Meta_asterisk": "",
	"Meta_at": "",
	"Meta_b": "",
	"Meta_B": "",
	"Meta_backslash": "",
	"Meta_BackSpace": "",
	"Meta_bar": "",
	"Meta_braceleft": "",
	"Meta_braceright": "",
	"Meta_bracketleft": "",
	"Meta_bracketright": "",
	"Meta_c": "",
	"Meta_C": "",
	"Meta_colon": "",
	"Meta_comma": "",
	"Meta_Control_a": "",
	"Meta_Control_asciicircum": "",
	"Meta_Control_b": "",
	"Meta_Control_backslash": "",
	"Meta_Control_bracketright": "",
	"Meta_Control_c": "",
	"Meta_Control_d": "",
	"Meta_Control_e": "",
	"Meta_Control_f": "",
	"Meta_Control_g": "",
	"Meta_Control_k": "",
	"Meta_Control_l": "",
	"Meta_Control_m": "",
	"Meta_Control_n": "",
	"Meta_Control_o": "",
	"Meta_Control_p": "",
	"Meta_Control_q": "",
	"Meta_Control_r": "",
	"Meta_Control_s": "",
	"Meta_Control_t": "",
	"Meta_Control_u": "",
	"Meta_Control_underscore": "",
	"Meta_Control_v": "",
	"Meta_Control_w": "",
	"Meta_Control_x": "",
	"Meta_Control_y": "",
	"Meta_Control_z": "",
	"Meta_d": "",
	"Meta_D": "",
	"Meta_Delete": "",
	"Meta_dollar": "",
	"Meta_e": "",
	"Meta_E": "",
	"Meta_eight": "",
	"Meta_equal": "",
	"Meta_Escape": "",
	"Meta_exclam": "",
	"Meta_f": "",
	"Meta_F": "",
	"Meta_five": "",
	"Meta_four": "",
	"Meta_g": "",
	"Meta_G": "",
	"Meta_greater": "",
	"Meta_h": "",
	"Meta_H": "",
	"Meta_i": "",
	"Meta_I": "",
	"Meta_j": "",
	"Meta_J": "",
	"Meta_k": "",
	"Meta_K": "",
	"Meta_l": "",
	"Meta_L": "",
	"Meta_less": "",
	"Meta_Linefeed": "",
	"Meta_m": "",
	"Meta_M": "",
	"Meta_minus": "",
	"Meta_n": "",
	"Meta_N": "",
	"Meta_nine": "",
	"Meta_nul": "",
	"Meta_numbersign": "",
	"Meta_o": "",
	"Meta_O": "",
	"Meta_one": "",
	"Meta_p": "",
	"Meta_P": "",
	"Meta_parenleft": "",
	"Meta_parenright": "",
	"Meta_percent": "",
	"Meta_period": "",
	"Meta_plus": "",
	"Meta_q": "",
	"Meta_Q": "",
	"Meta_question": "",
	"Meta_quotedbl": "",
	"Meta_r": "",
	"Meta_R": "",
	"Meta_s": "",
	"Meta_S": "",
	"Meta_semicolon": "",
	"Meta_seven": "",
	"Meta_six": "",
	"Meta_slash": "",
	"Meta_space": "",
	"Meta_t": "",
	"Meta_T": "",
	"Meta_Tab": "",
	"Meta_three": "",
	"Meta_two": "",
	"Meta_u": "",
	"Meta_U": "",
	"Meta_underscore": "",
	"Meta_v": "",
	"Meta_V": "",
	"Meta_w": "",
	"Meta_W": "",
	"Meta_x": "",
	"Meta_X": "",
	"Meta_y": "",
	"Meta_Y": "",
	"Meta_z": "",
	"Meta_Z": "",
	"Meta_zero": "",
	"minus": "-",
	"mu": "µ",
	"n": "n",
	"N": "N",
	"Next": "",
	"nine": "9",
	"nul": "",
	"numbersign": "#",
	"Num_Lock": "",
	"o": "o",
	"O": "O",
	"odiaeresis": "ö",
	"Odiaeresis": "Ö",
	"one": "1",
	"Ooblique": "Ø",
	"ordfeminine": "ª",
	"oslash": "ø",
	"p": "p",
	"P": "P",
	"parenleft": "(",
	"parenright": ")",
	"Pause": "",
	"percent": "%",
	"period": ".",
	"plus": "+",
	"Prior": "",
	"q": "q",
	"Q": "Q",
	"question": "?",
	"quotedbl": "\"",
	"r": "r",
	"R": "R",
	"Remove": "",
	"Return": "\n",
	"Right": "",
	"s": "s",
	"S": "S",
	"Scroll_Backward": "",
	"Scroll_Forward": "",
	"Scroll_Lock": "",
	"Select": "",
	"semicolon": ";",
	"seven": "7",
	"Shift": "",
	"Show_Memory": "",
	"Show_Registers": "",
	"Show_State": "",
	"six": "6",
	"slash": "/",
	"space": " ",
	"ssharp": "ß",
	"t": "t",
	"T": "T",
	"Tab": "\t",
	"thorn": "þ",
	"THORN": "Þ",
	"three": "3",
	"two": "2",
	"u": "u",
	"U": "U",
	"udiaeresis": "ü",
	"Udiaeresis": "Ü",
	"underscore": "_",
	"Up": "",
	"v": "v",
	"V": "V",
	"VoidSymbol": "",
	"w": "w",
	"W": "W",
	"x": "x",
	"X": "X",
	"y": "y",
	"Y": "Y",
	"z": "z",
	"Z": "Z",
	"zero": "0",
}
