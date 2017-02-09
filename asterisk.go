package main
import (
	"fmt"
	"os/exec"
)
func DefaultHandler(m map[string]string) {
	fmt.Printf("Event received: %v\n", m)
}
//Join Agent to room 8800+ext
func ast_login(agent string, ext string , campaignid string){
	conf_num:="8800"+ext
	if(ext>0){
		if(agents[agent]["ownchannel"]!=""){
			result, err := a.Action(map[string]string{"Action":"Hangup","Channel":agents[agent]["ownchannel"]})
			fmt.Println(result, err)
		}
		agents[agent]["id"]=agent
		agents[agent]["ext"]=ext
		agents[agent]["campaignid"]=campaignid
		agents[agent]["conf_num"]=conf_num
		agents[agent]["status"]="standby"
		//Callee number
		agents[agent]["callee"]=""
		agents[agent]["channel"]=""
		plog( "Login "+agent+", "+ext+", "+conf_num)
		result, err := a.Action(map[string]string{"Action":"Originate","Channel":"SIP/"+ext,"Context":"default","Exten":conf_num,"Priority":1})
		fmt.Println(result, err)
		db_log("standby",agent,ext,campaignid)
		db_getstate(campaignid)
	}else {
		plog( "Agent "+agent+" miss extension")
	}
	default_ratio(campaignid)
}
//Call to agent mobile phone and join to room 8800+ext
func ast_login_remote(agent string, ext string , campaignid string,dest string){
	conf_num:="8800"+ext
	if(ext>0) {
		if (agents[agent]["ownchannel"] != "") {
			result, err := a.Action(map[string]string{"Action":"Hangup", "Channel":agents[agent]["ownchannel"]})
			fmt.Println(result, err)
		}
		agents[agent]["id"]=agent
		agents[agent]["ext"]=ext
		agents[agent]["campaignid"]=campaignid
		agents[agent]["conf_num"]=conf_num
		agents[agent]["remote_num"]=dest
		agents[agent]["status"]="standby"
		//Callee number
		agents[agent]["callee"]=""
		agents[agent]["channel"]=""
		plog( "Login "+agent+", "+ext+", "+conf_num+", "+dest)
		result, err := a.Action(map[string]string{"Action":"Originate","Channel":"SIP/"+dest+"\\@siptrunk","Context":"default","Exten":conf_num,"Priority":1})
		fmt.Println(result, err)
		db_log("standby",agent,ext,campaignid)
		db_getstate(campaignid)

	}else {
		plog( "Agent "+agent+" miss extension")
	}
	default_ratio(campaignid)
}
func ast_logout(agent string){
	plog( "Logout agent: "+agent)
	if val, ok := agents[agent]; ok {
		plog( "At logout, hangup: "+agent+", "+val["ext"]+", "+val["conf_num"]+", "+val["channel"]+", "+val["ownchannel"]+", "+val["campaignid"]+"")
		result, err := a.Action(map[string]string{"Action":"Command","Command":"meetme kick "+val["conf_num"]+" all"})
		fmt.Println(result, err)
		plog("Do logout: "+agent+","+val["ext"]+","+val["conf_num"]);
		db_log("logout",agent,val["ext"],val["campaignid"])
		delete(agents,agent)
	}

}
//Change agent to new campaign 
func ast_chcamp(agent string, ext string , campaignid string){
	var status string
	tmpagent,check:=agents[agent]
	if(check){
		status=tmpagent["status"]
		agents[agent]["id"]=agent
		agents[agent]["ext"]=ext
		agents[agent]["campaignid"]=campaignid
		plog( "ast_chcamp changing campaign for ["+agent+"] with ext ["+ext+"] to campaign ["+campaignid+"]")
		db_log(status,agent,ext,campaignid)
		if(Len(ext)>0 && Len(agents[agent]["ownchannel"]>0)){
			result, err := a.Action(map[string]string{"Action": "Hangup", "Channel": agents[agent]["ownchannel"]})
			fmt.Println(result, err,agent,campaignid)
		}
		db_getstate(campaignid)
		if(db_ratio[campaignid]==""){
			db_ratio[campaignid]=default_ratio
		}
		if(cur_ratio[campaignid]==""){
			cur_ratio[campaignid]=1.0
		}
		if(ratio_up[campaignid]==""){
			ratio_up[campaignid]=default_ratio_up
		}
		if(ratio_down[campaignid]==""){
			ratio_down[campaignid]=default_ratio_down
		}

	}else{
		plog("ast_chcamp "+agent+" is not logged in")
	}
}
//Process hang up event
func ast_hangup_event(m map[string]string){

	plog("Hangup!  "+m["Channel"]+" "+ m["Uniqueid"]+" " + m["Callerid"]+"\n");
	var agent=uniqueid_list[m["Uniqueid"]]
	var channel=strings.Split(m["Channel"],"@")[0]
	plog("Agent: "+agent+"  hangup with UniqueID"+ m["Uniqueid"])
	delete(incall_cnarr,m["Channel"])
	if(mute_arr[channel]!=""){
		agent=mute_arr[channel]
		plog("########## Channel: "+channel+" AGENT: "+agent)
		usernum:=agents[agent]["usernum"]
		conf_num:=agents[agent]["conf_num"]
		mute(conf_num, usernum)

	}
}
//Set agent status = ready
func ast_ready(agent string){
	ratio:=1
	ext:=agents[agent]["ext"]
	conf_num:=agents[agent]["conf_num"]
	campaignid:=agents[agent]["campaignid"]
	status:=agents[agent]["status"]
	plog("Dials: "+dial_cnt+" Answers: "+ans_cnt+" Tapp: "+tapp_cnt+" Fail: "+fail_cnt)
	flashdata(campaignid)
	if (status=="standby"){
		plog("Agent "+agent+" is ready")
		agents[agent]["status"]="ready"
		agents[agent]["when"]=time.Now().Unix()
		ratio= calc_ratio(campaignid)
		plog("Ratio:"+ratio)
		plog("ast_ready "+agent+", "+ext+", "+conf_num)
		db_log("ready",agent,ext,campaignid)
		if(ratio==0){
			plog("No need to call")
		}else {
			db_dial(ratio,campaignid)
			numqueue[campaignid]+=ratio
		}
	}else{
		plog("Agent "+agent+" now in "+status+" not in stanby - ignore")
	}

}
//Set agent status = standby
func ast_standby(agent string){
	ext:=agents[agent]["ext"]
	conf_num:=agents[agent]["conf_num"]
	campaignid:=agents[agent]["campaignid"]
	channel:=agents[agent]["channel"]
	db_log("standby",agent,ext,campaignid)
	plog("ast_standby "+agent+", "+ext+", "+conf_num)
	if val, ok := agents[agent]; ok {
		if(val["status"]=="ready" ){
			agent_cnt[campaignid]--
			if(agent_cnt[campaignid]<0) {
				agent_cnt[campaignid] = 0
			}
		}
		if(val["status"]=="incall" ){
			agent_cnt[campaignid]--
			if(agent_cnt[campaignid]<0) {
				agent_cnt[campaignid] = 0
			}
			plog("Do hangup: "+agent+", "+channel+", "+conf_num)
			plog("Agent "+agent+" is standby")
			agents[agent]["ringcard_id"]=""
			agents[agent]["channel"]=""
			agents[agent]["callee"]=""
			result, err := a.Action(map[string]string{"Action": "Hangup", "Channel": agents[agent]["channel"],"Context":"default","Exten":conf_num,"Priority":1})
			fmt.Println(result, err)
		}
	}
	agents[agent]["status"]="standby"

}
//Start record call
func ast_rec_start(agent string,filename string, clientid string){
	if val, ok := agents[agent]; ok {
		if(val["status"]=="incall"){
			agents[agent]["clientid"]=clientid
			agents[agent]["logfile"]=filename
			agents[agent]["timestart"]=time.Now().Unix()
			agents[agent]["recstatus"]="recording"
			db_log("recon",agent,val["ext"],val["campaignid"])
			result, err := a.Action(map[string]string{"Action": "COMMAND", "COMMAND": "mixmonitor start "+val["channel"]+" "+filename})
			fmt.Println(result, err)
			plog("Record call for agent : "+agent+" with "+filename)
		}
	}
}
//Start record accept call
func ast_rec_accept_start(agent string,filename string, clientid string){
	if val, ok := agents[agent]; ok {
		if(val["status"]=="incall") {
			if(val["recstatus"]=="recording"){
				agents[agent]["acceptstart"]=time.Now().Unix()
				agents[agent]["logging"]="on"
				agents[agent]["acceptstatus"]="on"
				agents[agent]["acceptfile"]=filename
				plog("Mark accept start for agent : "+agent)
			}else{
				agents[agent]["clientid"]=clientid
				agents[agent]["timestart"]=time.Now().Unix()
				agents[agent]["recstatus"]="recording"
				agents[agent]["acceptstatus"]="on"
				agents[agent]["acceptstart"]=time.Now().Unix()
				agents[agent]["logging"]="off"
				db_log("recon",agent,val["ext"],val["campaignid"])
				result, err := a.Action(map[string]string{"Action": "COMMAND", "COMMAND": "mixmonitor start "+val["channel"]+" "+filename})
				fmt.Println(result, err)
				plog("Record call for agent : "+agent+" with "+filename)
			}
		}
	}
}
//Stop call record
func ast_rec_stop_mix(agent string, filename string,ringcard_id string){
	plog("Stop record call for agent : "+agent+" with "+filename)
	if val, ok := agents[agent]; ok {
		if(val["status"]=="incall"){
			db_log("recoff",agent,val["ext"],val["campaignid"])
			result, err := a.Action(map[string]string{"Action": "COMMAND", "COMMAND": "mixmonitor stop "+val["channel"]})
			fmt.Println(result, err)
			if(val["acceptstatus"]=="on"){
				cutstart:=int(val["acceptstart"])-int(val["timestart"])
				cutlen:=time.Now().Unix()-int(val["acceptstart"])
				db_log_rec(ringcard_id,val["campaignid"],val["clientid"],val["logfile"],0)
				db_log_rec(ringcard_id,val["campaignid"],val["clientid"],val["acceptfile"],1)
				cmdstr := "/home/trumpen/bin/mp3convertaccept.sh "+filename+" "+val["campaignid"]+" "+val["clientid"]+" "+val["acceptfile"]+" "+cutstart+" "+cutlen+" "+ringcard_id
				exec.Command(cmdstr).Start()
			}else{
				db_log_rec(ringcard_id,val["campaignid"],val["clientid"],val["logfile"],0)
				cmdstr := "/home/trumpen/bin/mp3convertdelayed.sh "+filename+" "+val["campaignid"]+" "+val["clientid"]
				exec.Command(cmdstr).Start()
			}
			agents[agent]["timestart"]=""
			agents[agent]["acceptfile"]=""
			agents[agent]["acceptstart"]=""
			agents[agent]["acceptstatus"]="off"
			agents[agent]["acceptstop"]=""
		}
	}
}
//Stop call accept record
func ast_rec_accept_stop(agent string, filename string){
	if val, ok := agents[agent]; ok{
		if(val["status"]=="incall") {
			if(agents[agent]["logging"]=="on"){
				agents[agent]["acceptstop"]=time.Now().Unix()
				agents[agent]["acceptfile"]=filename
				plog("Mark accept stop for agent : "+agent+" with "+filename)
			}else{
				db_log("recoff",agent,val["ext"],val["campaignid"])
				db_log_rec(val["ringcardid"],val["campaignid"],val["clientid"],val["logfile"],0)
				db_log_rec(val["ringcardid"],val["campaignid"],val["clientid"],val["acceptfile"],1)
				result, err := a.Action(map[string]string{"Action": "COMMAND", "COMMAND": "mixmonitor stop "+val["channel"]})
				fmt.Println(result, err)
			}
		}

	}
}
func mute(conf_num string ,user string){
	result, err := a.Action(map[string]string{"Action": "MeetmeMute", "Meetme": conf_num,"Usernum":user})
	fmt.Println(result, err)
	checkErr(err)
	plog("User "+user+" in conference "+conf_num+" MUTE");
}
func unmute(conf_num string ,user string){
	result, err := a.Action(map[string]string{"Action": "MeetmeUnMute", "Meetme": conf_num,"Usernum":user})
	fmt.Println(result, err)
	checkErr(err)
	plog("User "+user+" in conference "+conf_num+" UNMUTE");
}
//Calulate ratio
func calc_ratio(campaignid string)string {
	ratio :=int((cur_ratio[campaignid]*agent_cnt[campaignid]) - num_queue[campaignid])
	if(ratio>=2){
		return 2
	}else if(ratio<1){
		return 0
	}else{
		return 1
	}
}

