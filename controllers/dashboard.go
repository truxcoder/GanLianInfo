package controllers

import (
	"github.com/gin-gonic/gin"
)

type OrganTotal struct {
	OrganId string `json:"organId"`
	Total   int    `json:"total"`
}

func DashboardData(c *gin.Context) {

	var ageParams []interface{}
	var ageList []struct {
		OrganId               string `json:"organId"`
		OlderThanFifty        int    `json:"olderThanFifty"`
		BetweenFortyFifty     int    `json:"betweenFortyFifty"`
		BetweenThirtyForty    int    `json:"betweenThirtyForty"`
		YoungerThanThirty     int    `json:"youngerThanThirty"`
		YoungerThanThirtyFive int    `json:"youngerThanThirtyFive"`
	}
	var globalList struct {
		Total                 int `json:"total"`
		PartyMember           int `json:"partyMember"`
		Male                  int `json:"male"`
		OlderThanFifty        int `json:"olderThanFifty"`
		BetweenFortyFifty     int `json:"betweenFortyFifty"`
		BetweenThirtyForty    int `json:"betweenThirtyForty"`
		YoungerThanThirty     int `json:"youngerThanThirty"`
		YoungerThanThirtyFive int `json:"youngerThanThirtyFive"`
	}
	var genderList, politicalList, totalList []OrganTotal
	ageStr := "select organ_id, count(case when birthday < ? then 1 else null end) older_than_fifty" +
		",count(case when birthday >= ? and birthday <= ? then 1 else null end) between_forty_fifty" +
		",count(case when birthday >= ? and birthday <= ? then 1 else null end) between_thirty_forty" +
		",count(case when birthday > ? then 1 else null end) younger_than_thirty" +
		",count(case when birthday > ? then 1 else null end) younger_than_thirty_five" +
		" from personnels where user_type = 1 and status = 1 group by organ_id"
	politicalStr := "select organ_id, count(political) total from personnels where political = '中共党员' and user_type = 1 and status = 1 group by organ_id"
	genderStr := "select organ_id, count(gender) total from personnels where gender = '男' and user_type = 1 and status = 1 group by organ_id"
	totalStr := "select organ_id, count(1) total from personnels where user_type = 1 and status = 1 group by organ_id"
	globalStr := "select count(1) total" +
		",count(case when political = '中共党员' then 1 else null end) party_member" +
		",count(case when gender = '男' then 1 else null end) male" +
		",count(case when birthday < ? then 1 else null end) older_than_fifty" +
		",count(case when birthday >= ? and birthday <= ? then 1 else null end) between_forty_fifty" +
		",count(case when birthday >= ? and birthday <= ? then 1 else null end) between_thirty_forty" +
		",count(case when birthday > ? then 1 else null end) younger_than_thirty" +
		",count(case when birthday > ? then 1 else null end) younger_than_thirty_five" +
		" from personnels where user_type = 1 and status = 1"
	ageParams = append(ageParams, yearsAgo(50), yearsAgo(50), yearsAgo(40), yearsAgo(40), yearsAgo(30), yearsAgo(30), yearsAgo(35))

	db.Raw(ageStr, ageParams...).Scan(&ageList)
	db.Raw(genderStr).Scan(&genderList)
	db.Raw(politicalStr).Scan(&politicalList)
	db.Raw(totalStr).Scan(&totalList)
	db.Raw(globalStr, ageParams...).Scan(&globalList)
	r := gin.H{"code": 20000, "ageList": &ageList, "politicalList": &politicalList, "genderList": &genderList,
		"totalList": &totalList, "globalList": &globalList}
	c.JSON(200, r)
}
