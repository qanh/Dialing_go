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
)
func db_getstate(campaignid string){

	var t_ratio,t_ratio_up,t_ratio_down,t_wait_time,t_campNumber string
	err := db.QueryRow("SELECT ratio,wait_time, ratio_up, ratio_down ,campNumber from tCampaign where campaignID= ?",campaignid).Scan(&t_ratio,&t_wait_time,&t_ratio_up,&t_ratio_down,&t_campNumber)
	//if no row -> err !=nil
	checkErr(err)
	if(err == nil){
		//defer rows.Close()
		//rows, _ = stmt.Run(campaignid)
		//rows.Scan(&t_ratio,&t_wait_time,&t_ratio_up,&t_ratio_down,&t_campNumber)
		//if(t_ratio==nil){
		//	set_default_ratio(campaignid)
		//}else{
		ratioup, _ := strconv.ParseFloat(t_ratio_up, 64)
		ratiodown, _ := strconv.ParseFloat(t_ratio_down, 64)
		ratio, _ := strconv.ParseFloat(t_ratio, 64)
		wait_time, _ := strconv.Atoi(t_wait_time)
		plog("t_ratio:"+t_ratio,1)
		if (ratioup > -2 && ratioup < 2) {
			ratio_up[campaignid] = ratioup
			plog("Set ration up =" + t_ratio_up + " for campaign " + campaignid,1)
		}
		if (ratiodown > -2 && ratiodown < 2) {
			ratio_down[campaignid] = ratiodown
			plog("Set ration down = " + t_ratio_down + " for campaign " + campaignid,1)
		}
		if (wait_time > 10000 && wait_time < 90000) {
			dial_timeout = wait_time
			plog("Set dial timeout = " + t_wait_time + " for campaign " + campaignid,1)
		}
		if (ratio > 1 && ratio < 10) {
			db_ratio[campaignid] = ratio
			plog("Set ratio = " + t_ratio + " for campaign " + campaignid,1)
		}

		trunk_list[campaignid] = t_campNumber
		plog("Set trunk = " + t_campNumber + " for campaign " + campaignid,1)
		//}
	}else{
		set_default_ratio(campaignid)
	}

}
func db_log(status string, agent string, ext string, campaignid string){
	query, err :=db.Prepare("INSERT INTO log set state = ?, agentid = ?,extension = ?,kampanj = ?, tid = NOW()");
	checkErr(err)
	defer query.Close()
	//query.Raw.Bind(status,agent,ext,campaignid)
	_, err=query.Exec(status,agent,ext,campaignid)
	checkErr(err)
	plog( "db_log "+status+", "+agent+", "+ext+", "+campaignid,1)
}
func db_setstate(ringcardid string){
	stmt, err := db.Prepare("UPDATE tCampRingCards SET tapp = tapp + 1 where rID=?")
	checkErr(err)

	_, err = stmt.Exec( ringcardid)
	checkErr(err)
}

