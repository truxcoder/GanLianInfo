package controllers

import (
	"GanLianInfo/utils"
	"slices"
	"strings"
	"time"
)

type HighLevelPolice struct {
	ID         int64     `json:"id,string"`
	Name       string    `json:"name"`
	Birthday   time.Time `json:"birthday"`
	RetireDay  time.Time `json:"retireDay"`
	OrganID    string    `json:"organId"`
	Gender     string    `json:"gender"`
	G4StartDay time.Time `json:"g4StartDay"`
	G3StartDay time.Time `json:"g3StartDay"`
	G2StartDay time.Time `json:"g2StartDay"`
	G1StartDay time.Time `json:"g1StartDay"`
}
type posStruct struct {
	ID             int64               `json:"id"`
	Name           string              `json:"name"`
	Birthday       time.Time           `json:"birthday"`
	OrganID        string              `json:"organId"`
	Gender         string              `json:"gender"`
	CurrentRank    int64               `json:"currentRank"`
	RankName       string              `json:"rankName"`
	CurrentLevel   int64               `json:"currentLevel"`
	LevelName      string              `json:"levelName"`
	RankStartDay   time.Time           `json:"rankStartDay"`
	LevelStartDay  time.Time           `json:"levelStartDay"`
	FkStartDay     time.Time           `json:"fkStartDay"`
	ZkStartDay     time.Time           `json:"zkStartDay"`
	FcStartDay     time.Time           `json:"fcStartDay"`
	ZcStartDay     time.Time           `json:"zcStartDay"`
	RetireDay      time.Time           `json:"retireDay" gorm:"-"`
	LeaderPosts    []PostLevelPosition `json:"leaderPosts" gorm:"-"`
	NonLeaderPosts []PostLevelPosition `json:"nonLeaderPosts" gorm:"-"`
}

type PostLevelPosition struct {
	Id           int64     `json:"id"`
	PersonnelId  int64     `json:"personnelId"`
	StartDay     time.Time `json:"startDay"`
	EndDay       time.Time `json:"endDay"`
	IsLeader     int8      `json:"isLeader"`
	LevelName    string    `json:"levelName"`
	PositionName string    `json:"positionName"`
	Orders       int       `json:"orders"`
}

// 获取一高开始日期
func getG1StartDay(baseDay time.Time, p *posStruct) time.Time {
	var lst []time.Time
	// 条件一：任二高满三年，现曾任场所正处实职满三年
	if !p.ZcStartDay.IsZero() {
		d1 := utils.GetLatestDay(baseDay, getAttainTime(p.LeaderPosts, "正处级", 3))
		lst = append(lst, d1)
	}
	// 条件二：任二高满三年，现曾任场所副处实职满10年，年龄大于53岁
	// 条件三：任二高满三年，现曾任场所副处实职满8年，年龄大于55岁
	if !p.FcStartDay.IsZero() {
		d2 := utils.GetLatestDay(baseDay, getAttainTime(p.LeaderPosts, "正副处", 10), p.Birthday.AddDate(53, 0, 0))
		d3 := utils.GetLatestDay(baseDay, getAttainTime(p.LeaderPosts, "正副处", 8), p.Birthday.AddDate(55, 0, 0))
		lst = append(lst, d2, d3)
	}
	// 条件四：任二高满三年，离退休小于1年
	d4 := utils.GetLatestDay(baseDay, p.RetireDay.AddDate(-1, 0, 0))
	lst = append(lst, d4)

	return utils.GetEarliestDay(lst...)
}

// 获取二高开始日期
func getG2StartDay(baseDay time.Time, p *posStruct) time.Time {
	// 条件一：任三高满两年，现曾任场所副处实职满5年，年龄大于48岁
	// 条件二：任三高满两年，原副处满8年，离退休小于5年
	// 条件三：任三高满两年，曾任正科实职大于15年，离退休小于3年
	// 条件四：任三高满两年，曾任正科实职，任正副科大于20年，离退休小于3年
	// 条件五：任三高满两年，离退休小于1年
	var lst []time.Time
	if !p.FcStartDay.IsZero() {
		d1 := utils.GetLatestDay(baseDay, getAttainTime(p.LeaderPosts, "正副处", 5), p.Birthday.AddDate(48, 0, 0))

		lst = append(lst, d1)
	}
	d2 := utils.GetLatestDay(baseDay, getNonLeaderAttainTime(p.NonLeaderPosts, "正副处", 8), p.RetireDay.AddDate(-5, 0, 0))
	lst = append(lst, d2)
	if !p.ZkStartDay.IsZero() {
		// FixMe: 此处要确认是不是可以将正科级以上的经历算进去。
		d3 := utils.GetLatestDay(baseDay, getAttainTime(p.LeaderPosts, "正科级", 15), p.RetireDay.AddDate(-3, 0, 0))
		d4 := utils.GetLatestDay(baseDay, getAttainTime(p.LeaderPosts, "正副科", 20), p.RetireDay.AddDate(-3, 0, 0))
		lst = append(lst, d3, d4)
	}
	d5 := utils.GetLatestDay(baseDay, p.RetireDay.AddDate(-1, 0, 0))
	lst = append(lst, d5)

	return utils.GetEarliestDay(lst...)
}

