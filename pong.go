package main

// Homework
// make gameplay more interesting (I did this by adding acceleration to the ball, but I want to do something else too)
// 2 player vs playing computer
// Scale draws to resized window
// Load bitmaps for our paddles/ball

import (
	"math"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type gameState int

// Game states as enumerative iota
// https://blog.learngoprogramming.com/golang-const-type-enums-iota-bc4befd096d3
const (
	start gameState = iota
	play  gameState = iota
)

// default state
var state = start

// default color r, g, b
var gameColor = color{0, 255, 0}

// Current score can only go to 10 for the moment.
var nums = [][]byte{
	{1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 1, 1, 1}, // 0
	{1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 1, 1}, // 1
	{1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1}, // 2
	{1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1}, // 3
	{1, 0, 1, 1, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1}, // 4
	{1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 1}, // 5
	{1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1}, // 6
	{1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1}, // 7
	{1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1}, // 8
	{1, 1, 1, 1, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1}, // 9
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, // Win
}

// default window width and hieght. Might be a good idea to alow this to be resized.
const winWidth, winHeight int = 800, 600

// color struct defines the RGB spectrum of colors
type color struct {
	r, g, b byte
}

// pos structure for defining the X & Y coordinates on a grid
type pos struct {
	x, y float32
}

// ball struct defines the structure for the ball attributes
type ball struct {
	pos
	radius, xv, yv, accelleration float32
	color                         color
}

// lerp function provids linear interpolation, I use it to figure out a good
// position to place the scores
func lerp(a, b float32, pct float32) float32 {
	return a + pct*(b-a)
}

// drawBackground - Simple function to draw a dotted line down the middle of the screen
func drawBackground(pixels []byte) {
	x := int(winWidth) / 2.0

	var on bool // a flag that turns the drawing on and off
	// only draws a stripped line across middle of screen
	for y := 0; y < winHeight; y++ {
		if y%10 == 0 {
			on = !on
		}
		if on {
			setPixel(x, y, gameColor, pixels)
		}
	}
	// Maybe create a different background and include goals?
}

// drawNumbers draws the numbers given a certain position
func drawNumber(pos pos, color color, size, num int, pixels []byte) {
	startX := int(pos.x) - (size*3)/2
	startY := int(pos.y) - (size*5)/2

	for i, v := range nums[num] {
		if v == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		startX += size
		if (i+1)%3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
}

// draw function takes a ball receiver and pixel grid as input. Based where
// the ball was previously the ball is draw as a circle around its X & Y coords.
// This is not optimized since it needs to always draw a square and not a round circle.
func (ball *ball) draw(pixels []byte) {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x+x), int(ball.y+y), ball.color, pixels)
			}
		}
	}
}

// getCenter gets the center of the window from default
func getCenter() pos {
	return pos{float32(winWidth) / 2, float32(winHeight) / 2}
}

// updates the position of the ball in the pixels space based on velocity in x & y direction
// fixed screen rate
func (ball *ball) update(leftPaddle, rightPaddle *paddle, elapseTime float32) {
	ball.accelleration *= 1.0005
	ball.x += ball.xv * elapseTime * ball.accelleration
	ball.y += ball.yv * elapseTime * ball.accelleration

	// ball collissions for when the ball hits the top and bottom of the screen
	// takes into consideration the edge of the ball when it collides with the
	// wall.
	if ball.y-ball.radius < 0 {
		ball.yv = -ball.yv
		ball.y = ball.y + ball.radius/2.0 + ball.radius
	} else if ball.y+ball.radius > float32(winHeight) {
		ball.yv = -ball.yv
		ball.y = ball.y - ball.radius/2.0 - ball.radius
	}

	// if the ball reaches the other goal, updates score, resets ball pos
	// and sets state of game to the start.
	if ball.x < 0 {
		rightPaddle.score++
		ball.pos = getCenter()
		state = start
	} else if ball.x > float32(winWidth) {
		leftPaddle.score++
		ball.pos = getCenter()
		state = start
	}

	// Collision logic when ball collides with
	if (ball.x - ball.radius) < leftPaddle.x+leftPaddle.w/2 {
		if ball.y > leftPaddle.y-leftPaddle.h/2 && ball.y < leftPaddle.y+leftPaddle.h/2 {
			// Think I should add some logic to add different areas on the paddle
			// That make the ball fly in a new direction
			ball.xv = -ball.xv
			ball.x = leftPaddle.x + leftPaddle.w/2.0 + ball.radius
		}
	}
	// Same collision logic but for the right paddle
	if (ball.x + ball.radius) > rightPaddle.x-rightPaddle.w/2 {
		if ball.y > rightPaddle.y-rightPaddle.h/2 && ball.y < rightPaddle.y+rightPaddle.h/2 {
			ball.xv = -ball.xv
			ball.x = rightPaddle.x - rightPaddle.w/2.0 - ball.radius
		}
	}
}

// paddle structure for storing all information about a particular paddle.
type paddle struct {
	pos
	w, h, speed, acceleration float32
	score                     int
	color                     color
}

