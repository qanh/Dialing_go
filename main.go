package main
import (
	"log"
	"net/http"
	"os"
	"database/sql"
	"github.com/ivahaev/amigo"
	"github.com/kardianos/service"
	_ "github.com/go-sql-driver/mysql"
)
type program struct{}
//Init variable
var log_file="dialing.log"
var logger service.Logger
//database config
var db_string="user:password@/dbname"
var db *DB
//Asterisk variable
var settings = &amigo.Settings{Username: "trumpen", Password: "foobar", Host: "dev.dialingozone.com",Port:"1234"}
var a=amigo.New(settings)
var dial_timeout=25000
var agents=make(map[string]map[string]string)
var db_ratio map[string]float
var ratio_up map[string]int
var ratio_down map[string]int
var trunk_list map[string]string
var cur_ratio map[string]float
var agent_cnt map[string]int
var uniqueid_list map[string]string
var mute_arr map[string]string
var incall_cnarr map[string]string
//unknown what it is
var num_queue map[string]int
var default_ratio=1.5
var default_ratio_up=0.1
var default_ratio_down=0.2
var dial_cnt=0
var ans_cnt=0
var tapp_cnt=0
var fail_cnt=0

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	//listen asterisk event and request
	a.Connect()
	a.RegisterDefaultHandler(DefaultHandler)
	c := make(chan map[string]string, 100)
	a.SetEventChannel(c)
	//listen http request
	http.HandleFunc("/user_state", state_check) // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
	//Database mysql
	db, err = sql.Open("mysql", db_string)
	checkErr(err)
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	db.Close()
	return nil
}
func plog(str string){
	log.Println("LOG: ",str)
}
func init(){
	file, err := os.OpenFile(log_file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file",  ":", err)
	}
	defer file.Close()

	// assign it to the standard logger
	log.SetOutput(file)
}
func checkErr(err error) {
	if err != nil {
		plog(err)
	}
}
func main() {
	svcConfig := &service.Config{
		Name:        "DialingService",
		DisplayName: "Dialing Service",
		Description: "Dialing Asterisk app.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	checkErr(err)
	logger, err = s.Logger(nil)
	checkErr(err)
	err = s.Run()
	checkErr(err)
}