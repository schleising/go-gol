package main

import (
	"fmt"
	"math"
	"math/rand"
	"syscall/js"
)

var (
	window     = js.Global()
	canvas     js.Value
	context    js.Value
	windowSize struct{ width, height float64 }
	random     *rand.Rand
	newBoard   *BufferedBoard
)

func main() {
	newBoard = Initialise(25, 25)

	setupCanvas()
	setupRenderLoop()

	runForever := make(chan bool)
	<-runForever
}

func setupCanvas() {
	document := window.Get("document")
	fmt.Println("Doc    : ", document)

	canvas = document.Call("getElementById", "canvas")
	context = canvas.Call("getContext", "2d")

	fmt.Println("Canvas : ", canvas)
	fmt.Println("Context: ", context)

	updateWindowSizeJSCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resetWindowSize()
		return nil
	})
	window.Call("addEventListener", "resize", updateWindowSizeJSCallback)
	resetWindowSize()

	fmt.Println("Setup Canvas")
}

func resetWindowSize() {
	// https://stackoverflow.com/a/8486324/1478636
	squareWindowSide := math.Min(window.Get("innerWidth").Float(), window.Get("innerHeight").Float())
	windowSize.width = squareWindowSide
	windowSize.height = squareWindowSide
	canvas.Set("width", windowSize.width)
	canvas.Set("height", windowSize.height)
}

func setupRenderLoop() {
	var renderJSCallback js.Func
	renderJSCallback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		draw()
		update()
		window.Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			window.Call("requestAnimationFrame", renderJSCallback)
			return nil
		}), 75)
		return nil
	})
	window.Call("requestAnimationFrame", renderJSCallback)
	fmt.Println("Setup Render Loop")
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
	fmt.Println("Clear: ", windowSize.width, windowSize.height)
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
