package main

import (
	"github.com/alexflint/go-arg"
	"github.com/fossteams/fossteams-backend/internal/server"
	"github.com/sirupsen/logrus"
)

var args struct {
	Debug *bool `arg:"-d"`
}

func main() {
	arg.MustParse(&args)

	logger := logrus.New()
	s, err := server.New(logger)

	if err != nil {
		logger.Fatalf("unable to create server: %v", err)
	}

	logger.Fatal(s.Start("0.0.0.0:8050"))
}
