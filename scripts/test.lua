if KeyPressed("Space", true) then
    local x = GetFlag(GetId(), "x")
    x = x + 1
    SetFlag(GetId(), "x", x)
    print ("Space " .. x)
elseif KeyPressed("Space", false) then
    print ("Space")
end
