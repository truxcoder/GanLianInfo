package controllers

import (
	"GanLianInfo/models"
	"GanLianInfo/utils"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	log "github.com/truxcoder/truxlog"
	"strings"
	"time"
)

// 含当前级别名称的模型
type perPost struct {
	models.Personnel
	PostLevelName string `json:"postLevelName"`
	//PostLevel int `json:"postLevel"`
}

// 含职务开始日期的模型
type currentLevel struct {
	ID            int64     `json:"id,string"`
	OrganID       string    `json:"organId"`
	LevelStartDay time.Time `json:"levelStartDay"`
}

// LeaderTeam 班子成员模型
type LeaderTeam struct {
	models.Leader
	PersonnelName  string `json:"personnelName"`
	PoliceCode     string `json:"policeCode"`
	IdCode         string `json:"idCode"`
	OrganName      string `json:"organName"`
	OrganShortName string `json:"organShortName"`
}

func AutoTask() {
	if data, err := getPoliceTeamData(); err != nil {
		log.Error(err)
	} else {
		log.Success(data)
	}

}

// AnalysisPoliceTeamData 民警队伍数据分析数据获取
func AnalysisPoliceTeamData(c *gin.Context) {
	var (
		r gin.H
		d []models.Department
	)
	if data, err := getPoliceTeamData(); err != nil {
		log.Error(err)
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	} else {
		selectStr := "id,position"
		result := db.Table("departments").Select(selectStr).Where("dept_type = ?", 1).Order("sort desc,level_code asc").Find(&d)
		err = result.Error
		if err != nil {
			r = GetError(CodeServer)
			c.JSON(200, r)
			return
		}
		r = gin.H{"code": 20000, "data": data, "department": &d}
		c.JSON(200, r)
	}
}

// 获取民警队伍数据
func getPoliceTeamData() (data string, err error) {
	var (
		overview []struct { // 数据概述
			ID        string `json:"id"`        //单位id
			Headcount int    `json:"headcount"` //单位编制数
			Use       int    `json:"use"`       // 单位占用数
			Male      int    `json:"male"`      // 男民警数
			Female    int    `json:"female"`    // 女民警数
			Communist int    `json:"communist"` // 党员数
			Minority  int    `json:"minority"`  // 少数民族数
			Zk        int    `json:"zk"`        //正科
			Fk        int    `json:"fk"`        //副科
			Zc        int    `json:"zc"`        //正处
			Fc        int    `json:"fc"`        //副处
			Ft        int    `json:"ft"`        //副厅
		}
		personnels []perPost
	)

	// 先构建三个map，最终处理成json数据返回前端
	overviewMap := make(map[string]map[string]int) // 外层key是单位id，内层key是数据项名称，例如key为headcount，值为单位编制数
	ageMap := make(map[string]map[string]int)
	eduMap := make(map[string]map[string]int)
	selectStr := "personnels.*, levels.name post_level_name"
	joinStr := "left join levels on levels.id = personnels.current_level"
	// 在编人员
	if err = db.Model(&models.Personnel{}).Select(selectStr).Joins(joinStr).Where("status = ?", 1).Find(&personnels).Error; err != nil {
		return
	}

	// 各单位编制数和占用数
	selectStr = "departments.id,departments.headcount,(select count (1) from personnels where personnels.organ_id = departments.id and personnels.status = 1) use," +
		"(select count (1) from personnels where personnels.organ_id = departments.id and personnels.status = 1 and gender = '男') male," +
		"(select count (1) from personnels where personnels.organ_id = departments.id and personnels.status = 1 and gender = '女') female," +
		"(select count (1) from personnels where personnels.organ_id = departments.id and personnels.status = 1 and nation is not null and nation <> '汉族') minority," +
		"(select count (1) from personnels where personnels.organ_id = departments.id and personnels.status = 1 and political in ('中共党员', '中共预备党员') ) communist," +
		"(select count (1) from personnels where personnels.organ_id = departments.id and personnels.status = 1 and personnels.id in (select id from current_pos where level_name = '正科级') ) zk," +
		"(select count (1) from personnels where personnels.organ_id = departments.id and personnels.status = 1 and personnels.id in (select id from current_pos where level_name = '副科级') ) fk," +
		"(select count (1) from personnels where personnels.organ_id = departments.id and personnels.status = 1 and personnels.id in (select id from current_pos where level_name = '正处级') ) zc," +
		"(select count (1) from personnels where personnels.organ_id = departments.id and personnels.status = 1 and personnels.id in (select id from current_pos where level_name = '副处级') ) fc," +
		"(select count (1) from personnels where personnels.organ_id = departments.id and personnels.status = 1 and personnels.id in (select id from current_pos where level_name = '副厅级') ) ft"
	// 查询得到概述数据
	db.Table("departments").Select(selectStr).Where("dept_type = ?", 1).Find(&overview)

	// 遍历在编人员， 循环处理其年龄和学历
	for _, p := range personnels {
		genAge(ageMap, &p)
		genEdu(eduMap, &p)
	}

	// 遍历概述数据，以单位id为key生成概述map
	for _, o := range overview {
		overviewMap[o.ID] = map[string]int{
			"headcount": o.Headcount,
			"use":       o.Use,
			"male":      o.Male,
			"female":    o.Female,
			"communist": o.Communist,
			"minority":  o.Minority,
			"zk":        o.Zk,
			"fk":        o.Fk,
			"zc":        o.Zc,
			"fc":        o.Fc,
			"ft":        o.Ft,
		}
	}

	// 把三个map打包，生成返回数据
	dataMap := map[string]map[string]map[string]int{
		"overview": overviewMap,
		"age":      ageMap,
		"edu":      eduMap,
	}
	data, err = jsoniter.MarshalToString(dataMap)
	if err != nil {
		return "", err
	}
	return
}

