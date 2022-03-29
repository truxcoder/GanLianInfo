package controllers

import (
	"strconv"
	"time"
)

type SearchMod struct {
	Name                    string   `json:"name"`
	PoliceCode              string   `json:"policeCode"`
	AgeStart                int8     `json:"ageStart"`
	AgeEnd                  int8     `json:"ageEnd"`
	Gender                  string   `json:"gender"`
	Nation                  []string `json:"nation"`
	Political               []string `json:"political"`
	FullTimeEdu             []string `json:"fullTimeEdu"`
	FullTimeMajor           []string `json:"fullTimeMajor"`
	PartTimeEdu             []string `json:"partTimeEdu"`
	OrganID                 string   `json:"organId"`
	ProCert                 []string `json:"proCert"`
	IsSecret                string   `json:"isSecret"`
	PassExamDay             string   `json:"passExamDay"`
	HasPassport             string   `json:"hasPassport"`
	HasAppraisalIncompetent string   `json:"hasAppraisalIncompetent"`
	HasReport               string   `json:"hasReport"`
	Award                   []int8   `json:"award"`
	Punish                  []int8   `json:"punish"`
	FamilyAbroad            string   `json:"familyAbroad"`
	Level                   []string `json:"level"`
}

func yearsAgo(year int8) time.Time {
	now := time.Now()
	return now.AddDate(int(0-year), 0, 0)
}

func makeWhere(sm *SearchMod) (string, []interface{}) {
	var paramList []interface{}
	whereStr := " 1 = 1 "
	if sm.Name != "" {
		whereStr += " AND personnels.name LIKE ?"
		paramList = append(paramList, "%"+sm.Name+"%")
	}
	if sm.PoliceCode != "" {
		whereStr += " AND police_code = ?"
		paramList = append(paramList, sm.PoliceCode)
	}
	if sm.AgeStart != 0 && sm.AgeEnd != 0 {
		whereStr += " AND birthday >= ? AND birthday <= ?"
		paramList = append(paramList, yearsAgo(sm.AgeEnd), yearsAgo(sm.AgeStart))
	} else if sm.AgeStart != 0 {
		whereStr += " AND birthday <= ?"
		paramList = append(paramList, yearsAgo(sm.AgeStart))
	} else if sm.AgeEnd != 0 {
		whereStr += " AND birthday >= ?"
		paramList = append(paramList, yearsAgo(sm.AgeEnd))
	}
	if sm.Gender != "" {
		whereStr += " AND gender = ?"
		paramList = append(paramList, sm.Gender)
	}
	if len(sm.Nation) > 0 {
		hanFound := false
		otherFound := false
		for _, v := range sm.Nation {
			if v == "少数民族" {
				otherFound = true
			}
			if v == "汉族" {
				hanFound = true
			}
			if hanFound && otherFound {
				break
			}
		}
		if otherFound && !hanFound {
			whereStr += " AND nation <> ?"
			paramList = append(paramList, "汉族")
		} else if !otherFound {
			whereStr += " AND nation in ?"
			paramList = append(paramList, sm.Nation)
		}
	}
	if len(sm.Political) > 0 {
		whereStr += " AND political in ?"
		paramList = append(paramList, sm.Political)
	}
	if len(sm.FullTimeEdu) > 0 {
		whereStr += " AND full_time_edu in ?"
		paramList = append(paramList, sm.FullTimeEdu)
	}
	if len(sm.FullTimeMajor) > 0 {
		//whereStr += " AND full_time_major in ?"
		//paramList = append(paramList, sm.FullTimeMajor)
		for k, v := range sm.FullTimeMajor {
			if k == 0 {
				whereStr += " AND (full_time_major LIKE ?"
				paramList = append(paramList, "%"+v+"%")
			} else {
				whereStr += " OR full_time_major LIKE ?"
				paramList = append(paramList, "%"+v+"%")
			}
		}
		whereStr += ")"
	}
	if len(sm.PartTimeEdu) > 0 {
		whereStr += " AND part_time_edu in ?"
		paramList = append(paramList, sm.PartTimeEdu)
	}
	if sm.OrganID != "" {
		whereStr += " AND organ_id = ?"
		paramList = append(paramList, sm.OrganID)
	}
	if len(sm.ProCert) > 0 {
		//whereStr += " AND pro_cert in ?"
		//paramList = append(paramList, sm.ProCert)
		for k, v := range sm.ProCert {
			if k == 0 {
				whereStr += " AND pro_cert like ?"
			} else {
				whereStr += " OR pro_cert like ?"
			}

			paramList = append(paramList, "%"+v+"%")
		}
	}
	if sm.IsSecret != "" {
		if sm.IsSecret == "是" {
			whereStr += " AND is_secret = 2"
		} else {
			whereStr += " AND is_secret = 1"
		}
	}

	if sm.PassExamDay != "" {
		if sm.PassExamDay == "是" {
			whereStr += " AND pass_exam_day > ?"
		} else {
			whereStr += " AND (pass_exam_day < ? or pass_exam_day is null)"
		}
		paramList = append(paramList, yearsAgo(3))
	}
	if sm.HasPassport != "" {
		//whereStr += " AND has_passport = ?"
		//paramList = append(paramList, getBool(sm.HasPassport))
		if sm.HasPassport == "是" {
			whereStr += " AND (passport is not null and json_value(personnels.passport, '$[0]' RETURNING number) <> 0)"
		} else {
			whereStr += " AND (passport is null or passport = '' or json_value(personnels.passport, '$[0]' RETURNING number) = 0)"
		}
	}
	if sm.HasAppraisalIncompetent != "" {
		if sm.HasAppraisalIncompetent == "是" {
			whereStr += " AND personnels.id in (?)"
		} else {
			whereStr += " AND personnels.id not in (?)"
		}
		paramList = append(paramList, db.Table("appraisals").Select("personnel_id").Where("conclusion in ('不称职','不确定等次')"))
	}
	if sm.HasReport != "" {
		if sm.HasReport == "是" {
			whereStr += " AND personnels.id in (?)"
		} else {
			whereStr += " AND personnels.id not in (?)"
		}
		paramList = append(paramList, db.Table("person_reports").Select("personnel_id").Where("report_id in (?)",
			db.Table("reports").Select("id").Where("step <> ?", 99)))
	}
	if sm.FamilyAbroad != "" {
		if sm.FamilyAbroad == "是" {
			whereStr += " AND personnels.id in (?)"
		} else {
			whereStr += "AND personnels.id not in (?)"
		}
		paramList = append(paramList, db.Table("families").Select("personnel_id").Where("is_abroad = ?", 1))
	}
	if len(sm.Award) > 0 {
		whereStr += " AND personnels.id in (?)"
		paramList = append(paramList, db.Table("awards").Select("personnel_id").Where("grade in ?", sm.Award))
	}
	if len(sm.Punish) > 0 {
		whereStr += " AND personnels.id in (?)"
		paramList = append(paramList, db.Table("punishes").Select("personnel_id").Where("grade in ?", sm.Punish))
	}
	if len(sm.Level) > 0 {
		var level []int64
		for _, v := range sm.Level {
			_v, _ := strconv.Atoi(v)
			level = append(level, int64(_v))
		}
		var zero = time.Time{}
		whereStr += " AND personnels.id in (?)"
		paramList = append(paramList, db.Table("posts").Select("personnel_id").Where("level_id in ? and end_day = ?", level, zero))
	}
	return whereStr, paramList
}

func getBool(str string) int {
	if str == "是" {
		return 1
	}
	return 0
}
