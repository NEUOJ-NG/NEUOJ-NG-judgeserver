package form

type InternalError struct {
	Description  string `form:"description" binding:"required"`
	JudgehostLog string `form:"judgehostlog" binding:"required"`
	Disabled     string `form:"disabled" binding:"required"`
	CID          int    `form:"cid"`
	JudgingID    int    `form:"judgingid"`
}