func ast_ratio(ratio int, campaignid string,timeout int){
	if(timeout>10000 && timeout < 90000){
		dial_timeout=timeout
		plog("Set dial timeout = "+wait_time+" for campaign "+campaignid)
	}
	if(ratio > 1 && ratio< 10){
		db_ratio[campaignid]=ratio
		plog("Set dial timeout = "+wait_time+" for campaign "+campaignid)
	}

}
//Increase ratio
func ratio_up(campaignid string){
	if(db_ratio[campaignid]==""){
		db_ratio[campaignid]=default_ratio
	}
	if(cur_ratio[campaignid]<db_ratio[campaignid]){
		cur_ratio[campaignid]+=ratio_up[campaignid]
	}else{
		cur_ratio[campaignid]=db_ratio[campaignid]
	}
}
//Decrease ratio
func ratio_down(campaignid string){

	cur_ratio[campaignid]-=ratio_up[campaignid]
	if(cur_ratio[campaignid]<1){
		cur_ratio[campaignid]=1.0
	}
}
//Reset ratio
func ratio_reset(campaignid string){
	if(db_ratio[campaignid]==""){
		db_ratio[campaignid]=default_ratio
	}
	cur_ratio[campaignid]=1.0
	plog( "reset:Ratio for campaign "+campaignid);
}
//End of number
func ast_eon(campaignid string){
	num_queue[campaignid]--;
	if(num_queue[campaignid]<0){
		num_queue[campaignid]=0
	}
	plog( "decrease num_queue for campaign "+campaignid);
}
// Set default ratio
func default_ratio(campaignid string){
	if(db_ratio[campaignid]==0){
		db_ratio[campaignid]=default_ratio
	}
	if(cur_ratio[campaignid]==0){
		cur_ratio[campaignid]=1.0
	}
	if(ratio_up[campaignid]==0){
		ratio_up[campaignid]=default_ratio_up
	}
	if(ratio_down[campaignid]==0){
		ratio_down[campaignid]=default_ratio_down
	}
}