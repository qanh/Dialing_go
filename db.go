package main
import(
	"database/sql"
	//_ "github.com/go-sql-driver/mysql"
	//"github.com/ziutek/mymysql/mysql"
	//_ "github.com/ziutek/mymysql/native"
	"strconv"
 	"time"
	"os"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
)
func db_getstate(campaignid string){

	var t_ratio,t_ratio_up,t_ratio_down,t_wait_time,t_campNumber sql.NullString
	var ratioup,ratiodown,ratio float64
	var wait_time int
	err := db.QueryRow("select d.dmaxratio as ratio,d.dtry_time as wait_time ,d.dratio_up as ratio_up,d.dratio_down as ratio_down,t.campNumber from tCampaign t left join DialerSetting d on t.Dialer_Setting=d.dID where t.campaignID= ?",campaignid).Scan(&t_ratio,&t_wait_time,&t_ratio_up,&t_ratio_down,&t_campNumber)
	//if no row -> err !=nil
	checkErr(err)
	if(err == nil){
		//defer rows.Close()
		//rows, _ = stmt.Run(campaignid)
		//rows.Scan(&t_ratio,&t_wait_time,&t_ratio_up,&t_ratio_down,&t_campNumber)
		//if(t_ratio==nil){
		//	set_default_ratio(campaignid)
		//}else{
		if(t_ratio_up.Valid ) {
			ratioup, _ = strconv.ParseFloat(t_ratio_up.String, 64)
		}
		if(t_ratio_down.Valid ) {
			ratiodown, _ = strconv.ParseFloat(t_ratio_down.String, 64)
		}
		if(t_ratio.Valid ) {
			ratio, _ = strconv.ParseFloat(t_ratio.String, 64)
		}
		if(t_wait_time.Valid ) {
			wait_time, _ = strconv.Atoi(t_wait_time.String)
		}
		plog("t_ratio:"+t_ratio.String,1)
		if (ratioup > -2 && ratioup < 2) {
			ratio_up[campaignid] = ratioup
			plog("Set ration up =" + t_ratio_up.String + " for campaign " + campaignid,1)
		}
		if (ratiodown > -2 && ratiodown < 2) {
			ratio_down[campaignid] = ratiodown
			plog("Set ration down = " + t_ratio_down.String + " for campaign " + campaignid,1)
		}
		if (wait_time > 10000 && wait_time < 90000) {
			dial_timeout = wait_time
			plog("Set dial timeout = " + t_wait_time.String + " for campaign " + campaignid,1)
		}
		if (ratio > 1 && ratio < 10) {
			db_ratio[campaignid] = ratio
			plog("Set ratio = " + t_ratio.String + " for campaign " + campaignid,1)
		}

		trunk_list[campaignid] = t_campNumber.String
		plog("Set trunk = " + t_campNumber.String + " for campaign " + campaignid,1)
		//}
	}else{
		ast_set_default_ratio(campaignid)
	}

}
func db_log(status string, agent string, ext string, campaignid string){
	query, err :=db.Prepare("INSERT INTO log set state = ?, agentid = ?,extension = ?,kampanj = ?, tid = NOW()")
	checkErr(err)
	defer query.Close()
	_, err=query.Exec(status,agent,ext,campaignid)
	checkErr(err)
	plog( "db_log "+status+", "+agent+", "+ext+", "+campaignid,1)
}
func db_setstate(ringcardid string){
	stmt, err := db.Prepare("UPDATE tCampRingCards SET tapp = tapp + 1 where rID=?")
	checkErr(err)
	defer stmt.Close()
	_, err = stmt.Exec( ringcardid)
	checkErr(err)
}

