package controllers

import "time"

type LogCategory int8

const (
	DangerUpdate LogCategory = iota + 1
	DangerQuery
)
const (
	sortStr = `(case when length(d.level_code)>=3 then (select ii.sort from departments ii where ii.level_code = substring(d.level_code,1,3)) else null end) desc,
(case when length(d.level_code)>=6 then (select ii.sort from departments ii where ii.level_code = substring(d.level_code,1,6)) else null end) desc,
(case when length(d.level_code)>=9 then (select ii.sort from departments ii where ii.level_code = substring(d.level_code,1,9)) else null end) desc,
(case when length(d.level_code)>=12 then (select ii.sort from departments ii where ii.level_code = substring(d.level_code,1,12)) else null end) desc,
(case when length(d.level_code)>=15 then (select ii.sort from departments ii where ii.level_code = substring(d.level_code,1,15)) else null end) desc, 
personnels.sort desc nulls first`

	RbacModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`
)

var (
	LongLongDay = time.Date(2999, 12, 31, 0, 0, 0, 0, time.UTC)
)
