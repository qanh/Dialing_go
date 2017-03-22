package main
import (
	"fmt"
//	"os/exec"
//	"time"
)
func DefaultHandler(m map[string]string) {
	fmt.Printf("Event received: %v\n", m)
}
//Join Agent to room 8800+ext
func ast_login(agent string, ext string , campaignid string)string {
	conf_num:="8800"+ext
	if(len(ext)>0){
		if(agents[agent]["ownchannel"]!=""){
			result, err := a.Action(map[string]string{"Action":"Hangup","Channel":agents[agent]["ownchannel"]})
			fmt.Println(result, err)
		}
		if(agents[agent]== nil) {
			agents[agent] = make(map[string]string)
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
		result, err := a.Action(map[string]string{"Action":"Originate","Channel":"SIP/"+ext,"Context":"default","Exten":conf_num,"Priority":"1"})
		fmt.Println(result,err)
		fmt.Println("sdaf")
		db_log("standby",agent,ext,campaignid)
		db_getstate(campaignid)
		if(result["Message"]=="Error"){
			return result["Message"]
		}
		return "OK"
	}else {
		plog( "Agent "+agent+" miss extension")
		return "Agent "+agent+" miss extension"
	}
	//set_default_ratio(campaignid)

}
/*
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
//check lai trang thai cua reason code
func ast_originate_response(m map[string] string){
	actionID:=Split(m["ActionID"])
	callee:=actionID[0]
	ringcardid:=actionID[1]
	agent:=actionID[2]
	uid:=m["Uniqueid"]
	reason:=m["Reason"]
	fromchannel:=m["Channel"]
	campaignid:=camparr[callee]
	callarr2[callee+":"+uid]=ringcardid
	status:=""
	if agent != nil && agent!=""{
		uidarr[agent]=uid
	}
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
	plog("Originate result: "+callee+", "+uid+", "+reason+", "+status+", "+agent)
	flashdata(campaignid)
	if(m["Response"]=="Success"){
		idarr[uid]=callee
	}else{
		fail_cnt++
		fail_cntarr[campaignid]++
		ratio_up(campaignid)

	}
	failtext:="ejsvar"
	if(reason==1 || reason==8){
		failtext="trasigt"
	}
	db_set_num_status(campaignid,ringcardid,failtext,callee)
	plog("To DB ringcard id:"+ringcardid+" number:"+callee+" Campaign ID:"+campaignid)
	delete( idarr,uid)
	if(fromchannel[0:10]!="Local/8800"){
		num_queue[campaignid]--;
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
				break;
			}
		}
		//co the bo call arr
		//delete(callarr,callee+":"+ringcardid)

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
			num_queue[campaignid]+=ratio
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
func ast_rec_stop(agent string, filename string){
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
				//cmdstr := "/home/trumpen/bin/mp3convertaccept.sh "+filename+" "+val["campaignid"]+" "+val["clientid"]+" "+val["acceptfile"]+" "+cutstart+" "+cutlen+" "+ringcard_id
				//exec.Command(cmdstr).Start()
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
//Auto call to Dest and join to default room 8000000.
func ast_dial(dest string,ringcardid string,campaignid string ){
	conf:="8000000"
	trunkname:="siptrunk"
	callerid:="0852503080"
	if val, ok := trunk_list[campaignid]; ok {
		if(len(val)>5){
			trunkname="tr"+val
			callerid=val
		}
		plog("Trunk: "+ trunkname + " ,CLID: "+callerid)
	}
	dial_cnt++;
	dial_cntarr[campaignid]++;
	//callarr[dest+":"+ringcardid]=ringcardid
	camparr[dest]=campaignid
	plog("Dial "+conf+" to :"+dest+" Ringcard: "+ringcardid+" Campaign: "+campaignid )
	result, err := a.Action(map[string]string{"Action": "Originate",
		"Channel": 	"Local/"+dest+"@selecttrunk",
		"Context": 	"default",
		"Exten":	"8000000",
		"Timeout":	dial_timeout,
		"Callerid":	callerid,
		"Account": 	campaignid+":"+dest,
		"Async":	"1",
		"ActionID":	dest+":"+ringcardid,
		"Variable":	"__TRUNKNAME="+trunkname,
		"Priority":	"1"	})
	fmt.Println(result, err)
}
func ast_mdial_trunk(agent string,ext string,dest string,ringcardid string){
	ext=agents[agent]
	conf:="8800"+ext
	//callarr[dest+":"+ringcardid]=ringcardid
	mdialarr[conf]["ringcardid"]=ringcardid
	mdialarr[conf]["dest"]=dest
	plog("Mdial "+conf+" to :"+dest+" Ringcard: "+ringcardid)
	campaignid:=agents[agent]["campaignid"]
	trunkname:="tr0852503080"
	callerid:="0757575998"
	if val, ok := trunk_list[campaignid]; ok {
		if(len(val)>5){
			trunkname="tr"+val
			callerid=val
		}
		plog("Mdial Trunk: "+ trunkname + " ,CLID: "+callerid)
	}
	dial_cnt++;
	dial_cntarr[campaignid]++;
	//callarr[dest+":"+ringcardid]=ringcardid
	camparr[dest]=campaignid
	plog("Mdial "+conf+" to :"+dest+" Ringcard: "+ringcardid+" Campaign: "+campaignid )
	usernum:=agents[agent]["usernum"]
	unmute(conf,usernum)
	result, err := a.Action(map[string]string{"Action": "Originate",
		"Channel": 	"Local/"+conf,
		"Context": 	"mdialt",
		"Exten":	"SIP/"+dest,
		"Timeout":	dial_timeout,
		"Callerid":	callerid,
		"Account": 	campaignid+":"+dest,
		"Async":	"1",
		"ActionID":	dest+":"+ringcardid+":"+agent,
		"Variable":	"__myactionid="+dest+":"+ringcardid+":"+agent+":"+campaignid,
		"Variable":	"__TRUNKNAME="+trunkname,
		"Priority":	"1"	})
	fmt.Println(result, err)
}
func ast_hangup(agent string){
	plog("Hangup call for agent: "+agent)
	if val, ok := agents[agent]; ok {
		if (val["status"] == "incall"){
			ext := val["ext"]
			conf := val["conf_num"]
			channel := val["channel"]
			campaignid := val["campaignid"]
			plog("Do hangup: "+agent+" ,"+channel+" ,"+conf )
			plog("Do hangup: Agent "+agent+" is standby")
			agents[agent]["status"]="standby"
			agents[agent]["ringcardid"]=""
			usernum:=agents[agent]["usernum"]
			mute(conf,usernum)
			db_log("standby",agent,ext,campaignid)
			agents[agent]["channel"]=""
			agents[agent]["callee"]=""
			result, err := a.Action(map[string]string{"Action": "Originate",
				"Channel":	channel,
				"Context":	"default",
				"Exten":	conf,
				"Priority":	"1",
			})
			fmt.Println(result, err)
		}
	}
}
func ast_join(m map[string]string){
	channel:=m["Channel"]
	uid:=m["Uniqueid"]
	usernum:=m["Usernum"]
	tmpclid:=idarr[uid]
	callee:=tmpclid
	ringcardid:=callarr2[tmpclid+":"+uid]
	campaignid:=camparr[tmpclid]
	//none:=1
	conf:=""
	plog("Meetme Join!, "+callee+","+channel+" "+uid+" "+m["Meetme"]+" "+usernum);
	if(m["Meetme"]=="8000000"){
		ans_cnt++
		ans_cntarr["campaignid"]++
		num_queue["campaignid"]--
		oldwhen:=time.Now().Unix()
		when:=0
		nextagent:=0
		for key, value := range agents {
			if(value["status"]=="ready" && camparr[callee]==value["campaignid"]){
				if(nextagent==0){
					nextagent=key
				}
				when=value["when"]
				if(when<oldwhen){
					nextagent=key
					oldwhen=when
				}
			}
			plog("Ast_join: Search agent for call: "+key+" , "+when+" , "+oldwhen)
		}
		if(nextagent>0){
			agent:=nextagent
			ext:=agents[agent]["ext"]
			conf=agents[agent]["conf_num"]
			plog("Found agent "+agent+", "+ext+", "+conf+", campaign id: "+campaignid)
			plog("Redirect:"+channel+", "+conf)
			agents[agent]["status"]="incall"
			agents[agent]["callee"]=callee
			agent_cnt[campaignid]--
			if(agent_cnt[campaignid]<0) {
				agent_cnt[campaignid] = 0
			}
			db_log("incall",agent,ext,campaignid)
			agents[agent]["channel"]=channel
			ratio_down(campaignid)
			usernum=agents[agent]["usernum"]
			unmute(conf,usernum)
			a.Action(map[string]string{"Action": "Redirect",
				"Channel":	channel,
				"Context":	"default",
				"Exten":	conf,
				"Priority":	"1",
			})
			plog("Ringcard: "+ringcardid+", "+callee)
			db_log_soundfile(ringcardid,campaignid,agent)
		}else{
			plog("No agent for call with ringcard: "+ringcardid)
			plog("Do hangup:"+channel+", "+conf)
			//delete(callarr,callee+":"+ringcardid)
			delete(callarr2,callee+":"+uid)
			delete(idarr,uid)
			tapp_cnt++;
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
		if _, ok := idial_cnarr[channel]; ok {
			db_inbound_delete(channel)
		}
		if(channel[4:6]!="pseudo"){
			if(channel[4:3]!=conf[4:3]){
				if(Len(mdialarr[conf]["dest"])>5){
					for key, _ := range agents {
						if(agents[key]["conf_num"]==conf){
							agents[key]["channel"]=channel
							//agents[key]["ownchannel"]=channel
							agents[key]["status"]="incall"
							agents[key]["ringcardid"]=ringcardid
							ext:=mdialarr[conf]["dest"]
							db_log("incall",key,ext,campaignid)
							db_log_soundfile(ringcardid,campaignid,key)
							delete(mdialarr[conf],ringcardid)
							delete(mdialarr[conf],"dest")
							plog("Mdial Agent: "+key+" Ringcard: "+ringcardid+" Callee:"+callee+" Campaign: "+campaignid)
							break
						}
					}
				}
			}else{
				for key, _ := range agents {
					if(agents[key]["conf_num"]==conf){
						agents[key]["ownchannel"]=channel
						agents[key]["usernum"]=usernum
						mute(conf,usernum)
						db_user_connected(key,1)
						break
					}
				}
			}
		}
	}

}
func ast_leave(m map[string]string){
	channel:=m["Channel"]
	//usernum:=m["Usernum"]
	ext:=channel[4:3]
	hit:=0
	agent:=""
	plog("Ast_leave: "+channel+", "+ext)
	for key, value := range agents {
		if(value["ext"]==ext){
			agent=key
			hit=1
			status:=value["status"]
			current_channel:=value["channel"]
			db_user_connected(agent,0)
			agents[agent]["status"]="standby"
			if(status=="incall") {
				a.Action(map[string]string{"Action": "Hangup",
					"Channel":        current_channel,
					"Context":        "default",
					"Priority":        "1",
				})
			}
		}
	}
	//thua co the bo di
	if(hit==0){
		for key, value := range agents {
		 	if(value["channel"]==ext){
				conf:=value["conf"]
				usernum:=value["usernum"]
				mute(conf,usernum)
				agents[key]["status"]="standby"
				break
			}
		}
	}
}
func ast_setstate(campaignid string, tdb_ratio int,tratio_up float,tratio_down float, ttimeout int,tcampno string){
	if(tratio_down>-2.0 && tratio_down<2.0){
		ratio_down[campaignid]=tratio_down
		plog("Set ratio down ="+tratio_down+" for campaign:"+campaignid)
	}
	if(tratio_up>-2.0 && tratio_up<2.0){
		ratio_up[campaignid]=tratio_up
		plog("Set ratio up ="+tratio_up+" for campaign:"+campaignid)
	}
	if(ttimeout>10000 && ttimeout <90000){
		dial_timeout=ttimeout
		plog("Set dial timeout ="+ttimeout)
	}
	if(tdb_ratio>1 && tdb_ratio<10){
		db_ratio[campaignid]=tdb_ratio
		plog("Set ratio="+tdb_ratio+" for campaign:"+campaignid)
	}
	trunk_list[campaignid]=tcampno
	plog("Campaign: "+campaignid+" Trunk number:"+tcampno)
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
func ast_stepratio(tratio_up float,tratio_down float, campaignid string){
	if(tratio_down>-2.0 && tratio_down<2.0){
		ratio_down[campaignid]=tratio_down
		plog("Set ratio down ="+tratio_down+" for campaign:"+campaignid)
	}
	if(tratio_up>-2.0 && tratio_up<2.0){
		ratio_up[campaignid]=tratio_up
		plog("Set ratio up ="+tratio_up+" for campaign:"+campaignid)
	}
}
//End of number
func ast_eon(campaignid string){
	num_queue[campaignid]--;
	if(num_queue[campaignid]<0){
		num_queue[campaignid]=0
	}
	plog( "decrease num_queue for campaign "+campaignid);
}

func flashdata(campaignid string){
	os.Mkdir("/var/www/flash/"+campaignid)
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