// draw takes a paddle as a receiver, and a slice of bytes to draw
func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x - paddle.w/2)
	startY := int(paddle.y - paddle.h/2)

	for y := 0; y < int(paddle.h); y++ {
		for x := 0; x < int(paddle.w); x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}
	//
	numX := lerp(paddle.x, getCenter().x, 0.2)
	drawNumber(pos{numX, 35}, paddle.color, 10, paddle.score, pixels)
}

// update takes a paddle receiver and a keystate/controller axis parameter as well
// as elapseTime variable to determine the direction and speed of the paddle movement
func (paddle *paddle) update(keyState []uint8, controllerAxis int16, elapseTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 && (paddle.y-paddle.h/2) > 0 {
		paddle.y -= paddle.speed * elapseTime
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 && (paddle.y+paddle.h/2) < float32(winHeight) {
		paddle.y += paddle.speed * elapseTime
	}

	if math.Abs(float64(controllerAxis)) > 1500 {
		pct := float32(controllerAxis) / 32767.0
		paddle.y += paddle.speed * pct * elapseTime
	}
}

// aiUpdate is the computers logic for going after a ball. I made it good in the begining
// when the ball is at a slower speed, but the reset from the acceleration everytime
// the ball passes its field of vision keeps it from winning everytime.
func (paddle *paddle) aiUpdate(ball *ball, elapseTime float32) {
	if ball.y > (paddle.y+paddle.h/2) && paddle.y < float32(winHeight) {
		paddle.acceleration *= 1.04
		paddle.y += paddle.speed * elapseTime * paddle.acceleration
	} else if ball.y < (paddle.y-paddle.h/2) && paddle.y > 0 {
		paddle.acceleration *= 1.04
		paddle.y -= paddle.speed * elapseTime * paddle.acceleration
	} else {
		paddle.acceleration = 1
	}
}

// clear is a helper function to reset the pixel grid
func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

// setPixel - Takes coords and a color and a pixel grid and draws a pixel at
// a specific location.
func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func main() {

	// creates the window with the default window width and height, I 'OR' the flags to combine into one big flag
	window, err := sdl.CreateWindow("Pong", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	// gets a context from the window created
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	// creates a 2d texture matrix
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		panic(err)
	}
	defer tex.Destroy()

	// queries SDL for all the Game controllers connected to the system
	var controllerHandlers []*sdl.GameController
	for i := 0; i < sdl.NumJoysticks(); i++ {
		controllerHandlers = append(controllerHandlers, sdl.GameControllerOpen(i))
		defer controllerHandlers[i].Close()
	}

	// creates an in memory representation of the window space
	pixels := make([]byte, winWidth*winHeight*4)

	// Creates the player paddles, sets position, Width & Height, Speed, and color
	player1 := paddle{pos: pos{50, float32(winHeight / 2)}, w: 20, h: 100, speed: 500, color: gameColor} // Maybe I will play with player acceleration later
	player2 := paddle{pos: pos{float32(winWidth - 50), float32(winHeight / 2)}, speed: 300, acceleration: 1, w: 20, h: 100, color: gameColor}

	// creates a ball in the center of the screen, with a radius, velocity in the x & y direction, and color
	Ball := ball{getCenter(), 15, 300, 300, 1, gameColor}

	// Gets the state of the keys pressed on a keyboard
	keyState := sdl.GetKeyboardState()

	// Variables used to get a fixed rate for machines of all speeds based on time passed
	var frameStart time.Time
	var elapseTime float32

	// Variable to represent how hard someone holds a controller axis down
	var controllerAxis int16

	for {
		// Begining of game loop, grab the current time
		frameStart = time.Now()
		// MacOS requires that the events are taken care of or else program won't launch
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			// switch case for when someone quits out of application
			case *sdl.QuitEvent:
				println("Quit") // not necessary
				// decided with os.Exit since I was having issues when I just
				//broke the game loop and window wasn't closing properly
				os.Exit(0)
			}
		}

		// gets all input from all controllers, not great for more than on game controller
		for _, controller := range controllerHandlers {
			if controller != nil {
				controllerAxis = controller.Axis(sdl.CONTROLLER_AXIS_LEFTY)
			}
		}

		// Updating game states, currently after each score it waits for user spacebar press to start
		if state == play {
			Ball.update(&player1, &player2, elapseTime)
			player1.update(keyState, controllerAxis, elapseTime)
			player2.aiUpdate(&Ball, elapseTime)
		} else if state == start {
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.score == 10 || player2.score == 10 {
					player1.score = 0
					player2.score = 0
				}
				Ball.accelleration = 1
				state = play
			}
		} // need to add a clause for when a second player enters the game and deal with multiple inputs
		// Clear the game board, may be inefficient since it is O(n) speed
		clear(pixels)

		// Drawing pixels to screen
		//drawNumber(getCenter(), white, 15, 2, pixels)
		drawBackground(pixels)
		Ball.draw(pixels)
		player1.draw(pixels)
		player2.draw(pixels)

		// updating window
		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		// Frame Rate fix for all computer types
		elapseTime = float32(time.Since(frameStart).Seconds())
		if elapseTime < 0.004 {
			sdl.Delay(5 - uint32(elapseTime*10000))
			elapseTime = float32(time.Since(frameStart).Seconds())
		}
	}
}
