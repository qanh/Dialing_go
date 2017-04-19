package main
import (
	"fmt"
	"net/http"
	"strconv"
)

var listfile map[string]string
func state_check(w http.ResponseWriter, r *http.Request){
//fmt.Println("http request")
	r.ParseForm()
	switch action:=r.FormValue("action"); action{
		//call to agent anknytning then join it to room
		case "login":
			//http://dialern.televinken.se/user_state?agent=4711&anknytning=021&campaignID=5&action=login
			if ((r.FormValue("agent")=="") || (r.FormValue("anknytning")== "") || (r.FormValue("campaignID")== "")) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to login Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" CampaignID:"+r.FormValue("campaignID"))
				plog ("Missing argument to login Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" CampaignID:"+r.FormValue("campaignID"),1)
			} else {

				plog ("HTTP login Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" CampaignID:"+r.FormValue("campaignID"),1)

				code,message:=ast_login(r.FormValue("agent"),r.FormValue("anknytning"),r.FormValue("campaignID"))
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
				w.WriteHeader(http.StatusOK)
				plog ("http chcamp Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" CampaignID:"+r.FormValue("campaignID"),1);
				code,message:=ast_chcamp(r.FormValue("agent"),r.FormValue("campaignID"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
				//$poe_kernel->post( 'monitor', 'ast_chcamp', $agent, $anknytning, $kampanjid);
			}

		//call agent mobile phone number then join it to room
		case "loginremote":
			//http://dialern.televinken.se/user_state?agent=4711&anknytning=021&campaignID=5&action=loginremote
			if ((r.FormValue("agent")=="") || (r.FormValue("anknytning")== "") || (r.FormValue("campaignID")== "")|| (r.FormValue("dest")== "")) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to Login remote Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" CampaignID:"+r.FormValue("campaignID"))
				plog ("Missing argument to Login remote Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" CampaignID:"+r.FormValue("campaignID"),1)
			} else {
				w.WriteHeader(http.StatusOK)
				code,message:=ast_login_remote(r.FormValue("agent"),r.FormValue("anknytning"),r.FormValue("campaignID"),r.FormValue("dest"))
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
				w.WriteHeader(http.StatusOK)
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
				w.WriteHeader(http.StatusOK)
				code,message:=ast_mdial_trunk(r.FormValue("agent"),r.FormValue("anknytning"),r.FormValue("ringkort"),r.FormValue("dest"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		case "standby":
			//http://dialern.televinken.se/user_state?agent=4711&action=standby
			if ((r.FormValue("agent")=="") ) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to standby Agent:"+ r.FormValue("agent"))
				plog ("Missing argument to standby Agent:"+ r.FormValue("agent"))
			} else {
				w.WriteHeader(http.StatusOK)
				code,message:=ast_standby(r.FormValue("agent"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		case "ready":
			//http://dialern.televinken.se/user_state?agent=4711&action=ready
			if ((r.FormValue("agent")=="") ) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to ready Agent:"+ r.FormValue("agent"))
				plog ("Missing argument to ready Agent:"+ r.FormValue("agent"))
			} else {
				w.WriteHeader(http.StatusOK)
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
				w.WriteHeader(http.StatusOK)
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
				w.WriteHeader(http.StatusOK)
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
				plog ("Missing argument to hangup Agent:"+ r.FormValue("agent"))
			} else {
				w.WriteHeader(http.StatusOK)
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
				w.WriteHeader(http.StatusOK)
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
				w.WriteHeader(http.StatusOK)
				http_ratio,_:=strconv.Atoi(r.FormValue("ratio"))
				http_timout,_:=strconv.Atoi(r.FormValue("timeout"))
				code,message:=ast_ratio(http_ratio,r.FormValue("campaignID"),http_timout)
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		//listen specified call
		/*case "idial":
			//http://dialern.televinken.se/user_state?agent=4711&anknytning=021&ringkort=5&channel=SIP/123&action=idial
			if ((r.FormValue("agent")=="") || (r.FormValue("anknytning")== "") || (r.FormValue("ringkort")== "")|| (r.FormValue("dest")== "")|| (r.FormValue("channel")== "")) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to monitor dial Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" RingcardID:"+r.FormValue("ringkort")+" Channel:"+r.FormValue("channel"))
				plog ("Missing argument to monitor dial Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("anknytning")+" RingcardID:"+r.FormValue("ringkort")+" Channel:"+r.FormValue("channel"))
			} else {
				w.WriteHeader(http.StatusOK)
				code,message:=ast_idial(r.FormValue("agent"),r.FormValue("anknytning"),r.FormValue("dest"),r.FormValue("ringkort"),r.FormValue("channel"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}*/
		default:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, r.FormValue("action")+ " is not an allowed action" )


	}

}