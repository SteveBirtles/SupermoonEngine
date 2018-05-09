#if_new

print("Becoming pac man!...")

SetSprite(GetId(), 'N', 32)
Animate(GetId(), 'N', 32, 33, 8, false)
SetSprite(GetId(), 'E', 0)
Animate(GetId(), 'E', 0, 1, 8, false)
SetSprite(GetId(), 'S', 48)
Animate(GetId(), 'S', 48, 49, 8, false)
SetSprite(GetId(), 'W', 16)
Animate(GetId(), 'W', 16, 17, 8, false)

#always

local speed = 8

if KeyPressed("up", false) then
    SetVelocity(GetId(), 'N', speed, 1)
end

if KeyPressed("down", false) then
    SetVelocity(GetId(), 'S', speed, 1)
end

if KeyPressed("left", false) then
    SetVelocity(GetId(), 'W', speed, 1)
end

if KeyPressed("right", false) then
    SetVelocity(GetId(), 'E', speed, 1)
end

if KeyPressed("l", true) then
    SetClass(GetId(), "link")
    Reset(GetId())
end