func db_set_num_status(campaignid string , ringcardid string,reason string, number string){
	/*type Campaign struct {
		Phone1   string
		Phone2   string
		Phone3   string
		Phone4   string
		Phone5   string
		status1 int
		status2 int
		status3 int
		status4 int
		status5 int
	}
	//var campaign Campaign

	phone:=map [string]string{
		"Phone1":"",
		"Phone2":"",
		"Phone3":"",
		"Phone4":"",
		"Phone5":"",
	}*/
	update_query := "UPDATE tCampRingCards SET "
	var real_status int
	status:=map [string]int{
		"status1":0,
		"status2":0,
		"status3":0,
		"status4":0,
		"status5":0,
	}
	select_query := "SELECT Phone1,Phone2,Phone3,Phone4,Phone5,status1,status2,status3,status4,status5 from tCampRingCards WHERE rID ="+ringcardid
	row,err := db.Query(select_query)//.Scan(phone["Phone1"],phone["Phone2"],phone["Phone3"],phone["Phone4"],phone["Phone5"],status["status1"],status["status2"],status["status3"],status["status4"],status["status5"])
	checkErr(err)
	columnNames, _ := row.Columns()
	rc := NewMapStringScan(columnNames)
	row.Next()
	rc.Update(row)
	row.Close()
	i := 1
	index:=1
	callnote_status:="Result_No_Answer"
	for i < 6 {
		key:="Phone"+strconv.Itoa(index)
		if(number==rc.row[key]){
			i=6
		}else{
			index++
		}
		status["status"+strconv.Itoa(i)],_=strconv.Atoi(rc.row["status"+strconv.Itoa(i)])
		i++
	}
	//phone_key:="Phone"+index
	status_key:="status"+strconv.Itoa(index)
	if(reason=="trasigt"){
		real_status,_=strconv.Atoi(rc.row[status_key])
		real_status+=1000
	}else if(reason=="ejsvar"){
		real_status,_=strconv.Atoi(rc.row[status_key])
		real_status+=1
	}else{
		plog ("error: set_num_status () unknow reason\n", 1)
	}
	called:=1
	fail:=1
	for i := 1; i < 6; i++ {
		if(status["status"+strconv.Itoa(i)]==0 && len (rc.row["Phone"+strconv.Itoa(i)])>4){
			if(i!=index){
				called=0
			}
		}
		if(status["status"+strconv.Itoa(i)]<500 && len (rc.row["Phone"+strconv.Itoa(i)])>4){
			if(i!=index){
				fail=0
			}else if(reason=="ejsvar"){
				fail=0
			}
		}
	}
	if(fail==1){
		plog ("set_num_status: all number fail, update database for ringcard "+ringcardid,1)
		update_query="UPDATE tCampRingCards Set fail_try=fail_try+1 WHERE rID="+ringcardid
		callnote_status="Result_Connection_Failure"
	}else if(called==1){
		i:=1
		for i < 6 {
			phonestatus:="status"+strconv.Itoa(i)
			if(status[phonestatus]<500 && len(rc.row["Phone"+strconv.Itoa(i)])>4){
				update_query = update_query+" "+phonestatus+" = 0, "
			}else{
				update_query = update_query+" "+phonestatus+" = 1000, "
			}
			i++
		}
		tidsperiod := tidsperiod ()
		if(tidsperiod==1){
			update_query += "AM_try = AM_try + 1, AMdate=CURDATE(), lastcalldate=CURDATE(), "
		}else if(tidsperiod==2){
			update_query += "PM_try = PM_try + 1, PMdate=CURDATE(), lastcalldate=CURDATE(), "
		}else if(tidsperiod==3){
			update_query += "Evening_try = Evening_try + 1, Eveningdate=CURDATE(), lastcalldate=CURDATE(), "
		}else{
			update_query += "lastcalldate=CURDATE(),"
		}
		if(real_status<1000){
			update_query +=" no_try = no_try + 1 WHERE rID="+ringcardid
		}else{
			update_query +="no_try = no_try + 1,"+status_key+"="+strconv.Itoa(real_status)+" WHERE rID="+ringcardid
		}
	}else{
		tidsperiod := tidsperiod ()
		if(tidsperiod==1){
			update_query += "AM_try = AM_try + 1, AMdate=CURDATE(), lastcalldate=CURDATE(), "
		}else if(tidsperiod==2){
			update_query += "PM_try = PM_try + 1, PMdate=CURDATE(), lastcalldate=CURDATE(), "
		}else if(tidsperiod==3){
			update_query += "Evening_try = Evening_try + 1, Eveningdate=CURDATE(), lastcalldate=CURDATE(), "
		}else{
			update_query += "lastcalldate=CURDATE(),"
		}
		if(real_status<1000){
			update_query +=" no_try = no_try + 1 WHERE rID="+ringcardid
		}else{
			update_query +="no_try = no_try + 1,"+status_key+"="+strconv.Itoa(real_status)+" WHERE rID="+ringcardid
		}
	}
	db_callnote_fail(campaignid,ringcardid,number,callnote_status,"0")
	_, err = db.Exec(update_query)
	checkErr(err)

}
func db_dial(ratio int ,campaignid string ){
	plog("db_dial",1)
	tidsperiod:=tidsperiod ()
	query:="call PTakeActiveRingCard("+campaignid+","+strconv.Itoa(tidsperiod)+")"
	if(ratio>0){
		row, err := db.Query(query)

		if(err!=nil){
			checkErr(err)
			row.Close()
			ast_eon(campaignid)
		}else{
			db_dial_res(row,campaignid)
		}
		if(ratio>1){
			ratio--
			go db_dial(ratio,campaignid)
		}
	}
}
func db_dial_res(row *sql.Rows,campaignid string ){
	plog("db_dial_res",1)
	number:=""
	number_index:=0
	number_check:=1
	columnNames, err := row.Columns()
	checkErr(err)
	rc := NewMapStringScan(columnNames)
	if(!row.Next()){
		plog("Nu är det slut på telefonnummer i den här kampanjen",1)
		ast_eon(campaignid)
	}else {
		rc.Update(row)
		row.Close()
		ringcardid := rc.row["rID"]
		status1, _ := strconv.Atoi(rc.row["status1"])
		status2, _ := strconv.Atoi(rc.row["status2"])
		status3, _ := strconv.Atoi(rc.row["status3"])
		status4, _ := strconv.Atoi(rc.row["status4"])
		status5, _ := strconv.Atoi(rc.row["status5"])
		number1 := rc.row["Phone1"]
		number2 := rc.row["Phone2"]
		number3 := rc.row["Phone3"]
		number4 := rc.row["Phone4"]
		number5 := rc.row["Phone5"]
		plog("db_dial_res" + number1 + " " + number2, 1)
		if _, ok := list_ringcard[ringcardid]; ok {
			time.Sleep(1)
			go db_dial(1, campaignid)
		} else {
			list_ringcard["ringcardid"] = 1
			if (len(number1) > 4 && status1 == 0) {
				number = number1
				number_index = 1
			} else if (len(number2) > 4 && status2 == 0) {
				number = number2
				number_index = 2
			} else if (len(number3) > 4 && status3 == 0) {
				number = number3
				number_index = 3
			} else if (len(number4) > 4 && status4 == 0) {
				number = number4
				number_index = 4
			} else if (len(number5) > 4 && status5 == 0) {
				number = number5
				number_index = 5
			} else {
				number_ok := 0
				update_query := ""
				for i := 5; i > 0; i-- {
					status, _ := strconv.Atoi(rc.row["status" + strconv.Itoa(i)])
					if (status < 500) {
						number_ok = 1
					}
				}
				if (number_ok == 1) {
					tidsperiod := tidsperiod()
					update_query = "UPDATE tCampRingCards Set lastcalldate=CURDATE()"
					if (tidsperiod == 1) {
						update_query = update_query + ", AMdate=CURDATE()"
					} else if (tidsperiod == 2) {
						update_query = update_query + ", PMdate=CURDATE()"
					} else if (tidsperiod == 3) {
						update_query = update_query + ", Eveningdate=CURDATE()"
					}
					if ((len(number1) < 5) || (status1 > 500)) {
						update_query = update_query + ", status1 = 1000 "
					} else {
						update_query = update_query + ", status1 = 0 "
						number_check = 0
						number_index = 1
						number = number1
					}
					if ((len(number2) < 5) || (status2 > 500)) {
						update_query = update_query + ", status2 = 1000 "
					} else {
						update_query = update_query + ", status2 = 0 "
						number_check = 0
						number_index = 2
						number = number2
					}
					if ((len(number3) < 5) || (status3 > 500)) {
						update_query = update_query + ", status3 = 1000 "
					} else {
						update_query = update_query + ", status3 = 0 "
						number_check = 0
						number_index = 3
						number = number3
					}
					if ((len(number4) < 5) || (status4 > 500)) {
						update_query = update_query + ", status4 = 1000 "
					} else {
						update_query = update_query + ", status4 = 0 "
						number_check = 0
						number_index = 4
						number = number4
					}
					if ((len(number5) < 5) || (status5 > 500)) {
						update_query = update_query + ", status5 = 1000 "
					} else {
						update_query = update_query + ", status5 = 0 "
						number_check = 0
						number_index = 5
						number = number5
					}
					update_query = update_query + " WHERE rID=" + ringcardid
					if (number_check == 1) {
						plog("inga nummer med bra status och längd på detta ringkort $ringkort\n", 1)
						update_query = "UPDATE tCampRingCards Set userID=0,statusID=(select bortfall_status from tCampaign where campaignID=" + campaignid + "),subID=(select inget_nr_sub from tCampaign where campaignID=" + campaignid + "), closed_date=now() WHERE rID=" + ringcardid
					}
				} else {
					plog("inget nummer har bra status på detta ringkort $ringkort\n", 1)
					update_query = "UPDATE tCampRingCards Set userID=0,statusID=(select bortfall_status from tCampaign where campaignID=" + campaignid + "),subID=(select inget_nr_sub from tCampaign where campaignID=" + campaignid + "), closed_date=now() WHERE rID=" + ringcardid
				}

				_, err = db.Exec(update_query)
				checkErr(err)
				delete(list_ringcard, ringcardid)

				if (number_check != 0) {
					go db_dial(1, campaignid)
				}
			}
			plog(strconv.Itoa(number_index) + " " + strconv.Itoa(number_check), 1)
			if (number_index != 0 || number_check == 0) {
				update_query := "UPDATE tCampRingCards SET status" + strconv.Itoa(number_index) + " =  1 where rID=" + ringcardid
				_, err = db.Exec(update_query)
				checkErr(err)
				delete(list_ringcard, ringcardid)
				ast_dial(number, ringcardid, campaignid)
			}
		}
	}
}
func db_log_soundfile(ringcardid string ,campaignid string ,agent string,clientid string){
	//var clientid string
	//select_query := "SELECT clientID from tCampaign where campaignID="+campaignid
	//err:=db.QueryRow(select_query).Scan(&clientid)
	//if(err==nil){
		host,_:=os.Hostname()
		now := time.Now()
		datestring:=now.Format("20060102_030405")
		recname := "LOGG_"+datestring+"_u"+agent+"_c"+campaignid+"_"+ringcardid+"_"+clientid+".wav"

		go ast_rec_start(agent,recname,clientid)
		//db.Exec("INSERT INTO soundfile set rid = "+ringcardid+", userid = "+agent+", campaignid = "+campaignid+", clientid = "+clientid+", originfilename = "+recname+", filename = "+recname+", closed = 0, converted = 0, cut = 0, start = NOW(),asterisk_ip="+host)
		stmtIns, err := db.Prepare("INSERT INTO soundfile set rid = ?, userid = ?, campaignid = ?, clientid = ?, originfilename = ?, filename = ?, closed = 0, converted = 0, cut = 0, start = NOW(),asterisk_ip=?")
		checkErr(err)
		defer stmtIns.Close()
		_, err=stmtIns.Exec(ringcardid,agent,campaignid,clientid,recname,recname,host)
		plog("poe_kernel->post( monitor,ast_rec_start_mix, "+agent+", "+recname+","+clientid+")", 1 )
	//}
}

