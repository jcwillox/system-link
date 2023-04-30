package engine

import (
	"crypto/tls"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jcwillox/system-bridge/config"
)

func NewClientOptions() *mqtt.ClientOptions {
	url := fmt.Sprintf("%s:%s", config.Config.MQTT.Host, config.Config.MQTT.Port)
	if config.Config.MQTT.TLS {
		url = "ssl://" + url
	}

	opts := mqtt.NewClientOptions().
		AddBroker(url).
		SetClientID("system-bridge").
		SetUsername(config.Config.MQTT.Username).
		SetPassword(config.Config.MQTT.Password).
		SetConnectRetry(true).
		SetTLSConfig(&tls.Config{
			InsecureSkipVerify: false, // TODO: allow no verify / custom ssl
		})

	return opts

}
