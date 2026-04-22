var ticket = 4031713;

if (!plr.haveItem(ticket, 1)) {
    npc.sendOk("Hello, I am the ticket gate.");
} else {
    npc.sendSelection("Hello, I am the ticket gate. Which ticket do you want to use? You will be warped immediately.#b\r\n#L0##t4031713#");
    if (npc.selection() === 0) {
        if (!npc.sendYesNo("It seems like there is still plenty of room on this ride. Please keep your ticket ready so I can let you on. The journey may be long, but you will get to your destination safely. What do you think? Do you want to go on this ride?")) {
            npc.sendOk("You must have some business to take care of here, right?");
        } else {
            plr.gainItem(ticket, -1);
            plr.warp(103000100);
        }
    }
}