func db_reg_tapp(ringcardid string){
	_,err:=db.Exec("UPDATE tCampRingCards SET tapp = tapp + 1 where rID="+ringcardid)
	checkErr(err)
}

func db_inbound_delete(channel string){
	_,err:=db.Exec("delete from popurlar where channel="+channel)
	checkErr(err)
}

func db_user_connected(userid string ,connected int){
	_,err:=db.Exec("UPDATE tUsers SET connected = "+strconv.Itoa(connected)+" where userID="+userid)
	checkErr(err)
}
func db_log_rec(campaignid string ,clientid string,logfile string,va int){
	addr := "192.168.185.99"
	//_,err:=db.Exec("INSERT INTO tSoundFiles set campaignID = "+campaignid+", clientID = "+clientid+", filename = "+logfile+", createtime=NOW(), va = "+strconv.Itoa(va)+",addr = "+addr+", filesize = 0, status=0")
	//checkErr(err)
	stmtIns, err := db.Prepare("INSERT INTO tSoundFiles set campaignID = ?, clientID = ?, filename = ?, createtime=NOW(), va = ?,addr = ?, filesize = 0, status=0")
	checkErr(err)
	defer stmtIns.Close()
	_, err=stmtIns.Exec(campaignid,clientid,logfile,strconv.Itoa(va),addr)
}
func tidsperiod() int {
	now := time.Now()
	wday:=now.Weekday()
	hour:=now.Hour()
	if(wday.String()=="Sunday" || wday.String()=="Saturday"){return 4}
	if(hour>=6 && hour < 12){return 1}
	if(hour>=12 && hour < 17){return 2}
	if(hour>=17 && hour < 23){return 3}
	return 0
}
func db_get_file(fileid string)(int,string){
	var host_name string
	var path_sql string
	select_query := "select host_name,path from dialplan_voicefiles where id=="+fileid
	err:=db.QueryRow(select_query).Scan(&host_name,&path_sql)
	if(err==nil){
		ext := filepath.Ext(path_sql)
		path_sql=strings.Replace(path_sql,ext,".wav",-1)
		filename := path.Base(path_sql)
		out, err := os.Create("/var/lib/asterisk/sounds/dialplan/"+filename)
		if err != nil  {
			return 400,err.Error()
		}
		defer out.Close()
		resp, err := http.Get("http://"+host_name+"/"+path_sql)
		if err != nil {
			return 400,err.Error()
		}
		defer resp.Body.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil  {
			return 400,err.Error()
		}
	}
	return 200,"OK"
}

