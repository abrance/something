package ginroutev

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"text/template"
)

var routes []Route

type Route struct {
	Method string
	Path   string
}

//func recordRoutes() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		c.Next()
//		for _, routeInfo := range c.Engine.Routes() {
//			routes = append(routes, Route{
//				Method: routeInfo.Method,
//				Path:   routeInfo.Path,
//			})
//		}
//	}
//}

func Server() {
	r := gin.Default()
	//r.Use(recordRoutes())

	pprof.Register(r, "dev/pprof")
	r.GET("/example", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "example"})
	})

	r.GET("/visualize", func(c *gin.Context) {
		renderVisualization(c.Writer)
	})

	r.Run(":8080")
}

func renderVisualization(w http.ResponseWriter) {
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Gin Routes Visualization</title>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/d3/6.6.0/d3.min.js"></script>
	</head>
	<body>
		<svg width="800" height="600"></svg>
		<script>
			var routes = [
				{{- range . }}
				{ method: "{{.Method}}", path: "{{.Path}}" },
				{{- end }}
			];

			var svg = d3.select("svg"),
				width = +svg.attr("width"),
				height = +svg.attr("height");

			var simulation = d3.forceSimulation(routes)
				.force("x", d3.forceX(width / 2))
				.force("y", d3.forceY(height / 2))
				.force("collide", d3.forceCollide(50))
				.on("tick", ticked);

			function ticked() {
				var u = svg.selectAll("text")
					.data(routes);

				u.enter()
					.append("text")
					.merge(u)
					.attr("x", function(d) { return d.x; })
					.attr("y", function(d) { return d.y; })
					.text(function(d) { return d.method + " " + d.path; });

				u.exit().remove();
			}
		</script>
	</body>
	</html>
	`

	t, err := template.New("routes").Parse(tmpl)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	if err := t.Execute(w, routes); err != nil {
		fmt.Println("Error executing template:", err)
	}
}
