package main
import(
	//"database/sql"
	//_ "github.com/go-sql-driver/mysql"
	//"github.com/ziutek/mymysql/mysql"
	//_ "github.com/ziutek/mymysql/native"
	"strconv"
)
func db_getstate(campaignid string){

	rows,res, err := db.Query("SELECT ratio,wait_time, ratio_up, ratio_down ,campNumber from tCampaign where campaignID= "+campaignid)
	checkErr(err)
	//rows, _ = stmt.Run(campaignid)
	//rows.Next().Scan(&ratio,&wait_time,&ratioup,&ratiodown,&campNumber)
	if(len(rows)==0){
		set_default_ratio(campaignid)
	}else{
		if(float64(res.Map("ratio_up")) > -2 && float64(res.Map("ratio_up")) < 2){
			ratio_up[campaignid]=float64(res.Map("ratio_up"))
			plog("Set ration up ="+strconv.Itoa(res.Map("ratio_up"))+" for campaign "+campaignid)
		}
		if(res.Map("ratio_down") > -2 && res.Map("ratio_down") < 2){
			ratio_down[campaignid]=float64(res.Map("ratio_down"))
			plog("Set ration down = "+strconv.Itoa(res.Map("ratio_down"))+" for campaign "+campaignid)
		}
		if(res.Map("wait_time") > 10000 && res.Map("wait_time") < 90000){
			dial_timeout=res.Map("wait_time")
			plog("Set dial timeout = "+strconv.Itoa(res.Map("wait_time"))+" for campaign "+campaignid)
		}
		if(res.Map("ratio") > 10000 && res.Map("ratio") < 90000){
			db_ratio[campaignid]=float64(res.Map("ratio"))
			plog("Set ratio = "+strconv.Itoa(res.Map("ratio"))+" for campaign "+campaignid)
		}

		trunk_list[campaignid]=res.Map("campNumber")
		plog("Set trunk = "+res.Map("campNumber")+" for campaign "+campaignid)
	}


}
func db_log(status string, agent string, ext string, campaignid string){
	query, err :=db.Prepare("INSERT INTO log set state = ?, agentid = ?,extension = ?,kampanj = ?, tid = NOW()");
	checkErr(err)
	//defer query.Close()
	query.Raw.Bind(status,agent,ext,campaignid)
	_, _, err=query.Exec()
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