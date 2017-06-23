package cmq

import (
	"flag"
	"testing"

	//	"github.com/golang/glog"
)

func TestCreateQueue(t *testing.T) {
	flag.Set("logtostderr", "true")
	flag.Set("log_dir", "./")
	flag.Parse()
	c := Init("AKIDZ0FI9oVWZR7btywXiiXI6i9OQGD3vL0a", "b4OWRAg1A6qNiqpcGiz4weATZTAhATaV", "sh", false)

	//	c.CreateQueue("yp-test-queue-1")
	c.SendMessage("yp-test-queue-1", "hello")
	//	c.DeleteQueue("yp-test-queue-1")

}
