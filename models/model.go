package models

import (
	"GanLianInfo/utils"
	"time"

	"github.com/Insua/gorm-dm8/datatype"

	"gorm.io/gorm"
)

type Base struct {
	ID        int64 `json:"id,string" gorm:"autoIncrement:false;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int8
}

type BaseId struct {
	ID int64 `json:"id,string" gorm:"autoIncrement:false;primaryKey"`
}

func (b *BaseId) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = utils.GenId()
	return
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = utils.GenId()
	b.Version = 0
	return
}
func (b *Base) BeforeUpdate(tx *gorm.DB) (err error) {
	b.Version++
	return
}

type Personnel struct {
	Base
	UserId                  string    `json:"userId"`
	Name                    string    `json:"name" gorm:"size:50"`
	Gender                  string    `json:"gender" gorm:"size:50"`
	Nation                  string    `json:"nation" gorm:"size:50"`
	IdCode                  string    `json:"idCode" gorm:"size:50"`
	Birthday                time.Time `json:"birthday"`
	PoliceCode              string    `json:"policeCode" gorm:"size:50"`
	Political               string    `json:"political"`
	JoinPartyDay            time.Time `json:"joinPartyDay" update:"join_party_day"`
	JoinPartyPrePeriodStart time.Time `json:"joinPartyPrePeriodStart" update:"join_party_pre_period_start"`
	JoinPartyPrePeriodEnd   time.Time `json:"joinPartyPrePeriodEnd" update:"join_party_pre_period_end"`
	StartJobDay             time.Time `json:"startJobDay"  update:"start_job_day"`
	FullTimeEdu             string    `json:"fullTimeEdu"`
	FullTimeDegree          string    `json:"fullTimeDegree"`
	FullTimeMajor           string    `json:"fullTimeMajor"`
	FullTimeSchool          string    `json:"fullTimeSchool"`
	PartTimeEdu             string    `json:"partTimeEdu"`
	PartTimeDegree          string    `json:"partTimeDegree"`
	PartTimeMajor           string    `json:"partTimeMajor"`
	PartTimeSchool          string    `json:"partTimeSchool"`
	FinalEdu                string    `json:"finalEdu"`
	FinalDegree             string    `json:"finalDegree"`
	FinalMajor              string    `json:"finalMajor"`
	FinalSchool             string    `json:"finalSchool"`
	OrganID                 string    `json:"organId"`
	DepartmentId            string    `json:"departmentId"`
	BePoliceDay             time.Time `json:"bePoliceDay"`
	Training                string    `json:"training"`
	ProCert                 string    `json:"proCert" update:"pro_cert"`
	IsSecret                int8      `json:"isSecret" gorm:"default:0"`
	PassExamDay             time.Time `json:"passExamDay" update:"pass_exam_day"`
	Passport                string    `json:"passport"`
	Phone                   string    `json:"phone"`
	Photo                   string    `json:"photo"`
	Hometown                string    `json:"hometown"`
	Birthplace              string    `json:"birthplace"`
	Health                  string    `json:"health"`
	TechnicalTitle          string    `json:"technicalTitle"`
	Specialty               string    `json:"specialty"`
	Marriage                string    `json:"marriage"`
	Remark                  string    `json:"remark" update:"remark"`
	UserType                int8      `json:"userType"`
	DataStatus              int8      `json:"dataStatus"`
	Sort                    int       `json:"sort"`
	Status                  bool      `json:"status"`
	CreateTime              time.Time `json:"createTime"`
	UpdateTime              time.Time `json:"updateTime"`
}

type Department struct {
	ID          string `json:"id" gorm:"size:50;primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string    `json:"name" gorm:"size:50;not Null"`
	ShortName   string    `json:"shortName" gorm:"size:50"`
	DeptType    int       `json:"deptType"`
	DataStatus  int       `json:"dataStatus"`
	Code        string    `json:"code"`
	LevelCode   string    `json:"levelCode"`
	BusOrgCode  string    `json:"busOrgCode"`
	BusDeptCode string    `json:"busDeptCode"`
	ParentId    string    `json:"parentId"`
	Sort        int       `json:"sort"`
	Headcount   int       `json:"headcount"`
	Position    string    `json:"position"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

type Family struct {
	Base
	PersonnelId int64     `json:"personnelId,string"`
	Name        string    `json:"name"`
	Gender      string    `json:"gender"`
	Relation    string    `json:"relation"`
	Birthday    time.Time `json:"birthday"`
	Organ       string    `json:"organ"`
	Post        string    `json:"post" update:"post"`
	Political   string    `json:"political"`
	IsAbroad    bool      `json:"isAbroad" update:"is_abroad"`
}

// Resume ???????????????
type Resume struct {
	BaseId
	PersonnelId int64  `json:"personnelId,string"`
	Content     string `json:"content"`
}

// Training ???????????????
type Training struct {
	Base
	Title      string        `json:"title" gorm:"size:1000"`
	Intro      datatype.Clob `json:"intro"`
	Place      string        `json:"place"`
	Sponsor    string        `json:"sponsor"`   //????????????
	Organizer  string        `json:"organizer"` //????????????
	Category   int8          `json:"category"`
	Method     int8          `json:"method"`                   //????????????
	IsInner    int8          `json:"isInner" gorm:"default:2"` //????????????????????????
	IsFullTime int8          `json:"isFullTime"`               //????????????
	Period     int16         `json:"period"`                   //??????
	StartTime  time.Time     `json:"startTime"`
	EndTime    time.Time     `json:"endTime"`
}

// PersonTrain ???????????????
type PersonTrain struct {
	BaseId
	PersonnelId int64 `json:"personnelId,string"`
	TrainId     int64 `json:"trainId,string"`
}

// Appraisal ???????????????
type Appraisal struct {
	Base
	OrganId     string `json:"organId" gorm:"size:50"` //TODO: ???????????????
	Organ       string `json:"organ"`
	PersonnelId int64  `json:"personnelId,string"`
	Years       string `json:"years" gorm:"size:10"`
	Season      int8   `json:"season"`
	Conclusion  string `json:"conclusion"`
}

// Post ???????????????
type Post struct {
	Base
	Department  string    `json:"department" gorm:"size:200"`
	Organ       string    `json:"organ" gorm:"size:200"`
	StartDay    time.Time `json:"startDay"`
	EndDay      time.Time `json:"endDay" update:"end_day"`
	PositionId  int64     `json:"positionId,string"`
	LevelId     int64     `json:"levelId,string"`
	PersonnelId int64     `json:"personnelId,string" gorm:"size:50"`
}

// Position ???????????????
type Position struct {
	Base
	Name     string `json:"name"`
	IsLeader int8   `json:"isLeader"`
	IsChief  int8   `json:"isChief"`
	LevelId  int64  `json:"levelId,string"`
}

// Level ????????????????????????????????????????????????????????????
type Level struct {
	Base
	Name   string `json:"name"`
	Orders int    `json:"order" gorm:"not Null"`
}

// Appointment ???????????????
type Appointment struct {
	Base
	PersonnelId int64  `json:"personnelId,string"`
	AppointPost string `json:"appointPost"` //????????????
	RemovePost  string `json:"removePost"`  //????????????
	Reason      string `json:"reason"`      //????????????
	Opinion     string `json:"opinion"`     //??????????????????
}

// Award ???????????????
type Award struct {
	Base
	PersonnelId int64     `json:"personnelId,string" gorm:"size:50"`
	Category    int8      `json:"category"`
	GetTime     time.Time `json:"getTime"`
	Grade       int8      `json:"grade"`
	Content     string    `json:"content"`
	DocNumber   string    `json:"docNumber"`
	Organ       string    `json:"organ"`
}

// Punish ???????????????
type Punish struct {
	Base
	PersonnelId int64     `json:"personnelId,string" gorm:"size:50"`
	Category    int8      `json:"category"`
	GetTime     time.Time `json:"getTime"`
	Grade       int8      `json:"grade"`
	Content     string    `json:"content"`
	DocNumber   string    `json:"docNumber"`
	Organ       string    `json:"organ"`
}

type Module struct {
	Base
	Name      string `json:"name"`
	Title     string `json:"title"`
	Paths     string `json:"path"`
	Param     string `json:"param"`
	Rank      int8   `json:"rank"`
	Component string `json:"component"`
	Redirect  string `json:"redirect"`
	Icon      string `json:"icon"`
	Parent    int64  `json:"parent,string"`
	Orders    int8   `json:"order"`
}

// RoleDict ????????????
type RoleDict struct {
	Base
	Name  string `json:"name" gorm:"size:100"`  //?????????
	Title string `json:"title" gorm:"size:100"` //?????????
}

// PermissionDict ????????????
type PermissionDict struct {
	Base
	Name  string `json:"name" gorm:"size:20"`  //?????????
	Title string `json:"title" gorm:"size:20"` //?????????
}

// Discipline ???????????????
type Discipline struct {
	Base
	PersonnelId int64     `json:"personnelId,string"`
	Category    int8      `json:"category"`
	GetTime     time.Time `json:"getTime"`
	DictId      int64     `json:"dictId,string"`
	Content     string    `json:"content"`
	DocNumber   string    `json:"docNumber"`
	Deadline    time.Time `json:"deadline"`
	Organ       string    `json:"organ"`
}

// DisDict ?????????????????????
type DisDict struct {
	Base
	Name     string `json:"name" gorm:"size:20"`
	Category int8   `json:"category"`
	Term     int16  `json:"term"`
}

// EduDict ???????????????
type EduDict struct {
	Base
	Name     string `json:"name"`
	Category int8   `json:"category"`
	Sort     int16  `json:"sort"`
}

type Report struct {
	Base
	Title      string        `json:"title"`
	Step       int8          `json:"step"`
	ReportTime time.Time     `json:"reportTime"`
	Intro      datatype.Clob `json:"intro"`
	Steps      datatype.Clob `json:"steps"`
}

type PersonReport struct {
	BaseId
	PersonnelId int64 `json:"personnelId,string"`
	ReportId    int64 `json:"reportId,string"`
}

//type ReportStep struct {
//	BaseId
//	ReportId int64         `json:"reportId"`
//	Step     int8          `json:"step"`
//	StepTime time.Time     `json:"stepTime"`
//	Content  datatype.Clob `json:"content"`
//}

// EntryExit ????????????
type EntryExit struct {
	Base
	PersonnelId int64     `json:"personnelId,string"`
	Passport    int8      `json:"passport"`
	EnterTime   time.Time `json:"enterTime"`   //????????????
	ExitTime    time.Time `json:"exitTime"`    //????????????
	Destination string    `json:"destination"` //?????????
	Aim         string    `json:"aim"`         //????????????
	IsReport    int8      `json:"isReport"`    //????????????
}

// Affair ???????????????????????????
type Affair struct {
	Base
	PersonnelId int64         `json:"personnelId,string"`
	Title       string        `json:"title"`
	Category    int8          `json:"category"`
	Intro       datatype.Clob `json:"intro"`
}

// Talent ?????????
type Talent struct {
	Base
	PersonnelId   int64     `json:"personnelId,string"`
	Category      int8      `json:"category"`
	Skill         string    `json:"skill"`
	BeExaminerDay time.Time `json:"beExaminerDay"`
}

type TalentPick struct {
	Base
	PickerId int64         `json:"pickerId,string"`
	Category int8          `json:"category"`
	Title    string        `json:"title"`
	PickDate time.Time     `json:"pickDate"`
	Res      datatype.Clob `json:"res"`
}

// Custom ????????????????????????
type Custom struct {
	Base
	Name      string `json:"name"`
	AccountId string `json:"accountId"`
	Category  int8   `json:"category"`
	Content   string `json:"content"`
}

// Account ?????????
type Account struct {
	ID           string    `json:"id" gorm:"autoIncrement:false;primaryKey"`
	PersonnelId  int64     `json:"personnelId,string"`
	Name         string    `json:"name" gorm:"size:50"`
	Username     string    `json:"username"`
	IdCode       string    `json:"idCode" gorm:"size:50"`
	DepartmentId string    `json:"departmentId"`
	OrganID      string    `json:"organId"`
	UserType     int8      `json:"userType"`
	DataStatus   int8      `json:"dataStatus"`
	Sort         int       `json:"sort"`
	CreateTime   time.Time `json:"createTime"`
	UpdateTime   time.Time `json:"updateTime"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TextSize struct {
	Base
	Name datatype.Clob `json:"name" gorm:"size:8000"`
}

// Review ???????????????
type Review struct {
	Base
	PersonnelId int64         `json:"personnelId,string"`
	OrganID     string        `json:"organId"`
	Category    int8          `json:"category"`
	Reviewer    int64         `json:"reviewer,string"`
	Status      int8          `json:"status" gorm:"default:1"`
	Content     datatype.Clob `json:"content"`
}

// Feedback ???????????????
type Feedback struct {
	Base
	PersonnelId int64         `json:"personnelId,string"`
	OrganID     string        `json:"organId"`
	Category    int8          `json:"category"`
	Feedbacker  int64         `json:"feedbacker,string"`
	Status      int8          `json:"status" gorm:"default:1"`
	Content     datatype.Clob `json:"content"`
}

type Log struct {
	Base
	Category  int8          `json:"category"`
	IP        string        `json:"ip"`
	AccountId string        `json:"accountId"`
	Content   datatype.Clob `json:"content"`
}
