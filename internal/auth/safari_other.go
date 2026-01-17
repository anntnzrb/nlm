//go:build !darwin

package auth

//nolint:unused // retained for future browser detection enhancements
func detectSafari(debug bool) Browser {
	return Browser{Type: BrowserUnknown}
}
