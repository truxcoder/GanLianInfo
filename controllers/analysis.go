package controllers

import (
	"GanLianInfo/models"
	"GanLianInfo/utils"

	"github.com/gin-gonic/gin"

	jsoniter "github.com/json-iterator/go"
	log "github.com/truxcoder/truxlog"
)

func AutoTask() {
	if data, err := getPoliceTeamData(); err != nil {
		log.Error(err)
	} else {
		log.Success(data)
	}

}

func AnalysisPoliceTeamData(c *gin.Context) {
	var r gin.H
	if data, err := getPoliceTeamData(); err != nil {
		log.Error(err)
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	} else {
		r = gin.H{"code": 20000, "data": data}
		c.JSON(200, r)
	}
}

func getPoliceTeamData() (data string, err error) {
	var (
		headcount []struct {
			ID        string `json:"id"`
			Headcount int    `json:"headcount"`
			Use       int    `json:"use"`
		}
		personnels []struct {
			models.Personnel
			PostLevel int `json:"postLevel"`
		}
		//ageMap map[string]map[string]int
		//total  int
		//zero   time.Time
	)
	ageMap := make(map[string]map[string]int)
	selectStr := "personnels.*, (case when exists (select * from posts where posts.personnel_id = personnels.id and posts.level_id = (select levels.id from levels where levels.name = '正科级') and posts.position_id in (select id from positions where is_leader = 2) and posts.end_day ='0001-01-01 00:00:00.000000 +00:00') then 2 when exists (select * from posts where posts.personnel_id = personnels.id and posts.level_id = (select levels.id from levels where levels.name = '副科级') and posts.position_id in (select id from positions where is_leader = 2) and posts.end_day ='0001-01-01 00:00:00.000000 +00:00') then 1 else 0 end) as post_level"
	// 在编人员
	if err = db.Model(&models.Personnel{}).Select(selectStr).Where("status = ?", 1).Find(&personnels).Error; err != nil {
		return
	}
	//total = len(personnels)
	// 各单位编制数和占用数
	selectStr = "departments.id,departments.headcount,(select count (1) from personnels where personnels.organ_id = departments.id and personnels.status = 1) use"
	db.Table("departments").Select(selectStr).Where("dept_type = ?", 1).Find(&headcount)

	for _, p := range personnels {
		if age, ok := ageMap[p.OrganID]; ok {
			_age := utils.GetAgeFromBirthday(p.Birthday)
			age["num"] += 1
			age["total"] += _age
			if p.PostLevel == 2 {
				age["zkNum"] += 1
				age["zkTotal"] += _age
				if _age > age["zkOldest"] {
					age["zkOldest"] = _age
				}
				if (age["zkYoungest"] != 0 && _age < age["zkYoungest"]) || age["zkYoungest"] == 0 {
					age["zkYoungest"] = _age
				}
			}
			if p.PostLevel == 1 {
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
				if p.PostLevel == 2 {
					age["zk-20-29"] += 1
				}
				if p.PostLevel == 1 {
					age["fk-20-29"] += 1
				}
			}
			if _age >= 25 && _age < 30 {
				age["25-29"] += 1
				if p.PostLevel == 2 {
					age["zk-20-29"] += 1
				}
				if p.PostLevel == 1 {
					age["fk-20-29"] += 1
				}
			}
			if _age >= 30 && _age < 35 {
				age["30-34"] += 1
				if p.PostLevel == 2 {
					age["30-34"] += 1
				}
				if p.PostLevel == 1 {
					age["fk-30-34"] += 1
				}
			}
			if _age >= 35 && _age < 40 {
				age["35-39"] += 1
				if p.PostLevel == 2 {
					age["35-39"] += 1
				}
				if p.PostLevel == 1 {
					age["fk-35-39"] += 1
				}
			}
			if _age >= 40 && _age < 45 {
				age["40-44"] += 1
				if p.PostLevel == 2 {
					age["40-44"] += 1
				}
				if p.PostLevel == 1 {
					age["fk-40-44"] += 1
				}
			}
			if _age >= 45 && _age < 50 {
				age["45-49"] += 1
				if p.PostLevel == 2 {
					age["45-49"] += 1
				}
				if p.PostLevel == 1 {
					age["fk-45-49"] += 1
				}
			}
			if _age >= 50 && _age < 55 {
				age["50-54"] += 1
				if p.PostLevel == 2 {
					age["50-54"] += 1
				}
				if p.PostLevel == 1 {
					age["fk-50-54"] += 1
				}
			}
			if _age >= 55 && _age < 60 {
				age["55-59"] += 1
				if p.PostLevel == 2 {
					age["55-59"] += 1
				}
				if p.PostLevel == 1 {
					age["fk-55-59"] += 1
				}
			}

		} else {
			_age := utils.GetAgeFromBirthday(p.Birthday)

			temp := map[string]int{"num": 1, "maleNum:": 0, "femaleNum": 0, "zkNum": 0, "fkNum": 0, "total": _age, "zkTotal": 0, "fkTotal": 0, "oldest": _age, "youngest": _age, "zkOldest": 0, "fkOldest": 0, "zkYoungest": 0, "fkYoungest": 0, "maleTotal": 0, "femaleTotal": 0, "maleOldest": 0, "maleYoungest": 0, "femaleOldest": 0, "femaleYoungest": 0, "willRetire": 0, "20-24": 0, "zk-20-29": 0, "fk-20-29": 0, "25-29": 0, "30-34": 0, "zk-30-34": 0, "fk-30-34": 0, "35-39": 0, "zk-35-39": 0, "fk-35-39": 0, "40-44": 0, "zk-40-44": 0, "fk-40-44": 0, "45-49": 0, "zk-45-49": 0, "fk-45-49": 0, "50-54": 0, "zk-50-54": 0, "fk-50-54": 0, "55-59": 0, "zk-55-59": 0, "fk-55-59": 0}
			if p.PostLevel == 2 {
				temp["zkNum"] = 1
				temp["zkOldest"] = _age
				temp["zkYoungest"] = _age
			}
			if p.PostLevel == 1 {
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
				if p.PostLevel == 2 {
					temp["zk-20-29"] = 1
				}
				if p.PostLevel == 1 {
					temp["fk-20-29"] = 1
				}
			}
			if _age >= 25 && _age < 30 {
				temp["25-29"] = 1
				if p.PostLevel == 2 {
					temp["zk-20-29"] = 1
				}
				if p.PostLevel == 1 {
					temp["fk-20-29"] = 1
				}
			}
			if _age >= 30 && _age < 35 {
				temp["30-34"] = 1
				if p.PostLevel == 2 {
					temp["30-34"] = 1
				}
				if p.PostLevel == 1 {
					temp["fk-30-34"] = 1
				}
			}
			if _age >= 35 && _age < 40 {
				temp["35-39"] = 1
				if p.PostLevel == 2 {
					temp["35-39"] = 1
				}
				if p.PostLevel == 1 {
					temp["fk-35-39"] = 1
				}
			}
			if _age >= 40 && _age < 45 {
				temp["40-44"] = 1
				if p.PostLevel == 2 {
					temp["40-44"] = 1
				}
				if p.PostLevel == 1 {
					temp["fk-40-44"] = 1
				}
			}
			if _age >= 45 && _age < 50 {
				temp["45-49"] = 1
				if p.PostLevel == 2 {
					temp["45-49"] = 1
				}
				if p.PostLevel == 1 {
					temp["fk-45-49"] = 1
				}
			}
			if _age >= 50 && _age < 55 {
				temp["50-54"] = 1
				if p.PostLevel == 2 {
					temp["50-54"] = 1
				}
				if p.PostLevel == 1 {
					temp["fk-50-54"] = 1
				}
			}
			if _age >= 55 && _age < 60 {
				temp["55-59"] = 1
				if p.PostLevel == 2 {
					temp["55-59"] = 1
				}
				if p.PostLevel == 1 {
					temp["fk-55-59"] = 1
				}
			}
			ageMap[p.OrganID] = temp
		}
	}

	data, err = jsoniter.MarshalToString(ageMap)
	if err != nil {
		return "", err
	}
	return
}
