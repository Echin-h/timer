package Z

import (
	"go.uber.org/zap/zapcore"
	"time"
)

var Z Zap

type Zap struct {
	Prefix       string `default:"gin-vue-admin" yaml:"prefix" json:"prefix"` // 日志时间前缀
	Level        string `default:"panic" yaml:"level" json:"level"`           // 级别
	Path         string `default:"" yaml:"path" json:"path"`                  // path记得改一下，我之前乱写的
	Director     string `default:"director" yaml:"director" json:"director"`
	EncoderLevel string `default:"LowercaseLevelEncoder" yaml:"encoder_level" json:"encoder_level"` // 编码级
	MaxAge       int    `default:"7" yaml:"max_age" json:"max_age"`                                 // 日志留存时间
	ShowLine     bool   `default:"true" yaml:"show_line" json:"show_line"`                          // 是否显示行
	LogInConsole bool   `default:"true" yaml:"log_in_console" json:"log_in_console"`                // 是否输出控制台
	Encoder      string `default:"console" yaml:"encoder" json:"encoder"`                           // 输出形式 Json/console
}

func (z *Zap) GetLevelEncoder() zapcore.LevelEncoder {
	switch {
	case z.EncoderLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case z.EncoderLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case z.EncoderLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case z.EncoderLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

func (z *Zap) CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(Z.Prefix + t.Format("2006/01/02 - 15:04:05.000"))
}

func (z *Zap) TransmitLvl() zapcore.Level {
	switch z.Level {
	case "debug":
		return zapcore.DebugLevel
	case "Info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "Dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.DebugLevel
	}
}
