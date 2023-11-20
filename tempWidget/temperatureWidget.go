package tempWidget

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type TemperatureWidget struct {
	widget.BaseWidget // Erbt Basisfunktionen eines Widgets
	Temperature       float64
}

func NewTemperatureWidget() *TemperatureWidget {
	t := &TemperatureWidget{}
	t.ExtendBaseWidget(t) // Wichtig, um sicherzustellen, dass das Widget korrekt initialisiert wird
	return t
}

type temperatureRenderer struct {
	widget *TemperatureWidget
	rect   *canvas.Rectangle
	text   *canvas.Text
	frame  *canvas.Rectangle
}

// MinSize implements fyne.WidgetRenderer.
func (t *temperatureRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, 25)

}

// Refresh implements fyne.WidgetRenderer.
func (r *temperatureRenderer) Refresh() {
	r.text.Text = fmt.Sprintf("%.1f°", r.widget.Temperature)
	r.text.Color = color.RGBA{255, 255, 0, 255} // Textfarbe

	temp := r.widget.Temperature
	if temp < -60 {
		temp = -60
	} else if temp > 60 {
		temp = 60
	}

	// Aktualisieren Sie die Farbe des Balkens basierend auf der Temperatur
	if temp < 0 {
		r.rect.FillColor = color.RGBA{0, 0, 255, 255} // Blau für kalte Temperaturen
	} else {
		r.rect.FillColor = color.RGBA{255, 0, 0, 255} // Rot für warme Temperaturen
	}

	// Aktualisieren Sie das Erscheinungsbild des Rahmens
	r.frame.FillColor = color.Transparent                // Transparenter Hintergrund für den Rahmen
	r.frame.StrokeColor = color.RGBA{100, 100, 100, 255} // Weiße Farbe für die Linie des Rahmens
	r.frame.StrokeWidth = 1                              // Breite der Linie des Rahmens

	canvas.Refresh(r.rect)
	canvas.Refresh(r.text)
	canvas.Refresh(r.frame)
}

func (t *TemperatureWidget) CreateRenderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(color.RGBA{200, 200, 200, 255})
	text := canvas.NewText("", color.RGBA{255, 255, 0, 255})
	frame := canvas.NewRectangle(color.RGBA{100, 100, 100, 255})

	return &temperatureRenderer{widget: t, rect: rect, text: text, frame: frame}
}

func (r *temperatureRenderer) Layout(size fyne.Size) {
	padding := float32(2)                                                      // Randbreite des Rahmens
	innerSize := fyne.NewSize(size.Width-(padding*2), size.Height-(padding*2)) // Innere Größe nach Abzug des Randes

	r.frame.Resize(size)
	r.frame.Move(fyne.NewPos(0, 0)) // Rahmen soll bei (0, 0) beginnen

	// Berechnen Sie die Breite des Temperaturbalkens basierend auf der Temperatur
	progress := (r.widget.Temperature + 60) / 120
	rectWidth := float32(progress) * innerSize.Width
	r.rect.Resize(fyne.NewSize(rectWidth, innerSize.Height)) // Setzen Sie die Größe des Temperaturbalkens
	r.rect.Move(fyne.NewPos(padding, padding))               // Bewegen Sie den Temperaturbalken, um den Rand zu berücksichtigen

	r.text.TextSize = 14 // Anpassen der Textgröße
	r.text.Alignment = fyne.TextAlignCenter
	r.text.Resize(innerSize)
	r.text.Move(fyne.NewPos(padding, padding))

	r.Refresh()
}
func (r *temperatureRenderer) Update() {
	// Diese Methode sollte nur die Eigenschaften des Widgets aktualisieren
	// und keine Refresh-Aufrufe enthalten, um Rekursion zu vermeiden.
}

func (r *temperatureRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.frame, r.rect, r.text}
}

func (r *temperatureRenderer) Destroy() {}
