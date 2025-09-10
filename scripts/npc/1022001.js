// Regular Cab (Perion) stateless refactor
var maps = [104000000, 101000000, 103000000, 100000000];
var prices = [1200, 1000, 800, 1000];

var isBeginner = (plr.job() === 0);

npc.sendNext("Hi! I drive the Regular Cab. If you want to go from town to town safely and fast, then ride our cab. We'll gladly take you to your destination with an affordable price.");

// build destination menu
var text = "";
if (isBeginner) {
    text += "We have a special 90% discount for beginners. ";
}
text += "Choose your destination, for fees will change from place to place.#b";

for (var i = 0; i < maps.length; i++) {
    var price = isBeginner ? Math.floor(prices[i] / 10) : prices[i];
    text += "\r\n#L" + i + "##m" + maps[i] + "# (" + price + " mesos)#l";
}
npc.sendSelection(text);

var selection = npc.selection();
if (selection < 0 || selection >= maps.length) {
    selection = 0;
}

var map = maps[selection];
var price = isBeginner ? Math.floor(prices[selection] / 10) : prices[selection];

if (npc.sendYesNo("You don't have anything else to do here, huh? Do you really want to go to #b#m" + map + "##k? It'll cost you #b" + price + " mesos#k.")) {
    if (plr.mesos() >= price) {
        plr.giveMesos(-price);
        plr.warp(map);
    } else {
        npc.sendNext("You don't have enough mesos. Sorry to say this, but without them, you won't be able to ride the cab.");
    }
} else {
    npc.sendNext("There's a lot to see in this town, too. Come back and find us when you need to go to a different town.");
}