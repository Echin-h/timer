package Z

import "testing"

func TestInit(t *testing.T) {
	l := Init()
	l.Debug("log Debug成功")
	l.Warn("log Warn成功")
	l.Info("log Info成功")
}
