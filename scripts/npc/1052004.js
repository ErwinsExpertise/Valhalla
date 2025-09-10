var coupon = 5152057;
var face;

if (plr.getGender() < 1) {
    face = [20000, 20001, 20002, 20003, 20004, 20005, 20006, 20008, 20012, 20014, 20015, 20022, 20028];
} else {
    face = [21000, 21001, 21002, 21003, 21004, 21005, 21006, 21007, 21008, 21012, 21013, 21023, 21026];
}

for (var i = 0; i < face.length; i++) {
    face[i] += Math.floor(plr.getFace() / 100) % 10 * 100;
}

var sel = npc.sendAvatar("Let's see...for #b#t5152057##k, you can get a new face. That's right, I can completely transform your face! Wanna give it a shot? Please consider your choice carefully.", face);

if (itemCount(coupon) > 0) {
    removeItemsByID(coupon, 1);
    plr.setFace(face[sel]);
    npc.sendOk("Alright, it's all done! Check yourself out in the mirror. Well, aren't you lookin' marvelous? Haha! If you're sick of it, just give me another call, alright?");
} else {
    npc.sendNext("Hmm ... it looks like you don't have the coupon specifically for this place. Sorry to say this, but without the coupon, there's no plastic surgery for you...");
}