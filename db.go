package main
import(
	//"database/sql"
	//_ "github.com/go-sql-driver/mysql"
	//"github.com/ziutek/mymysql/mysql"
	//_ "github.com/ziutek/mymysql/native"
	"strconv"
)
func db_getstate(campaignid string){

	var t_ratio,t_ratio_up,t_ratio_down,t_wait_time,t_campNumber string
	err := db.QueryRow("SELECT ratio,wait_time, ratio_up, ratio_down ,campNumber from tCampaign where campaignID= ?",campaignid).Scan(&t_ratio,&t_wait_time,&t_ratio_up,&t_ratio_down,&t_campNumber)
	checkErr(err)
	strconv.ParseFloat("3.1415", 64)
	//defer rows.Close()
	//rows, _ = stmt.Run(campaignid)
	//rows.Scan(&t_ratio,&t_wait_time,&t_ratio_up,&t_ratio_down,&t_campNumber)
	//if(t_ratio==nil){
	//	set_default_ratio(campaignid)
	//}else{
		ratioup,_:=strconv.ParseFloat(t_ratio_up, 64)
		ratiodown,_:=strconv.ParseFloat(t_ratio_down, 64)
		ratio,_:=strconv.ParseFloat(t_ratio, 64)
		wait_time,_:=strconv.Atoi(t_wait_time)
		if(ratioup > -2 && ratioup < 2){
			ratio_up[campaignid]=ratioup
			plog("Set ration up ="+t_ratio_up+" for campaign "+campaignid)
		}
		if(ratiodown > -2 && ratiodown < 2){
			ratio_down[campaignid]=ratiodown
			plog("Set ration down = "+t_ratio_down+" for campaign "+campaignid)
		}
		if(wait_time> 10000 && wait_time < 90000){
			dial_timeout=wait_time
			plog("Set dial timeout = "+t_wait_time+" for campaign "+campaignid)
		}
		if(ratio > 10000 && ratio < 90000){
			db_ratio[campaignid]=ratio
			plog("Set ratio = "+t_ratio+" for campaign "+campaignid)
		}

		trunk_list[campaignid]=t_campNumber
		plog("Set trunk = "+t_campNumber+" for campaign "+campaignid)
	//}


}
func db_log(status string, agent string, ext string, campaignid string){
	query, err :=db.Prepare("INSERT INTO log set state = ?, agentid = ?,extension = ?,kampanj = ?, tid = NOW()");
	checkErr(err)
	defer query.Close()
	//query.Raw.Bind(status,agent,ext,campaignid)
	_, err=query.Exec(status,agent,ext,campaignid)
	checkErr(err)
	plog( "db_log "+status+", "+agent+", "+ext+", "+campaignid)
}
/*func ast_setstate(ringcardid string){
	stmt, err := db.Prepare("UPDATE tCampRingCards SET tapp = tapp + 1 where rID=?")
	checkErr(err)

	_, err = stmt.Exec( ringcardid)
	checkErr(err)
}

func set_num_status(campaignid string , ringcardid string,reason string, number string){
	select_query := "SELECT * from tCampRingCards WHERE rID ="+ringcardid
	rows, err := db.Query(select_query)
	checkErr(err)
	defer rows.Close()
	//chua xong
}*/