package pattern
/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

type House struct {
  glassAmount int
  height int
  rooms int
}

type HouseBuilderI interface {
  GlassAmount(g int) HouseBuilderI
  Height(h int) HouseBuilderI
  Rooms(r int) HouseBuilderI

  Build() struct {glassAmount, height, rooms int}
}

type HouseBuilderS struct {
  glassAmount, height, rooms int
}

func (b *HouseBuilderS)  GlassAmount(g int) HouseBuilderI {
  b.glassAmount = g
  return b
}

func (b *HouseBuilderS)  Height(h int) HouseBuilderI {
  b.height = h
  return b
}

func (b *HouseBuilderS)  Rooms(r int) HouseBuilderI {
  b.rooms = r
  return b
}

func (b *HouseBuilderS) Build() struct {glassAmount, height, rooms int} {
  v := struct {
    glassAmount, height, rooms int
  }{b.glassAmount, b.height, b.rooms}
  return v 
}

