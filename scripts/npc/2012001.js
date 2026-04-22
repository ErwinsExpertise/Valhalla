var TICKET_ID = 4031047;

if (!npc.sendYesNo("It seems like there is still room on this ride. Please have your ticket ready so you can get on. The journey may be long, but you will get to your destination safely. What do you think? Do you want to go on this trip?")) {
    npc.sendOk("You must have some business to take care of here, right?");
} else if (plr.haveItem(TICKET_ID, 1)) {
    plr.gainItem(TICKET_ID, -1);
    plr.warp(101000300);
} else {
    npc.sendOk("Make sure you have a Ellinia ticket to travel on this boat.");
}