// 处理民警学历， 传两个参数，一个是学历map， 一个是当前要处理的人员指针
func genEdu(eduMap map[string]map[string]int, p *perPost) {
	// 如果学历map里还没有这个单位， 则生成这个单位
	if _, ok := eduMap[p.OrganID]; !ok {
		eduMap[p.OrganID] = map[string]int{"yjs": 0, "yjs_zk": 0, "yjs_fk": 0, "qrzyjs": 0, "qrzyjs_zk": 0, "qrzyjs_fk": 0, "bk": 0, "bk_zk:": 0, "bk_fk": 0, "qrzbk": 0, "qrzbk_zk": 0, "qrzbk_fk": 0, "zk": 0, "zk_zk:": 0, "zk_fk": 0, "qrzzk": 0, "qrzzk_zk": 0, "qrzzk_fk": 0}
	}
	// 取出当前单位信息， 当前单位也是一个内层map
	edu := eduMap[p.OrganID]
	if strings.Contains(p.FinalEdu, "研究生") {
		edu["yjs"] += 1
		if p.PostLevelName == "正科级" {
			edu["yjs_zk"] += 1
		}
		if p.PostLevelName == "副科级" {
			edu["yjs_fk"] += 1
		}
	}
	if strings.Contains(p.FullTimeEdu, "研究生") {
		edu["qrzyjs"] += 1
		if p.PostLevelName == "正科级" {
			edu["qrzyjs_zk"] += 1
		}
		if p.PostLevelName == "副科级" {
			edu["qrzyjs_fk"] += 1
		}
	}
	if strings.Contains(p.FinalEdu, "大学") && p.FinalEdu != "大学普通班" {
		edu["bk"] += 1
		if p.PostLevelName == "正科级" {
			edu["bk_zk"] += 1
		}
		if p.PostLevelName == "副科级" {
			edu["bk_fk"] += 1
		}
	}
	if strings.Contains(p.FullTimeEdu, "大学") && p.FullTimeEdu != "大学普通班" {
		edu["qrzbk"] += 1
		if p.PostLevelName == "正科级" {
			edu["qrzbk_zk"] += 1
		}
		if p.PostLevelName == "副科级" {
			edu["qrzbk_fk"] += 1
		}
	}
	if strings.Contains(p.FinalEdu, "大专") || p.FinalEdu == "大学普通班" {
		edu["zk"] += 1
		if p.PostLevelName == "正科级" {
			edu["zk_zk"] += 1
		}
		if p.PostLevelName == "副科级" {
			edu["zk_fk"] += 1
		}
	}
	if strings.Contains(p.FullTimeEdu, "大专") || p.FullTimeEdu == "大学普通班" {
		edu["qrzzk"] += 1
		if p.PostLevelName == "正科级" {
			edu["qrzzk_zk"] += 1
		}
		if p.PostLevelName == "副科级" {
			edu["qrzzk_fk"] += 1
		}
	}
}

