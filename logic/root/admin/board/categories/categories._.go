package categories

type CreateInput struct {
	Name string `json:"name" form:"name" validate:"required,min=1,max=32"`
	Slug string `json:"slug" form:"slug" validate:"required,slug,min=1,max=32"`
}

type UpdateInput struct {
	ID   string `json:"id" form:"id" validate:"omitempty,uuid"`
	Name string `json:"name" form:"name" validate:"required,min=1,max=32"`
	Slug string `json:"slug" form:"slug" validate:"required,slug,min=1,max=32"`
}

type DestroyInput struct {
	ID string `json:"id" form:"id" validate:"required,uuid"`
}

type MoveInput struct {
	ID string `json:"id" form:"id" validate:"required,uuid"`
}
