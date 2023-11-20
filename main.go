package main

import (
	"LocalWeather/tempWidget"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var weather LocalWeather
var selectedDay = 0
var weatherImages = []*canvas.Image{}

const (
	windowWidth  = 380
	windowHeight = 700
)

func main() {
	// Load all weather images from resource
	weatherImages = append(weatherImages, canvas.NewImageFromResource(resWeather_0))
	weatherImages = append(weatherImages, canvas.NewImageFromResource(resWeather_1_2_3))
	weatherImages = append(weatherImages, canvas.NewImageFromResource(resWeather_45_48))
	weatherImages = append(weatherImages, canvas.NewImageFromResource(resWeather_51_53_55))
	weatherImages = append(weatherImages, canvas.NewImageFromResource(resWeather_56_57))
	weatherImages = append(weatherImages, canvas.NewImageFromResource(resWeather_61_63_65))
	weatherImages = append(weatherImages, canvas.NewImageFromResource(resWeather_66_67))
	weatherImages = append(weatherImages, canvas.NewImageFromResource(resWeather_71_73_75))
	weatherImages = append(weatherImages, canvas.NewImageFromResource(resWeather_77))
	weatherImages = append(weatherImages, canvas.NewImageFromResource(resWeather_80_81_82))
	weatherImages = append(weatherImages, canvas.NewImageFromResource(resWeather_85_86))
	weatherImages = append(weatherImages, canvas.NewImageFromResource(resWeather_95))
	weatherImages = append(weatherImages, canvas.NewImageFromResource(resWeather_96_99))

	// Load Logo image
	logo := canvas.NewImageFromResource(resLogo)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(283, 157))

	// Set all images to contain mode and minimum size
	for _, v := range weatherImages {
		v.FillMode = canvas.ImageFillContain
		v.SetMinSize(fyne.NewSize(380, 380))
	}

	// Background gradient
	gradient := canvas.NewLinearGradient(color.Black, color.RGBA{0, 0, 255, 255}, 0)

	weather.InitLocalWeather()

	a := app.NewWithID("com.lennart.localweather")
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("LocalWeather (c)2023 by Lennart Martens")
	w.Resize(fyne.NewSize(windowWidth, windowHeight))
	w.SetFixedSize(true)
	w.CenterOnScreen()

	// Labels fo 7 days
	var labelDate, labelRain, labelWeather, labelCity, labelCountry [7]*widget.Label

	for day := 0; day < 7; day++ {
		labelDate[day] = widget.NewLabel(weather.localDateEu[day])
		labelDate[day].Alignment = fyne.TextAlignCenter
		labelDate[day].TextStyle = fyne.TextStyle{Bold: true}

		labelRain[day] = widget.NewLabel(fmt.Sprintf("Regenwahrscheinlichkeit : %d%%", weather.Daily.PrecipitationProbabilityMean[day]))
		labelRain[day].Alignment = fyne.TextAlignCenter
		labelWeather[day] = widget.NewLabel(weather.weather[day])
		labelWeather[day].Alignment = fyne.TextAlignCenter
		labelCity[day] = widget.NewLabel(fmt.Sprintf("Stadt: %s %s", weather.location.Postcode, weather.location.City))
		labelCity[day].Alignment = fyne.TextAlignCenter
		labelCountry[day] = widget.NewLabel(fmt.Sprintf("Land: %s", weather.location.CountryName))
		labelCountry[day].Alignment = fyne.TextAlignCenter
	}

	// Label space
	labelSpace := widget.NewLabel("           ")

	// Label minTemp
	labelMinTemp := widget.NewLabel("Min. Temperatur :")
	labelMinTemp.Alignment = fyne.TextAlignCenter

	// Label maxTemp
	labelMaxTemp := widget.NewLabel("Max. Temperatur :")
	labelMaxTemp.Alignment = fyne.TextAlignCenter

	// Buttons
	buttonNext := widget.NewButtonWithIcon("          ", theme.NewThemedResource(theme.MediaFastForwardIcon()), func() {})
	buttonPrevious := widget.NewButtonWithIcon("          ", theme.NewThemedResource(theme.MediaFastRewindIcon()), func() {})
	buttonHelp := widget.NewButtonWithIcon("    ", theme.NewThemedResource(theme.HelpIcon()), func() {})

	// Progress bars for 7 days
	var progressRain [7]*widget.ProgressBar
	for day := 0; day < 7; day++ {
		progressRain[day] = widget.NewProgressBar()
		progressRain[day].SetValue(float64(weather.Daily.PrecipitationProbabilityMean[day]) / 100)
	}

	// Min, max temperatur Widget for 7 days
	var minTemp [7]*tempWidget.TemperatureWidget
	var maxTemp [7]*tempWidget.TemperatureWidget
	for day := 0; day < 7; day++ {
		minTemp[day] = tempWidget.NewTemperatureWidget()
		minTemp[day].Temperature = weather.Daily.Temperature2MMin[day]
		maxTemp[day] = tempWidget.NewTemperatureWidget()
		maxTemp[day].Temperature = weather.Daily.Temperature2MMax[day]
	}

	// Containers for 7 days
	var vBox []*fyne.Container
	var hBoxDateButtons []*fyne.Container
	for day := 0; day < 7; day++ {
		hBoxDateButtons = append(hBoxDateButtons, container.NewHBox(labelSpace, buttonPrevious, labelDate[day], buttonNext))
		vBox = append(vBox, container.NewVBox(hBoxDateButtons[day], labelMinTemp, minTemp[day], labelMaxTemp, maxTemp[day], labelRain[day], progressRain[day], labelWeather[day], getWeatherImg(day), labelCity[day], labelCountry[day]))
	}

	hBoxDateButtons[0].Add(buttonHelp)

	// Buttons next and previous
	buttonNext.OnTapped = func() {
		if selectedDay < 6 {
			selectedDay++
			w.SetContent(container.NewStack(gradient, vBox[selectedDay]))
		}
	}

	buttonPrevious.OnTapped = func() {
		if selectedDay > 0 {
			selectedDay--
			w.SetContent(container.NewStack(gradient, vBox[selectedDay]))
		}
	}
	buttonHelp.OnTapped = func() {
		dlgHelp := dialog.NewCustom("(c)2023 by Lennart Martens\nLicense: MIT\nWetter API: Open-Meteo\nLocation API: BigDataCloud", "Ok", logo, w)
		dlgHelp.Show()
	}

	w.SetContent(container.NewStack(gradient, vBox[selectedDay]))
	w.ShowAndRun()

}

func getWeatherImg(day int) *canvas.Image {
	switch weather.Daily.WeatherCode[day] {
	case 0:
		return weatherImages[0]
	case 1, 2, 3:
		return weatherImages[1]
	case 45, 48:
		return weatherImages[2]
	case 51, 53, 55:
		return weatherImages[3]
	case 56, 57:
		return weatherImages[4]
	case 61, 63, 65:
		return weatherImages[5]
	case 66, 67:
		return weatherImages[6]
	case 71, 73, 75:
		return weatherImages[7]
	case 77:
		return weatherImages[8]
	case 80, 81, 82:
		return weatherImages[9]
	case 85, 86:
		return weatherImages[10]
	case 95:
		return weatherImages[11]
	case 96, 99:
		return weatherImages[12]
	default:
		return weatherImages[0]
	}
}