// 处理年龄信息
func genAge(ageMap map[string]map[string]int, p *perPost) {
	if age, ok := ageMap[p.OrganID]; ok {
		_age := utils.GetAgeFromBirthday(p.Birthday)
		age["num"] += 1
		age["total"] += _age
		if p.PostLevelName == "正科级" {
			age["zkNum"] += 1
			age["zkTotal"] += _age
			if _age > age["zkOldest"] {
				age["zkOldest"] = _age
			}
			if (age["zkYoungest"] != 0 && _age < age["zkYoungest"]) || age["zkYoungest"] == 0 {
				age["zkYoungest"] = _age
			}
		}
		if p.PostLevelName == "副科级" {
			age["fkNum"] += 1
			age["fkTotal"] += _age
			if _age > age["fkOldest"] {
				age["fkOldest"] = _age
			}
			if (age["fkYoungest"] != 0 && _age < age["fkYoungest"]) || age["fkYoungest"] == 0 {
				age["fkYoungest"] = _age
			}
		}
		if _age > age["oldest"] {
			age["oldest"] = _age
		}
		if _age < age["youngest"] {
			age["youngest"] = _age
		}
		if p.Gender == "男" {
			age["maleNum"] += 1
			age["maleTotal"] += _age
			if _age > age["maleOldest"] {
				age["maleOldest"] = _age
			}
			if (age["maleYoungest"] != 0 && _age < age["maleYoungest"]) || age["maleYoungest"] == 0 {
				age["maleYoungest"] = _age
			}
			if _age >= 58 {
				age["willRetire"] += 1
			}
		}
		if p.Gender == "女" {
			age["femaleNum"] += 1
			age["femaleTotal"] += _age
			if _age > age["femaleOldest"] {
				age["femaleOldest"] = _age
			}
			if (age["femaleYoungest"] != 0 && _age < age["femaleYoungest"]) || age["femaleYoungest"] == 0 {
				age["femaleYoungest"] = _age
			}
			if _age >= 53 {
				age["willRetire"] += 1
			}
		}
		if _age >= 20 && _age < 25 {
			age["20-24"] += 1
			if p.PostLevelName == "正科级" {
				age["zk-20-24"] += 1
			}
			if p.PostLevelName == "副科级" {
				age["fk-20-24"] += 1
			}
		}
		if _age >= 25 && _age < 30 {
			age["25-29"] += 1
			if p.PostLevelName == "正科级" {
				age["zk-25-29"] += 1
			}
			if p.PostLevelName == "副科级" {
				age["fk-25-29"] += 1
			}
		}
		if _age >= 30 && _age < 35 {
			age["30-34"] += 1
			if p.PostLevelName == "正科级" {
				age["zk-30-34"] += 1
			}
			if p.PostLevelName == "副科级" {
				age["fk-30-34"] += 1
			}
		}
		if _age >= 35 && _age < 40 {
			age["35-39"] += 1
			if p.PostLevelName == "正科级" {
				age["zk-35-39"] += 1
			}
			if p.PostLevelName == "副科级" {
				age["fk-35-39"] += 1
			}
		}
		if _age >= 40 && _age < 45 {
			age["40-44"] += 1
			if p.PostLevelName == "正科级" {
				age["zk-40-44"] += 1
			}
			if p.PostLevelName == "副科级" {
				age["fk-40-44"] += 1
			}
		}
		if _age >= 45 && _age < 50 {
			age["45-49"] += 1
			if p.PostLevelName == "正科级" {
				age["zk-45-49"] += 1
			}
			if p.PostLevelName == "副科级" {
				age["fk-45-49"] += 1
			}
		}
		if _age >= 50 && _age < 55 {
			age["50-54"] += 1
			if p.PostLevelName == "正科级" {
				age["zk-50-54"] += 1
			}
			if p.PostLevelName == "副科级" {
				age["fk-50-54"] += 1
			}
		}
		if _age >= 55 && _age < 60 {
			age["55-59"] += 1
			if p.PostLevelName == "正科级" {
				age["zk-55-59"] += 1
			}
			if p.PostLevelName == "副科级" {
				age["fk-55-59"] += 1
			}
		}

	} else {
		_age := utils.GetAgeFromBirthday(p.Birthday)

		temp := map[string]int{"num": 1, "maleNum:": 0, "femaleNum": 0, "zkNum": 0, "fkNum": 0, "total": _age, "zkTotal": 0, "fkTotal": 0, "oldest": _age, "youngest": _age, "zkOldest": 0, "fkOldest": 0, "zkYoungest": 0, "fkYoungest": 0, "maleTotal": 0, "femaleTotal": 0, "maleOldest": 0, "maleYoungest": 0, "femaleOldest": 0, "femaleYoungest": 0, "willRetire": 0, "20-24": 0, "zk-20-24": 0, "fk-20-24": 0, "25-29": 0, "zk-25-29": 0, "fk-25-29": 0, "30-34": 0, "zk-30-34": 0, "fk-30-34": 0, "35-39": 0, "zk-35-39": 0, "fk-35-39": 0, "40-44": 0, "zk-40-44": 0, "fk-40-44": 0, "45-49": 0, "zk-45-49": 0, "fk-45-49": 0, "50-54": 0, "zk-50-54": 0, "fk-50-54": 0, "55-59": 0, "zk-55-59": 0, "fk-55-59": 0}
		if p.PostLevelName == "正科级" {
			temp["zkNum"] = 1
			temp["zkOldest"] = _age
			temp["zkYoungest"] = _age
		}
		if p.PostLevelName == "副科级" {
			temp["fkNum"] = 1
			temp["fkOldest"] = _age
			temp["fkYoungest"] = _age
		}
		if p.Gender == "男" {
			temp["maleNum"] = 1
			temp["maleTotal"] = _age
			temp["maleOldest"] = _age
			temp["maleYoungest"] = _age
			if _age >= 58 {
				temp["willRetire"] = 1
			}
		}
		if p.Gender == "女" {
			temp["femaleNum"] = 1
			temp["femaleTotal"] += _age
			temp["femaleOldest"] = _age
			temp["femaleYoungest"] = _age
			if _age >= 53 {
				temp["willRetire"] = 1
			}
		}
		if _age >= 20 && _age < 25 {
			temp["20-24"] = 1
			if p.PostLevelName == "正科级" {
				temp["zk-20-24"] = 1
			}
			if p.PostLevelName == "副科级" {
				temp["fk-20-24"] = 1
			}
		}
		if _age >= 25 && _age < 30 {
			temp["25-29"] = 1
			if p.PostLevelName == "正科级" {
				temp["zk-25-29"] = 1
			}
			if p.PostLevelName == "副科级" {
				temp["fk-25-29"] = 1
			}
		}
		if _age >= 30 && _age < 35 {
			temp["30-34"] = 1
			if p.PostLevelName == "正科级" {
				temp["zk-30-34"] = 1
			}
			if p.PostLevelName == "副科级" {
				temp["fk-30-34"] = 1
			}
		}
		if _age >= 35 && _age < 40 {
			temp["35-39"] = 1
			if p.PostLevelName == "正科级" {
				temp["zk-35-39"] = 1
			}
			if p.PostLevelName == "副科级" {
				temp["fk-35-39"] = 1
			}
		}
		if _age >= 40 && _age < 45 {
			temp["40-44"] = 1
			if p.PostLevelName == "正科级" {
				temp["zk-40-44"] = 1
			}
			if p.PostLevelName == "副科级" {
				temp["fk-40-44"] = 1
			}
		}
		if _age >= 45 && _age < 50 {
			temp["45-49"] = 1
			if p.PostLevelName == "正科级" {
				temp["zk-45-49"] = 1
			}
			if p.PostLevelName == "副科级" {
				temp["fk-45-49"] = 1
			}
		}
		if _age >= 50 && _age < 55 {
			temp["50-54"] = 1
			if p.PostLevelName == "正科级" {
				temp["zk-50-54"] = 1
			}
			if p.PostLevelName == "副科级" {
				temp["fk-50-54"] = 1
			}
		}
		if _age >= 55 && _age < 60 {
			temp["55-59"] = 1
			if p.PostLevelName == "正科级" {
				temp["zk-55-59"] = 1
			}
			if p.PostLevelName == "副科级" {
				temp["fk-55-59"] = 1
			}
		}
		ageMap[p.OrganID] = temp
	}
}

