package webhandlers

import "net/http"

func ChatWebHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/web/chat" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "/app/uiassignment/home.html")
}
