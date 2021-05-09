package scour

import (
	"github.com/grines/scour/pkg/awsdata"
	"github.com/grines/scour/pkg/completion"
)

func Start() {
	tid := awsdata.SetTracking()
	completion.Start(tid)
}
