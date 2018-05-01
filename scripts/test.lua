SetFocus(GetId(), true)

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

--if KeyPressed("Space", true) then
--    local x = GetFlag(GetId(), "x")
--    x = x + 1
--    SetFlag(GetId(), "x", x)
--    print ("Space JUST PRESSED (" .. x .. ")")
--elseif KeyPressed("Space", false) then
--    print ("Space")
--end
