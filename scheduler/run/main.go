package main

import scheduler "minik8s/scheduler/src"

func main() {
	scheduler := scheduler.New()
	scheduler.Start()
}
