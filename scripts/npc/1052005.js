npc.sendYesNo("If you use the regular coupon, you may end up with a random new look for your face...do you still want to do it using #b#t5152056##k?");

if (npc.sendYesNo("Proceed with the surgery?")) {
    if (plr.itemCount(5152056) <= 0) {
        npc.sendNext("Hmm ... it looks like you don't have the coupon specifically for this place. Sorry to say this, but without the coupon, there's no plastic surgery for you...");
    } else {
        var face;
        if (plr.gender() == 0) {
            var maleBase = [20000, 20005, 20008, 20012, 20016, 20022, 20032];
            face = maleBase[Math.floor(Math.random() * maleBase.length)];
        } else {
            var femaleBase = [21000, 21002, 21008, 21014, 21020, 21024, 21029];
            face = femaleBase[Math.floor(Math.random() * femaleBase.length)];
        }
        face += Math.floor((plr.getFace() / 100) % 10) * 100;
        
        plr.removeItemsByID(5152056, 1);
        plr.setFace(face);
        npc.sendNext("Okay, the surgery's done. Here's a mirror--check it out. What a masterpiece, no? Haha! If you ever get tired of this look, please feel free to come visit me again.");
    }
} else {
    npc.sendNext("I see...take your time and see if you really want it. Let me know when you've decided.");
}