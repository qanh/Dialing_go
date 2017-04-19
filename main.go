package main
import (
	"log"
	"net/http"
	"os"
	"fmt"
	"database/sql"
	"path"
	"github.com/ivahaev/amigo"
	"github.com/kardianos/service"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/bradfitz/gomemcache/memcache"
)
type program struct{}
//Init variable
var log_file="dialing.log"
var file os.File
var port="8000"
//var logger service.Logger
//database config

//var db_host="127.0.0.1:3306"
//var db_user="dialing"
//var db_pass="Dl@fj1ra"
//var db_name="dialingdb"
//var db =autorc.New("tcp", "", db_host, db_user, db_pass, db_name)
var db *sql.DB
var db_string="dialing:Dl@fj1ra@/dialingdb"
//Asterisk variable
var settings = &amigo.Settings{Username: "trumpen", Password: "foobar", Host: "dev.dialingozone.com",Port:"1234"}
//memcache
//var mc := memcache.New("127.0.0.1:11211")
var a=amigo.New(settings)
var dial_timeout=25000
var agents=make(map[string]map[string]string)
var db_ratio =make(map[string]float64)
var ratio_up =make(map[string]float64)
var ratio_down =make(map[string]float64)
var trunk_list =make(map[string]string)
var cur_ratio =make(map[string]float64)
var agent_cnt =make(map[string]int)
var list_ringcard=make(map[string]int)
//var mute_arr =make(map[string]string)
//var incall_cnarr =make(map[string]string)
var dial_cntarr =make(map [string]int)
//var callarr =make(map [string]string)
//var mdialarr =make(map [string]string)
var inbound_arr =make(map [string]int)
var ans_cntarr =make(map [string]int)
var call_arr  =make(map [string]map[string]string)
//var idarr =make(map [string]string)
var tapp_cntarr  =make(map [string]int)
//unknown what it is
var num_queue =make(map[string]int)
var default_ratio=1.5
var default_ratio_up=0.1
var default_ratio_down=0.2
var dial_cnt=0
var ans_cnt=0
var tapp_cnt=0
var fail_cnt=0
var fail_cntarr=make(map [string]int)

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	//db.Close()
	file.Close()
	return nil
}
func plog(str string,level ...int){
	debug:=4
	if(level[0]<=debug) {
		log.Println("LOG: ", str)
	}
}
func init(){
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := path.Dir(ex)

	file, err := os.OpenFile(exPath+"/"+log_file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file",  ":", err)
	}
	//defer file.Close()
	// assign it to the standard logger
	log.SetOutput(file)
	db, err = sql.Open("mysql", db_string)
	//db=mysql.New("tcp", "", db_host, db_user, db_pass, db_name)
	if err != nil {
		plog("DB error",1)
	}else{
		plog("DB connected",1)
	}
	//listen asterisk event and request
	go a.Connect()
	//register asterisk event listener
	//a.RegisterDefaultHandler(DefaultHandler)
	a.RegisterHandler("Hangup",ast_hangup_event)
	a.RegisterHandler("MeetmeJoin",ast_join)
	a.RegisterHandler("MeetmeLeave",ast_leave)
	a.RegisterHandler("OriginateResponse",ast_originate_response)
	c := make(chan map[string]string, 100)
	a.SetEventChannel(c)
	//listen http request

	http.HandleFunc("/user_state", state_check) // set router
	err = http.ListenAndServe(":"+port, nil) // set listen port

	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
		plog("ListenAndServe Error",1)
	}else{
		fmt.Println("ListenAndServe on port "+port,1)
	}

}
func checkErr(err error) {
	if err != nil {
		plog(err.Error(),1)
	}
}
func main() {
	svcConfig := &service.Config{
		Name:        "DialingService",
		DisplayName: "Dialing Service",
		Description: "Dialing Asterisk app.",
	}
	fmt.Println("Start")
	//Database mysql


	prg := &program{}
	s, err := service.New(prg, svcConfig)
	checkErr(err)
	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	//logger, err = s.Logger(nil)
	//checkErr(err)
	err = s.Run()
	checkErr(err)
}