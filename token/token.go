package token

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/sandertv/gophertunnel/minecraft/auth"
	"golang.org/x/oauth2"
)

// tokenSource returns a token source for using with a gophertunnel client. It either reads it from the
// token.tok file if cached or requests logging in with a device code.
func TokenSource() oauth2.TokenSource {
	check := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	token := new(oauth2.Token)
	tokenData, err := ioutil.ReadFile("token.tok")
	if err == nil {
		_ = json.Unmarshal(tokenData, token)
	} else {
		token, err = auth.RequestLiveToken()
		check(err)
	}
	src := auth.RefreshTokenSource(token)
	_, err = src.Token()
	if err != nil {
		// The cached refresh token expired and can no longer be used to obtain a new token. We require the
		// user to log in again and use that token instead.
		token, err = auth.RequestLiveToken()
		check(err)
		src = auth.RefreshTokenSource(token)
	}
	go func() {
		log.Info().Msgf("Writing token file")
		c := make(chan os.Signal, 3)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
		<-c

		tok, _ := src.Token()
		b, _ := json.Marshal(tok)
		_ = ioutil.WriteFile("token.tok", b, 0644)
		os.Exit(0)
	}()
	return src
}
