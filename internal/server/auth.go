package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/davecgh/go-spew/spew"
)

func (s *server) handleGetAuthCallback(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	code, state, err := s.parseCodeAndStateFromURL(r.URL)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, err)
		return
	}

	spew.Dump(code, state)

	err = s.member.Login(ctx, code, state)
	if err != nil {
		s.logger.WithError(err).Errorln()
	}

	// err = s.user.VerifyUserRegistrationByToken(ctx, bearer)
	// if err != nil {
	// 	s.logger.WithError(err).Error("failed to verify user")
	// 	s.writeResponse(w, http.StatusBadRequest, zrule.AuthStatus{
	// 		Status: zrule.StatusInvalid,
	// 	})
	// 	return
	// }

	// s.redis.Set(
	// 	ctx,
	// 	fmt.Sprintf(zrule.CACHE_ZRULE_AUTH_TOKEN, state),
	// 	bearer.AccessToken,
	// 	time.Minute*5,
	// )

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte(`
		<!DOCTYPE html>
		<html>
			<title>ZRule EVE SSO Auth Callback</title>
			<script>
				setTimeout(function() {
					window.close()
				}, 250)

			</script>
		</html>
	`))

}

// func (s *server) handleGetAuthLogin(w http.ResponseWriter, r *http.Request) {

// 	var ctx = r.Context()

// 	state := r.URL.Query().Get("state")
// 	if state == "" {
// 		return
// 	}

// 	uri := s.auth.AuthorizationURI(ctx, state)

// 	_, _ = w.Write([]byte(fmt.Sprintf(`
// 	<html>
// 		<body>
// 			<a href="" onclick="popupCenter(600, 800)">Click Here To Login To CCP</a>

// 			<script>
// 				function popupCenter(w, h) {
// 					var left = (screen.width/2)-(w/2);
// 					var top = (screen.height/2)-(h/2);
// 					window.open("%s", "", 'toolbar=no, location=no, directories=no, status=no, menubar=no, scrollbars=no, resizable=no, copyhistory=no, width='+w+', height='+h+', top='+top+', left='+left);
// 				}
// 			</script>
// 		</body>
// 	</html>
// 	`, uri)))

// }

func (s *server) parseCodeAndStateFromURL(uri *url.URL) (code, state string, err error) {

	code = uri.Query().Get("code")
	state = uri.Query().Get("state")
	if code == "" || state == "" {
		return "", "", fmt.Errorf("required paramter missing from request")
	}

	return code, state, nil

}
