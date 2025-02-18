package main

import (
	_ "awesomeProject/docs"
	"context"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"os"
	"time"
)

// @Summary Генерация графика
// @Accept  json
// @Produce  png
// @Param X0 query number true "X0"
// @Param Y0 query number true "Y0"
// @Param A	 query number true "A"
// @Param B	 query number true "B"
// @Param E	 query number true "E"
// @Param D	 query number true "D"
// @Success 200 {file} png "Генерируемый график"
// @Router /api/chart [get]
func chartHandler(c *gin.Context) {
	var model Model

	_ = c.ShouldBind(&model)

	x, y, a, b, d, e := model.X0, model.Y0, model.A, model.B, model.D, model.E
	data := make([]opts.LineData, 0)

	for i := 0; i < 150; i++ {
		dx := (e-a*y)*x + x
		dy := (d*dx-b)*y + y

		x = dx
		y = dy

		data = append(data, opts.LineData{Value: []float64{x, y}})
	}

	x, y = model.X0, model.Y0-5
	for i := 0; i < 150; i++ {
		dx := (e-a*y)*x + x
		dy := (d*dx-b)*y + y

		x = dx
		y = dy

		data = append(data, opts.LineData{Value: []float64{x, y}})
	}

	x, y = model.X0, model.Y0-10
	for i := 0; i < 150; i++ {
		dx := (e-a*y)*x + x
		dy := (d*dx-b)*y + y

		x = dx
		y = dy

		data = append(data, opts.LineData{Value: []float64{x, y}})
	}

	line := charts.NewLine()
	line.AddSeries("count", data)
	htmlContent := line.RenderContent()

	tmpFile, _ := os.CreateTemp("", "chart-*.html")
	defer os.Remove(tmpFile.Name())

	_, _ = tmpFile.Write(htmlContent)
	tmpFile.Close()

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	_ = chromedp.Run(ctx,
		chromedp.Navigate("file://"+tmpFile.Name()),
		chromedp.Sleep(2*time.Second), // Ждем рендеринг
		chromedp.FullScreenshot(&buf, 90),
	)

	outputPath := "chart.png"

	os.WriteFile(outputPath, buf, 0644)
	c.FileAttachment(outputPath, "chart.png")

	defer os.Remove(outputPath)
}

// @title			Хищник-Жертва
func main() {
	r := gin.Default()

	r.GET("/api/chart", chartHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
