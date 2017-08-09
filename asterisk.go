package main
import (
	"fmt"
	"strings"
	"strconv"
	"os/exec"
	"time"
	"github.com/bradfitz/gomemcache/memcache"
)
func DefaultHandler(m map[string]string) {
	fmt.Printf("Event received: %v\n\n\n", m)
}
//Join Agent to room 8800+ext
func ast_login(agent string, ext string , campaignid string,clientid string)(int , string) {
	conf_num:="8800"+ext
	if(agents[agent]["ownchannel"]!=""){
		a.Action(map[string]string{"Action":"Hangup","Channel":agents[agent]["ownchannel"]})

	}
	if(agents[agent]== nil) {
		agents[agent] = make(map[string]string)
	}
	agents[agent]["id"]=agent
	agents[agent]["ext"]=ext
	agents[agent]["campaignid"]=campaignid
	agents[agent]["conf_num"]=conf_num
	agents[agent]["status"]="standby"
	agents[agent]["clientid"]=clientid
	//Callee number
	agents[agent]["callee"]=""
	agents[agent]["channel"]=""
	plog( "Login "+agent+", "+ext+", "+conf_num,1)
	result, _ := a.Action(map[string]string{"Action":"Originate","Channel":"SIP/"+ext,"Context":"default","Exten":conf_num,"Priority":"1"})

	if(result["Response"]=="Error"){
		return 406,result["Message"]
	}
	go db_log("standby",agent,ext,campaignid)
	go db_getstate(campaignid)
	_,check:=cur_ratio[campaignid]
	if(!check){
		cur_ratio[campaignid]=1.0
	}
	return 200,"OK"
}

