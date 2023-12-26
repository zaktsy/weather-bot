package main

import (
	"encoding/json"
	"flag"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var cache = make(map[string]string)
var cacheLife = make(map[string]time.Time)

func main() {
	token := flag.String("token", "", "")
	weatherApi := flag.String("weatherApi", "", "")
	geocoderApi := flag.String("geocoderApi", "", "")

	flag.Parse()
	if *token == "" || *weatherApi == "" {
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(*token)
	if err != nil {
		fmt.Println(err)
		return
	}
	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() && update.Message.Command() == "start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Напишите ваш город, или место в вашем городе, я скажу, какая там погода")
			bot.Send(msg)
			continue
		}

		rawCity := strings.TrimSpace(update.Message.Text)
		cityMessage, err := GetNormalizedCityMessage(rawCity, *geocoderApi)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Извините, не смог понять, погода какого города вам нужна"))
			continue
		}

		weatherMessage := getWeatherMessage(cityMessage, *weatherApi)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, weatherMessage)
		bot.Send(msg)
	}
}

func GetNormalizedCityMessage(rawCity string, geocoderApi string) (string, error) {
	urlapi := fmt.Sprintf("https://geocode-maps.yandex.ru/1.x?apikey=%s&geocode=%s&kind=locality&format=json", geocoderApi, url.QueryEscape(rawCity))

	resp, err := http.Get(urlapi)
	if err != nil {
		return fmt.Sprintf("Извините, я не знаю города %s", rawCity), err
	}

	if resp.StatusCode == 404 {
		return fmt.Sprintf("Извините, я не знаю города %s", rawCity), err
	}

	defer resp.Body.Close()
	var geocoder Geocoder
	data, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &geocoder); err != nil {
		log.Println(err)
		return fmt.Sprintf("Извините, не смог узнать погоду в городе %s", rawCity), err
	}

	return geocoder.Response.GeoObjectCollection.FeatureMember[0].GeoObject.Name, nil
}

func getWeatherMessage(city string, weatherApi string) string {
	if weather, ok := cache[city]; ok && time.Now().Before(cacheLife[city]) {
		return weather
	}

	urlapi := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?appid=%s&q=%s&lang=ru&units=metric", weatherApi, url.QueryEscape(city))

	resp, err := http.Get(urlapi)
	if err != nil {
		return fmt.Sprintf("Извините, не смог узнать погоду в городе %s", city)
	}

	if resp.StatusCode == 404 {
		return fmt.Sprintf("Извините, не смог узнать погоду в городе %s", city)
	}
	defer resp.Body.Close()
	var weather WeatherC
	data, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &weather); err != nil {
		log.Println(err)
		return fmt.Sprintf("Извините, не смог узнать погоду в городе %s", city)
	}

	output := fmt.Sprintf("%s, %s\n\nПогода\nОписание: %s\nТемпература: %.2f °C\nОщущается как: %.2f °C\nСкорость ветра: %.2f км/ч\nВлажность: %.2f",
		weather.Name, weather.Sys.Country, weather.Weather[0].Description, weather.Main.Temp, weather.Main.FeelsLike, weather.Wind.Speed, weather.Main.Humidity)
	return output
}
