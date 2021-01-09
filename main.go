package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"strconv"
	"syscall/js"
)

var (
	window     = js.Global()
	canvas     js.Value
	context    js.Value
	windowSize struct{ width, height float64 }
	random     *rand.Rand
	newBoard   *BufferedBoard
	messages   chan string
)

func main() {
	messages = make(chan string)
	go func() {
		for message := range messages {
			fmt.Println(message)
		}
	}()
	messages <- "GoL::main"
	newBoard = Initialise(25, 25)

	frameInterval := setupCanvas()
	setupRenderLoop(frameInterval)

	messages <- "GoL::Running"
	runForever := make(chan bool)
	<-runForever
}

func setupCanvas() int64 {
	document := window.Get("document")

	canvas = document.Call("getElementById", "canvas")
	context = canvas.Call("getContext", "2d")

	pageURL := document.Get("location").Get("href").String()
	params := parseURLQueryParams(pageURL)
	messages <- fmt.Sprintf("WASM::setupCanvas Params: %+v", params)

	updateWindowSizeJSCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resetWindowSize()
		return nil
	})
	window.Call("addEventListener", "resize", updateWindowSizeJSCallback)
	resetWindowSize()

	messages <- "GoL::Setup Canvas"

	return params["frameInterval"]
}

func parseURLQueryParams(pageURL string) map[string]int64 {
	params := map[string]int64{
		"frameInterval": 75,
	}
	parse, e := url.Parse(pageURL)
	if e != nil {
		return params
	}
	for paramKey, paramValues := range parse.Query() {
		if len(paramValues) > 0 {
			if value, e := strconv.ParseInt(paramValues[0], 10, 64); e == nil {
				params[paramKey] = value
				messages <- fmt.Sprintln("Key: ", paramKey)
				messages <- fmt.Sprintln("Val: ", value)
			} else {
				params[paramKey] = -1
			}
		}
	}
	return params
}

func resetWindowSize() {
	// https://stackoverflow.com/a/8486324/1478636
	squareWindowSide := math.Min(window.Get("innerWidth").Float(), window.Get("innerHeight").Float())
	windowSize.width = squareWindowSide
	windowSize.height = squareWindowSide
	canvas.Set("width", windowSize.width)
	canvas.Set("height", windowSize.height)
}

func setupRenderLoop(frameInterval int64) {
	var renderJSCallback js.Func
	renderJSCallback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		draw()
		update()
		window.Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			window.Call("requestAnimationFrame", renderJSCallback)
			return nil
		}), frameInterval)
		return nil
	})
	window.Call("requestAnimationFrame", renderJSCallback)
	messages <- "GoL::Setup Render Loop"
}

func draw() {
	clearCanvas()
	color := "#ffffff"
	strokeStyle(color)
	fillStyle(color)
	lineWidth(0.75)
	padding := float64(4)

	squareSize := math.Min(windowSize.width/float64(newBoard.cols), windowSize.height/float64(newBoard.rows))
	side := squareSize - padding*2

	for row := 0; row < newBoard.rows; row++ {
		for column := 0; column < newBoard.cols; column++ {
			x := float64(column)*squareSize + padding
			y := float64(row)*squareSize + padding
			drawStrokeRect(x, y, side, side)
			if newBoard.GetState(row, column) {
				drawFillRect(x, y, side, side)
			}
		}
	}
}

func update() {
	newBoard.Iterate()
}

func clearCanvas() {
	context.Call("clearRect", 0, 0, windowSize.width, windowSize.height)
}

func strokeStyle(style string) {
	context.Set("strokeStyle", style)
}

func fillStyle(style string) {
	context.Set("fillStyle", style)
}

func lineWidth(width float64) {
	context.Set("lineWidth", width)
}

func drawStrokeRect(x, y, width, height float64) {
	context.Call("strokeRect", x, y, width, height)
}

func drawFillRect(x, y, width, height float64) {
	context.Call("fillRect", x, y, width, height)
}