//Call to agent mobile phone and join to room 8800+ext
func ast_login_remote(agent string, ext string , campaignid string,dest string,clientid string)(int , string){
	conf_num:="8800"+ext
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
	agents[agent]["clientid"]=clientid
	//Callee number
	agents[agent]["callee"]=""
	agents[agent]["channel"]=""
	plog( "Login "+agent+", "+ext+", "+conf_num+", "+dest,1)
	result, _ := a.Action(map[string]string{"Action":"Originate","Channel":"SIP/"+dest+"\\@siptrunk","Context":"default","Exten":conf_num,"Priority":"1"})
	if(result["Response"]=="Error"){
		return 406,result["Message"]
	}
	db_log("standby",agent,ext,campaignid)
	db_getstate(campaignid)
	_,check:=cur_ratio[campaignid]
	if(!check){
		cur_ratio[campaignid]=1.0
	}
	return 200,"OK"

}
func ast_logout(agent string)(int , string){
	//plog( "Logout agent: "+agent,1)
	if val, ok := agents[agent]; ok {
		plog( "At logout, hangup: "+agent+", "+val["ext"]+", "+val["conf_num"]+", "+val["channel"]+", "+val["ownchannel"]+", "+val["campaignid"]+"",1)
		result, _ := a.Action(map[string]string{"Action":"Command","Command":"meetme kick "+val["conf_num"]+" all"})
		if(result["Response"]=="Error"){
			return 406,result["Message"]
		}
		db_log("logout",agent,val["ext"],val["campaignid"])
		delete(agents,agent)

	}
	return 200,"OK"
}
//Change agent to new campaign 
func ast_chcamp(agent string,  campaignid string)(int , string){
	var status string
	_,check:=agents[agent]
	if(check){
		status=agents[agent]["status"]
		agents[agent]["id"]=agent
		agents[agent]["campaignid"]=campaignid
		plog( "ast_chcamp changing campaign for ["+agent+"] to campaign ["+campaignid+"]",1)
		db_log(status,agent,agents[agent]["ext"],campaignid)
		db_getstate(campaignid)
	}else{
		plog("ast_chcamp "+agent+" is not logged in",1)
		return 400,"Agent is not logged in"
	}
	_,check=cur_ratio[campaignid]
	if(!check){
		cur_ratio[campaignid]=1.0
	}
	return 200,"OK"
}
//Process hang up event
func ast_hangup_event(m map[string]string){

	plog("Hangup!  "+m["Channel"]+" "+ m["Uniqueid"]+" " + m["Callerid"],1)
	var agent=call_arr[m["Uniqueid"]]["agent"]
	//var channel=strings.Split(m["Channel"],"@")[0]

	delete(inbound_arr,m["Channel"])
	delete(call_arr,m["Uniqueid"])
	if(agent !="") {
		usernum := agents[agent]["usernum"]
		conf_num := agents[agent]["conf_num"]
		mute(conf_num, usernum, agent)
	}

}
//check lai trang thai cua reason code
func ast_originate_response(m map[string] string){
	var agent string
	actionID:=strings.Split(m["ActionID"],":")
	callee:=actionID[0]
	ringcardid:=actionID[1]
	campaignid:=actionID[2]
	if(len(actionID)>3) {
		agent = actionID[3]
	}
	uid:=m["Uniqueid"]
	reason,_:=strconv.Atoi(m["Reason"])
	fromchannel:=m["Channel"]
	if(call_arr[uid]== nil) {
		call_arr[uid] = make(map[string]string)
	}
	call_arr[uid]["ringcardid"]=ringcardid
	if(agent!="") {
		call_arr[uid]["agent"] = agent
	}
	call_arr[uid]["callee"]=callee
	call_arr[uid]["campaignid"]=campaignid
	status:=""
	//if agent != nil && agent!=""{
	//	uidarr[agent]=uid
	//}
	switch reason{
	case 0:
		status="Failure"
	case 1:
		status="No Answer"
	case 4:
		status="OK"
	case 5:
		status="BUSY"
	case 8:
		status="Congested"
	default:
		status="Unknown fail"
	}
	plog("Originate result: "+callee+", "+uid+", "+strconv.Itoa(reason)+", "+status+", "+agent,1)
	//Use socket to control reatime status of call
	// /flashdata(campaignid)
	if(m["Response"]!="Success"){
		fail_cnt++
		fail_cntarr[campaignid]++
		ast_ratio_up(campaignid)

	}
	failtext:="ejsvar"
	if(reason==1 || reason==8){
		failtext="trasigt"
	}
	db_set_num_status(campaignid,ringcardid,failtext,callee)
	plog("To DB ringcard id:"+ringcardid+" number:"+callee+" Campaign ID:"+campaignid,1)
	if(fromchannel[0:10]!="Local/8800"){
		num_queue[campaignid]--
		localcount:=0
		for _, value := range agents {
			if (value["status"] == "ready" && value["campaignid"] == campaignid) {
				localcount++
			}
		}
		agent_cnt[campaignid]=localcount
		for _, value := range agents {
			if(value["status"]=="ready" && value["campaignid"]==campaignid){
				ratio:=calc_ratio(campaignid)
				num_queue[campaignid]+=ratio
				if(ratio>0){
					db_dial(ratio,campaignid)
				}
				break
			}
		}
	}
}