// AnalysisLeaderTeamData 获取数据分析班子数据
func AnalysisLeaderTeamData(c *gin.Context) {
	var (
		r   gin.H
		d   []models.Department
		ld  []LeaderTeam
		err error
	)
	selectStr := "leaders.*,per.name as personnel_name, per.police_code as police_code, per.id_code as id_code," +
		"d.name as organ_name, d.short_name as organ_short_name"
	joinStr := "left join personnels as per on leaders.personnel_id = per.id " +
		"left join departments as d on leaders.organ_id = d.id "
	result := db.Table("leaders").Select(selectStr).Joins(joinStr).Order("organ_id, sort").Find(&ld)
	err = result.Error
	if err != nil {
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	}

	if data, err := getLeaderTeamData(); err != nil {
		log.Error(err)
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	} else {
		selectStr = "id,position"
		result = db.Table("departments").Select(selectStr).Where("dept_type = ?", 1).Order("sort desc,level_code asc").Find(&d)
		err = result.Error
		if err != nil {
			r = GetError(CodeServer)
			c.JSON(200, r)
			return
		}
		r = gin.H{"code": 20000, "data": data, "leaders": &ld, "department": &d}
		c.JSON(200, r)
	}
}

func getLeaderTeamData() (data string, err error) {

	var (
		overview []struct {
			ID       string `json:"id"`       //单位id
			Male     int    `json:"male"`     //班子成员男性数量
			Female   int    `json:"female"`   //班子成员女性数量
			Minority int    `json:"minority"` //班子成员少数民族数量
			Zc       int    `json:"zc"`       //正处占用
			Fc       int    `json:"fc"`       //副处占用
		}
		personnels []perPost
		cl         []currentLevel
	)

	overviewMap := make(map[string]map[string]int)
	ageMap := make(map[string]map[string]int)
	eduMap := make(map[string]map[string]int)
	postMap := make(map[string]map[string]int)
	//selectStr := "personnels.*, (case when exists (select * from posts where posts.personnel_id = personnels.id and posts.level_id = (select levels.id from levels where levels.name = '正科级') and posts.position_id in (select id from positions where is_leader = 2) and posts.end_day ='0001-01-01 00:00:00.000000 +00:00') then 2 when exists (select * from posts where posts.personnel_id = personnels.id and posts.level_id = (select levels.id from levels where levels.name = '副科级') and posts.position_id in (select id from positions where is_leader = 2) and posts.end_day ='0001-01-01 00:00:00.000000 +00:00') then 1 else 0 end) as post_level"
	selectStr := "personnels.*, levels.name post_level_name"
	joinStr := "left join levels on levels.id = personnels.current_level"
	// 找出全局所有班子人员
	if err = db.Model(&models.Personnel{}).Select(selectStr).Joins(joinStr).Where("status = ? and personnels.id in (select personnel_id from leaders)", 1).Find(&personnels).Error; err != nil {
		return
	}
	selectStr = "current_pos.id, current_pos.level_start_day,leaders.organ_id"
	joinStr = "left join leaders on current_pos.id = leaders.personnel_id"
	// 副处级任职时间
	if err = db.Table("current_pos").Select(selectStr).Joins(joinStr).Where("level_name = ? and current_pos.id in (select personnel_id from leaders)", "副处级").Find(&cl).Error; err != nil {
		return
	}
	// 各单位班子总体情况
	selectStr = "departments.id," +
		"(select count (1) from leaders where leaders.organ_id = departments.id and (select gender from personnels where personnels.id = leaders.personnel_id) = '男') male," +
		"(select count (1) from leaders where leaders.organ_id = departments.id and (select gender from personnels where personnels.id = leaders.personnel_id) = '女') female," +
		"(select count (1) from leaders where leaders.organ_id = departments.id and (select nation from personnels where personnels.id = leaders.personnel_id)  is not null and (select nation from personnels where personnels.id = leaders.id)  <> '汉族') minority," +
		"(select count (1) from personnels where personnels.organ_id = departments.id and personnels.status = 1 and personnels.id in (select id from current_pos where level_name = '正处级') ) zc," +
		"(select count (1) from personnels where personnels.organ_id = departments.id and personnels.status = 1 and personnels.id in (select id from current_pos where level_name = '副处级') ) fc "
	db.Table("departments").Select(selectStr).Where("dept_type = ? and short_name <> ?", 1, "局机关").Find(&overview)

	genLeaderPost(postMap, cl)
	for _, p := range personnels {
		genLeaderAge(ageMap, &p)
		genLeaderEdu(eduMap, &p)
	}

	for _, o := range overview {
		overviewMap[o.ID] = map[string]int{
			"male":     o.Male,
			"female":   o.Female,
			"minority": o.Minority,
			"zc":       o.Zc,
			"fc":       o.Fc,
		}
	}

	dataMap := map[string]map[string]map[string]int{
		"overview": overviewMap,
		"age":      ageMap,
		"edu":      eduMap,
		"post":     postMap,
	}
	data, err = jsoniter.MarshalToString(dataMap)
	if err != nil {
		return "", err
	}
	return
}

