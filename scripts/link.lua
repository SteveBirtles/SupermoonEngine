#if_new

SetFocus(GetId(), true)

SetSprite(GetId(), 'N', 12)
Animate(GetId(), 'N', 40, 49, 15, false)
SetSprite(GetId(), 'E', 13)
Animate(GetId(), 'E', 50, 59, 15, false)
SetSprite(GetId(), 'S', 10)
Animate(GetId(), 'S', 20, 29, 15, false)
SetSprite(GetId(), 'W', 11)
Animate(GetId(), 'W', 30, 39, 15, false)

#if_focus

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