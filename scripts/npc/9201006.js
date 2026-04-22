var mapId = plr.mapID();

if (mapId === 680000210) {
    npc.sendSelection("\r\n#L0##bWhen does the wedding begin?#l\r\n#L1#I want out!#l");
    var selection = npc.selection();

    if (selection === 0) {
        npc.sendOk("We will wait until the bride and groom are ready. Please wait a few minutes!");
    } else if (selection === 1) {
        plr.removeAll(5251100);
        plr.warp(680000000);
    } else {
        npc.sendOk("Bye");
    }
} else if (mapId === 680000200) {
    npc.sendOk("Uhh, sorry for the delay. Father John went to do something very quickly. It should not be long; please wait for us to start.");
} else {
    npc.sendOk("Bye");
}
