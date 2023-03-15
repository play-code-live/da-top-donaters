package main

import (
	"He110/donation-report-manager/internal/business/domain/configs"
	"He110/donation-report-manager/internal/business/domain/user"
	getConfigUseCase "He110/donation-report-manager/internal/business/use_cases/get_config_use_case"
	saveConfigUseCase "He110/donation-report-manager/internal/business/use_cases/save_config_use_case"
	saveTokenUseCase "He110/donation-report-manager/internal/business/use_cases/save_token_use_case"
	donationClient "He110/donation-report-manager/internal/pkg/donation-alerts-client"
	"He110/donation-report-manager/internal/web"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	host := "http://localhost"
	port := 8092
	socketHost := "ws://localhost:8092"
	client, err := donationClient.NewClient("10386", "x9Auz25j1PULNJXl4FScvSnnEKzJIf95oXXYPgvq", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	uStorage := user.NewStorage()
	cStorage := configs.NewStorage()

	useCaseGetConfig := getConfigUseCase.New(uStorage, cStorage)
	useCaseSaveToken := saveTokenUseCase.New(uStorage)
	useCaseSaveConfig := saveConfigUseCase.New(uStorage, cStorage)

	container := web.NewTemplateContainer("templates/base/*.gohtml")
	if err := container.FindAndRegister("templates/pages/"); err != nil {
		panic(err)
	}

	errGroup, ctx := errgroup.WithContext(context.Background())
	app := web.NewApp(container, client, useCaseGetConfig, useCaseSaveToken)

	router := mux.NewRouter()

	router.HandleFunc("/redirect", app.HandlerRedirect())
	router.HandleFunc("/redirect/{channelId}", app.HandlerChanneledRedirect())

	router.Path("/config/{channelId}").Methods(http.MethodGet).HandlerFunc(app.HandlerGetConfig(socketHost))
	router.Path("/config/{channelId}").Methods(http.MethodPost).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		channelId := vars["channelId"]
		if err = r.ParseForm(); err != nil {
			http.Error(w, "Cannot parse form data", http.StatusBadRequest)
			return
		}
		count, _ := strconv.Atoi(r.FormValue("donaters_count"))

		if _, err = useCaseSaveConfig.Perform(saveConfigUseCase.Parameters{
			ChannelId:     channelId,
			Title:         r.FormValue("panel_title"),
			DonatersCount: count,
		}); err != nil {
			http.Error(w, "Cannot save configs", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/config/"+channelId, http.StatusSeeOther)
	})
	router.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		if err = container.MustGet("config_anonymous").Execute(w, nil); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	router.HandleFunc("/socket/{channelId}", app.SocketBridge())

	ctx, _ = signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	// #### Web Server
	errGroup.Go(func() error {
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		}
		fmt.Println("Web server has started")
		return server.ListenAndServe()
	})

	<-ctx.Done()
	log.Println("service gracefully shutdown")
}
