package permissions

import (
	"xn--gckvb8fzb.com/hyperuplink/models/category"
)

type PermissionRow struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	Level        string `json:"level"`
}

type GroupView struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Rows        []PermissionRow     `json:"rows"`
	AddableCats []category.Category `json:"addable_categories"`
}

type View struct {
	DefaultLevel  string      `json:"default_level"`
	Groups        []GroupView `json:"groups"`
	HasCategories bool        `json:"has_categories"`
}

type ApplyInput struct {
	GroupID    string `json:"group_id" form:"group_id" validate:"omitempty,slug,max=32"`
	CategoryID string `json:"category_id" form:"category_id" validate:"omitempty,uuid"`
	Level      string `json:"level" form:"level" validate:"required,oneof=none read read_write read_write_moderate"`
}

type GroupCreateInput struct {
	ID   string `json:"id" form:"id" validate:"required,slug,min=1,max=32"`
	Name string `json:"name" form:"name" validate:"required,min=1,max=32"`
}

type GroupDestroyInput struct {
	ID string `json:"id" form:"id" validate:"required,slug,min=1,max=32"`
}

type RemoveInput struct {
	GroupID    string `json:"group_id" form:"group_id" validate:"required,slug,max=32"`
	CategoryID string `json:"category_id" form:"category_id" validate:"required,uuid"`
}
