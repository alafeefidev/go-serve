package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

//TODO: change
const TEST_JS_NAME   = "BIG_GLOB_GOLD.js"

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
					script := []byte(fmt.Sprintf(`<script type="text/javascript">%s</script>`, getScript()))
					body = bytes.Replace(body, []byte(rep), append(script, []byte(rep)...), 1)
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

func getScript() string {
	f, err := os.ReadFile(TEST_JS_NAME)
	if errors.Is(err, os.ErrNotExist) {
		return `alert("cannot find script file");`
	}
	if err != nil {
		return fmt.Sprintf(`alert("error: %v")`, err)
	}
	return string(f)
}