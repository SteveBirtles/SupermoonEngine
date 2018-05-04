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

if KeyPressed("space", true) then
    SetPosition(GetId(), 0, 0, 0)
end

if KeyPressed("c", true) then
    local a, b, c = GetPosition(GetId())
    Create(a + 2, b, c, "test")
end

if KeyPressed("s", true) then
    local s = GetScript(GetId())
    print(s)
end

if KeyPressed("k", true) then
    Delete(1)
end

if KeyPressed("m", true) then
    SetClassActive("test", 5)
end

if KeyPressed("enter", true) then
    local ids = Nearby(GetId(), 10)
    print("Entities:")
    for i = 1, #ids do
        print (ids[i])
    end
end

local direction, velocity, distance = GetVelocity(GetId())

if direction == 'N' then
    if distance == 0 then
        SetSprite(GetId(), 12)
    else
        Animate(GetId(), 40, 49, 15)
    end
elseif direction == 'E' then
    if distance == 0 then
        SetSprite(GetId(), 13)
    else
        Animate(GetId(), 50, 59, 15)
    end
elseif direction == 'S' then
    if distance == 0 then
        SetSprite(GetId(), 10)
    else
        Animate(GetId(), 20, 29, 15)
    end
elseif direction == 'W' then
    if distance == 0 then
        SetSprite(GetId(), 11)
    else
        Animate(GetId(), 30, 39, 15)
    end
end