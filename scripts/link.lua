#if_new

print("Becoming link!...")

SetFocus(GetId(), true)

SetSprite(GetId(), 'N', 96)
Animate(GetId(), 'N', 97, 106, 15, false)
SetSprite(GetId(), 'E', 112)
Animate(GetId(), 'E', 113, 122, 15, false)
SetSprite(GetId(), 'S', 64)
Animate(GetId(), 'S', 65, 74, 15, false)
SetSprite(GetId(), 'W', 80)
Animate(GetId(), 'W', 81, 90, 15, false)

SetFlag(GetId(), "score", "0")

#always

local speed = 5
local slow = KeyPressed("leftshift", false)

if slow then
    speed = 2
end

if KeyPressed("up", slow) then
    SetVelocity(GetId(), 'N', speed, 1)
end

if KeyPressed("down", slow) then
    SetVelocity(GetId(), 'S', speed, 1)
end

if KeyPressed("left", slow) then
    SetVelocity(GetId(), 'W', speed, 1)
end

if KeyPressed("right", slow) then
    SetVelocity(GetId(), 'E', speed, 1)
end

if KeyPressed("j", true) then
    PlaySound("jump.wav")
end

if KeyPressed("b", true) then
    PlayMusic("bennyhill.mp3")
end

if KeyPressed("p", true) then
    SetClass(GetId(), "pacman")
    Reset(GetId())
end

if KeyPressed("c", true) then
  local x, y, z = GetPosition(GetId())
  Create(x+1, y, z, "pacman")
end
