package main
import (
	"log"
	"net/http"
	"os"
	"fmt"
	//"database/sql"
	"github.com/ivahaev/amigo"
	"github.com/kardianos/service"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	//"github.com/bradfitz/gomemcache/memcache"
)
type program struct{}
//Init variable
var log_file="dialing.log"
var file os.File
var logger service.Logger
//database config

var db_host="127.0.0.1:3306"
var db_user="dialing"
var db_pass="Dl@fj1ra"
var db_name="dialingdb"
var db mysql.Conn
//Asterisk variable
var settings = &amigo.Settings{Username: "trumpen", Password: "foobar", Host: "dev.dialingozone.com",Port:"1234"}
//memcache
//var mc := memcache.New("127.0.0.1:11211")
var a=amigo.New(settings)
var dial_timeout=25000
var agents=make(map[string]map[string]string)
var db_ratio map[string]float64
var ratio_up map[string]float64
var ratio_down map[string]float64
var trunk_list map[string]string
var cur_ratio map[string]float64
var agent_cnt map[string]int
var uniqueid_list map[string]string
var mute_arr map[string]string
var incall_cnarr map[string]string
var dial_cntarr map [string]int
var callarr map [string]string
var camparr map [string]string
var mdialarr map [string]string
var idial_cnarr map [string]int
var ans_cntarr map [string]int
var callarr2  map [string]string
var idarr map [string]string
var tapp_cntarr  map [string]int
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
	//register asterisk event listener
	//a.RegisterDefaultHandler(DefaultHandler)
	//a.RegisterHandler("Hangup",ast_hangup_event)
	//a.RegisterHandler("MeetmeJoin",ast_join)
	//a.RegisterHandler("MeetmeLeave",ast_leave)
	//a.RegisterHandler("OriginateResponse",ast_originate_response)
	c := make(chan map[string]string, 100)
	a.SetEventChannel(c)
	//listen http request
	http.HandleFunc("/user_state", state_check) // set router
	err := http.ListenAndServe(":8001", nil) // set listen port
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
	//Database mysql
	//db, err = sql.Open("mysql", db_string)
	db=mysql.New("tcp", "", db_host, db_user, db_pass, db_name)
	err = db.Connect()
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
	file.Close()
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
	//defer file.Close()

	// assign it to the standard logger
	log.SetOutput(file)
	log.Println("This is a test log entry")
}
func checkErr(err error) {
	if err != nil {
		plog(err.Error())
	}
}
func main() {
	svcConfig := &service.Config{
		Name:        "DialingService",
		DisplayName: "Dialing Service",
		Description: "Dialing Asterisk app.",
	}
	fmt.Println("Init Amigo")
	if(db== nil){
		plog("Start")
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	checkErr(err)
	/*if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}*/
	logger, err = s.Logger(nil)
	checkErr(err)
	err = s.Run()
	checkErr(err)
}