func db_user_wrapup(userid string){
	_,err:=db.Exec("UPDATE tUsers SET status ='Wrapup' where userID="+userid)
	checkErr(err)
}

func db_callnote_fail(campaignid string, ringcardid string, number string,status string,userid string){
	if(userid==""){
		userid="0"
	}
	callnote:="{\"status\":\""+status+"\",\"phone\":\""+number+"\"}"
	//_,err:=db.Exec("INSERT INTO tCampRingCards_callnote set campaignid ="+campaignid+", cardid = "+ringcardid+", userid = "+userid+", time=NOW(), callnote = '"+callnote+"',action = 'perlapp', operator = 'perlapp'")
	//checkErr(err)
	stmtIns, err := db.Prepare("INSERT INTO tCampRingCards_callnote set campaignid =? ,cardid =? , userid =?, time=NOW(), callnote =? ,action = 'perlapp', operator = 'perlapp'")
	checkErr(err)
	defer stmtIns.Close()
	_, err=stmtIns.Exec(campaignid,ringcardid,userid,callnote)
	checkErr(err)

}
//chua lam robocaller

func db_robo_call(id string, maxcall string , percent string)(int,string){
	select_query:="select r.campaign_id,r.segments,r.voices,t.campNumber from robocaller r inner join tCampaign t on r.campaign_id=t.campaignID where r.id="+id+" and r.status=1"
	row,err := db.Query(select_query)
	checkErr(err)
	columnNames, _ := row.Columns()
	rc := NewMapStringScan(columnNames)
	row.Next()
	rc.Update(row)
	row.Close()
	//convert field voice (JSON string) to array maps
	var voices []map[string]string
	json.Unmarshal([]byte(rc.row["voices"]), &voices)
	fmt.Println(voices)
	mc.Set(&memcache.Item{Key: "robo_call", Value: []byte(strconv.Itoa(0))})
	for i:=0;i<len(voices);i++ {
		where:=""
		if voices[i]["segmentid"] !=""{
			where +="and t.segmentID in ("+voices[i]["segmentid"]+")"
		}
		if voices[i]["statusid"] !=""{
			where +="and t.statusID in ("+voices[i]["statusid"]+")"
		}
		if voices[i]["subid"] !=""{
			where +="and t.subID in ("+voices[i]["subid"]+")"
		}
		select_query="select t.rID,t.Phone1,t.Phone2,t.Phone3,t.Phone4,t.Phone5,t.statusID,t.subID,d.path from tCampRingCards t inner join dialplan_voicefiles d on d.id="+voices[i]["soundid"]+" where t.campaignID="+rc.row["campaign_id"]+" "+where
		row,err = db.Query(select_query)
		checkErr(err)
		max, _ := strconv.Atoi(maxcall)
		go db_robo_call_process(row,max,percent,id,rc.row["campNumber"],rc.row["campaign_id"])
	}
	return 200,"OK"

}
func db_robo_call_process(rows *sql.Rows ,maxcall int, percent string, taskid string, trunk string, campaignid string){
	columnNames, _ := rows.Columns()
	rc := NewMapStringScan(columnNames)
	for rows.Next() {
		item,_ := mc.Get("robo_call");
		count, _ := strconv.Atoi(string(item.Value))

		//fmt.Println(count)
		rc.Update(rows)
		for count>= maxcall{
			time.Sleep(time.Second * 20)
			item,_ = mc.Get("robo_call");
			count,_= strconv.Atoi(string(item.Value))
		}
		count++
		mc.Set(&memcache.Item{Key: "robo_call", Value: []byte(strconv.Itoa(count))})
		phonenum:=rc.row["Phone1"]+":"+rc.row["Phone2"]+":"+rc.row["Phone3"]+":"+rc.row["Phone4"]+":"+rc.row["Phone5"]
		fName := filepath.Base(rc.row["path"])
		extName := filepath.Ext(rc.row["path"])
		bName := fName[:len(fName)-len(extName)]
		go db_robo_callnote(campaignid,rc.row["rID"])
		go ast_robo_call(phonenum,bName,trunk,taskid,rc.row["rID"],percent)
	}
}