func ast_join(m map[string]string){
	channel:=m["Channel"]
	uid:=m["Uniqueid"]
	usernum:=m["Usernum"]
	//tmpclid:=idarr[uid]

	callee:=call_arr[uid]["callee"]
	ringcardid:=call_arr[uid]["ringcardid"]
	campaignid:=call_arr[uid]["campaignid"]
	//none:=1
	conf:=""
	plog("Meetme Join!, "+callee+","+channel+" "+uid+" "+m["Meetme"]+" "+usernum,1)
	if(m["Meetme"]=="8000000"){
		ans_cnt++
		ans_cntarr["campaignid"]++
		num_queue["campaignid"]--
		oldwhen:=time.Now().Unix()
		when:=0
		nextagent:=0
		for key, value := range agents {
			if(value["status"]=="ready" && campaignid==value["campaignid"]){
				if(nextagent==0){
					nextagent,_=strconv.Atoi(key)
				}
				when,_=strconv.Atoi(value["when"])
				if(int64(when)<oldwhen){
					nextagent,_=strconv.Atoi(key)
					oldwhen=int64(when)
				}
			}
			plog("Ast_join: Search agent for call: "+key+" , "+strconv.Itoa(when)+" , "+strconv.FormatInt(oldwhen,10),1)
		}
		if(nextagent>0){
			agent:=strconv.Itoa(nextagent)
			ext:=agents[agent]["ext"]
			conf=agents[agent]["conf_num"]
			plog("Found agent "+agent+", "+ext+", "+conf+", campaign id: "+campaignid,1)
			plog("Redirect:"+channel+", "+conf,1)
			agents[agent]["status"]="incall"
			agents[agent]["callee"]=callee
			agent_cnt[campaignid]--
			if(agent_cnt[campaignid]<0) {
				agent_cnt[campaignid] = 0
			}
			db_log("incall",agent,ext,campaignid)
			agents[agent]["channel"]=channel
			ast_ratio_down(campaignid)
			usernum=agents[agent]["usernum"]
			unmute(conf,usernum,agent)
			a.Action(map[string]string{"Action": "Redirect",
				"Channel":	channel,
				"Context":	"default",
				"Exten":	conf,
				"Priority":	"1",
			})
			//fmt.Println(result)
			plog("Ringcard: "+ringcardid+", "+callee,1)
			db_log_soundfile(ringcardid,campaignid,agent)
		}else{
			plog("No agent for call with ringcard: "+ringcardid,1)
			plog("Do hangup:"+channel+", "+conf,1)
			//delete(callarr,callee+":"+ringcardid)
			//delete(callarr2,callee+":"+uid)
			//delete(idarr,uid)
			tapp_cnt++
			tapp_cntarr[campaignid]++
			ratio_reset(campaignid)
			db_reg_tapp(ringcardid)
			a.Action(map[string]string{"Action": "Hangup",
				"Channel":	channel,
				"Context":	"default",
				"Priority":	"1",
			})
		}
	}else{
		conf=m["Meetme"]
		if _, ok := inbound_arr[channel]; ok {
			db_inbound_delete(channel)
		}
		if(channel[4:6]!="pseudo"){
			if(channel[4:7]!=conf[4:7]){
				//if(len(mdialarr[conf]["dest"])>5){
					for key, _ := range agents {
						if(agents[key]["conf_num"]==conf){
							agents[key]["channel"]=channel
							//agents[key]["ownchannel"]=channel
							agents[key]["status"]="incall"
							agents[key]["ringcardid"]=ringcardid
							ext:=agents[key]["ext"]
							db_log("incall",key,ext,campaignid)
							db_log_soundfile(ringcardid,campaignid,key)
							call_arr[uid]["agent"]=key
							url:="/dialing/card/"+ringcardid+"?dialnumber="+callee
							mc.Set(&memcache.Item{Key: "redirect_"+agents[key]["clientid"]+"_"+key, Value: []byte(url)})
							break
						}
					}
				//}
			}else{
				for key, _ := range agents {
					if(agents[key]["conf_num"]==conf){
						agents[key]["ownchannel"]=channel
						agents[key]["usernum"]=usernum
						mute(conf,usernum,key)
						db_user_connected(key,1)
						break
					}
				}
			}
		}
	}

}
//Set agent status = ready
func ast_ready(agent string)(int , string){
	ratio:=1
	ext:=agents[agent]["ext"]
	conf_num:=agents[agent]["conf_num"]
	campaignid:=agents[agent]["campaignid"]
	status:=agents[agent]["status"]
	plog("Dials: "+strconv.Itoa(dial_cnt)+" Answers: "+strconv.Itoa(ans_cnt)+" Tapp: "+strconv.Itoa(tapp_cnt)+" Fail: "+strconv.Itoa(fail_cnt),1)
	//use websocket to update status call not do yet
	//flashdata(campaignid)
	if (status=="standby"){
		plog("Agent "+agent+" is ready",1)
		agents[agent]["status"]="ready"
		agents[agent]["when"]=strconv.FormatInt(time.Now().Unix(),10)
		agent_cnt[campaignid]++
		ratio= calc_ratio(campaignid)
		plog("Ratio:"+strconv.Itoa(ratio),1)
		plog("ast_ready "+agent+", "+ext+", "+conf_num,1)
		db_log("ready",agent,ext,campaignid)
		if(ratio==0){
			plog("No need to call",1)
		}else {
			go db_dial(ratio,campaignid)
			num_queue[campaignid]+=ratio
		}
	}else{
		plog("Agent "+agent+" now in "+status+" not in stanby - ignore",1)
	}
	return 200,"OK"

}
//Set agent status = standby
func ast_standby(agent string)(int , string){
	if val, ok := agents[agent]; ok {
		ext:=agents[agent]["ext"]
		conf_num:=agents[agent]["conf_num"]
		campaignid:=agents[agent]["campaignid"]
		channel:=agents[agent]["channel"]
		db_log("standby",agent,ext,campaignid)
		plog("ast_standby "+agent+", "+ext+", "+conf_num,1)
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
			plog("Do hangup: "+agent+", "+channel+", "+conf_num,1)
			plog("Agent "+agent+" is standby",1)
			a.Action(map[string]string{"Action": "Hangup", "Channel": agents[agent]["channel"],"Context":"default","Exten":conf_num,"Priority":"1"})
			agents[agent]["ringcard_id"]=""
			agents[agent]["channel"]=""
			agents[agent]["callee"]=""
			//fmt.Println(result, err)
		}
		agents[agent]["status"]="standby"
		return 200,"OK"
	}
	return 400,"Agent not avaiable"

}
//Start record call
func ast_rec_start(agent string,filename string, clientid string) (int , string){
	if val, ok := agents[agent]; ok {
		if(val["status"]=="incall"){
			agents[agent]["clientid"]=clientid
			agents[agent]["logfile"]=filename
			agents[agent]["timestart"]=strconv.FormatInt(time.Now().Unix(),10)
			agents[agent]["recstatus"]="recording"
			db_log("recon",agent,val["ext"],val["campaignid"])
			result, _ := a.Action(map[string]string{"Action": "COMMAND", "COMMAND": "mixmonitor start "+val["channel"]+" "+filename})
			if(result["Response"]=="Error"){
				return 406,result["Message"]
			}
			plog("Record call for agent : "+agent+" with "+filename,1)

		}
	}
	return 200,"OK"
}

