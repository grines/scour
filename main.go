package main

import (
	"github.com/grines/scour/pkg/awsdata"
	"github.com/grines/scour/pkg/completion"
)

func main() {
	tid := awsdata.SetTracking()
	completion.Start(tid)
}
