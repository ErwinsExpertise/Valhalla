if (plr.getQuestStatus(3230) < 1) {
    npc.sendOk("We are the toy soldiers here guarding this room, preventing anyone else from entering. I cannot inform you of the reasoning behind this policy. Now, if you'll excuse me, I am working here.")
} else if (plr.getQuestStatus(3230) > 1) {
    npc.sendOk("Thanks to #h0#, we got the #bPendulum#k back and destroyed the monster from the other dimension. Thankfully we haven't found one like that since. I can't thank you enough for helping us out. Hope you enjoy your stay here at Ludibrium!")
} else {
    var proceed = npc.sendYesNo("Hmmm...I've heard a lot about you through #b#p2040001##k. You got him a bunch of #bTasty Walnut#k so he can fight off boredom at work. Well ... alright, then. There's a dangerous, dangerous monster inside. I want to ask you for help in regards to locating it. Would you like to help me out?")
    if (!proceed) {
        npc.sendOk("I see. It's very understandable, considering the fact that you'll be facing a very dangerous monster inside. If you ever feel a change of heart, then please come talk to me. I sure can use help from someone like you.")
    }
    npc.sendNext("Thank you so much. Actually, #b#p2040001##k asked you to get #bTasty Walnuts#k as a way of testing your abilities to see if you can handle this, so don't think of it as a random request. I think someone like you can handle adversity well.")
    var enter = npc.sendYesNo("A while ago, a monster came here from another dimension thanks to a crack in dimensions, and it stole the pendulum of the clock. It hid itself inside the room over there camouflaged as a dollhouse. It all looks the same to me, so there's no way to find it. Would you help us locate it?")
    if (!enter) {
        npc.sendOk("I see. Please talk to me when you're ready to take on this task. I advise you not to take too much time, though, for the monster may turn into something totally different. We have to act like we don't know anything.")
    }
    npc.sendNext("Alright! I'll take you to a room, where you'll find a number of dollhouses all over the place. One of them will look slightly different from the others. Your job is to locate it and break its door. If you locate it correctly, you'll find #bPendulum#k. If you break a wrong dollhouse, however, you'll be sent out here without warning, so please be careful on that.")
    npc.sendBackNext("You'll also find monsters in there, and they have gotten so powerful thanks to the monster from the other dimension that you won't be able to take them down. Please find #bPendulum#k within the time limit and then notify #b#p2040028##k, who should be inside. Let's get this started!", true, true)
    if (plr.countPlayersOnMap(922000010) < 1) {
        // Reset map
        plr.resetMapEntry(922000010)
        // Remove 4031094
        plr.removeItemsByID(4031094, plr.itemCount(4031094))
        // Warp player
        plr.warpPartyToInstance(922000010, 0)
        // Start timer 600 seconds to map 221023200
        plr.setMapTimeLimit(600, 221023200)
    } else {
        npc.sendNext("Someone else must be inside looking for the dollhouse. Unfortunately I can only let in one person at a time, so please wait for your turn.")
    }
}