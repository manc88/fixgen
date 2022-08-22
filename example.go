package main

type Driver struct {
	name string
}

func (d *Driver) String() string {
	return d.name + "driver"
}

type Car struct {
	speed  int
	price  int64
	model  string
	pass   []string
	props  map[int]string
	driver Driver
}
