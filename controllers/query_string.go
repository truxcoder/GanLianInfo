package controllers

var (
	willUpInThreeMonth = "AND exists (select 1 from posts where posts.personnel_id = personnels.id and exists(select 1 from levels where posts.level_id = levels.id and levels.orders > 3) and exists(select 1 from positions where positions.id = posts.position_id and is_leader = 1 and positions.name <> '一级警长')  and ((posts.end_day = ?  and posts.start_day <= ?)  or (posts.end_day <> ?  and posts.start_day<= ? and exists (select *  from posts as po  where po.position_id = posts.position_id and po.personnel_id = posts.personnel_id and po.end_day = ?))))"
)
