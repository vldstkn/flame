package models

type Swipe struct {
	UserId1      int64 `db:"user_id1"`
	UserId2      int64 `db:"user_id2"`
	UserIsLiked1 *bool `db:"user_is_liked1"`
	UserIsLiked2 *bool `db:"user_is_liked2"`
}