func db_robo_callnote(campaignid string, ringcardid string){
	callnote:="Robocaller has called this card"
	//_,err:=db.Exec("INSERT INTO tCampRingCards_callnote set campaignid ="+campaignid+", cardid = "+ringcardid+", userid = 0, time=NOW(), callnote = '"+callnote+"',type = 0, operator = 'perlapp'")
	//checkErr(err)
	stmtIns, err := db.Prepare("INSERT INTO tCampRingCards_callnote set campaignid =? ,cardid =? , userid =0, time=NOW(), callnote =? ,type = 0, operator = 'perlapp'")
	checkErr(err)
	defer stmtIns.Close()
	_, err=stmtIns.Exec(campaignid,ringcardid,callnote)
}

func db_robo_call_status(m map[string]string){
	//jsonString, _ := json.Marshal(m)
	plog ("ast_robo_call_event ",1);
	query:="insert into  robocaller_log set `taskid`="+m["TaskID"]+" ,`rid` = "+m["CardID"]+",`status`='"+m["Status"]+"', `reason`='"+m["Reason"]+"' "
	if(m["Length"]!=""){
		query+=" ,`voice_length` ="+m["Length"]
	}
	if(m["Duration"]!=""){
		query+=" ,`duration` ="+m["Duration"]
	}
	if(m["Phone"]!=""){
		query+=" ,`phone` ="+m["Phone"]
	}
	_,err:=db.Exec(query)
	checkErr(err)
}

