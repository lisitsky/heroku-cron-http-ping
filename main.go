package main

import (
	"net/http"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)


func main()  {
	var config = getConfig()

	var timeout = time.Duration(config.Http.Timeout) * time.Second

	var wg sync.WaitGroup

	for i := 0; i < config.Ping.Count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ping(config.Ping.URL, timeout)
		}()
		time.Sleep(time.Duration(config.Ping.Delay) * time.Millisecond)
	}
	wg.Wait()
}

func ping(url string, timeout time.Duration)  {
	var client = http.Client{
		Timeout: timeout,
	}
	log.Debugf("starting client with url: %s and timeout: %v", url, timeout)
	resp, err := client.Get(url)
	if err != nil {
		log.Errorf("ping error: %s", err)
		return
	}
	log.Debugf("ping response: %s", resp.Status)
}

func getConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	var config Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("unable to decode config: %v", err)
	}

	logLevel, err := log.ParseLevel(config.Log.Level)
	if err != nil {
		log.Errorf("cannot parse log level: %s, log level set to error", config.Log.Level)
		logLevel = log.ErrorLevel
	}
	log.SetLevel(logLevel)

	log.Debugf("parsed config: %+v", config)

	return &config
}