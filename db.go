package main
import(
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)
func db_getstate(campaignid string){
	var (
		ratio int
		wait_time int
		ratioup int
		ratiodown int
		campNumber int
	)
	rows, err := db.Query("SELECT ratio,wait_time, ratio_up, ratio_down ,campNumber from tCampaign where campaignID= ?", campaignid).Scan(&ratio,&wait_time,&ratioup,&ratiodown,&campNumber)
	checkErr(err)
	defer rows.Close()
	if(ratioup > -2 && ratioup < 2){
		ratio_up[campaignid]=ratioup
		plog("Set ration up ="+ratioup+" for campaign "+campaignid)
	}
	if(ratiodown > -2 && ratiodown < 2){
		ratio_down[campaignid]=ratiodown
		plog("Set ration down = "+ratiodown+" for campaign "+campaignid)
	}
	if(wait_time > 10000 && wait_time < 90000){
		dial_timeout=wait_time
		plog("Set dial timeout = "+wait_time+" for campaign "+campaignid)
	}
	if(ratio > 10000 && ratio < 90000){
		db_ratio[campaignid]=ratio
		plog("Set ratio = "+ratio+" for campaign "+campaignid)
	}
	trunk_list[campaignid]=campNumber
	plog("Set trunk = "+campNumber+" for campaign "+campaignid)

}
