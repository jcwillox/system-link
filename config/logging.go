package config

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jcwillox/system-bridge/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"time"
)

type MQTTLogger struct {
	Lvl zerolog.Level
}

func (l MQTTLogger) Println(v ...interface{}) {
	log.WithLevel(l.Lvl).Msg("[MQTT] " + fmt.Sprint(v...))
}
func (l MQTTLogger) Printf(format string, v ...interface{}) {
	log.WithLevel(l.Lvl).Msgf("[MQTT] "+format, v...)
}

func SetupLogging() {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}

	fileWriter := zerolog.ConsoleWriter{
		Out: &lumberjack.Logger{
			Filename:   strings.TrimSuffix(utils.ExePath, ".exe") + ".log",
			MaxSize:    1,
			MaxBackups: 3,
		},
		NoColor:    true,
		TimeFormat: time.DateTime,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(i.(string))
		},
	}

	log.Logger = log.Output(zerolog.MultiLevelWriter(consoleWriter, fileWriter))
}

func SetLogLevels() {
	// set log level
	level, err := zerolog.ParseLevel(Config.LogLevel)
	log.Info().Str("lvl", Config.LogLevel).Msg("log level requested")
	if err != nil {
		log.Fatal().Err(err).Msg("fatal error setting log level")
	} else {
		log.Info().Int8("lvl", int8(level)).Str("name", zerolog.GlobalLevel().String()).Msg("log level set")
	}
	zerolog.SetGlobalLevel(level)

	// set mqtt log level
	mqtt.CRITICAL = MQTTLogger{zerolog.FatalLevel}
	if Config.LogLevelMqtt == "critical" {
		return
	}
	mqtt.ERROR = MQTTLogger{zerolog.ErrorLevel}
	if Config.LogLevelMqtt == "error" {
		return
	}
	mqtt.WARN = MQTTLogger{zerolog.WarnLevel}
	if Config.LogLevelMqtt == "warn" {
		return
	}
	mqtt.DEBUG = MQTTLogger{zerolog.DebugLevel}
	if Config.LogLevelMqtt == "debug" {
		return
	}
}