// 获取三高开始日期
func getG3StartDay(baseDay time.Time, p *posStruct) time.Time {
	// 条件一：任四高满两年，现曾任场所副处实职满2年
	// 条件二：任四高满两年，原副处满5年，年龄大于50岁
	// 条件三：任四高满两年，现任正科实职,任正科大于10年，年龄大于50岁
	// 条件四：任四高满两年，现任正科实职，任正副科大于15年，年龄大于50岁
	// 条件五：任四高满两年，曾任正副科实职大于15年，离退休小于3年
	// 条件六：任四高满两年，曾任正副科实职大于20年，离退休小于5年
	// 条件七：任四高满两年，离退休小于1年
	var d1, d2, d3, d4, d5, d6, d7 time.Time
	var lst []time.Time
	if !p.FcStartDay.IsZero() {
		d1 = utils.GetLatestDay(baseDay, getAttainTime(p.LeaderPosts, "副处级", 2))

		lst = append(lst, d1)
		//log.Successf("getNonLeaderAttainTime=%v\n", getNonLeaderAttainTime(p.NonLeaderPosts, "正副处", 5))
	}
	d2 = utils.GetLatestDay(baseDay, getNonLeaderAttainTime(p.NonLeaderPosts, "正副处", 5), p.Birthday.AddDate(50, 0, 0))
	lst = append(lst, d2)
	if p.LevelName == "正科级" {
		d3 = utils.GetLatestDay(baseDay, getAttainTime(p.LeaderPosts, "正科级", 10), p.Birthday.AddDate(50, 0, 0))
		d4 = utils.GetLatestDay(baseDay, getAttainTime(p.LeaderPosts, "正副科", 15), p.Birthday.AddDate(50, 0, 0))
		lst = append(lst, d3, d4)
	}
	if !p.FkStartDay.IsZero() {
		d5 = utils.GetLatestDay(baseDay, getAttainTime(p.LeaderPosts, "正副科", 15), p.RetireDay.AddDate(-3, 0, 0))
		d6 = utils.GetLatestDay(baseDay, getAttainTime(p.LeaderPosts, "正副科", 20), p.RetireDay.AddDate(-5, 0, 0))
		lst = append(lst, d5, d6)
	}
	d7 = utils.GetLatestDay(baseDay, p.RetireDay.AddDate(-1, 0, 0))
	lst = append(lst, d7)
	return utils.GetEarliestDay(lst...)
}

// 获取四高开始日期
func getG4StartDay(baseDay time.Time, p *posStruct) time.Time {
	// 条件一：一级警长满两年，现任正科实职，任正科满5年，年龄大于43岁
	// 条件二：一级警长满两年，现任正科实职，任正副科满10年，年龄大于45岁
	// 条件三：一级警长满两年，现任副科实职，任副科满15年，年龄大于48岁
	// 条件四：一级警长满两年，现任副科实职，任副科满13年，年龄大于50岁
	// 条件五：一级警长满两年，曾任正副科实职，任正副科满15年，离退休小于5年
	// 条件六：一级警长满两年，离退休小于3年
	var d1, d2, d3, d4, d5, d6 time.Time
	var lst []time.Time
	//log.Successf("getNonLeaderAttainTime=%v\n", getNonLeaderAttainTime(p.NonLeaderPosts, "正副处", 5))
	if p.LevelName == "正科级" {
		d1 = utils.GetLatestDay(baseDay, getAttainTime(p.LeaderPosts, "正科级", 5), p.Birthday.AddDate(43, 0, 0))
		d2 = utils.GetLatestDay(baseDay, getAttainTime(p.LeaderPosts, "正副科", 10), p.Birthday.AddDate(45, 0, 0))
		lst = append(lst, d1, d2)
		//if p.ID == 10485224311040 {
		//	log.Successf("%v---%v---%v\n", baseDay.Local(), d1, d2)
		//}
	}
	if p.LevelName == "副科级" {
		d3 = utils.GetLatestDay(baseDay, getAttainTime(p.LeaderPosts, "副科级", 15), p.Birthday.AddDate(48, 0, 0))
		d4 = utils.GetLatestDay(baseDay, getAttainTime(p.LeaderPosts, "副科级", 13), p.Birthday.AddDate(50, 0, 0))
		lst = append(lst, d3, d4)
		//log.Successf("%s:正科5年原%v现%v，正副科10年原%v，现%v\n", p.Name, p.ZkStartDay.AddDate(5, 0, 0), getAttainTime(p.LeaderPosts, "正科级", 5), p.FkStartDay.AddDate(10, 0, 0), getAttainTime(p.LeaderPosts, "正副科", 10))
	}
	if !p.ZkStartDay.IsZero() || !p.FkStartDay.IsZero() {
		d5 = utils.GetLatestDay(baseDay, getAttainTime(p.LeaderPosts, "正副科", 15), p.RetireDay.AddDate(-5, 0, 0))
		lst = append(lst, d5)
	}
	d6 = utils.GetLatestDay(baseDay, p.RetireDay.AddDate(-3, 0, 0))
	lst = append(lst, d6)
	return utils.GetEarliestDay(lst...)
}

