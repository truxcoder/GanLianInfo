package controllers

import (
	"strconv"
	"time"

	log "github.com/truxcoder/truxlog"

	sq "github.com/Masterminds/squirrel"
)

type SearchMod struct {
	Name                    string   `json:"name"`
	PoliceCode              string   `json:"policeCode"`
	AgeStart                int8     `json:"ageStart"`
	AgeEnd                  int8     `json:"ageEnd"`
	JobAgeStart             int8     `json:"JobAgeStart"`
	JobAgeEnd               int8     `json:"JobAgeEnd"`
	Gender                  string   `json:"gender"`
	Nation                  []string `json:"nation"`
	Political               []string `json:"political"`
	FullTimeEdu             []string `json:"fullTimeEdu"`
	FullTimeDegree          []string `json:"fullTimeDegree"`
	FullTimeMajor           []string `json:"fullTimeMajor"`
	FinalEdu                []string `json:"finalEdu"`
	FinalDegree             []string `json:"finalDegree"`
	FinalMajor              []string `json:"finalMajor"`
	OrganID                 []string `json:"organId"`
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
	LevelId                 []string `json:"levelId"`
	PositionId              []string `json:"positionId"`
	IsChief                 string   `json:"isChief"`
	IsLeader                string   `json:"isLeader"`
	Current                 string   `json:"current"`
	TwoPost                 []string `json:"twoPost"`
	Status                  bool     `json:"status"`
	Extra                   string   `json:"extra"`
}

func yearsAgo(year int8) time.Time {
	now := time.Now()
	return now.AddDate(int(0-year), 0, 0)
}

func monthAgo(month int16) time.Time {
	now := time.Now()
	return now.AddDate(0, int(0-month), 0)
}

