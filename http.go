package main
import (
	"fmt"
	"net/http"
)

var listfile map[string]string
func state_check(w http.ResponseWriter, r *http.Request){
//fmt.Println("http request")
	r.ParseForm()
	switch action:=r.FormValue("action"); action{
		//call to agent ext then join it to room
		case "login":
			//http://dialern.televinken.se/user_state?agent=4711&ext=021&campaignID=5&action=login
			if ((r.FormValue("agent")=="") || (r.FormValue("ext")== "") || (r.FormValue("campaignid")== "")) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to login Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("ext")+" CampaignID:"+r.FormValue("campaignid"))
				plog ("Missing argument to login Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("ext")+" CampaignID:"+r.FormValue("campaignid"))
			} else {

				plog ("HTTP login Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("ext")+" CampaignID:"+r.FormValue("campaignid"))

				code,message:=ast_login(r.FormValue("agent"),r.FormValue("ext"),r.FormValue("campaignid"))
				w.WriteHeader(code)
				fmt.Fprintf(w, message)
			}
		/*//change campaign
		case "chcamp":
			//http://dialern.televinken.se/user_state?agent=4711&ext=021&campaignid=5&action=chcamp
			if ((r.FormValue("agent")=="") || (r.FormValue("ext")== "") || (r.FormValue("campaignid")== "")) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to chcamp Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("ext")+" CampaignID:"+r.FormValue("campaignid"))
				plog ("Missing argument to chcamp Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("ext")+" CampaignID:"+r.FormValue("campaignid"))
			} else {
				w.WriteHeader(http.StatusOK)
				plog ("http chcamp Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("ext")+" CampaignID:"+r.FormValue("campaignid"));
				ast_chcamp(r.FormValue("agent"),r.FormValue("ext"),r.FormValue("campaignid"))
				//$poe_kernel->post( 'monitor', 'ast_chcamp', $agent, $anknytning, $kampanjid);
			}

		//call agent mobile phone number then join it to room
		case "loginremote":
			//http://dialern.televinken.se/user_state?agent=4711&ext=021&campaignID=5&action=loginremote
			if ((r.FormValue("agent")=="") || (r.FormValue("ext")== "") || (r.FormValue("campaignid")== "")|| (r.FormValue("dest")== "")) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to Login remote Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("ext")+" CampaignID:"+r.FormValue("campaignid"))
				plog ("Missing argument to Login remote Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("ext")+" CampaignID:"+r.FormValue("campaignid"))
			} else {
				w.WriteHeader(http.StatusOK)
				ast_login_remote(r.FormValue("agent"),r.FormValue("ext"),r.FormValue("campaignid"),r.FormValue("dest"))
			}
		/*case "dial":
			//http://dialern.televinken.se/user_state?agent=4711&ext=021&ringcardid=5&action=dial
			if ((r.FormValue("agent")=="") || (r.FormValue("ext")== "") || (r.FormValue("ringcardid")== "")|| (r.FormValue("dest")== "")) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to manual dial Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("ext")+" RingcardID:"+r.FormValue("ringcardid"))
				plog ("Missing argument to manual dial Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("ext")+" RingcardID:"+r.FormValue("ringcardid"))
			} else {
				w.WriteHeader(http.StatusOK)
				ast_mdial(r.FormValue("agent"),r.FormValue("ext"),r.FormValue("ringcardid"),r.FormValue("dest"))
			}
		case "tdial":
			//http://dialern.televinken.se/user_state?agent=4711&ext=021&ringcardid=5&action=tdial
			if ((r.FormValue("agent")=="") || (r.FormValue("ext")== "") || (r.FormValue("ringcardid")== "")|| (r.FormValue("dest")== "")) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to manual trunk dial Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("ext")+" RingcardID:"+r.FormValue("ringcardid"))
				plog ("Missing argument to manual trunk dial Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("ext")+" RingcardID:"+r.FormValue("ringcardid"))
			} else {
				w.WriteHeader(http.StatusOK)
				ast_mdial_trunk(r.FormValue("agent"),r.FormValue("ext"),r.FormValue("ringcardid"),r.FormValue("dest"))
			}
		case "standby":
			//http://dialern.televinken.se/user_state?agent=4711&action=standby
			if ((r.FormValue("agent")=="") ) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to standby Agent:"+ r.FormValue("agent"))
				plog ("Missing argument to standby Agent:"+ r.FormValue("agent"))
			} else {
				w.WriteHeader(http.StatusOK)
				ast_standby(r.FormValue("agent"))
			}
		case "ready":
			//http://dialern.televinken.se/user_state?agent=4711&action=ready
			if ((r.FormValue("agent")=="") ) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to ready Agent:"+ r.FormValue("agent"))
				plog ("Missing argument to ready Agent:"+ r.FormValue("agent"))
			} else {
				w.WriteHeader(http.StatusOK)
				ast_ready(r.FormValue("agent"))
			}
		case "rec_start":
			//http://dialern.televinken.se/user_state?agent=4711&clientid&recname=fghdfg&action=rec_start
			if ((r.FormValue("agent")=="") || Len(r.FormValue("recname"))<2) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to start record file sound Agent:"+ r.FormValue("agent")+" Record file name:"+r.FormValue("recname"))
				plog ("Missing argument to start record file sound Agent:"+ r.FormValue("agent")+" Record file name:"+r.FormValue("recname"))
			} else {
				w.WriteHeader(http.StatusOK)
				listfile[r.FormValue("agent")]=recname
				ast_rec_start(r.FormValue("agent"),r.FormValue("recname"),r.FormValue("clientid"))
			}
		case "rec_stop":
			//http://dialern.televinken.se/user_state?agent=4711&recname=fghdfg&action=rec_stop
			var recname string
			if(Len(r.FormValue("agent"))<2) {
				recname = listfile[r.FormValue("agent")]
			}
			if ((r.FormValue("agent")=="") || Len(recname)<2) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to stop record file  Agent:"+ r.FormValue("agent")+" Record file name:"+recname)
				plog ("Missing argument to stop record file sound Agent:"+ r.FormValue("agent")+" Record file name:"+recname)
			} else {
				w.WriteHeader(http.StatusOK)
				listfile[r.FormValue("agent")]=recname
				ast_rec_stop(r.FormValue("agent"),r.FormValue("recname"))
			}
			delete(listfile,r.FormValue("agent"))
		case "logout":
			//http://dialern.televinken.se/user_state?agent=4711&action=logout
			if ((r.FormValue("agent")=="") ) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to logout Agent:"+ r.FormValue("agent"))
				plog ("Missing argument to logout Agent:"+ r.FormValue("agent"))
			} else {
				w.WriteHeader(http.StatusOK)
				ast_logout(r.FormValue("agent"))
			}
		case "hangup":
			//http://dialern.televinken.se/user_state?agent=4711&action=hangup
			if ((r.FormValue("agent")=="") ) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to hangup Agent:"+ r.FormValue("agent"))
				plog ("Missing argument to hangup Agent:"+ r.FormValue("agent"))
			} else {
				w.WriteHeader(http.StatusOK)
				ast_hangup(r.FormValue("agent"))
			}
		case "setratio":
			//http://dialern.televinken.se/user_state?agent=4711&ratio=1&campaignid=234&timeout=12300&action=setratio
			if ((r.FormValue("agent")=="4711")||(r.FormValue("ratio")== "") || (r.FormValue("campaignid")== "") ) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to set ratio Ratio:"+ r.FormValue("ratio")+"Campaign ID:"+ r.FormValue("campaignid"))
				plog ("Missing argument to set ratio Agent:"+ r.FormValue("agent"))
			} else {
				w.WriteHeader(http.StatusOK)
				ast_ratio(r.FormValue("ratio"),r.FormValue("campaignid"),r.FormValue("timeout"))
			}
		case "ratiostep":
			//http://dialern.televinken.se/user_state?agent=4711&ratioup=1&rationdown&campaignid=234&timeout=12300&action=ratiostep
			if ((r.FormValue("agent")=="4711")||(r.FormValue("ratioup")== "") || (r.FormValue("ratiodown")== "")|| (r.FormValue("campaignid")== "") ) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to set step ratio Ratio Up:"+ r.FormValue("ratioup")+"Ratio Down:"+ r.FormValue("ratiodown")+"Campaign ID:"+ r.FormValue("campaignid"))
				plog ("MMissing argument to set step ratio Ratio Up:"+ r.FormValue("ratioup")+"Ratio Down:"+ r.FormValue("ratiodown")+"Campaign ID:"+ r.FormValue("campaignid"))
			} else {
				w.WriteHeader(http.StatusOK)
				ast_ratio(r.FormValue("ratio"),r.FormValue("campaignid"),r.FormValue("timeout"))
			}
		//listen specified call
		case "idial":
			//http://dialern.televinken.se/user_state?agent=4711&ext=021&ringcardid=5&channel=SIP/123&action=idial
			if ((r.FormValue("agent")=="") || (r.FormValue("ext")== "") || (r.FormValue("ringcardid")== "")|| (r.FormValue("dest")== "")|| (r.FormValue("channel")== "")) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Missing argument to monitor dial Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("ext")+" RingcardID:"+r.FormValue("ringcardid")+" Channel:"+r.FormValue("channel"))
				plog ("Missing argument to monitor dial Agent:"+ r.FormValue("agent")+" Ext:"+r.FormValue("ext")+" RingcardID:"+r.FormValue("ringcardid")+" Channel:"+r.FormValue("channel"))
			} else {
				w.WriteHeader(http.StatusOK)
				ast_idial(r.FormValue("agent"),r.FormValue("ext"),r.FormValue("dest"),r.FormValue("ringcardid"),r.FormValue("channel"))
			}*/
		default:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, r.FormValue("action")+ " is not an allowed action" )


	}

}