//Stop call record
func ast_rec_stop(agent string, filename string)(int , string){
	plog("Stop record call for agent : "+agent+" with "+filename,1)
	if val, ok := agents[agent]; ok {
		if(val["status"]=="incall"){
			db_log("recoff",agent,val["ext"],val["campaignid"])
			result, _ := a.Action(map[string]string{"Action": "COMMAND", "COMMAND": "mixmonitor stop "+val["channel"]})
			if(result["Response"]=="Error"){
				return 406,result["Message"]
			}

			db_log_rec(val["campaignid"],val["clientid"],val["logfile"],0)
			agents[agent]["timestart"]=""
			agents[agent]["acceptfile"]=""
			agents[agent]["acceptstart"]=""
			agents[agent]["acceptstatus"]="off"
			agents[agent]["acceptstop"]=""
		}
	}
	return 200,"OK"
}

//Auto call to Dest and join to default room 8000000.
func ast_dial(dest string,ringcardid string,campaignid string )(int , string){
	conf:="8000000"
	trunkname:="siptrunk"
	callerid:=""
	if val, ok := trunk_list[campaignid]; ok {
		if(len(val)>5){
			trunkname="tr"+val
			callerid=val
		}
		plog("Trunk: "+ trunkname,1)
	}
	dial_cnt++
	dial_cntarr[campaignid]++
	actionID:=dest+":"+ringcardid+":"+campaignid
	//callarr[dest+":"+ringcardid]=ringcardid
	//camparr[dest]=campaignid
	plog("Dial "+conf+" to :"+dest+" Ringcard: "+ringcardid+" Campaign: "+campaignid ,1)
	result, _ := a.Action(map[string]string{"Action": "Originate",
		"Channel": 	"Local/"+dest+"@selecttrunk/n",
		"Context": 	"default",
		"Exten":	"8000000",
		"Timeout":	strconv.Itoa(dial_timeout),
		"Callerid":	callerid,
		"Async":	"1",
		"ActionID":	actionID,
		"Variable":	"__myactionid="+actionID+",__TRUNKNAME="+trunkname,
		"Priority":	"1"	})
	if(result["Response"]=="Error"){
		return 406,result["Message"]
	}
	return 200,"OK"
}
func ast_mdial_trunk(agent string,ext string,dest string,ringcardid string)(int , string){
	ext=agents[agent]["ext"]
	conf:="8800"+ext
	//callarr[dest+":"+ringcardid]=ringcardid
	//mdialarr[conf]["ringcardid"]=ringcardid
	//mdialarr[conf]["dest"]=dest
	plog("Mdial "+conf+" to :"+dest+" Ringcard: "+ringcardid,1)
	campaignid:=agents[agent]["campaignid"]
	trunkname:="siptrunk"
	//callerid:=""
	if val, ok := trunk_list[campaignid]; ok {
		if(len(val)>5){
			trunkname="tr"+val
			//callerid=val
		}
		plog("Mdial Trunk: "+ trunkname ,1)
	}
	dial_cnt++
	dial_cntarr[campaignid]++
	//callarr[dest+":"+ringcardid]=ringcardid
	//camparr[dest]=campaignid
	plog("Mdial "+conf+" to :"+dest+" Ringcard: "+ringcardid+" Campaign: "+campaignid ,1)
	usernum:=agents[agent]["usernum"]
	unmute(conf,usernum,agent)
	actionID:=dest+":"+ringcardid+":"+campaignid+":"+agent
	result, _ := a.Action(map[string]string{"Action": "Originate",
		"Channel": 	"Local/"+conf,
		"Context": 	"mdialt",
		"Exten":	"SIP/"+dest,
		"Timeout":	strconv.Itoa(dial_timeout),
		//"Callerid":	callerid,
		"Async":	"1",
		"ActionID":	actionID,
		"Variable":	"__myactionid="+actionID+",__TRUNKNAME="+trunkname,
		//"Variable":	"__TRUNKNAME="+trunkname,
		"Priority":	"1"	})
	if(result["Response"]=="Error"){
		return 406,result["Message"]
	}
	return 200,"OK"
}
func ast_hangup(agent string)(int , string){
	plog("Hangup call for agent: "+agent,1)
	if val, ok := agents[agent]; ok {
		if (val["status"] == "incall"){
			ext := val["ext"]
			conf := val["conf_num"]
			channel := val["channel"]
			campaignid := val["campaignid"]
			plog("Do hangup: "+agent+" ,"+channel+" ,"+conf ,1)
			plog("Do hangup: Agent "+agent+" is standby",1)
			agents[agent]["status"]="standby"
			agents[agent]["ringcardid"]=""
			usernum:=agents[agent]["usernum"]
			mute(conf,usernum,agent)
			db_log("standby",agent,ext,campaignid)
			agents[agent]["channel"]=""
			agents[agent]["callee"]=""
			result, _ := a.Action(map[string]string{"Action": "Hangup",
				"Channel":	channel,
				"Context":	"default",
				"Exten":	conf,
				"Priority":	"1",
			})
			if(result["Response"]=="Error"){
				return 406,result["Message"]
			}

		}
	}
	return 200,"OK"
}

