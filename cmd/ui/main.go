package main

import (
	"net/http"

	"github.com/chasefleming/elem-go"
	"github.com/chasefleming/elem-go/attrs"
)

func main() {
	http.HandleFunc("/", dashboardHandler)
	http.ListenAndServe(":8080", nil)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	html := elem.Html(attrs.Props{attrs.Lang: "en"},
		elem.Head(nil,
			elem.Meta(attrs.Props{attrs.Charset: "UTF-8"}),
			elem.Meta(attrs.Props{
				attrs.Name:    "viewport",
				attrs.Content: "width=device-width, initial-scale=1.0",
			}),
			elem.Title(nil, elem.Text("Nyla Analytics Dashboard")),
			elem.Script(attrs.Props{attrs.Src: "https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"}),
		),
		elem.Body(attrs.Props{"class": "bg-gray-50 min-h-screen"},
			// Header
			elem.Header(attrs.Props{"class": "bg-white shadow px-6 py-4 flex items-center justify-between"},
				elem.Div(attrs.Props{"class": "text-2xl font-bold text-indigo-700"}, elem.Text("Nyla Analytics")),
				elem.Nav(nil,
					elem.A(attrs.Props{"href": "#", "class": "text-gray-600 hover:text-indigo-700 px-3"}, elem.Text("Dashboard")),
					elem.A(attrs.Props{"href": "#", "class": "text-gray-600 hover:text-indigo-700 px-3"}, elem.Text("Settings")),
				),
			),
			// Main flex container
			elem.Div(attrs.Props{"class": "flex"},
				// Sidebar
				elem.Aside(attrs.Props{"class": "w-64 bg-white border-r min-h-screen p-6 hidden md:block"},
					elem.Nav(attrs.Props{"class": "space-y-4"},
						elem.A(attrs.Props{"href": "#", "class": "block text-indigo-700 font-semibold"}, elem.Text("Overview")),
						elem.A(attrs.Props{"href": "#", "class": "block text-gray-600 hover:text-indigo-700"}, elem.Text("Pages")),
						elem.A(attrs.Props{"href": "#", "class": "block text-gray-600 hover:text-indigo-700"}, elem.Text("Visitors")),
						elem.A(attrs.Props{"href": "#", "class": "block text-gray-600 hover:text-indigo-700"}, elem.Text("Settings")),
					),
				),
				// Main content
				elem.Main(attrs.Props{"class": "flex-1 p-8"},
					elem.H1(attrs.Props{"class": "text-3xl font-bold mb-6 text-gray-900"}, elem.Text("Dashboard")),
					elem.Div(attrs.Props{"class": "grid grid-cols-1 md:grid-cols-3 gap-6 mb-8"},
						// Metric Cards
						elem.Div(attrs.Props{"class": "bg-white rounded-lg shadow p-6"},
							elem.Div(attrs.Props{"class": "text-sm text-gray-500"}, elem.Text("Total Pageviews")),
							elem.Div(attrs.Props{"class": "text-2xl font-bold text-indigo-700 mt-2"}, elem.Text("--")),
						),
						elem.Div(attrs.Props{"class": "bg-white rounded-lg shadow p-6"},
							elem.Div(attrs.Props{"class": "text-sm text-gray-500"}, elem.Text("Unique Visitors")),
							elem.Div(attrs.Props{"class": "text-2xl font-bold text-indigo-700 mt-2"}, elem.Text("--")),
						),
						elem.Div(attrs.Props{"class": "bg-white rounded-lg shadow p-6"},
							elem.Div(attrs.Props{"class": "text-sm text-gray-500"}, elem.Text("Active Users")),
							elem.Div(attrs.Props{"class": "text-2xl font-bold text-indigo-700 mt-2"}, elem.Text("--")),
						),
					),
					// Chart placeholder
					elem.Div(attrs.Props{"class": "bg-white rounded-lg shadow p-6 h-64 flex items-center justify-center text-gray-400"},
						elem.Text("[Traffic Chart Placeholder]"),
					),
				),
			),
		),
	).Render()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