// 获取实职达到时间
func getAttainTime(posts []PostLevelPosition, level string, year int) (att time.Time) {
	var sub time.Duration
	var start, end, original time.Time
	var key string
	if strings.Contains(level, "正副") {
		key = strings.Trim(level, "正副")
	}
	for _, post := range posts {
		if (key == "" && post.LevelName == level) || (key != "" && post.LevelName == "正"+key+"级") || (key != "" && post.LevelName == "副"+key+"级") {
			//fmt.Printf("start:%v----end:%v\n", post.Start, post.End)
			if start.IsZero() {
				start = post.StartDay
				end = post.EndDay
				original = start.AddDate(year, 0, 0)
				if end.IsZero() {
					break
				}
				//fmt.Printf("start.IsZero(), then: start = %v, end = %v, original = %v\n", start, end, original)
			} else if post.StartDay.After(end) {
				sub = sub + post.StartDay.Sub(end)
				end = post.EndDay
				//fmt.Printf("post.Start.After(end) && !end.IsZero(), sub = %v天\n", sub.Hours()/24)
				if end.IsZero() {
					break
				}
			} else if post.EndDay.After(end) {
				end = post.EndDay
				//fmt.Printf("post.End.After(end) && !end.IsZero()， end = %v\n", end)
			}
			if post.EndDay.IsZero() {
				end = post.EndDay
				break
			}
		}
	}
	att = original.Add(sub)
	if !end.IsZero() && att.After(end) {
		//att = time.Time{}
		att = LongLongDay
	}
	return att
}

// 获取非实职达到时间
func getNonLeaderAttainTime(posts []PostLevelPosition, level string, year int) (att time.Time) {
	var sub time.Duration
	var start, end, original time.Time
	var key string
	var positions []string
	if level == "正处级" {
		positions = append(positions, "调研员", "一级警长")
	}
	if level == "副处级" {
		positions = append(positions, "副调研员", "二级警长", "助理调研员")
	}
	if level == "正副处" {
		positions = append(positions, "副调研员", "二级警长", "调研员", "一级警长", "助理调研员")
	}
	if strings.Contains(level, "正副") {
		key = strings.Trim(level, "正副")
	}
	for _, post := range posts {
		if key == "" && post.LevelName == level && slices.Contains(positions, post.PositionName) || key != "" && (post.LevelName == "正"+key+"级" || post.LevelName == "副"+key+"级") && slices.Contains(positions, post.PositionName) {
			if start.IsZero() {
				start = post.StartDay
				end = post.EndDay
				original = start.AddDate(year, 0, 0)
				if end.IsZero() {
					break
				}
				//fmt.Printf("start.IsZero(), then: start = %v, end = %v, original = %v\n", start, end, original)
			} else if post.StartDay.After(end) {
				sub = sub + post.StartDay.Sub(end)
				end = post.EndDay
				//fmt.Printf("post.Start.After(end) && !end.IsZero(), sub = %v天\n", sub.Hours()/24)
				if end.IsZero() {
					break
				}
			} else if post.EndDay.After(end) {
				end = post.EndDay
				//fmt.Printf("post.End.After(end) && !end.IsZero()， end = %v\n", end)
			}
			if post.EndDay.IsZero() {
				end = post.EndDay
				break
			}
		}
	}
	if original.IsZero() {
		return LongLongDay
	}
	att = original.Add(sub)
	if !end.IsZero() && att.After(end) {
		//att = time.Time{}
		att = LongLongDay
	}
	return att
}
