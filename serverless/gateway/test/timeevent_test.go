package test

import (
	"fmt"
	"testing"

	"github.com/robfig/cron"
)

func TestTimeEvent(t *testing.T){
	fmt.Println("Test Start")
	c := cron.New()
	c.AddFunc("*/1 * * * *", func() {
		fmt.Println("Every Second")
	})
	c.Start()
	select {}
}