func makeWhere(sm *SearchMod) (string, []interface{}) {
	var paramList []interface{}
	var zero time.Time
	whereStr := " 1 = 1 "

	if sm.Status {
		whereStr += " AND personnels.status = 1"
	} else {
		whereStr += " AND personnels.status = 0"
	}
	if len(sm.OrganID) > 0 {
		whereStr += " AND organ_id in ?"
		paramList = append(paramList, sm.OrganID)
	}
	if sm.Extra != "" {
		switch sm.Extra {
		case "fc":
			var thisYear = time.Now().Year()
			var threeYears = []string{strconv.Itoa(thisYear - 1), strconv.Itoa(thisYear - 2), strconv.Itoa(thisYear - 3)}
			whereStr += " AND start_job_day is not null AND start_job_day <> ? AND start_job_day <= ?" +
				" AND personnels.id in (select personnel_id from (select personnel_id, count(1) total from posts where position_id in (select positions.id from positions where is_leader = 2) group by personnel_id) where total > 1)" +
				" AND personnels.id in (select personnel_id from posts where posts.level_id = (select id from levels where levels.name = '?????????') and posts.start_day <= ? and posts.position_id in (select id from positions where is_leader = 2))" +
				" AND personnels.id not in (select personnel_id from posts where posts.level_id in (select id from levels where levels.orders < 4) and posts.position_id in (select id from positions where is_leader = 2))" +
				" AND final_edu in (select name from edu_dicts where sort > 20)" +
				" AND pass_exam_day is not null AND pass_exam_day >= ?" +
				" AND personnels.id in (select personnel_id from (select personnel_id, count(1) total from appraisals where appraisals.season = 100 and appraisals.years in ? and appraisals.conclusion not in ('?????????','???????????????') group by personnel_id) where total = 3)" +
				" AND political = '????????????'" +
				" AND join_party_day <= ?"
			paramList = append(paramList, zero, yearsAgo(5), yearsAgo(3), yearsAgo(3), threeYears, yearsAgo(3))
		case "zc":
			var thisYear = time.Now().Year()
			var threeYears = []string{strconv.Itoa(thisYear - 1), strconv.Itoa(thisYear - 2), strconv.Itoa(thisYear - 3)}
			whereStr += " AND start_job_day is not null AND start_job_day <> ? AND start_job_day <= ?" +
				" AND personnels.id in (select personnel_id from (select personnel_id, count(1) total from posts where position_id in (select positions.id from positions where is_leader = 2) group by personnel_id) where total > 1)" +
				" AND personnels.id in (select personnel_id from posts where posts.level_id = (select id from levels where levels.name = '?????????') and posts.start_day <= ? and posts.position_id in (select id from positions where is_leader = 2))" +
				" AND personnels.id not in (select personnel_id from posts where posts.level_id in (select id from levels where levels.orders < 3) and posts.position_id in (select id from positions where is_leader = 2))" +
				" AND final_edu in (select name from edu_dicts where sort > 20)" +
				" AND pass_exam_day is not null AND pass_exam_day >= ?" +
				" AND personnels.id in (select personnel_id from (select personnel_id, count(1) total from appraisals where appraisals.season = 100 and appraisals.years in ? and appraisals.conclusion not in ('?????????','???????????????') group by personnel_id) where total = 3)" +
				" AND political = '????????????'" +
				" AND join_party_day <= ?"
			paramList = append(paramList, zero, yearsAgo(5), yearsAgo(2), yearsAgo(3), threeYears, yearsAgo(3))
		case "willUpInSixMonth":
			// ????????????????????????????????????????????????????????????????????????????????????
			// 1???end_day?????????start_day?????????????????????
			// 2???end_day????????????start_day??????????????????????????????posts????????????????????????????????????
			// ??????position_id?????????????????????position_id?????????personnel_id?????????????????????personnel_id?????????end_day?????????
			whereStr += " AND personnels.id in (select personnel_id from posts where posts.level_id in (select id from levels where levels.orders > 3) and posts.position_id in (select id from positions where is_leader = 1 and positions.name <> '????????????') and ((posts.end_day = ? and posts.start_day <= ?) or (posts.end_day <> ? and posts.start_day<= ? and exists (select * from posts as po where po.position_id = posts.position_id and po.personnel_id = posts.personnel_id and po.end_day = ?))))"
			paramList = append(paramList, zero, monthAgo(18), zero, monthAgo(18), zero)
		case "willUpInThreeMonth":
			whereStr += " AND personnels.id in (select personnel_id from posts where posts.level_id in (select id from levels where levels.orders > 3) and posts.position_id in (select id from positions where is_leader = 1 and positions.name <> '????????????') and ((posts.end_day = ? and posts.start_day <= ?) or (posts.end_day <> ? and posts.start_day<= ? and exists (select * from posts as po where po.position_id = posts.position_id and po.personnel_id = posts.personnel_id and po.end_day = ?))))"
			paramList = append(paramList, zero, monthAgo(21), zero, monthAgo(21), zero)
		case "willRetireInTwoYear":
			whereStr += " AND ((personnels.birthday <= ? and personnels.gender = '???') or (personnels.birthday <= ? and personnels.gender = '???'))"
			paramList = append(paramList, yearsAgo(58), yearsAgo(53))
		}
		return whereStr, paramList
	}
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
	if sm.JobAgeStart != 0 && sm.JobAgeEnd != 0 {
		whereStr += " AND start_job_day is not null AND start_job_day <> ? AND start_job_day >= ? AND start_job_day <= ?"
		paramList = append(paramList, zero, yearsAgo(sm.JobAgeEnd), yearsAgo(sm.JobAgeStart))
	} else if sm.JobAgeStart != 0 {
		whereStr += " AND start_job_day is not null AND start_job_day <> ? AND start_job_day <= ?"
		paramList = append(paramList, zero, yearsAgo(sm.JobAgeStart))
	} else if sm.JobAgeEnd != 0 {
		whereStr += " AND start_job_day is not null AND start_job_day <> ? AND start_job_day >= ?"
		paramList = append(paramList, zero, yearsAgo(sm.JobAgeEnd))
	}
	if sm.Gender != "" {
		whereStr += " AND gender = ?"
		paramList = append(paramList, sm.Gender)
	}
	if len(sm.Nation) > 0 {
		hanFound := false
		otherFound := false
		for _, v := range sm.Nation {
			if v == "????????????" {
				otherFound = true
			}
			if v == "??????" {
				hanFound = true
			}
			if hanFound && otherFound {
				break
			}
		}
		if otherFound && !hanFound {
			whereStr += " AND nation <> ?"
			paramList = append(paramList, "??????")
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
	if len(sm.FullTimeDegree) > 0 {
		whereStr += " AND full_time_degree in ?"
		paramList = append(paramList, sm.FullTimeDegree)
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
	if len(sm.FinalEdu) > 0 {
		whereStr += " AND final_edu in ?"
		paramList = append(paramList, sm.FinalEdu)
	}
	if len(sm.FinalDegree) > 0 {
		whereStr += " AND final_degree in ?"
		paramList = append(paramList, sm.FinalDegree)
	}
	if len(sm.FinalMajor) > 0 {
		for k, v := range sm.FullTimeMajor {
			if k == 0 {
				whereStr += " AND (final_major LIKE ?"
				paramList = append(paramList, "%"+v+"%")
			} else {
				whereStr += " OR final_major LIKE ?"
				paramList = append(paramList, "%"+v+"%")
			}
		}
		whereStr += ")"
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
		if sm.IsSecret == "???" {
			whereStr += " AND is_secret = 2"
		} else {
			whereStr += " AND is_secret = 1"
		}
	}

	if sm.PassExamDay != "" {
		if sm.PassExamDay == "???" {
			whereStr += " AND pass_exam_day > ?"
		} else {
			whereStr += " AND (pass_exam_day < ? or pass_exam_day is null)"
		}
		paramList = append(paramList, yearsAgo(3))
	}
	if sm.HasPassport != "" {
		//whereStr += " AND has_passport = ?"
		//paramList = append(paramList, getBool(sm.HasPassport))
		if sm.HasPassport == "???" {
			whereStr += " AND (passport is not null and json_value(personnels.passport, '$[0]' RETURNING number) <> 0)"
		} else {
			whereStr += " AND (passport is null or passport = '' or json_value(personnels.passport, '$[0]' RETURNING number) = 0)"
		}
	}
	if sm.HasAppraisalIncompetent != "" {
		if sm.HasAppraisalIncompetent == "???" {
			whereStr += " AND personnels.id in (?)"
		} else {
			whereStr += " AND personnels.id not in (?)"
		}
		paramList = append(paramList, db.Table("appraisals").Select("personnel_id").Where("conclusion in ('?????????','???????????????')"))
	}
	if sm.HasReport != "" {
		if sm.HasReport == "???" {
			whereStr += " AND personnels.id in (?)"
		} else {
			whereStr += " AND personnels.id not in (?)"
		}
		paramList = append(paramList, db.Table("person_reports").Select("personnel_id").Where("report_id in (?)",
			db.Table("reports").Select("id").Where("step <> ?", 99)))
	}
	if sm.FamilyAbroad != "" {
		if sm.FamilyAbroad == "???" {
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
		//whereStr += " AND personnels.id in (?)"
		whereStr += " AND personnels.id in (select personnel_id from posts where level_id in ? and end_day = ? and posts.position_id in (select id from positions where is_leader = 2))"
		//paramList = append(paramList, db.Table("posts").Select("personnel_id").Where("level_id in ? and end_day = ?", level, zero))
		paramList = append(paramList, level, zero)
	}
	if !isPostConditionEmpty(sm) {
		stmt := sq.Select("personnel_id").From("posts")
		if len(sm.LevelId) > 0 {
			var level []int64
			for _, v := range sm.LevelId {
				_v, _ := strconv.Atoi(v)
				level = append(level, int64(_v))
			}
			stmt = stmt.Where(sq.Eq{"level_id": level})
		}
		if len(sm.PositionId) > 0 {
			var positions []int64
			for _, v := range sm.PositionId {
				_v, _ := strconv.Atoi(v)
				positions = append(positions, int64(_v))
			}
			stmt = stmt.Where(sq.Eq{"position_id": positions})
		}
		if sm.IsChief != "" {
			if sm.IsChief == "??????" {
				stmt = stmt.Where("position_id in (select positions.id from positions where is_chief = 2)")
			} else {
				stmt = stmt.Where("position_id in (select positions.id from positions where is_chief = 1)")
			}
		}
		if sm.IsLeader != "" {
			if sm.IsLeader == "???" {
				stmt = stmt.Where("position_id in (select positions.id from positions where is_leader= 2)")
			} else {
				stmt = stmt.Where("position_id in (select positions.id from positions where is_leader= 1)")
			}
		}
		if sm.Current != "" {
			if sm.Current == "???" {
				stmt = stmt.Where("position_id in (select positions.id from positions where end_day = ?)", zero)

			} else {
				stmt = stmt.Where("position_id in (select positions.id from positions where end_day <> ?)", zero)
			}
		}

		sql, args, _ := stmt.ToSql()
		log.Successf("sql:%s\nargs:%v\n", sql, args)
		whereStr += " AND personnels.id in (" + sql + ")"
		paramList = append(paramList, args)
	}
	if len(sm.TwoPost) > 0 {
		var leaderList []int64
		if sm.TwoPost[0] == "??????" {
			whereStr += " AND personnels.id in (select personnel_id from (select personnel_id, count(1) total from posts where position_id in (select positions.id from positions where is_leader = 2) group by personnel_id) where total > 1)"
		} else {
			for _, v := range sm.TwoPost {
				_v, _ := strconv.Atoi(v)
				leaderList = append(leaderList, int64(_v))
			}
			whereStr += " AND personnels.id in (select personnel_id from (select personnel_id, count(1) total from posts where position_id in (select positions.id from positions where is_leader = 2 and positions.level_id in ?) group by personnel_id) where total > 1)"
			paramList = append(paramList, leaderList)
		}
	}
	return whereStr, paramList
}

func getBool(str string) int {
	if str == "???" {
		return 1
	}
	return 0
}

func isPostConditionEmpty(sm *SearchMod) bool {
	return len(sm.LevelId) == 0 && len(sm.PositionId) == 0 && sm.IsChief == "" && sm.IsLeader == "" && sm.Current == ""
}
