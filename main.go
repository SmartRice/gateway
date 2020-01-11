package main

import (
	"errors"
	"fmt"
	"github.com/SmartRice/gateway/model"
	"os"
	"runtime"
	"sync"
	"time"
)

type App struct {
	Name       string
	ServerList []model.APIServer
	launched   bool
	hostname   string
}

// NewApp Wrap application
func NewApp(name string) *App {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "undefined"
	}
	app := &App{
		Name:       name,
		ServerList: []model.APIServer{},
		launched:   false,
		hostname:   hostname,
	}
	return app
}

// Launch Launch app
func (app *App) Launch() error {

	if app.launched {
		return nil
	}

	app.launched = true

	name := app.Name + " / " + app.hostname
	fmt.Println("[ App " + name + " ] Launching ...")
	var wg = sync.WaitGroup{}

	// start connect to DB
	//for _, db := range app.DBList {
	//	err := db.Connect()
	//	if err != nil {
	//		fmt.Println("Connect DB error " + err.Error())
	//		return err
	//	}
	//}

	//fmt.Println("[ App " + name + " ] DBs connected.")

	// start servers
	for _, s := range app.ServerList {
		wg.Add(1)
		go s.Start(&wg)
	}
	fmt.Println("[ App " + name + " ] Servers started.")

	// start workers
	//for _, wk := range app.WorkerList {
	//	wg.Add(1)
	//	go wk.Execute()
	//}
	//fmt.Println("[ App " + name + " ] Workers started.")

	fmt.Println("[ App " + name + " ] Totally launched!")
	go callGCManually()
	wg.Wait()

	return nil
}

// callGCManually
func callGCManually() {
	for {
		time.Sleep(2 * time.Minute)
		runtime.GC()
	}
}

func (app *App) SetupAPIServer(t string) (model.APIServer, error) {
	var newID = len(app.ServerList) + 1
	var server model.APIServer
	server = model.NewHTTPAPIServer(newID, app.hostname)

	if server == nil {
		return nil, errors.New("server type " + t + " is invalid (HTTP/THRIFT)")
	}
	app.ServerList = append(app.ServerList, server)
	return server, nil
}

func main() {
	//var newID = len(app.ServerList) + 1
	var app = NewApp("SmartLife API Gateway")
	var server, _ = app.SetupAPIServer("HTTP")

	// ============= health check ==========================//
	now := time.Now()
	server.SetHandler(model.APIMethod.GET, "/gateway/v1/health-check", func(req model.APIRequest, res model.APIResponder) error {
		info := map[string]interface{}{}
		info["serviceName"] = app.Name
		info["startTime"] = now
		info["status"] = "OK"

		return res.Respond(&model.APIResponse{Status: model.APIStatus.Ok, Message: "OK", Data: []interface{}{info}})
	})

	server.Expose(80)

	app.Launch()
}
