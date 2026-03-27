package middleware

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func Serve(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := newResponseWriterWrapper(w)
		next.ServeHTTP(ww, r)

		contentType := w.Header().Get("Content-Type")
		body := ww.buf.Bytes()

		if strings.Contains(contentType, "text/html") {
			if bytes.Contains(body, []byte("<html")) && bytes.Contains(body, []byte("</html")) {
				rep := "</head>"
				if bytes.Contains(body, []byte(rep)) {
					extra := []byte(`<script type="text/javascript">alert("gando!");</script>`)
					body = bytes.Replace(body, []byte(rep), append(extra, []byte(rep)...), 1)
				} else {
					log.Println("no head tag found")
				}
			} else {
				log.Println("no html tag found")
			}
		} else {
			log.Println("response is not html")
		}
		w.Write(body)
	}
	return http.HandlerFunc(fn)
}