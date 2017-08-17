package main
import (
	"fmt"
	"net/http"
	"strconv"
	"io/ioutil"
	"bytes"
)

var listfile map[string]string
func state_check(w http.ResponseWriter, r *http.Request){
//fmt.Println("http request")
	r.ParseForm()
	//fmt.Println(r)
	switch action:=r.FormValue("action"); action{
		//call to agent anknytning then join it to room
		case "login":
			//http://dialern.televinken.se/user_state?agent=4711&anknytning=021&campaignID=5&action=login
			if ((r.FormValue("agent")=="") || (r.FormValue("anknytning")== "") || (r.FormValue("campaignID")== "") || (r.FormValue("clientid")== "")) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to login Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" CampaignID:"+r.FormValue("campaignID")+" ClientID:"+r.FormValue("clientid"))
				plog ("Missing argument to login Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" CampaignID:"+r.FormValue("campaignID")+" ClientID:"+r.FormValue("clientid"),1)
			} else {
				inbound:=""
				if val, ok := r.FormValue("inbound"); ok {
					inbound	=val
				}
				plog ("HTTP login Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" CampaignID:"+r.FormValue("campaignID"),1)

				code,message:=ast_login(r.FormValue("agent"),r.FormValue("anknytning"),r.FormValue("campaignID"),r.FormValue("clientid"),inbound)
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		//change campaign
		case "chcamp":
			//http://dialern.televinken.se/user_state?agent=4711&anknytning=021&campaignID=5&action=chcamp
			if ((r.FormValue("agent")=="") || (r.FormValue("anknytning")== "") || (r.FormValue("campaignID")== "")) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to chcamp Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" CampaignID:"+r.FormValue("campaignID"))
				plog ("Missing argument to chcamp Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" CampaignID:"+r.FormValue("campaignID"),1)
			} else {
				
				plog ("http chcamp Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" CampaignID:"+r.FormValue("campaignID"),1)
				inbound:=""
				if val, ok := r.FormValue("inbound"); ok {
					inbound	=val
				}
				code,message:=ast_chcamp(r.FormValue("agent"),r.FormValue("campaignID"),inbound)
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
				//$poe_kernel->post( 'monitor', 'ast_chcamp', $agent, $anknytning, $kampanjid)
			}

		//call agent mobile phone number then join it to room
		case "loginremote":
			//http://dialern.televinken.se/user_state?agent=4711&anknytning=021&campaignID=5&action=loginremote
			if ((r.FormValue("agent")=="") || (r.FormValue("anknytning")== "") || (r.FormValue("campaignID")== "")|| (r.FormValue("dest")== "") || (r.FormValue("clientid")== "")) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to Login remote Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" CampaignID:"+r.FormValue("campaignID")+" ClientID:"+r.FormValue("clientid"))
				plog ("Missing argument to Login remote Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" CampaignID:"+r.FormValue("campaignID")+" ClientID:"+r.FormValue("clientid"),1)
			} else {
				inbound:=""
				if val, ok := r.FormValue("inbound"); ok {
					inbound	=val
				}
				code,message:=ast_login_remote(r.FormValue("agent"),r.FormValue("anknytning"),r.FormValue("campaignID"),r.FormValue("dest"),r.FormValue("clientid"),inbound)
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		/*case "dial":
			//http://dialern.televinken.se/user_state?agent=4711&anknytning=021&ringkort=5&action=dial
			if ((r.FormValue("agent")=="") || (r.FormValue("anknytning")== "") || (r.FormValue("ringkort")== "")|| (r.FormValue("dest")== "")) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to manual dial Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" RingcardID:"+r.FormValue("ringkort"))
				plog ("Missing argument to manual dial Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" RingcardID:"+r.FormValue("ringkort"))
			} else {
				
				code,message:=ast_dial(r.FormValue("agent"),r.FormValue("anknytning"),r.FormValue("ringkort"),r.FormValue("dest"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}*/
		case "tdial":
			//http://dialern.televinken.se/user_state?agent=4711&anknytning=021&ringkort=5&action=tdial
			if ((r.FormValue("agent")=="") || (r.FormValue("anknytning")== "") || (r.FormValue("ringkort")== "")|| (r.FormValue("dest")== "")) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to manual trunk dial Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" RingcardID:"+r.FormValue("ringkort"))
				plog ("Missing argument to manual trunk dial Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" RingcardID:"+r.FormValue("ringkort"),1)
			} else {
				
				code,message:=ast_mdial_trunk(r.FormValue("agent"),r.FormValue("anknytning"),r.FormValue("ringkort"),r.FormValue("dest"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		case "standby":
			//http://dialern.televinken.se/user_state?agent=4711&action=standby
			if ((r.FormValue("agent")=="") ) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to standby Agent:"+ r.FormValue("agent"))
				plog ("Missing argument to standby Agent:"+ r.FormValue("agent"),1)
			} else {
				
				code,message:=ast_standby(r.FormValue("agent"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		case "ready":
			//http://dialern.televinken.se/user_state?agent=4711&action=ready
			if ((r.FormValue("agent")=="") ) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to ready Agent:"+ r.FormValue("agent"))
				plog ("Missing argument to ready Agent:"+ r.FormValue("agent"),1)
			} else {
				
				code,message:=ast_ready(r.FormValue("agent"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		case "rec_start":
			//http://dialern.televinken.se/user_state?agent=4711&clientid&recname=fghdfg&action=rec_start
			if ((r.FormValue("agent")=="") || len(r.FormValue("recname"))<2) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to start record file sound Agent:"+ r.FormValue("agent")+" Record file name:"+r.FormValue("recname"))
				plog ("Missing argument to start record file sound Agent:"+ r.FormValue("agent")+" Record file name:"+r.FormValue("recname"),1)
			} else {
				
				listfile[r.FormValue("agent")]=r.FormValue("recname")
				code,message:=ast_rec_start(r.FormValue("agent"),r.FormValue("recname"),r.FormValue("clientid"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		case "rec_stop":
			//http://dialern.televinken.se/user_state?agent=4711&recname=fghdfg&action=rec_stop
			var recname string
			if(len(r.FormValue("agent"))<2) {
				recname = listfile[r.FormValue("agent")]
			}
			if ((r.FormValue("agent")=="") || len(recname)<2) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to stop record file  Agent:"+ r.FormValue("agent")+" Record file name:"+recname)
				plog ("Missing argument to stop record file sound Agent:"+ r.FormValue("agent")+" Record file name:"+recname,1)
			} else {
				
				listfile[r.FormValue("agent")]=recname
				code,message:=ast_rec_stop(r.FormValue("agent"),r.FormValue("recname"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
			delete(listfile,r.FormValue("agent"))
		case "logout":
			//http://dialern.televinken.se/user_state?agent=4711&action=logout
			if ((r.FormValue("agent")=="") ) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to logout Agent:"+ r.FormValue("agent"))
				plog ("Missing argument to logout Agent:"+ r.FormValue("agent"),1)
			} else {
				code,message:=ast_logout(r.FormValue("agent"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		case "hangup":
			//http://dialern.televinken.se/user_state?agent=4711&action=hangup
			if ((r.FormValue("agent")=="") ) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to hangup Agent:"+ r.FormValue("agent"))
				plog ("Missing argument to hangup Agent:"+ r.FormValue("agent"),1)
			} else {
				
				code,message:=ast_hangup(r.FormValue("agent"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		case "setratio":
			//http://dialern.televinken.se/user_state?agent=4711&ratio=1&campaignID=234&timeout=12300&action=setratio
			if ((r.FormValue("agent")=="4711")||(r.FormValue("ratio")== "") || (r.FormValue("campaignID")== "") ) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to set ratio Ratio:"+ r.FormValue("ratio")+"Campaign ID:"+ r.FormValue("campaignID"))
				plog ("Missing argument to set ratio Agent:"+ r.FormValue("agent"),1)
			} else {
				
				http_ratio,_:=strconv.Atoi(r.FormValue("ratio"))
				http_timout,_:=strconv.Atoi(r.FormValue("timeout"))
				code,message:=ast_ratio(http_ratio,r.FormValue("campaignID"),http_timout)
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		case "ratiostep":
			//http://dialern.televinken.se/user_state?agent=4711&rup=1&rationdown&campaignID=234&timeout=12300&action=ratiostep
			if ((r.FormValue("agent")=="4711")||(r.FormValue("rup")== "") || (r.FormValue("rner")== "")|| (r.FormValue("campaignID")== "") ) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to set step ratio Ratio Up:"+ r.FormValue("rup")+"Ratio Down:"+ r.FormValue("rner")+"Campaign ID:"+ r.FormValue("campaignID"))
				plog ("MMissing argument to set step ratio Ratio Up:"+ r.FormValue("rup")+"Ratio Down:"+ r.FormValue("rner")+"Campaign ID:"+ r.FormValue("campaignID"),1)
			} else {

				rup,_:=strconv.ParseFloat(r.FormValue("rup"),  64)
				rner,_:=strconv.ParseFloat(r.FormValue("rner"), 64)
				code,message:=ast_stepratio(rup,rner,r.FormValue("campaignID"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		//listen specified call
		case "idial":
			//http://dialern.televinken.se/user_state?agent=4711&anknytning=021&ringkort=5&channel=SIP/123&action=idial
			if ((r.FormValue("agent")=="") || (r.FormValue("anknytning")== "") || (r.FormValue("ringkort")== "")|| (r.FormValue("dest")== "")|| (r.FormValue("channel")== "")) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to monitor dial Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" RingcardID:"+r.FormValue("ringkort")+" Channel:"+r.FormValue("channel"))
				plog ("Missing argument to monitor dial Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" RingcardID:"+r.FormValue("ringkort")+" Channel:"+r.FormValue("channel"),1)
			} else {
				
				code,message:=ast_idial(r.FormValue("agent"),r.FormValue("anknytning"),r.FormValue("dest"),r.FormValue("ringkort"),r.FormValue("channel"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		case "copyFile":
			if(r.FormValue("fileID")==""){
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w,"Missing argument to download file ")
				plog ("Missing argument to download file ", 1)
			}else{
				code,message:=db_get_file(r.FormValue("fileID"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		case "transfer":
			if(r.FormValue("phonenumber")=="" ||r.FormValue("agent")=="" || r.FormValue("to_agent")==""){
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w,"Missing argument to transfer call: "+r.FormValue("agent")+","+r.FormValue("to_agent")+","+r.FormValue("phonenumber"))
				plog ("Missing argument to transfer call: "+r.FormValue("agent")+","+r.FormValue("to_agent")+","+r.FormValue("phonenumber"), 1)
			}else{
				code,message:=ast_transfer(r.FormValue("agent"),r.FormValue("to_agent"),r.FormValue("phonenumber"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		case "record":
			if(r.FormValue("phonenumber")=="" ||r.FormValue("recname")=="" || r.FormValue("trunk")==""){
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w,"Missing argument to record voice promts: "+r.FormValue("phonenumber")+","+r.FormValue("recname")+","+r.FormValue("trunk"))
				plog ("Missing argument to record voice promts: "+r.FormValue("phonenumber")+","+r.FormValue("recname")+","+r.FormValue("trunk"), 1)
			}else{
				code,message:=ast_record(r.FormValue("phonenumber"),r.FormValue("recname"),r.FormValue("trunk"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		case "recordcancel":
			if(r.FormValue("phonenumber")=="" ||r.FormValue("recname")=="" ){
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w,"Missing argument to stop record voice promts: "+r.FormValue("phonenumber")+","+r.FormValue("recname"))
				plog ("Missing argument to stop record voice promts: "+r.FormValue("phonenumber")+","+r.FormValue("recname"), 1)
			}else{
				code,message:=ast_record_stop(r.FormValue("phonenumber"),r.FormValue("recname"),1)
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		case "recordfinish":
			if(r.FormValue("phonenumber")=="" ||r.FormValue("recname")=="" ){
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w,"Missing argument to stop record voice promts: "+r.FormValue("phonenumber")+","+r.FormValue("recname"))
				plog ("Missing argument to stop record voice promts: "+r.FormValue("phonenumber")+","+r.FormValue("recname"), 1)
			}else{
				code,message:=ast_record_stop(r.FormValue("phonenumber"),r.FormValue("recname"),0)
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		case "recordget":
			if(r.FormValue("recname")==""){
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w,"Missing argument to get record voice promts: "+r.FormValue("recname"))
				plog ("Missing argument to get record voice promts: "+r.FormValue("recname"), 1)
			}else{
				file:="/var/lib/asterisk/sounds/dialplan/"+r.FormValue("recname")+".wav"
				streamPDFbytes, err := ioutil.ReadFile(file)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w,err.Error())

				}else{
					b := bytes.NewBuffer(streamPDFbytes)

					// stream straight to client(browser)
					w.Header().Set("Content-type", "application/octet-stream")
					w.WriteHeader(200)
					if _, err := b.WriteTo(w); err != nil { // <----- here!
						fmt.Fprintf(w, "%s", err.Error())
					}
				}

			}
		//case "peerstatus":
			//code,message:=ast_peerstatus(r.FormValue("peer"))
			//w.WriteHeader(code)
			//fmt.Fprintf(w, message)
		case "peerdelete":
			code,message:=ast_delete_peercache()
			w.WriteHeader(code)
			fmt.Fprintf(w, message)
		default:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, r.FormValue("action")+ " is not an allowed action" )


	}

}