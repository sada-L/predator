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

// @Summary	Генерация графика
// @Accept		json
// @Produce	png
// @Param		X0	query		number	true	"X0"
// @Param		Y0	query		number	true	"Y0"
// @Param		A	query		number	true	"A"
// @Param		B	query		number	true	"B"
// @Param		E	query		number	true	"E"
// @Param		D	query		number	true	"D"
// @Success	200	{file}		png		"Генерируемый график"
// @Failure	400	{string}	string
// @Failure	500	{string}	string
// @Router		/api/chart [get]
func chartHandler(c *gin.Context) {
	var model Model

	err := c.ShouldBind(&model)
	if err != nil {
		c.JSON(400, gin.H{"massage": "bind error"})
		return
	}

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

	tmpFile, err := os.CreateTemp("", "chart-*.html")
	if err != nil {
		c.JSON(500, gin.H{"massage": "create temp file error"})
		return
	}

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(htmlContent)
	if err != nil {
		c.JSON(500, gin.H{"massage": "write temp file error"})
		return
	}
	tmpFile.Close()

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	err = chromedp.Run(ctx,
		chromedp.Navigate("file://"+tmpFile.Name()),
		chromedp.Sleep(10*time.Second), // Ждем рендеринг
		chromedp.FullScreenshot(&buf, 90),
	)
	if err != nil {
		c.JSON(500, gin.H{"massage": "full screen shot error"})
		return
	}

	outputPath := "chart.png"

	os.WriteFile(outputPath, buf, 0644)
	c.FileAttachment(outputPath, "chart.png")
	c.JSON(200, gin.H{"message": "ok"})

	defer os.Remove(outputPath)
}

// @title	Хищник-Жертва
func main() {
	r := gin.Default()

	r.GET("/api/chart", chartHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
