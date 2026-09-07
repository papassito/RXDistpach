package webassets

import (
	"embed"
	"net/http"
)

//go:embed static/*
var staticFS embed.FS

// RegisterAssetRoutes monta los vectores estáticos en el enrutador HTTP del Gateway/Edge
func RegisterAssetRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/assets/logo.svg", func(w http.ResponseWriter, r *http.Request) {
		data, err := staticFS.ReadFile("static/logo.svg")
		if err != nil {
			// Si no ha sido copiado aún en tiempo de compilación, devolvemos un estado adecuado
			http.Error(w, "Asset no encontrado", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Write(data)
	})

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		// Servimos el SVG con soporte moderno. Para máxima compatibilidad,
		// se puede configurar como image/svg+xml.
		data, err := staticFS.ReadFile("static/icon.svg")
		if err != nil {
			http.Error(w, "Favicon no encontrado", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Write(data)
	})
}
