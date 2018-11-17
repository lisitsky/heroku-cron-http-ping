package main

type Config struct {
	Ping struct {
		URL    string
		Cookie string
		Count  int
		Delay  int
	}

	Http struct {
		Timeout int
	}

	Log struct {
		Level string
	}
}