// 处理班子年龄信息
func genLeaderAge(ageMap map[string]map[string]int, p *perPost) {
	_1960, _ := time.Parse(time.RFC3339, "1960-01-01T00:00:00Z")
	_1970, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")
	_1980, _ := time.Parse(time.RFC3339, "1980-01-01T00:00:00Z")
	_1990, _ := time.Parse(time.RFC3339, "1990-01-01T00:00:00Z")
	_2000, _ := time.Parse(time.RFC3339, "2000-01-01T00:00:00Z")
	if age, ok := ageMap[p.OrganID]; ok {
		_age := utils.GetAgeFromBirthday(p.Birthday)
		age["num"] += 1
		age["total"] += _age
		if p.PostLevelName == "正处级" {
			age["zcNum"] += 1
			age["zcTotal"] += _age
			if _age > age["zcOldest"] {
				age["zcOldest"] = _age
			}
			if (age["zcYoungest"] != 0 && _age < age["zcYoungest"]) || age["zcYoungest"] == 0 {
				age["zcYoungest"] = _age
			}
			if p.Birthday.After(_1960) && p.Birthday.Before(_1970) {
				age["zc-60er"] += 1
			}
			if p.Birthday.After(_1970) && p.Birthday.Before(_1980) {
				age["zc-70er"] += 1
			}
			if p.Birthday.After(_1980) && p.Birthday.Before(_1990) {
				age["zc-80er"] += 1
			}
			if p.Birthday.After(_1990) && p.Birthday.Before(_2000) {
				age["zc-90er"] += 1
			}
		}
		if p.PostLevelName == "副处级" {
			age["fcNum"] += 1
			age["fcTotal"] += _age
			if _age > age["fcOldest"] {
				age["fcOldest"] = _age
			}
			if (age["fcYoungest"] != 0 && _age < age["fcYoungest"]) || age["fcYoungest"] == 0 {
				age["fcYoungest"] = _age
			}
			if p.Birthday.After(_1960) && p.Birthday.Before(_1970) {
				age["fc-60er"] += 1
			}
			if p.Birthday.After(_1970) && p.Birthday.Before(_1980) {
				age["fc-70er"] += 1
			}
			if p.Birthday.After(_1980) && p.Birthday.Before(_1990) {
				age["fc-80er"] += 1
			}
			if p.Birthday.After(_1990) && p.Birthday.Before(_2000) {
				age["fc-90er"] += 1
			}
		}
		if _age > age["oldest"] {
			age["oldest"] = _age
		}
		if _age < age["youngest"] {
			age["youngest"] = _age
		}

		if _age >= 30 && _age < 35 {
			age["30-34"] += 1
			if p.PostLevelName == "正处级" {
				age["zc-30-34"] += 1
			}
			if p.PostLevelName == "副处级" {
				age["fc-30-34"] += 1
			}
		}
		if _age >= 35 && _age < 40 {
			age["35-39"] += 1
			if p.PostLevelName == "正处级" {
				age["zc-35-39"] += 1
			}
			if p.PostLevelName == "副处级" {
				age["fc-35-39"] += 1
			}
		}
		if _age >= 40 && _age < 45 {
			age["40-44"] += 1
			if p.PostLevelName == "正处级" {
				age["zc-40-44"] += 1
			}
			if p.PostLevelName == "副处级" {
				age["fc-40-44"] += 1
			}
		}
		if _age >= 45 && _age < 50 {
			age["45-49"] += 1
			if p.PostLevelName == "正处级" {
				age["zc-45-49"] += 1
			}
			if p.PostLevelName == "副处级" {
				age["fc-45-49"] += 1
			}
		}
		if _age >= 50 && _age < 55 {
			age["50-54"] += 1
			if p.PostLevelName == "正处级" {
				age["zc-50-54"] += 1
			}
			if p.PostLevelName == "副处级" {
				age["fc-50-54"] += 1
			}
		}
		if _age >= 55 && _age < 60 {
			age["55-59"] += 1
			if p.PostLevelName == "正处级" {
				age["zc-55-59"] += 1
			}
			if p.PostLevelName == "副处级" {
				age["fc-55-59"] += 1
			}
		}

	} else {
		_age := utils.GetAgeFromBirthday(p.Birthday)

		temp := map[string]int{"num": 1, "zcNum": 0, "fcNum": 0, "total": _age, "zcTotal": 0, "fcTotal": 0, "oldest": _age, "youngest": _age, "zcOldest": 0, "fcOldest": 0, "zcYoungest": 0, "fcYoungest": 0, "30-34": 0, "zc-30-34": 0, "fc-30-34": 0, "35-39": 0, "zc-35-39": 0, "fc-35-39": 0, "40-44": 0, "zc-40-44": 0, "fc-40-44": 0, "45-49": 0, "zc-45-49": 0, "fc-45-49": 0, "50-54": 0, "zc-50-54": 0, "fc-50-54": 0, "55-59": 0, "zc-55-59": 0, "fc-55-59": 0, "zc-60er": 0, "fc-60er": 0, "zc-70er": 0, "fc-70er": 0, "zc-80er": 0, "fc-80er": 0, "zc-90er": 0, "fc-90er": 0}
		if p.PostLevelName == "正处级" {
			temp["zcNum"] = 1
			temp["zcOldest"] = _age
			temp["zcYoungest"] = _age
			if p.Birthday.After(_1960) && p.Birthday.Before(_1970) {
				temp["zc-60er"] = 1
			}
			if p.Birthday.After(_1970) && p.Birthday.Before(_1980) {
				temp["zc-70er"] = 1
			}
			if p.Birthday.After(_1980) && p.Birthday.Before(_1990) {
				temp["zc-80er"] = 1
			}
			if p.Birthday.After(_1990) && p.Birthday.Before(_2000) {
				temp["zc-90er"] = 1
			}
		}
		if p.PostLevelName == "副处级" {
			temp["fcNum"] = 1
			temp["fcOldest"] = _age
			temp["fcYoungest"] = _age
			if p.Birthday.After(_1960) && p.Birthday.Before(_1970) {
				temp["fc-60er"] = 1
			}
			if p.Birthday.After(_1970) && p.Birthday.Before(_1980) {
				temp["fc-70er"] = 1
			}
			if p.Birthday.After(_1980) && p.Birthday.Before(_1990) {
				temp["fc-80er"] = 1
			}
			if p.Birthday.After(_1990) && p.Birthday.Before(_2000) {
				temp["fc-90er"] = 1
			}
		}

		if _age >= 30 && _age < 35 {
			temp["30-34"] = 1
			if p.PostLevelName == "正处级" {
				temp["zc-30-34"] = 1
			}
			if p.PostLevelName == "副处级" {
				temp["fc-30-34"] = 1
			}
		}
		if _age >= 35 && _age < 40 {
			temp["35-39"] = 1
			if p.PostLevelName == "正处级" {
				temp["zc-35-39"] = 1
			}
			if p.PostLevelName == "副处级" {
				temp["fc-35-39"] = 1
			}
		}
		if _age >= 40 && _age < 45 {
			temp["40-44"] = 1
			if p.PostLevelName == "正处级" {
				temp["zc-40-44"] = 1
			}
			if p.PostLevelName == "副处级" {
				temp["fc-40-44"] = 1
			}
		}
		if _age >= 45 && _age < 50 {
			temp["45-49"] = 1
			if p.PostLevelName == "正处级" {
				temp["zc-45-49"] = 1
			}
			if p.PostLevelName == "副处级" {
				temp["fc-45-49"] = 1
			}
		}
		if _age >= 50 && _age < 55 {
			temp["50-54"] = 1
			if p.PostLevelName == "正处级" {
				temp["zc-50-54"] = 1
			}
			if p.PostLevelName == "副处级" {
				temp["fc-50-54"] = 1
			}
		}
		if _age >= 55 && _age < 60 {
			temp["55-59"] = 1
			if p.PostLevelName == "正处级" {
				temp["zc-55-59"] = 1
			}
			if p.PostLevelName == "副处级" {
				temp["fc-55-59"] = 1
			}
		}
		ageMap[p.OrganID] = temp
	}
}

