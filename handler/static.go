package handler

import "net/http"

// StaticJS is http handler that handle static JS assets
func (h Handler) StaticJS() http.Handler {
	fs := http.FileServer(http.Dir("files/var/www/assets/creative/js"))
	return fs
}

// StaticCSS is http handler that handle static JS assets
func (h Handler) StaticCSS() http.Handler {
	fs := http.FileServer(http.Dir("files/var/www/assets/creative/css"))
	return fs
}

// StaticImages is http handler that handle static JS assets
func (h Handler) StaticImages() http.Handler {
	fs := http.FileServer(http.Dir("files/var/www/assets/creative/images"))
	return fs
}

// StaticScript is http handler that handle static JS assets
func (h Handler) StaticScript() http.Handler {
	fs := http.FileServer(http.Dir("files/var/www/assets/creative/scripts"))
	return fs
}

// StaticStyles is http handler that handle static JS assets
func (h Handler) StaticStyles() http.Handler {
	fs := http.FileServer(http.Dir("files/var/www/assets/creative/styles"))
	return fs
}

// StaticData is http handler that handle static exposes data to public
func (h Handler) StaticData() http.Handler {
	fs := http.FileServer(http.Dir("files/data"))
	return fs
}
