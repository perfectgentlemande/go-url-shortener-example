package api

import "net/http"

func (c *Controller) resolve(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("resolve"))
}
func (c *Controller) shorten(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("shorten"))
}