func ast_leave(m map[string]string){
	channel:=m["Channel"]
	//usernum:=m["Usernum"]
	ext:=channel[4:7]
	agent:=""
	plog("Ast_leave: "+channel+", "+ext,1)
	for key, value := range agents {
		if(value["ext"]==ext){
			agent=key
			status:=value["status"]
			current_channel:=value["channel"]
			db_user_connected(agent,0)
			agents[agent]["status"]="standby"
			if(status=="incall") {
				a.Action(map[string]string{"Action": "Hangup",
					"Channel":        current_channel,
				})
			}
		}
	}

}

func mute(conf_num string ,user string,agent string){
	_, err := a.Action(map[string]string{"Action": "MeetmeMute", "Meetme": conf_num,"Usernum":user})
	//fmt.Println(result, err)
	checkErr(err)
	plog("User "+user+" in conference "+conf_num+" MUTE",1)
}
func unmute(conf_num string ,user string,agent string){
	_, err := a.Action(map[string]string{"Action": "MeetmeUnMute", "Meetme": conf_num,"Usernum":user})
	//fmt.Println(result, err)
	checkErr(err)
	plog("User "+user+" in conference "+conf_num+" UNMUTE",1)
}
//Calulate ratio
func calc_ratio(campaignid string)int {
	plog("current ratio:"+strconv.FormatFloat(cur_ratio[campaignid], 'E', -1, 64)+ " agent_cnt:"+strconv.Itoa(agent_cnt[campaignid]) +" num queue:"+strconv.Itoa(num_queue[campaignid]),1)
	ratio :=int((cur_ratio[campaignid] * float64(agent_cnt[campaignid])) - float64(num_queue[campaignid]))
	if(ratio>=2){
		return 2
	}else if(ratio<1){
		return 0
	}else{
		return 1
	}
}

