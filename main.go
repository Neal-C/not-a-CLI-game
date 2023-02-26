//lint:file-ignore ST1006 because I like my code rustic, and I name my parameters however I want
//lint:file-ignore ST1003 forgot one this one does but it must be important

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const (
	NOTHING            = 0
	WALL               = 1
	PLAYER             = 70
	BLANK_SPACE string = " "
	BLOCK_LEVEL string = "H"
	MAX_SAMPLES        = 100
)

//*****************************
//* Input * /
//*****************************

type Input struct {
	pressedKey byte
}

func (self *Input) update() {

	self.pressedKey = 0;


	// tick := time.NewTicker(time.Millisecond * 2);





	// free: for {
	// 	select {
	// 	case <- tick.C:
	// 		break free;
		
	// 	default:

	// 		b := make([]byte, 1);

	// 		os.Stdin.Read(b); //blocks until stdin has stuff in buffer
	// 		self.pressedKey = b[0];
	// 	}
	// }





	// ch := make(chan byte);

	// // ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 4);
	// tick := time.NewTicker(time.Millisecond * 2);

	// // defer cancel();

	// go func(){
	// 	b := make([]byte, 1);

	// 	os.Stdin.Read(b); //blocks until stdin has stuff in buffer
	// 	ch <- b[0];
	// }();

	// select {
	// 	case key := <- ch:{
	// 		self.pressedKey = key;
	// 	}
	// 	case <- tick.C:{
	// 		return
	// 	}
	// }


}

//*****************************
//* Position * /
//*****************************


type Position struct {
	x int
	y int
}

//*****************************
//* Player * /
//*****************************


type Player struct {
	position         Position
	level            *Level
	reverseDirection bool
	input            *Input
}

func (self *Player) update() {

	if self.reverseDirection {
		self.position.x -= 1
		if self.position.x == 2 {
			self.position.x += 1
			self.reverseDirection = false
		}
		return
	}

	self.position.x += 1
	if self.position.x == (self.level.width - 2) {
		self.position.x -= 1
		self.reverseDirection = true
	}
}

//*****************************
//* Stats * /
//*****************************


type Stats struct {
	start  time.Time
	frames int
	fps    float64
}

//checked
func newStats() *Stats {
	return &Stats{
		start: time.Now(),
	}
}

//checked
func (self *Stats) update() {
	self.frames++
	if self.frames == MAX_SAMPLES {
		self.fps = float64(self.frames) / time.Since(self.start).Seconds()
		self.frames = 0
		self.start = time.Now()
	}
}

//*****************************
//* Level * /
//*****************************


type Level struct {
	width  int
	height int
	data   [][]int
}

//checked
func newLevel(width int, height int) *Level {
	data := make([][]int, height)

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			data[h] = make([]int, width)
		}
	}

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			if h == 0 {
				data[h][w] = WALL
			}

			if w == 0 {
				data[h][w] = WALL
			}

			if w == (width - 1) {
				data[h][w] = WALL
			}

			if h == (height - 1) {
				data[h][w] = WALL
			}

		}
	}
	return &Level{
		width:  width,
		height: height,
		data:   data,
	}
}

// func (self *Level) x(){

// }

func (self *Level) set(position Position, value int) {
	self.data[position.y][position.x] = value
}

//*****************************
//* Game * /
//*****************************


type Game struct {
	isRunning  bool
	level      *Level
	stats      *Stats
	player     *Player
	input      *Input
	drawBuffer *bytes.Buffer
}

//checked
func newGame(width int, height int) *Game {

	//copy pasted magic
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run();
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run();

	var (
		level = newLevel(width, height)

		input = &Input{}
	)

	return &Game{
		level:      level,
		drawBuffer: new(bytes.Buffer),
		stats:      newStats(),
		input:      input,
		player: &Player{
			input: input,
			level: level,
			position: Position{
				x: 2,
				y: 5,
			},
		},
	}
}

//checked
func (self *Game) start() {
	self.isRunning = true;
	self.loop();
}

//checked
func (self *Game) loop() {
	for self.isRunning {
		self.input.update();
		self.update();
		self.render();
		self.stats.update();
		time.Sleep(time.Millisecond * 2) //limit FPS // Virtually no glitch at 2
	}
}

//checked
func (self *Game) update() {

	self.level.set(self.player.position, NOTHING);
	self.player.update();
	self.level.set(self.player.position, 70);
}

//checked
func (self *Game) renderLevel() {
	for h := 0; h < self.level.height; h++ {
		for w := 0; w < self.level.width; w++ {

			if self.level.data[h][w] == NOTHING {
				self.drawBuffer.WriteString(" ");
			}

			if self.level.data[h][w] == WALL {
				self.drawBuffer.WriteString("H");
			}

			if self.level.data[h][w] == PLAYER {
				self.drawBuffer.WriteString("😎");
			}
		}
		self.drawBuffer.WriteString("\n")
	}
}

//checked
func (self *Game) render() {

	self.drawBuffer.Reset()
	//copy pasted black magic, that resets or empties/flushes the buffer
	//copy pasted black magic that clears the terminal
	fmt.Fprint(os.Stdout, "\033[2J\033[1;1H")

	self.renderLevel()
	self.renderStats()

	fmt.Fprint(os.Stdout, self.drawBuffer.String())
	//listen inputs
	//blocking
}

//checked
func (self *Game) renderStats() {
	self.drawBuffer.WriteString("-- STATS \n");
	self.drawBuffer.WriteString(fmt.Sprintf("FPS: %.2f \n", self.stats.fps));
	self.drawBuffer.WriteString(fmt.Sprintf("KEYPRESS: %v \n", self.input.pressedKey));
}

func main() {
	width := 80
	height := 18

	game := newGame(width, height)
	game.start()
}
