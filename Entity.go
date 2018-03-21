package main

type Entity struct {

  id uint32
  active bool

  x float64   // position
  y float64
  z float64

  dx float64  // velocity
  dy float64
  dz float64
  dn int      // number of squares to continue at that velocity

  class string  // corrisponds to a lua file
  script string // their lua script
  properties map[string]string   // entity properties map
  timers map[string]uint16

  sprite int
  animated bool
  startSprite int
  endSprite int
  animationSpeed float64

}