func db_set_num_status(campaignid string , ringcardid string,reason string, number string){
	type Campaign struct {
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
	update_query := "UPDATE tCampRingCards SET ";
	var real_status int
	status:=map [string]int{
		"status1":0,
		"status2":0,
		"status3":0,
		"status4":0,
		"status5":0,
	}
	phone:=map [string]string{
		"Phone1":"",
		"Phone2":"",
		"Phone3":"",
		"Phone4":"",
		"Phone5":"",
	}
	select_query := "SELECT Phone1,Phone2,Phone3,Phone4,Phone5,status1,status2,status3,status4,status5 from tCampRingCards WHERE rID ="+ringcardid
	err := db.QueryRow(select_query).Scan(phone["Phone1"],phone["Phone2"],phone["Phone3"],phone["Phone4"],phone["Phone5"],status["status1"],status["status2"],status["status3"],status["status4"],status["status5"])
	checkErr(err)
	i := 1
	index:=1
	for i < 6 {
		key:="Phone"+strconv.Itoa(index)
		if(number==phone[key]){
			i=6
		}else{
			index++
		}
		i++
	}
	//phone_key:="Phone"+index
	status_key:="status"+strconv.Itoa(index)
	if(reason=="trasigt"){
		real_status=status[status_key]+1000
	}else if(reason=="ejsvar"){
		real_status=status[status_key]+1
	}else{
		plog ("error: set_num_status () unknow reason\n", 1);
	}
	called:=1
	fail:=1
	for i := 1; i < 6; i++ {
		if(status["status"+strconv.Itoa(i)]==0 && len (phone["Phone"+strconv.Itoa(i)])>4){
			if(i!=index){
				called=0
			}
		}
		if(status["status"+strconv.Itoa(i)]<500 && len (phone["Phone"+strconv.Itoa(i)])>4){
			if(i!=index){
				fail=0
			}else if(reason=="ejsvar"){
				fail=0
			}
		}
	}
	if(fail==1){
		plog ("set_num_status: all number fail, update database for ringcard "+ringcardid,1);
		update_query="UPDATE tCampRingCards Set fail_try=fail_try+1 WHERE rID="+ringcardid
	}else if(called==1){
		i:=1
		for i < 6 {
			phonestatus:="status"+strconv.Itoa(i)
			if(status[phonestatus]<500 && len(phone["Phone"+strconv.Itoa(i)])>4){
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
	_, err = db.Exec(update_query)
	checkErr(err)

}
func db_dial(ratio int ,campaignid string ){
	plog("db_dial",1)
	tidsperiod:=tidsperiod ()
	query:="call PTakeActiveRingCard("+campaignid+","+strconv.Itoa(tidsperiod)+")"
	if(ratio>0){
		row, err := db.Query(query)
		defer row.Close()
		if(err!=nil){
			checkErr(err)
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
		fmt.Println("Nu är det slut på telefonnummer i den här kampanjen")
	}
	rc.Update(row)
	fmt.Printf("%#v\n\n", rc.row)
	ringcardid:=rc.row["rID"]
	status1,_:=strconv.Atoi(rc.row["status1"])
	status2,_:=strconv.Atoi(rc.row["status2"])
	status3,_:=strconv.Atoi(rc.row["status3"])
	status4,_:=strconv.Atoi(rc.row["status4"])
	status5,_:=strconv.Atoi(rc.row["status5"])
	number1:=rc.row["Phone1"]
	number2:=rc.row["Phone2"]
	number3:=rc.row["Phone3"]
	number4:=rc.row["Phone4"]
	number5:=rc.row["Phone5"]
	plog("db_dial_res"+number1+" "+number2,1)
	if _, ok := list_ringcard[ringcardid]; ok{
		time.Sleep(1)
		go db_dial(1,campaignid)
	}else{
		list_ringcard["ringcardid"]=1
		if(len(number1)>4 && status1==0){
			number=number1
			number_index=1
		}else if(len(number2)>4 && status2==0){
			number=number2
			number_index=2
		}else if(len(number3)>4 && status3==0){
		number=number3
		number_index=3
		}else if(len(number4)>4 && status4==0){
		number=number4
		number_index=4
		}else if(len(number5)>4 && status5==0){
		number=number5
		number_index=5
		}else{
			number_ok:=0
			update_query:=""
			for i:=5 ;i>0;i-- {
				status,_:=strconv.Atoi(rc.row["status"+strconv.Itoa(i)])
				if(status<500){
					number_ok=0
				}
			}
			if(number_ok==1){
				tidsperiod := tidsperiod ()
				update_query="UPDATE tCampRingCards Set lastcalldate=CURDATE()"
				if (tidsperiod == 1) {
					update_query = update_query+", AMdate=CURDATE()"
				}else if (tidsperiod == 2) {
					update_query = update_query+", PMdate=CURDATE()";
				} else if (tidsperiod == 3) {
					update_query = update_query+", Eveningdate=CURDATE()";
				}
				if ((len(number1) < 5) || (status1 > 500)) {
					update_query = update_query+", status1 = 1000 ";
				} else {
					update_query = update_query+", status1 = 0 ";
					number_check = 0;
					number_index = 1;
					number = number1;
				}
				if ((len(number2) < 5) || (status2 > 500)) {
					update_query = update_query+", status2 = 1000 ";
				} else {
					update_query = update_query+", status2 = 0 ";
					number_check = 0;
					number_index = 2;
					number = number2;
				}
				if ((len(number3) < 5) || (status3 > 500)) {
					update_query = update_query+", status3 = 1000 ";
				} else {
					update_query = update_query+", status3 = 0 ";
					number_check = 0;
					number_index = 3;
					number = number3;
				}
				if ((len(number4) < 5) || (status4 > 500)) {
					update_query = update_query+", status4 = 1000 ";
				} else {
					update_query = update_query+", status4 = 0 ";
					number_check = 0;
					number_index = 4;
					number = number4;
				}
				if ((len(number5) < 5) || (status5 > 500)) {
					update_query = update_query+", status5 = 1000 ";
				} else {
					update_query = update_query+", status5 = 0 ";
					number_check = 0;
					number_index = 5;
					number = number5;
				}
				update_query = update_query+" WHERE rID="+ringcardid
				if (number_check==1) {
					plog ("inga nummer med bra status och längd på detta ringkort $ringkort\n", 1);
					update_query = "UPDATE tCampRingCards Set userID=0,statusID=(select bortfall_status from tCampaign where campaignID="+campaignid+"),subID=(select inget_nr_sub from tCampaign where campaignID="+campaignid+"), closed_date=now() WHERE rID="+ringcardid;
				}
			}else{
				plog ("inget nummer har bra status på detta ringkort $ringkort\n", 1);
				update_query = "UPDATE tCampRingCards Set userID=0,statusID=(select bortfall_status from tCampaign where campaignID="+campaignid+"),subID=(select inget_nr_sub from tCampaign where campaignID="+campaignid+"), closed_date=now() WHERE rID="+ringcardid;
			}

			_,err=db.Exec(update_query)
			checkErr(err)
			delete(list_ringcard,ringcardid)

			//if(number_check!=0){
			//	go db_dial(1,campaignid)
			//}
		}
		plog(strconv.Itoa(number_index)+" "+strconv.Itoa(number_check),1)
		if(number_index!=0 || number_check==0){
			update_query:="UPDATE tCampRingCards SET status"+strconv.Itoa(number_index)+" =  1 where rID="+ringcardid
			_,err=db.Exec(update_query)
			checkErr(err)
			delete(list_ringcard,ringcardid)
			ast_dial(number,ringcardid,campaignid)
		}
	}
}
func db_log_soundfile(ringcardid string ,campaignid string ,agent string){
	var clientid string
	select_query := "SELECT clientID from tCampaign where campaignID="+campaignid
	err:=db.QueryRow(select_query).Scan(&clientid)
	if(err!=nil){
		host,_:=os.Hostname()
		now := time.Now()
		datestring:=now.Format("20060102_030405")
		recname := "LOGG_"+datestring+"_u"+agent+"_c"+campaignid+"_"+ringcardid+"_"+clientid+".wav"
		ast_rec_start(agent,recname,clientid)
		db.Exec("INSERT INTO soundfile set rid = "+ringcardid+", userid = "+agent+", campaignid = "+campaignid+", clientid = "+clientid+", originfilename = "+recname+", filename = "+recname+", closed = 0, converted = 0, cut = 0, start = NOW(),asterisk_ip="+host)
	}
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
	_,err:=db.Exec("INSERT INTO tSoundFiles set campaignID = "+campaignid+", clientID = "+clientid+", filename = "+logfile+", createtime=NOW(), va = "+strconv.Itoa(va)+",addr = "+addr+", filesize = 0, status=0")
	checkErr(err)
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