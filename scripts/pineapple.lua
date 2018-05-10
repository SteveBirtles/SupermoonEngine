#if_new

SetSprite(GetId(), 'NESW', 261)
SetActive("pineapple", 4)

#always

if Proximity(GetId(), GetFocus()) < 3 then
    SetSprite(GetId(), 'NESW', 379)
else
    SetSprite(GetId(), 'NESW', 261)
end

if Proximity(GetId(), GetFocus()) < 1 then

    local oldscore = GetFlag(GetFocus(), "score")
    if oldscore ~= "" then
        local newscore = oldscore + 1
        print (newscore)
        SetFlag(GetFocus(), "score", newscore)
        Delete(GetId())
    end

end
