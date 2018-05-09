#if_new

SetFocus(GetId(), true)

SetSprite(GetId(), 'N', 96)
Animate(GetId(), 'N', 97, 106, 15, false)
SetSprite(GetId(), 'E', 112)
Animate(GetId(), 'E', 113, 122, 15, false)
SetSprite(GetId(), 'S', 64)
Animate(GetId(), 'S', 65, 74, 15, false)
SetSprite(GetId(), 'W', 80)
Animate(GetId(), 'W', 81, 90, 15, false)

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


if KeyPressed("p", true) then
    print("Becoming pac man!...")
    SetClass(GetId(), "pacman")
    Reset(GetId())
end