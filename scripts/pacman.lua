#if_new

SetFocus(GetId(), true)

SetSprite(GetId(), 'N', 32)
Animate(GetId(), 'N', 32, 33, 8, false)
SetSprite(GetId(), 'E', 0)
Animate(GetId(), 'E', 0, 1, 8, false)
SetSprite(GetId(), 'S', 48)
Animate(GetId(), 'S', 48, 49, 8, false)
SetSprite(GetId(), 'W', 16)
Animate(GetId(), 'W', 16, 17, 8, false)

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