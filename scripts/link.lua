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

#always

if KeyPressed("up", false) then
    SetVelocity(GetId(), 'N', 5, 1)
end

if KeyPressed("down", false) then
    SetVelocity(GetId(), 'S', 5, 1)
end

if KeyPressed("left", false) then
    SetVelocity(GetId(), 'W', 5, 1)
end

if KeyPressed("right", false) then
    SetVelocity(GetId(), 'E', 5, 1)
end