// 处理班子学历信息
func genLeaderEdu(eduMap map[string]map[string]int, p *perPost) {
	if _, ok := eduMap[p.OrganID]; !ok {
		eduMap[p.OrganID] = map[string]int{"yjs": 0, "yjs_zc": 0, "yjs_fc": 0, "qrzyjs": 0, "qrzyjs_zc": 0, "qrzyjs_fc": 0, "bk": 0, "bk_zc:": 0, "bk_fc": 0, "qrzbk": 0, "qrzbk_zc": 0, "qrzbk_fc": 0, "zk": 0, "zk_zc:": 0, "zk_fc": 0, "qrzzk": 0, "qrzzk_zc": 0, "qrzzk_fc": 0}
	}
	edu := eduMap[p.OrganID]
	if strings.Contains(p.FinalEdu, "研究生") {
		edu["yjs"] += 1
		if p.PostLevelName == "正处级" {
			edu["yjs_zc"] += 1
		}
		if p.PostLevelName == "副处级" {
			edu["yjs_fc"] += 1
		}
	}
	if strings.Contains(p.FullTimeEdu, "研究生") {
		edu["qrzyjs"] += 1
		if p.PostLevelName == "正处级" {
			edu["qrzyjs_zc"] += 1
		}
		if p.PostLevelName == "副处级" {
			edu["qrzyjs_fc"] += 1
		}
	}
	if strings.Contains(p.FinalEdu, "大学") && p.FinalEdu != "大学普通班" {
		edu["bk"] += 1
		if p.PostLevelName == "正处级" {
			edu["bk_zc"] += 1
		}
		if p.PostLevelName == "副处级" {
			edu["bk_fc"] += 1
		}
	}
	if strings.Contains(p.FullTimeEdu, "大学") && p.FullTimeEdu != "大学普通班" {
		edu["qrzbk"] += 1
		if p.PostLevelName == "正处级" {
			edu["qrzbk_zc"] += 1
		}
		if p.PostLevelName == "副处级" {
			edu["qrzbk_fc"] += 1
		}
	}
	if strings.Contains(p.FinalEdu, "大专") || p.FinalEdu == "大学普通班" {
		edu["zk"] += 1
		if p.PostLevelName == "正处级" {
			edu["zk_zc"] += 1
		}
		if p.PostLevelName == "副处级" {
			edu["zk_fc"] += 1
		}
	}
	if strings.Contains(p.FullTimeEdu, "大专") || p.FullTimeEdu == "大学普通班" {
		edu["qrzzk"] += 1
		if p.PostLevelName == "正处级" {
			edu["qrzzk_zc"] += 1
		}
		if p.PostLevelName == "副处级" {
			edu["qrzzk_fc"] += 1
		}
	}
}

