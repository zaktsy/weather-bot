package main

type WeatherC struct {
	Name    string
	Weather []Weather
	Wind    Wind
	Main    Main
	Sys     Sys
}

type Weather struct {
	Main        string `json:"main"`
	Description string `json:"description"`
}

type Wind struct {
	Speed float32 `json:"speed"`
}

type Main struct {
	Temp      float32 `json:"temp"`
	Pressure  float32 `json:"pressure"`
	Humidity  float32 `json:"humidity"`
	FeelsLike float32 `json:"feels_like"`
}

type Sys struct {
	Country string `json:"country"`
}
