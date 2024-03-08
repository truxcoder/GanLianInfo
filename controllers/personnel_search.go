package controllers

import (
	"strconv"
	"time"

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
				" AND personnels.id in (select personnel_id from posts where posts.level_id = (select id from levels where levels.name = '正科级') and posts.start_day <= ? and posts.position_id in (select id from positions where is_leader = 2))" +
				" AND personnels.id not in (select personnel_id from posts where posts.level_id in (select id from levels where levels.orders < 4) and posts.position_id in (select id from positions where is_leader = 2))" +
				" AND final_edu in (select name from edu_dicts where sort > 20)" +
				" AND pass_exam_day is not null AND pass_exam_day >= ?" +
				" AND personnels.id in (select personnel_id from (select personnel_id, count(1) total from appraisals where appraisals.season = 100 and appraisals.years in ? and appraisals.conclusion not in ('不称职','不确定等次') group by personnel_id) where total = 3)" +
				" AND political = '中共党员'" +
				" AND join_party_day <= ?"
			paramList = append(paramList, zero, yearsAgo(5), yearsAgo(3), yearsAgo(3), threeYears, yearsAgo(3))
		case "zc":
			var thisYear = time.Now().Year()
			var threeYears = []string{strconv.Itoa(thisYear - 1), strconv.Itoa(thisYear - 2), strconv.Itoa(thisYear - 3)}
			whereStr += " AND start_job_day is not null AND start_job_day <> ? AND start_job_day <= ?" +
				" AND personnels.id in (select personnel_id from (select personnel_id, count(1) total from posts where position_id in (select positions.id from positions where is_leader = 2) group by personnel_id) where total > 1)" +
				" AND personnels.id in (select personnel_id from posts where posts.level_id = (select id from levels where levels.name = '副处级') and posts.start_day <= ? and posts.position_id in (select id from positions where is_leader = 2))" +
				" AND personnels.id not in (select personnel_id from posts where posts.level_id in (select id from levels where levels.orders < 3) and posts.position_id in (select id from positions where is_leader = 2))" +
				" AND final_edu in (select name from edu_dicts where sort > 20)" +
				" AND pass_exam_day is not null AND pass_exam_day >= ?" +
				" AND personnels.id in (select personnel_id from (select personnel_id, count(1) total from appraisals where appraisals.season = 100 and appraisals.years in ? and appraisals.conclusion not in ('不称职','不确定等次') group by personnel_id) where total = 3)" +
				" AND political = '中共党员'" +
				" AND join_party_day <= ?"
			paramList = append(paramList, zero, yearsAgo(5), yearsAgo(2), yearsAgo(3), threeYears, yearsAgo(3))
		case "willUpInSixMonth":
			// 此处查询逻辑很绕。对于任职信息的判断有两种情况符合条件。
			// 1，end_day为空，start_day小于指定日期。
			// 2，end_day不为空，start_day小于指定日期。同时在posts表里能找到这样一条记录，
			// 它的position_id等于当前记录的position_id，它的personnel_id等于当前记录的personnel_id，它的end_day为空。
			//whereStr += " AND personnels.id in (select personnel_id from posts where posts.level_id in (select id from levels where levels.orders > 3) and posts.position_id in (select id from positions where is_leader = 1 and positions.name <> '一级警长') and ((posts.end_day = ? and posts.start_day <= ?) or (posts.end_day <> ? and posts.start_day<= ? and exists (select * from posts as po where po.position_id = posts.position_id and po.personnel_id = posts.personnel_id and po.end_day = ?))))"
			whereStr += willUpInSixMonth
			//paramList = append(paramList, zero, monthAgo(18), zero, monthAgo(18), zero)
			paramList = append(paramList, monthAgo(18))
		case "willUpInThreeMonth":
			//whereStr += " AND personnels.id in (select personnel_id from posts where posts.level_id in (select id from levels where levels.orders > 3) and posts.position_id in (select id from positions where is_leader = 1 and positions.name <> '一级警长') and ((posts.end_day = ? and posts.start_day <= ?) or (posts.end_day <> ? and posts.start_day<= ? and exists (select * from posts as po where po.position_id = posts.position_id and po.personnel_id = posts.personnel_id and po.end_day = ?))))"
			whereStr += willUpInThreeMonth
			//paramList = append(paramList, zero, monthAgo(21), zero, monthAgo(21), zero)
			paramList = append(paramList, monthAgo(21))
		case "willRetireInTwoYear":
			whereStr += " AND ((personnels.birthday <= ? and personnels.gender = '男') or (personnels.birthday <= ? and personnels.gender = '女'))"
			paramList = append(paramList, yearsAgo(58), yearsAgo(53))
		case "willUp1", "willUp2d", "willUp3d", "willUp4d", "willUp1z", "willUp2z", "willUp3z", "willUp4z", "willUp1k", "willUp1g", "willUp2g", "willUp3g", "willUp4g", "willUp1j", "willUp2j", "willUp3j", "willUp4j", "willUp1y":
			whereStr += buildWillUpStr(&paramList, sm.Extra)
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
		for k, v := range sm.FinalMajor {
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
			if sm.IsChief == "正职" {
				stmt = stmt.Where("position_id in (select positions.id from positions where is_chief = 2)")
			} else {
				stmt = stmt.Where("position_id in (select positions.id from positions where is_chief = 1)")
			}
		}
		if sm.IsLeader != "" {
			if sm.IsLeader == "是" {
				stmt = stmt.Where("position_id in (select positions.id from positions where is_leader= 2)")
			} else {
				stmt = stmt.Where("position_id in (select positions.id from positions where is_leader= 1)")
			}
		}
		if sm.Current != "" {
			if sm.Current == "是" {
				stmt = stmt.Where("position_id in (select positions.id from positions where end_day = ?)", zero)

			} else {
				stmt = stmt.Where("position_id in (select positions.id from positions where end_day <> ?)", zero)
			}
		}

		sql, args, _ := stmt.ToSql()
		//log.Successf("sql:%s\nargs:%v\n", sql, args)
		whereStr += " AND personnels.id in (" + sql + ")"
		// 对空参数进行判断，如果不传参数，则不加入参数列表。
		if len(args) != 0 {
			paramList = append(paramList, args)
		}

	}
	if len(sm.TwoPost) > 0 {
		var leaderList []int64
		if sm.TwoPost[0] == "不限" {
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

func buildWillUpStr(paramList *[]interface{}, extra string) (whereStr string) {
	//var zero time.Time
	//var preRankMap = map[string]string{
	//	"willUp2d": "一级调研员", "willUp3d": "四级调研员", "willUp4d": "一级主任科员", "willUp1z": "二级主任科员", "willUp2z": "三级主任科员",
	//	"willUp3z": "四级主任科员", "willUp4z": "一级科员", "willUp1k": "二级科员", "willUp1j": "二级警长", "willUp2j": "三级警长",
	//	"willUp3j": "四级警长", "willUp4j": "一级警员", "willUp1y": "二级警员",
	//}
	var preRankMap = map[string][]string{
		"willUp2d": {"三级调研员", "三级高级警长"}, "willUp3d": {"四级调研员", "四级高级警长"}, "willUp4d": {"一级主任科员", "一级警长"}, "willUp1z": {"二级主任科员", "二级警长"}, "willUp2z": {"三级主任科员", "三级警长"},
		"willUp3z": {"四级主任科员", "四级警长"}, "willUp4z": {"一级科员", "一级警员"}, "willUp1k": {"二级科员", "二级警员"}, "willUp1j": {"二级主任科员", "二级警长"}, "willUp2j": {"三级主任科员", "三级警长"},
		"willUp3j": {"四级主任科员", "四级警长"}, "willUp4j": {"一级科员", "一级警员"}, "willUp1y": {"二级科员", "二级警员"},
	}
	var preLevelMap = map[string]string{
		"willUp3d": "副处级", "willUp1z": "正科级", "willUp3z": "副科级", "willUp1j": "正科级", "willUp3j": "副科级",
	}
	//now := time.Now().Local()
	//配偶子女移居国境外
	whereStr += " AND personnels.id not in (select DISTINCT personnel_id from families where families.is_abroad = 1) "
	//处分期未满
	whereStr += " AND personnels.id not in (select personnel_id from disciplines where disciplines.deadline > CURDATE()) "
	switch extra {
	case "willUp2d", "willUp4d", "willUp2z", "willUp4z", "willUp1k":
		whereStr += " AND personnels.organ_id in (select id from departments where short_name = '局机关')"
		whereStr += " AND personnels.id in (select current_pos.id from current_pos where rank_name in ? and ADD_MONTHS(rank_start_day,(24 + rank_add_month)) <= CURDATE())"
		*paramList = append(*paramList, preRankMap[extra])
	case "willUp2j", "willUp4j", "willUp1y":
		whereStr += " AND personnels.organ_id not in (select id from departments where short_name = '局机关')"
		whereStr += " AND personnels.id in (select current_pos.id from current_pos where rank_name in ? and ADD_MONTHS(rank_start_day,(24 + rank_add_month)) <= CURDATE())"
		*paramList = append(*paramList, preRankMap[extra])
	case "willUp3d", "willUp1z", "willUp3z":
		whereStr += " AND personnels.organ_id in (select id from departments where short_name = '局机关')"
		whereStr += " AND personnels.id in (select current_pos.id from current_pos where rank_name in ? and (ADD_MONTHS(rank_start_day,(24 + rank_add_month)) <= CURDATE() or (level_name = ? and ADD_MONTHS(level_start_day,(24 + level_add_month)) <= CURDATE())))"
		*paramList = append(*paramList, preRankMap[extra], preLevelMap[extra])
	case "willUp1j", "willUp3j":
		whereStr += " AND personnels.organ_id not in (select id from departments where short_name = '局机关')"
		whereStr += " AND personnels.id in (select current_pos.id from current_pos where rank_name in ? and (ADD_MONTHS(rank_start_day,(24 + rank_add_month)) <= CURDATE() or (level_name = ? and ADD_MONTHS(level_start_day,(24 + level_add_month)) <= CURDATE())))"
		*paramList = append(*paramList, preRankMap[extra], preLevelMap[extra])
	case "willUp4g":
		whereStr += " AND personnels.organ_id not in (select id from departments where short_name = '局机关')"
		//whereStr += " AND personnels.id in (select current_pos.id from current_pos,used_pos where rank_name in ('一级警长', '一级主任科员') and ADD_MONTHS(rank_start_day,(24 + rank_add_month)) <= CURDATE() and current_pos.id = used_pos.id AND ((level_name = '正科级' and ADD_MONTHS(level_start_day,5*12) <= CURDATE() and ADD_MONTHS(personnels.birthday, 43*12) <= CURDATE()) or (level_name = '副科级' and ADD_MONTHS(level_start_day,15*12) <= CURDATE() and ADD_MONTHS(personnels.birthday, 48*12) <= CURDATE()) or (ADD_MONTHS(IFNULL(fk_start_day, zk_start_day),(15*12+IFNULL(fk_add_month, zk_add_month))) <= CURDATE() and ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 55 else 50 end)*12) <= CURDATE()) or (level_name = '正科级' and ADD_MONTHS(IFNULL(fk_start_day, zk_start_day),10*12) <= CURDATE() and ADD_MONTHS(personnels.birthday, 45*12) <= CURDATE()) or (level_name = '副科级' and ADD_MONTHS(fk_start_day,13*12) <= CURDATE() and ADD_MONTHS(personnels.birthday, 50*12) <= CURDATE()) or ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 57 else 52 end)*12) <= CURDATE()))"
		whereStr += " AND exists (select current_pos.id from current_pos,used_pos where current_pos.id = personnels.id and used_pos.id = personnels.id and rank_name in ('一级警长', '一级主任科员') and ADD_MONTHS(rank_start_day,(24 + rank_add_month)) <= CURDATE()  AND ((level_name = '正科级' and ADD_MONTHS(level_start_day,5*12) <= CURDATE() and ADD_MONTHS(personnels.birthday, 43*12) <= CURDATE()) or (level_name = '副科级' and ADD_MONTHS(level_start_day,15*12) <= CURDATE() and ADD_MONTHS(personnels.birthday, 48*12) <= CURDATE()) or (ADD_MONTHS(IFNULL(fk_start_day, zk_start_day),(15*12+IFNULL(fk_add_month, zk_add_month))) <= CURDATE() and ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 55 else 50 end)*12) <= CURDATE()) or (level_name = '正科级' and ADD_MONTHS(IFNULL(fk_start_day, zk_start_day),10*12) <= CURDATE() and ADD_MONTHS(personnels.birthday, 45*12) <= CURDATE()) or (level_name = '副科级' and ADD_MONTHS(fk_start_day,13*12) <= CURDATE() and ADD_MONTHS(personnels.birthday, 50*12) <= CURDATE()) or ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 57 else 52 end)*12) <= CURDATE()))"
	case "willUp3g":
		whereStr += " AND personnels.organ_id not in (select id from departments where short_name = '局机关')"
		whereStr += " AND personnels.id in (with fd_list as ( select id from personnels where ( ((ADD_MONTHS(IFNULL((select min(start_day) from posts where personnel_id = personnels.id and position_id = (select id from positions where name = '副调研员')),'3000-01-01 00:00:00.000000 +00:00' ),5*12)<= CURDATE()) and ADD_MONTHS(personnels.birthday, 50*12) <= CURDATE())))select personnels.id from personnels left join fd_list on fd_list.id = personnels.id where personnels.id in (select current_pos.id from current_pos,used_pos where rank_name in('四级高级警长', '四级调研员') and ADD_MONTHS(rank_start_day,(24 + rank_add_month)) <= CURDATE() and current_pos.id = used_pos.id AND (ADD_MONTHS(IFNULL(fc_start_day, '3000-01-01 00:00:00.000000 +00:00' ),(2*12+IFNULL(fc_add_month, 0))) <= CURDATE() or current_pos.id = fd_list.id or (level_name = '正科级' and ADD_MONTHS(level_start_day,10*12) <= CURDATE() and ADD_MONTHS(personnels.birthday, 50*12) <= CURDATE()) or (level_name = '正科级' and ADD_MONTHS(IFNULL(fk_start_day, zk_start_day),15*12) <= CURDATE() and ADD_MONTHS(personnels.birthday, 50*12) <= CURDATE()) or (ADD_MONTHS(IFNULL(fk_start_day, zk_start_day),(15*12+IFNULL(fk_add_month, zk_add_month))) <= CURDATE() and ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 57 else 52 end)*12) <= CURDATE()) or (ADD_MONTHS(IFNULL(fk_start_day, zk_start_day),(20*12+IFNULL(fk_add_month, zk_add_month))) <= CURDATE() and ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 55 else 50 end)*12) <= CURDATE()) or ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 59 else 54 end)*12) <= CURDATE() )))"
	case "willUp2g":
		whereStr += " AND personnels.organ_id not in (select id from departments where short_name = '局机关')"
		whereStr += " AND personnels.id in (with fd_list as (select id from personnels where (((ADD_MONTHS(IFNULL((select min(start_day) from posts where personnel_id = personnels.id and position_id = (select id from positions where name = '副调研员')),'3000-01-01 00:00:00.000000 +00:00' ),8*12)<= CURDATE()) and ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 55 else 50 end)*12) <= CURDATE()))) select personnels.id from personnels left join fd_list on fd_list.id = personnels.id where  personnels.id in (select current_pos.id from current_pos,used_pos where rank_name in('三级高级警长', '三级调研员') and ADD_MONTHS(rank_start_day,(24 + rank_add_month)) <= CURDATE() and current_pos.id = used_pos.id AND ((ADD_MONTHS(IFNULL(fc_start_day, '3000-01-01 00:00:00.000000 +00:00' ),(2*12+IFNULL(fc_add_month, 0))) <= CURDATE() and ADD_MONTHS(personnels.birthday, 48*12) <= CURDATE()) or current_pos.id = fd_list.id or (ADD_MONTHS(IFNULL(zk_start_day, '3000-01-01 00:00:00.000000 +00:00'),(15*12+IFNULL(zk_add_month, 0))) <= CURDATE() and ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 57 else 52 end)*12) <= CURDATE())\nor (ADD_MONTHS(IFNULL(zk_start_day, '3000-01-01 00:00:00.000000 +00:00'),(20*12+IFNULL(fk_add_month, zk_add_month))) <= CURDATE() and ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 57 else 52 end)*12) <= CURDATE()) or ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 59 else 54 end)*12) <= CURDATE())))"
	}

	return whereStr
}

func getBool(str string) int {
	if str == "是" {
		return 1
	}
	return 0
}

func isPostConditionEmpty(sm *SearchMod) bool {
	return len(sm.LevelId) == 0 && len(sm.PositionId) == 0 && sm.IsChief == "" && sm.IsLeader == "" && sm.Current == ""
}