// 处理班子任职信息
func genLeaderPost(postMap map[string]map[string]int, cl []currentLevel) {

	// 一个包含人员id和人员曾任职务数量的模型
	type postCountType struct {
		PersonnelId int64 `json:"personnelId"`
		Total       int64 `json:"total"`
	}

	var (
		ids       []int64 // 人员id列表
		posIds    []int64 // 实职副处职务id列表
		postCount []postCountType
	)
	postCountMap := make(map[int64]postCountType)

	// 循环遍历传过来的待处理人员列表，生成人员id列表，便于向查询语句传参
	for _, c := range cl {
		ids = append(ids, c.ID)
	}
	// 查询得到实职副处职务id列表
	db.Model(models.Position{}).Where("is_leader = 2 and level_id in (select id from levels where levels.name = ?)", "副处级").Pluck("id", &posIds)
	// 查询得到postCountType列表
	db.Model(models.Post{}).Select("personnel_id, count(1) total").Where("posts.personnel_id in ? and posts.position_id in ?", ids, posIds).Group("personnel_id").Find(&postCount)

	// 遍历postCount，得到一个map,以人员id为key, 通过人员id，可以得到这个结构体，从而得到人员曾任副处级职务的经历数量
	for _, pc := range postCount {
		postCountMap[pc.PersonnelId] = pc
	}

	for _, c := range cl {
		postTotal := postCountMap[c.ID].Total
		if _, ok := postMap[c.OrganID]; !ok {
			postMap[c.OrganID] = map[string]int{"year1": 0, "year2": 0, "year5": 0, "year10": 0, "post1": 0, "post2": 0, "post3": 0}
		}
		if c.LevelStartDay.After(yearsAgo(2)) {
			postMap[c.OrganID]["year1"]++
		}
		if c.LevelStartDay.Before(yearsAgo(2)) && c.LevelStartDay.After(yearsAgo(5)) {
			postMap[c.OrganID]["year2"]++
		}
		if c.LevelStartDay.Before(yearsAgo(5)) && c.LevelStartDay.After(yearsAgo(10)) {
			postMap[c.OrganID]["year5"]++
		}
		if c.LevelStartDay.Before(yearsAgo(10)) {
			postMap[c.OrganID]["year10"]++
		}
		if postTotal == 1 {
			postMap[c.OrganID]["post1"]++
		}
		if postTotal == 2 {
			postMap[c.OrganID]["post2"]++
		}
		if postTotal > 2 {
			postMap[c.OrganID]["post3"]++
		}
	}

}

