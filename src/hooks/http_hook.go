package hooks

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/tftp-go-team/hooktftp/src/config"
	"github.com/tftp-go-team/libgotftp/src"
)

var HTTPHook = HookComponents{
	func(url string, _ tftp.Request) (*HookResult, error) {
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) == 0 {
					return nil
				}

				for key, val := range via[0].Header {
					if key != "Referer" {
						req.Header[key] = val
					}
				}

				return nil
			},
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", "hooktftp v0.9.1")

		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			res.Body.Close()
			return nil, fmt.Errorf("Bad response '%v' from %v", res.Status, url)
		}

		return newHookResult(res.Body, nil, int(res.ContentLength), nil), nil

	},
	func(s string, extra config.HookExtraArgs) (string, error) {
		shouldUrlDecode := extra["urldecode"].(bool)

		if !shouldUrlDecode {
			return s, nil
		}

		unescaped, err := url.PathUnescape(s)
		if err != nil {
			return "", err
		}
		return unescaped, nil
	},
}
