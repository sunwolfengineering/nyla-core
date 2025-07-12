package handlers

import (
	"net/http"

	"github.com/chasefleming/elem-go"
	"github.com/chasefleming/elem-go/attrs"
	"github.com/chasefleming/elem-go/htmx"
)

type UIHandlers struct {
	APIBaseURL string
}

func (h *UIHandlers) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	statsURL := h.APIBaseURL + "/v1/stats/realtime"
	html := elem.Html(attrs.Props{attrs.Lang: "en"},
		elem.Head(nil,
			elem.Meta(attrs.Props{attrs.Charset: "UTF-8"}),
			elem.Meta(attrs.Props{
				attrs.Name:    "viewport",
				attrs.Content: "width=device-width, initial-scale=1.0",
			}),
			elem.Title(nil, elem.Text("Nyla Analytics Dashboard")),
			elem.Script(attrs.Props{attrs.Src: "https://unpkg.com/htmx.org@1.9.12"}),
			elem.Script(attrs.Props{attrs.Src: "https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"}),
		),
		elem.Body(attrs.Props{attrs.Class: "bg-gray-50 min-h-screen"},
			// Header
			elem.Header(attrs.Props{attrs.Class: "bg-white shadow px-6 py-4 flex items-center justify-between"},
				elem.Div(attrs.Props{attrs.Class: "text-2xl font-bold text-indigo-700"}, elem.Text("Nyla Analytics")),
				elem.Nav(nil,
					elem.A(attrs.Props{attrs.Href: "#", attrs.Class: "text-gray-600 hover:text-indigo-700 px-3"}, elem.Text("Dashboard")),
					elem.A(attrs.Props{attrs.Href: "#", attrs.Class: "text-gray-600 hover:text-indigo-700 px-3"}, elem.Text("Settings")),
				),
			),
			// Main flex container
			elem.Div(attrs.Props{attrs.Class: "flex"},
				// Sidebar
				elem.Aside(attrs.Props{attrs.Class: "w-64 bg-white border-r min-h-screen p-6 hidden md:block"},
					elem.Nav(attrs.Props{attrs.Class: "space-y-4"},
						elem.A(attrs.Props{attrs.Href: "#", attrs.Class: "block text-indigo-700 font-semibold"}, elem.Text("Overview")),
						elem.A(attrs.Props{attrs.Href: "#", attrs.Class: "block text-gray-600 hover:text-indigo-700"}, elem.Text("Pages")),
						elem.A(attrs.Props{attrs.Href: "#", attrs.Class: "block text-gray-600 hover:text-indigo-700"}, elem.Text("Visitors")),
						elem.A(attrs.Props{attrs.Href: "#", attrs.Class: "block text-gray-600 hover:text-indigo-700"}, elem.Text("Settings")),
					),
				),
				// Main content
				elem.Main(attrs.Props{attrs.Class: "flex-1 p-8"},
					elem.H1(attrs.Props{attrs.Class: "text-3xl font-bold mb-6 text-gray-900"}, elem.Text("Dashboard")),
					elem.Div(attrs.Props{attrs.Class: "grid grid-cols-1 md:grid-cols-3 gap-6 mb-8"},
						// Metric Cards
						elem.Div(attrs.Props{
							attrs.Class:    "bg-white rounded-lg shadow p-6",
							htmx.HXGet:     statsURL,
							htmx.HXTrigger: "load, every 30s",
							htmx.HXSwap:    "innerHTML",
						},
							elem.Text("Loading..."),
						),
						elem.Div(attrs.Props{attrs.Class: "bg-white rounded-lg shadow p-6"},
							elem.Div(attrs.Props{attrs.Class: "text-sm text-gray-500"}, elem.Text("Unique Visitors")),
							elem.Div(attrs.Props{attrs.Class: "text-2xl font-bold text-indigo-700 mt-2"}, elem.Text("--")),
						),
						elem.Div(attrs.Props{attrs.Class: "bg-white rounded-lg shadow p-6"},
							elem.Div(attrs.Props{attrs.Class: "text-sm text-gray-500"}, elem.Text("Active Users")),
							elem.Div(attrs.Props{attrs.Class: "text-2xl font-bold text-indigo-700 mt-2"}, elem.Text("--")),
						),
					),
					// Chart placeholder
					elem.Div(attrs.Props{attrs.Class: "bg-white rounded-lg shadow p-6 h-64 flex items-center justify-center text-gray-400"},
						elem.Text("[Traffic Chart Placeholder]"),
					),
				),
			),
		),
	).Render()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