func ast_ratio(ratio int, campaignid string,timeout int)(int , string){
	if(timeout>10000 && timeout < 90000){
		dial_timeout=timeout
		plog("Set dial timeout = "+strconv.Itoa(timeout)+" for campaign "+campaignid,1)
	}
	if(ratio > 1 && ratio< 10){
		db_ratio[campaignid]=float64(ratio)
		plog("Set dial timeout = "+strconv.Itoa(timeout)+" for campaign "+campaignid,1)
	}
	return 200,"OK"
}
//Increase ratio
func ast_ratio_up(campaignid string){
	if(db_ratio[campaignid]==0){
		db_ratio[campaignid]=default_ratio
	}
	if(cur_ratio[campaignid]<db_ratio[campaignid]){
		cur_ratio[campaignid]+=ratio_up[campaignid]
	}else{
		cur_ratio[campaignid]=db_ratio[campaignid]
	}
}
//Decrease ratio
func ast_ratio_down(campaignid string){

	cur_ratio[campaignid]-=ratio_up[campaignid]
	if(cur_ratio[campaignid]<1){
		cur_ratio[campaignid]=1.0
	}
}
//Reset ratio
func ratio_reset(campaignid string){
	if(db_ratio[campaignid]==0){
		db_ratio[campaignid]=default_ratio
	}
	cur_ratio[campaignid]=1.0
	plog( "reset:Ratio for campaign "+campaignid,1)
}
func ast_stepratio(tratio_up float64,tratio_down float64, campaignid string)(int , string){
	if(tratio_down>-2.0 && tratio_down<2.0){
		ratio_down[campaignid]=tratio_down
		plog("Set ratio down ="+strconv.FormatFloat(tratio_down,'E', -1, 64)+" for campaign:"+campaignid,1)
	}
	if(tratio_up>-2.0 && tratio_up<2.0){
		ratio_up[campaignid]=tratio_up
		plog("Set ratio up ="+strconv.FormatFloat(tratio_up,'E', -1, 64)+" for campaign:"+campaignid,1)
	}
	return 200,"OK"
}
//End of number
func ast_eon(campaignid string){
	num_queue[campaignid]--
	if(num_queue[campaignid]<0){
		num_queue[campaignid]=0
	}
	plog( "decrease num_queue for campaign "+campaignid,1)
}
/*
func flashdata(campaignid string){
	//os.Mkdir("/var/www/flash/"+campaignid,"755")
	file, err := os.Create("/var/www/flash/"+campaignid+"/variables.txt")
	if err != nil {
		log.Fatalln("Failed to open file data",  ":", err)
	}
	defer file.Close()
	file.WriteString(Sprintln("ratio=%0.1f&waiting=%d&tapp=%d&totaltparingda=%d&totaltanswers=%d&slask=1\n",cur_ratio[campaignid],agent_cnt[campaignid],tapp_cntarr[campaignid],dial_cntarr[campaignid],ans_cntarr[campaignid]))
}
*/
// Set default ratio
func set_default_ratio(campaignid string){
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

func ast_idial(agent string,ext string,dest string,ringcard string, channel string)(int , string){
	usernum:=agents[agent]["usernum"]
	campaignid:=agents[agent]["campaignid"]
	conf:="8800"+ext
	unmute(conf,usernum,agent)
	dial_cnt++
	dial_cntarr[campaignid]++
	agents[agent]["channel"]=channel
	plog("Inbound from "+dest+" to "+conf+" Ringcard:"+ringcard,1)
	result, _:=a.Action(map[string]string{"Action": "Redirect",
		"Channel":channel,
		"Context":"default",
		"Exten":conf,
		"Priority":"1"})
	if(result["Response"]=="Error"){
		return 406,result["Message"]
	}
	return 200,"OK"
}
func ast_transfer(agent string,toagent string, phonenum string)(int , string){
	cur_conf:="8800"+agents[agent]["ext"]
	conf:="8800"+agents[toagent]["ext"]
	ringcardid:=agents[agent]["ringcardid"]
	channel:=agents[agent]["channel"]
	campaignid:=agents[agent]["campaignid"]
	usernum:=agents[agent]["usernum"]
	agents[agent]["logging"]="off"
	db_log("recoff",agent,agents[agent]["ext"],campaignid)
	a.Action(map[string]string{"Action": "COMMAND",
		"Command":"mixmonitor stop "+channel	})
	//plog("standby: Agent "+agent+" is standby",1)
	agents[agent]["status"]="standby"
	db_log("standby",agent,agents[agent]["ext"],campaignid)
	agents[agent]["ringcardid"]=""
	agents[agent]["channel"]=""
	agents[agent]["callee"]=""
	mute(cur_conf,usernum,agent)
	unmute(conf, agents[toagent]["usernum"],toagent)
	plog("Transfer call from "+agent+" , "+campaignid+" , "+phonenum+" , "+ringcardid+" , "+channel+" to "+toagent,1)
	db_user_connected(agent,0)
	result, _:=a.Action(map[string]string{"Action": "Redirect",
		"Channel":channel,
		"Context":"default",
		"Exten":conf,
		"Priority":"1"})
	if(result["Response"]=="Error"){
		return 406,result["Message"]
	}
	return 200,"OK"
}
func ast_record(phonenum string, recfile string,trunk string)(int,string){
	channel:="SIP/"+phonenum
	mc.Set(&memcache.Item{Key: "record_"+phonenum, Value: []byte(channel)})
	if(len(phonenum)>5){
		channel="SIP/tr"+trunk+"/"+phonenum
	}
	result, _ := a.Action(map[string]string{"Action": "Originate",
		"Channel": 	channel,
		"Context": 	"record",
		"Exten":	"s",
		"Timeout":	strconv.Itoa(dial_timeout),
		"Async":	"1",
		"ActionID":	"record_"+phonenum,
		"Variable":	"__recfile="+recfile,
		"Priority":	"1"	})
	if(result["Response"]=="Error"){
		return 406,result["Message"]
	}
	return 200,"OK"

}
func ast_record_stop(phonenum string, recfile string,delete int)(int,string){
	channel, _ := mc.Get("record_"+phonenum)
	mc.Delete("record_"+phonenum)
	a.Action(map[string]string{"Action": "Hangup",
		"Channel":string(channel.Value)	})
	if(delete == 1){
		cmd :=exec.Command("rm"," /var/lib/asterisk/sounds/dialplan/"+recfile+".wav")
		cmd.Run()
	}
	return 200,"OK"
}

func ast_peerstatus(peer string)(int,string){
	result, _ := a.Action(map[string]string{"Action": "SIPshowpeer",
		"Peer": 	peer	})
	status:="FAIL"
	if(result["Response"]=="Error"){
		return 406,result["Message"]
	}else{
		plog(result["Status"],1)
		if(result["Status"][0:2]=="OK"){
			status="OK"
		}
	}
	plog(status,1)
	return 200,"OK"
}
func ast_delete_peercache()(int,string){
	result, _ := a.Action(map[string]string{"Action": "SIPpeerstatus","ActionID":"allpeers"})
	if(result["Response"]=="Error"){
		return 406,result["Message"]
	}
	cmd,_ :=exec.Command("sudo /usr/sbin/asterisk -rx 'core show channels concise'|grep '@selecttrunk'|awk -F '!' '$2 ~ /dial-out/ && $5 ~ /Ring/ && $9 ~ /"+"1133"+":/'|wc -l").Output()
	fmt.Println(cmd)
	return 200,"OK"
}
func ast_peer_status(m map[string]string){
	if m["ActionID"]=="allpeers"{
		if len(m["Peer"])==7 {
			//fmt.Println(m["Peer"][4:])
			mc.Delete("peer_"+m["Peer"][4:])
		}
	}
}
func ast_channel(m map[string]string){
	if m["ActionID"]=="allchannel" && m["Context"]=="dial-out" && strings.Contains(m["Channel"],"selecttrunk"){

	}
}
func checknumqueue(){
	cmd,err :=exec.Command("sudo /usr/sbin/asterisk -rx 'core show channels concise'|wc -l").Output()
	fmt.Println(cmd,err)
}