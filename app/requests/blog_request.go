package requests

type CreateBlogRequest struct {
	Title    string `json:"title" validate:"required"`
	Content  string `json:"content" validate:"required"`
	AuthorID int    `json:"author_id" validate:"required"`
}

type UpdateBlogRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}
