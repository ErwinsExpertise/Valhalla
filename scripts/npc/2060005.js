if (plr.getQuestStatus(6002) > 1) {
    npc.sendOk("You already guarded the pig, and you did just fine.");

}

if (plr.getQuestStatus(6002) < 1) {
    npc.sendOk("What pig? Where did you hear about that?");

}

if (plr.itemCount(4031508) > 5 && plr.itemCount(4031507) > 5) {
    npc.sendOk("I don't need another one of #bKenta's Reports#k and I'm all stocked up on #bPheromone#k. You don't need to go in.");

}

var em = npc.getEventManager("q6002");
var prop = em.getProperty("state");

if (prop == null || prop == 0) {
    plr.removeItemsByID(4031507, plr.itemCount(4031507));
    plr.removeItemsByID(4031508, plr.itemCount(4031508));
    em.startInstance(plr);
} else {
    npc.sendNotice(5, "Someone is attempting to protect the Watch Hog already. Please try again later.");
}