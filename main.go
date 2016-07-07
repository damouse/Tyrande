package main

func runpipe() {
	p := NewPipeline()
	p.run(open("0.png"))
	p.save()
}

func main() {
	// adjustments()
	// sandbox()
	// sand2()

	// stripImage()
	// edgy("hue.png")
	runpipe()
}
