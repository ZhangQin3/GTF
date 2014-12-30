package log

import (
	"time"
)

type toggleText struct {
	Title  string
	Switch string
	TextId int64
	Text   string
}

type toggleImage struct {
	Title  string
	Switch string
	ImgId  int64
	ImgSrc string
}

// swtch: toggle switch state, the value is "on", "off"
func ToggleText(title, text, swtch string) {
	if text == "" {
		return
	}
	if swtch == "on" {
		swtch = "block"
	} else {
		swtch = "none"
	}
	log.Output("TOGGLE_TEXT",
		LOnlyFile,
		toggleText{
			title,
			swtch,
			time.Now().UnixNano(),
			text,
		})
}

// swtch: toggle switch state, the value is "on", "off"
func ToggleImage(title, imgSrc, swtch string) {
	if swtch == "on" {
		swtch = "block"
	} else {
		swtch = "none"
	}
	log.Output("TOGGLE_IMAGE",
		LOnlyFile,
		toggleImage{
			title,
			swtch,
			time.Now().UnixNano(),
			imgSrc,
		})
}
