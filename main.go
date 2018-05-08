package main

import (
	"log"
	"os"
	"text/template"

	awsapi "./awsapi"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	owm "github.com/briandowns/openweathermap" // "owm" for easier use
	"github.com/gin-gonic/gin"
)

const URL = "http://ip-api.com/json"

const weatherTemplate = `Current weather for {{.Name}}:
    Conditions: {{range .Weather}} {{.Description}} {{end}}
    Now:         {{.Main.Temp}} {{.Unit}}
    High:        {{.Main.TempMax}} {{.Unit}}
    Low:         {{.Main.TempMin}} {{.Unit}}
	`

var ginLambda *ginadapter.GinLambda

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	gin.SetMode(gin.ReleaseMode)
	if ginLambda == nil {
		log.Printf("Gin cold start")
		r := gin.Default()
		r.GET("/wether", getWether)

		ginLambda = ginadapter.New(r)
	}

	return ginLambda.Proxy(req)
}

func main() {
	lambda.Start(Handler)
}

func getWether(c *gin.Context) {
	token, err := awsapi.GetSsm("SLACK_TOKEN_FOR_WETHER", "ap-northeast-1")
	if err != nil {
		c.String(400, "Ssm error")
		return
	}
	key, err := awsapi.GetSsm("OWM_API_KEY", "ap-northeast-1")
	if err != nil {
		c.String(400, "Ssm error")
		return
	}
	if token != c.Query("token") {
		c.String(400, "Incorrect token")
		return
	}
	wether, err := getWetherMap(key)
	if err != nil {
		c.String(400, "Get wether error")
		return
	}
	c.String(200, wether)
}

func getCurrent(l, u, lang string, key string) *owm.CurrentWeatherData {
	w, err := owm.NewCurrent(u, lang, key)
	if err != nil {
		log.Fatalln(err)
	}
	w.CurrentByName(l)
	return w
}

func getWetherMap(key string) (string, error) {
	w := getCurrent("Tokyo", "c", "en", key)
	tmpl, err := template.New("weather").Parse(weatherTemplate)
	if err != nil {
		log.Println(err)
		return "", err
	}

	err = tmpl.Execute(os.Stdout, w)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return tmpl.Name(), nil
}
