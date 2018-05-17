#if_new

SetSprite(GetId(), 'NESW', 261)
SetActive("pineapple", 4)
SetFlag(GetId(), "sun", "false")

#always

if Proximity(GetId(), GetFocus()) < 3 then
    SetSprite(GetId(), 'NESW', 379)
    if GetFlag(GetId(), "sun") == "false" then
        PlaySound("on.wav")
        SetFlag(GetId(), "sun", "true")
    end
else
    SetSprite(GetId(), 'NESW', 261)
    if GetFlag(GetId(), "sun") == "true" then
        PlaySound("off.wav")
        SetFlag(GetId(), "sun", "false")
    end
end

if Proximity(GetId(), GetFocus()) < 1 then

    local oldscore = GetFlag(GetFocus(), "score")
    if oldscore ~= "" then
        local newscore = oldscore + 1
        print (newscore)
        SetFlag(GetFocus(), "score", newscore)
        PlaySound("pickup.wav")
        Delete(GetId())
    end

end
