var maps = [251000000, 200000100];
var mapNames = ["Herb Town", "Orbis"];

if (plr.mapID() === 251000000) {
    maps = [250000100];
    mapNames = ["Mu Lung"];
}

var text = "Where do you want to go today?";
for (var i = 0; i < maps.length; i++) {
    text += "\r\n#L" + i + "# " + mapNames[i] + "#l";
}

npc.sendSelection(text);
var selection = npc.selection();

if (selection >= 0 && selection < maps.length) {
    npc.sendNext("All right, see you next time.");
    plr.warp(maps[selection]);
}
