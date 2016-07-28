package main

// Main model object
type Sentinal struct {
	closing chan bool
}

func NewSentinal() *Sentinal {
	return &Sentinal{
		closing: make(chan bool, 0),
	}
}

// Main run loop. Does not return.
func (s *Sentinal) start() {

}

// func (s *Sentinal) hunt() {
// 	// Load the image
// 	p := open("lowsett.png")

// 	// Line detection
// 	start := time.Now()
// 	// chunks, lines := hunt(p, colors, COLOR_THRESHOLD, LINE_WIDTH)
// 	fmt.Printf("Hunt completed in: %s\n", time.Since(start))

// 	// Model detection

// 	// Update movement logic

// 	p = output(p.Bounds(), chunks, lines)
// 	save(p, "huntress.png")
// }

// Wait for the run to complete
func (s *Sentinal) wait() {
	<-s.closing
}
