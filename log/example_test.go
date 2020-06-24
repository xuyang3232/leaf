package log_test

import (
	"github.com/name5566/leaf/log"
	l "log"
)

func Example() {
	name := "Leaf"

	log.Debugf("My name is %v", name)
	log.Infof("My name is %v", name)
	log.Errorf("My name is %v", name)
	// log.Fatalf("My name is %v", name)

	logger, err := log.NewLoggerLeaf("release", "", l.LstdFlags)
	if err != nil {
		return
	}
	defer logger.Close()

	logger.Debug("will not print")
	logger.Infof("My name is %v", name)

	log.Export(logger)

	log.Debug("will not print")
	log.Infof("My name is %v", name)
}
