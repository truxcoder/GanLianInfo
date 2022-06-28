package controllers

import (
	"GanLianInfo/models"
	"strconv"
	"time"

	log "github.com/truxcoder/truxlog"

	"github.com/gin-gonic/gin"
)

type Appointment struct {
	models.Appointment
	OrganShortName string `json:"organShortName"`
	PersonnelName  string `json:"personnelName"`
	PoliceCode     string `json:"policeCode"`
}

// AppointmentList 任免登记列表
func AppointmentList(c *gin.Context) {
	var mos []Appointment
	var mo struct {
		PersonnelId int64 `json:"personnelId,string"`
	}
	selectStr := "appointments.*,p.name as personnel_name,p.police_code as police_code,d.short_name as organ_short_name," +
		"p.name as personnel_name,p.police_code as police_code"
	joinStr := "left join personnels as p on appointments.personnel_id = p.id " +
		"left join departments as d on p.organ_id = d.id "
	getList(c, "appointments", &mo, &mos, &selectStr, &joinStr)

}

// AppointmentTableDetail 导出干审表
func AppointmentTableDetail(c *gin.Context) {
	var (
		err        error
		r          gin.H
		zero       time.Time
		awards     []models.Award
		punishes   []models.Punish
		appraisals []models.Appraisal
	)

	var mo struct {
		PersonnelId int64 `json:"personnelId,string"`
	}
	var mos struct {
		Name           string    `json:"name"`
		Gender         string    `json:"gender"`
		Birthday       string    `json:"birthday"`
		Nation         string    `json:"nation"`
		JoinPartyDay   time.Time `json:"joinPartyDay"`
		StartJobDay    time.Time `json:"startJobDay"`
		Hometown       string    `json:"hometown"`
		Birthplace     string    `json:"birthplace"`
		Health         string    `json:"health"`
		TechnicalTitle string    `json:"technicalTitle"`
		Specialty      string    `json:"specialty"`
		FullTimeEdu    string    `json:"fullTimeEdu"`
		FullTimeDegree string    `json:"fullTimeDegree"`
		FullTimeMajor  string    `json:"fullTimeMajor"`
		FullTimeSchool string    `json:"fullTimeSchool"`
		PartTimeEdu    string    `json:"partTimeEdu"`
		PartTimeDegree string    `json:"partTimeDegree"`
		PartTimeMajor  string    `json:"partTimeMajor"`
		PartTimeSchool string    `json:"partTimeSchool"`
		Resume         string    `json:"resume"`
	}
	var posts []struct {
		Department   string `json:"department"`
		Organ        string `json:"organ"`
		PositionName string `json:"positionName"`
		IsLeader     int8   `json:"isLeader"`
	}
	var family []struct {
		Name      string    `json:"name"`
		Relation  string    `json:"relation"`
		Birthday  time.Time `json:"birthday"`
		Organ     string    `json:"organ"`
		Post      string    `json:"post"`
		Political string    `json:"political"`
	}
	var thisYear = time.Now().Year()
	var threeYears = []string{strconv.Itoa(thisYear - 1), strconv.Itoa(thisYear - 2), strconv.Itoa(thisYear - 3)}
	var selectStr = "name,gender,birthday,nation,join_party_day,start_job_day,hometown,birthplace,health,technical_title,specialty,full_time_edu,full_time_degree,full_time_major,full_time_school,part_time_edu,part_time_degree,part_time_major,part_time_school,resumes.content as resume"
	var joinStr = "left join resumes on resumes.personnel_id = personnels.id"
	var postSelect = "posts.department, posts.organ, positions.name as position_name, positions.is_leader as is_leader"
	var postJoin = "left join positions on positions.id = posts.position_id"
	if err = c.BindJSON(&mo); err != nil {
		log.Error(err)
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}
	db.Table("personnels").Select(selectStr).Joins(joinStr).Where("personnels.id = ?", mo.PersonnelId).Limit(1).Find(&mos)
	db.Table("posts").Select(postSelect).Joins(postJoin).Where("personnel_id = ? and end_day = ?", mo.PersonnelId, zero).Find(&posts)
	db.Table("families").Where("personnel_id = ?", mo.PersonnelId).Find(&family)
	db.Table("awards").Where("personnel_id = ?", mo.PersonnelId).Find(&awards)
	db.Table("punishes").Where("personnel_id = ?", mo.PersonnelId).Find(&punishes)
	db.Table("appraisals").Where("personnel_id = ? and season = 100 and years in ?", mo.PersonnelId, threeYears).Find(&appraisals)
	r = gin.H{"code": 20000, "data": &mos, "posts": &posts, "family": &family, "awards": &awards, "punishes": &punishes, "appraisals": &appraisals}
	c.JSON(200, r)
}