func buildHighLevelStr(paramList *[]interface{}, level string, date time.Time) (whereStr string) {
	//var zero time.Time
	//now := time.Now().Local()

	//处分期未满
	//whereStr += " AND personnels.id not in (select personnel_id from disciplines where disciplines.deadline > CURDATE()) "

	switch level {
	case "g4":
		whereStr += " AND exists (select current_pos.id from current_pos,used_pos where current_pos.id = personnels.id and used_pos.id = personnels.id and rank_name in ('一级警长', '一级主任科员') and ADD_MONTHS(rank_start_day,(24 + rank_add_month)) <= ?  AND ((level_name = '正科级' and ADD_MONTHS(level_start_day,5*12) <= ? and ADD_MONTHS(personnels.birthday, 43*12) <= ?) or (level_name = '副科级' and ADD_MONTHS(level_start_day,15*12) <= ? and ADD_MONTHS(personnels.birthday, 48*12) <= ?) or (ADD_MONTHS(IFNULL(fk_start_day, zk_start_day),(15*12+IFNULL(fk_add_month, zk_add_month))) <= ? and ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 55 else 50 end)*12) <= ?) or (level_name = '正科级' and ADD_MONTHS(IFNULL(fk_start_day, zk_start_day),10*12) <= ? and ADD_MONTHS(personnels.birthday, 45*12) <= ?) or (level_name = '副科级' and ADD_MONTHS(fk_start_day,13*12) <= ? and ADD_MONTHS(personnels.birthday, 50*12) <= ?) or ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 57 else 52 end)*12) <= ?))"
		*paramList = append(*paramList, date, date, date, date, date, date, date, date, date, date, date, date)
	case "g3":
		whereStr += " AND personnels.id in (with fd_list as ( select id from personnels where ( ((ADD_MONTHS(IFNULL((select min(start_day) from posts where personnel_id = personnels.id and position_id = (select id from positions where name = '副调研员')),'3000-01-01 00:00:00.000000 +00:00' ),5*12)<= ?) and ADD_MONTHS(personnels.birthday, 50*12) <= ?)))select personnels.id from personnels left join fd_list on fd_list.id = personnels.id where personnels.id in (select current_pos.id from current_pos,used_pos where rank_name in('四级高级警长', '四级调研员') and ADD_MONTHS(rank_start_day,(24 + rank_add_month)) <= ? and current_pos.id = used_pos.id AND (ADD_MONTHS(IFNULL(fc_start_day, '3000-01-01 00:00:00.000000 +00:00' ),(2*12+IFNULL(fc_add_month, 0))) <= ? or current_pos.id = fd_list.id or (level_name = '正科级' and ADD_MONTHS(level_start_day,10*12) <= ? and ADD_MONTHS(personnels.birthday, 50*12) <= ?) or (level_name = '正科级' and ADD_MONTHS(IFNULL(fk_start_day, zk_start_day),15*12) <= ? and ADD_MONTHS(personnels.birthday, 50*12) <= ?) or (ADD_MONTHS(IFNULL(fk_start_day, zk_start_day),(15*12+IFNULL(fk_add_month, zk_add_month))) <= ? and ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 57 else 52 end)*12) <= ?) or (ADD_MONTHS(IFNULL(fk_start_day, zk_start_day),(20*12+IFNULL(fk_add_month, zk_add_month))) <= ? and ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 55 else 50 end)*12) <= ?) or ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 59 else 54 end)*12) <= ? )))"
		*paramList = append(*paramList, date, date, date, date, date, date, date, date, date, date, date, date, date)
	case "g2":
		whereStr += " AND personnels.id in (with fd_list as (select id from personnels where (((ADD_MONTHS(IFNULL((select min(start_day) from posts where personnel_id = personnels.id and position_id = (select id from positions where name = '副调研员')),'3000-01-01 00:00:00.000000 +00:00' ),8*12)<= CURDATE()) and ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 55 else 50 end)*12) <= CURDATE()))) select personnels.id from personnels left join fd_list on fd_list.id = personnels.id where  personnels.id in (select current_pos.id from current_pos,used_pos where rank_name in('三级高级警长', '三级调研员') and ADD_MONTHS(rank_start_day,(24 + rank_add_month)) <= CURDATE() and current_pos.id = used_pos.id AND ((ADD_MONTHS(IFNULL(fc_start_day, '3000-01-01 00:00:00.000000 +00:00' ),(2*12+IFNULL(fc_add_month, 0))) <= CURDATE() and ADD_MONTHS(personnels.birthday, 48*12) <= CURDATE()) or current_pos.id = fd_list.id or (ADD_MONTHS(IFNULL(zk_start_day, '3000-01-01 00:00:00.000000 +00:00'),(15*12+IFNULL(zk_add_month, 0))) <= CURDATE() and ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 57 else 52 end)*12) <= CURDATE())\nor (ADD_MONTHS(IFNULL(zk_start_day, '3000-01-01 00:00:00.000000 +00:00'),(20*12+IFNULL(fk_add_month, zk_add_month))) <= CURDATE() and ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 57 else 52 end)*12) <= CURDATE()) or ADD_MONTHS(personnels.birthday, (case personnels.gender when '男' then 59 else 54 end)*12) <= CURDATE())))"
	}

	return whereStr
}
