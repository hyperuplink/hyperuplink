package data

import "github.com/gofiber/fiber/v3"

type Data struct {
	c fiber.Ctx

	store map[string]interface{}
}

func New(c fiber.Ctx) (f *Data) {
	f = new(Data)
	f.c = c

	f.store = make(map[string]interface{})

	return f
}

func (d *Data) Set(key string, val interface{}) {
	d.store[key] = val
}

func (d *Data) Get(key string) (val interface{}) {
	return d.store[key]
}
