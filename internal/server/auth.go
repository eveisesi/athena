package server

import (
	"fmt"
	"net/http"
	"net/url"
)

func (s *server) handleGetAuthCallback(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	code, state, err := s.parseCodeAndStateFromURL(r.URL)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, err)
		return
	}

	err = s.member.Login(ctx, code, state)
	if err != nil {
		s.logger.WithError(err).Errorln()
		s.writeError(ctx, w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte(`
		<!DOCTYPE html>
		<html>
			<title>Athena EVE SSO Auth Callback</title>
			<style>
			body {
				background-color: #585858;
			}
			</style>
			<body>
				<h2>Athena EVE SSO Auth Callback</h2>
			</body>
			<script>
				setTimeout(function() {
					window.close()
				}, 1000)
			</script>
		</html>
	`))

}

// func (s *server) handleGetAuthLogin(w http.ResponseWriter, r *http.Request) {

// 	// var ctx = r.Context()

// 	// state := r.URL.Query().Get("state")
// 	// if state == "" {
// 	// 	return
// 	// }

// 	// uri := s.auth.AuthorizationURI(ctx, state)

// 	// _, _ = w.Write([]byte(fmt.Sprintf(`
// 	// <html>
// 	// 	<body>
// 	// 		<a href="" onclick="popupCenter(600, 800)">Click Here To Login To CCP</a>

// 	// 		<script>
// 	// 			function popupCenter(w, h) {
// 	// 				var left = (screen.width/2)-(w/2);
// 	// 				var top = (screen.height/2)-(h/2);
// 	// 				window.open("%s", "", 'toolbar=no, location=no, directories=no, status=no, menubar=no, scrollbars=no, resizable=no, copyhistory=no, width='+w+', height='+h+', top='+top+', left='+left);
// 	// 			}
// 	// 		</script>
// 	// 	</body>
// 	// </html>
// 	// `, uri)))

// }

func (s *server) parseCodeAndStateFromURL(uri *url.URL) (code, state string, err error) {

	code = uri.Query().Get("code")
	state = uri.Query().Get("state")
	if code == "" || state == "" {
		return "", "", fmt.Errorf("required paramter missing from request")
	}

	return code, state, nil

}
