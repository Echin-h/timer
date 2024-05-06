package main

import "Timer/Z"

func main() {
	l := Z.Init()
	l.Debug("log Debug成功")
	l.Warn("log Warn成功")
	l.Info("log Info成功")
}
