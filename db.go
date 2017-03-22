package main
import(
	//"database/sql"
	//_ "github.com/go-sql-driver/mysql"
	//"github.com/ziutek/mymysql/mysql"
	//_ "github.com/ziutek/mymysql/native"
	"strconv"
)
func db_getstate(campaignid string){

	var t_ratio,t_ratio_up,t_ratio_down float64
	var t_wait_time int
	var t_campNumber string
	rows, err := db.Query("SELECT ratio,wait_time, ratio_up, ratio_down ,campNumber from tCampaign where campaignID= "+campaignid)//.Scan(&t_ratio,&t_wait_time,&t_ratio_up,&t_ratio_down,&t_campNumber)
	checkErr(err)
	//defer rows.Close()
	//rows, _ = stmt.Run(campaignid)
	rows.Scan(&t_ratio,&t_wait_time,&t_ratio_up,&t_ratio_down,&t_campNumber)
	//if(t_ratio==nil){
	//	set_default_ratio(campaignid)
	//}else{
		if(t_ratio_up > -2 && t_ratio_up < 2){
			ratio_up[campaignid]=t_ratio_up
			plog("Set ration up ="+strconv.FormatFloat(t_ratio_up, 'E', -1, 64)+" for campaign "+campaignid)
		}
		if(t_ratio_down > -2 && t_ratio_down < 2){
			ratio_down[campaignid]=t_ratio_down
			plog("Set ration down = "+strconv.FormatFloat(t_ratio_down, 'E', -1, 64)+" for campaign "+campaignid)
		}
		if(t_wait_time> 10000 && t_wait_time < 90000){
			dial_timeout=t_wait_time
			plog("Set dial timeout = "+strconv.Itoa(t_wait_time)+" for campaign "+campaignid)
		}
		if(t_ratio > 10000 && t_ratio < 90000){
			db_ratio[campaignid]=t_ratio
			plog("Set ratio = "+strconv.FormatFloat(t_ratio, 'E', -1, 64)+" for campaign "+campaignid)
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