package web

import (
	projectErrors "He110/donation-report-manager/internal/business/errors"
	getConfigUseCase "He110/donation-report-manager/internal/business/use_cases/get_config_use_case"
	saveTokenUseCase "He110/donation-report-manager/internal/business/use_cases/save_token_use_case"
	donationClient "He110/donation-report-manager/internal/pkg/donation-alerts-client"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type App struct {
	container      *TemplateContainer
	ucGetConfig    *getConfigUseCase.UseCase
	ucSaveToken    *saveTokenUseCase.UseCase
	daClient       *donationClient.Client
	sessionStorage *sessions.CookieStore
	connections    map[string]*websocket.Conn
}

func NewApp(c *TemplateContainer, daClient *donationClient.Client, ucGetConfig *getConfigUseCase.UseCase, ucSaveToken *saveTokenUseCase.UseCase) *App {
	return &App{
		container:      c,
		ucGetConfig:    ucGetConfig,
		ucSaveToken:    ucSaveToken,
		daClient:       daClient,
		sessionStorage: sessions.NewCookieStore(securecookie.GenerateRandomKey(32)),
		connections:    make(map[string]*websocket.Conn, 0),
	}
}

func (a *App) HandlerPanel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = a.container.MustGet("panel").Execute(w, nil)
	}
}

func (a *App) HandlerGetConfig(socketHost string) http.HandlerFunc {
	cfgTpl := a.container.MustGet("config")

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		channelId, exists := vars["channelId"]
		if !exists {
			http.Error(w, "Channel id is required", http.StatusBadRequest)
			return
		}

		response := GetConfigData{
			IsAuthorized: false,
			ChannelId:    channelId,
			SocketHost:   socketHost,
		}

		config, err := a.ucGetConfig.Perform(channelId)
		if err != nil && errors.Is(err, projectErrors.NotAuthorizedError{}) {
			_ = cfgTpl.Execute(w, response)
			return
		} else if err != nil {
			http.Error(w, "Cannot fetch extension config: "+err.Error(), http.StatusInternalServerError)
			return
		}
		response.IsAuthorized = true
		response.Config = config
		if err = cfgTpl.Execute(w, response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (a *App) HandlerChanneledRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		channelId := vars["channelId"]
		session, _ := a.sessionStorage.Get(r, SessionKey)
		session.Values[SessionKeyChannelId] = channelId
		if err := session.Save(r, w); err != nil {
			http.Error(w, "Cannot store channel id", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/redirect", http.StatusSeeOther)
	}
}

func (a *App) HandlerRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := a.sessionStorage.Get(r, SessionKey)
		channelIdRaw, exists := session.Values[SessionKeyChannelId]
		if !exists {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		channelId, ok := channelIdRaw.(string)
		if !ok {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		if r.URL.Query().Has("code") {
			// We got the code, so it's second step
			code := r.URL.Query().Get("code")
			token, err := a.daClient.ObtainAccessToken(code)
			if err != nil {
				http.Error(w, "Cannot obtain access token: "+err.Error(), http.StatusInternalServerError)
				return
			}

			_, err = a.ucSaveToken.Perform(saveTokenUseCase.Parameters{
				ChannelId:    channelId,
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
				ExpiresAfter: time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
			})
			if err != nil {
				http.Error(w, "Cannot save access token", http.StatusInternalServerError)
			}

			a.refreshConfigPage(channelId)
			_ = a.container.MustGet("close").Execute(w, nil)
			return
		}

		http.Redirect(w, r, a.daClient.GetAuthLink(), http.StatusSeeOther)
	}
}

func (a *App) SocketBridge() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		channelId := vars["channelId"]
		log.Println(fmt.Sprintf("Config page of %s has been connected", channelId))
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		a.connections[channelId] = conn
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (a *App) refreshConfigPage(channelId string) {
	conn, exists := a.connections[channelId]
	if !exists {
		return
	}
	log.Println("Refresh command has been sent")
	if err := conn.WriteMessage(websocket.BinaryMessage, []byte("refresh")); err != nil {
		if err.Error() != "websocket: close sent" {
			log.Println(err)
		}
		delete(a.connections, channelId)
	}
}
