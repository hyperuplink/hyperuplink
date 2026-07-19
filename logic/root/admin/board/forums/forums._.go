package forums

import (
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
)

type CategoryWithForums struct {
	Category category.Category `json:"category"`
	Forums   []forum.Forum     `json:"forums"`
}

type CreateInput struct {
	Name       string `json:"name" form:"name" validate:"required,min=1,max=32"`
	Slug       string `json:"slug" form:"slug" validate:"required,slug,min=1,max=32"`
	CategoryID string `json:"category_id" form:"category_id" validate:"required,uuid"`
}

type UpdateInput struct {
	ID          string `json:"id" form:"id" validate:"omitempty,uuid"`
	Name        string `json:"name" form:"name" validate:"required,min=1,max=32"`
	Slug        string `json:"slug" form:"slug" validate:"required,slug,min=1,max=32"`
	CategoryID  string `json:"category_id" form:"category_id" validate:"required,uuid"`
	Description string `json:"description" form:"description" validate:"min=0,max=128"`
}

type DestroyInput struct {
	ID string `json:"id" form:"id" validate:"required,uuid"`
}

type MoveInput struct {
	ID string `json:"id" form:"id" validate:"required,uuid"`
}
