var TICKET_ID = 4031045;

if (!npc.sendYesNo("We are preparing to depart for Orbis. If you need to take care of some things, I suggest you do that first before getting on board. Do you still wish to get on?")) {
    npc.sendOk("You must have some business to take care of here, right?");
} else if (plr.haveItem(TICKET_ID, 1)) {
    plr.gainItem(TICKET_ID, -1);
    plr.warp(200000100);
} else {
    npc.sendOk("Oh, no... You do not have a ticket with you. I can't let you on without one. Please buy a ticket at the ticket sales guide...");
}