func db_voicedrop_callnote(campaignid string,ringcardid string, agent string, callnote string){
	stmtIns, err := db.Prepare("INSERT INTO tCampRingCards_callnote set campaignid =? ,cardid =? , userid =?, time=NOW(), callnote =? , operator = 'perlapp'")
	checkErr(err)
	defer stmtIns.Close()
	_, err=stmtIns.Exec(campaignid,ringcardid,agent,callnote)
	//_,err:=db.Exec("INSERT INTO tCampRingCards_callnote set campaignid ="+campaignid+", cardid = "+ringcardid+", userid = "+agent+", time=NOW(), callnote = '"+callnote+"', operator = 'perlapp'")
	checkErr(err)
}
/**
  using a map
*/
type mapStringScan struct {
	// cp are the column pointers
	cp []interface{}
	// row contains the final result
	row      map[string]string
	colCount int
	colNames []string
}

func NewMapStringScan(columnNames []string) *mapStringScan {
	lenCN := len(columnNames)
	s := &mapStringScan{
		cp:       make([]interface{}, lenCN),
		row:      make(map[string]string, lenCN),
		colCount: lenCN,
		colNames: columnNames,
	}
	for i := 0; i < lenCN; i++ {
		s.cp[i] = new(sql.RawBytes)
	}
	return s
}

func (s *mapStringScan) Update(rows *sql.Rows) error {
	if err := rows.Scan(s.cp...); err != nil {
		return err
	}

	for i := 0; i < s.colCount; i++ {
		if rb, ok := s.cp[i].(*sql.RawBytes); ok {
			s.row[s.colNames[i]] = string(*rb)
			*rb = nil // reset pointer to discard current value to avoid a bug
		} else {
			return fmt.Errorf("Cannot convert index %d column %s to type *sql.RawBytes", i, s.colNames[i])
		}
	}
	return nil
}