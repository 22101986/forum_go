package main

import (
	"net/http"
	"time"
)

// Constants for the Window fields
const CountMaxReq int = 20
const WindowSize time.Duration = 1 * time.Second

type Window struct {
	MaxReq int
	Req    []time.Time
	Size   time.Duration
}

// create the Sliding Window
func NewWindow(CountMaxReq int, WindowSize time.Duration) *Window {
	return &Window{
		MaxReq: CountMaxReq,
		Req:    []time.Time{},
		Size:   WindowSize,
	}
}

// function that use the window and allow, or not, the request
func (win *Window) Allow() bool {
	now := time.Now()
	for len(win.Req) > 0 && now.Sub(win.Req[0]) > win.Size { // delete all old enough requests in the Window
		win.Req = win.Req[1:]
	}
	if len(win.Req) < win.MaxReq { // add the new request in the Window
		win.Req = append(win.Req, now)
		return true
	}
	return false
}

// check the number of requests before redirecting on the handler
func VerifReq(win *Window, HandlerFunc func(http.ResponseWriter, *http.Request, *StructSquare, *StructData), w http.ResponseWriter, r *http.Request, sSquared *StructSquare, sData *StructData) {
	if !win.Allow() {
		ErrorHandler(w, r, http.StatusTooManyRequests, "too many requests")
	} else {
		HandlerFunc(w, r, sSquared, sData)
	}
}
