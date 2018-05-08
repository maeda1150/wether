package main

import (
	owm "github.com/briandowns/openweathermap" // "owm" for easier use
	"log"
	"os"
	"text/template"
)

const URL = "http://ip-api.com/json"

const weatherTemplate = `Current weather for {{.Name}}:
    Conditions: {{range .Weather}} {{.Description}} {{end}}
    Now:         {{.Main.Temp}} {{.Unit}}
    High:        {{.Main.TempMax}} {{.Unit}}
    Low:         {{.Main.TempMin}} {{.Unit}}
`

func getCurrent(l, u, lang string) *owm.CurrentWeatherData {
	w, err := owm.NewCurrent(u, lang, os.Getenv("OWM_API_KEY"))
	if err != nil {
		log.Fatalln(err)
	}
	w.CurrentByName(l)
	return w
}

func main() {
	w := getCurrent("Tokyo", "c", "en")
	tmpl, err := template.New("weather").Parse(weatherTemplate)
	if err != nil {
		log.Fatalln(err)
	}

	err = tmpl.Execute(os.Stdout, w)
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(0)
}
