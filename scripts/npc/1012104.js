// Brittany the assistant hair NPC – stateless version
const couponHaircut = 5150052
const couponDye = 5151035
const maleHaircuts = [30310,30330,30060,30150,30410,30210,30140,30120,30200,30560,30510,30610,30470]
const femaleHaircuts = [31150,31310,31300,31160,31100,31410,31030,31080,31070,31610,31350,31510,31740]

const choice = npc.sendMenu(
    "I'm Brittany the assistant. If you have #b#t5150052##k or #b#t5151035##k by any chance, then how about letting me change your hairdo?",
    "#bHaircut (REG coupon)",
    "Dye your hair (REG coupon)"
)

// Haircut branch
if (choice === 0) {
    let hair
    if (plr.gender() < 1) {
        hair = maleHaircuts[Math.floor(Math.random() * maleHaircuts.length)]
    } else {
        hair = femaleHaircuts[Math.floor(Math.random() * femaleHaircuts.length)]
    }
    hair += plr.getHairStyle() % 10

    const ok = npc.sendYesNo("If you use the REG coupon your hair will change RANDOMLY with a chance to obtain a new experimental style that even you didn't think was possible. Are you going to use #b#t5150052##k and really change your hairstyle?")
    if (ok) {
        if (plr.itemCount(couponHaircut) > 0) {
            plr.giveItem(couponHaircut, -1)
            plr.setHairStyle(hair)
            npc.sendNext("Hey, here's the mirror. What do you think of your new haircut? I know it wasn't the smoothest of all, but didn't it come out pretty nice? Come back later when you need to change it up again!")
        } else {
            npc.sendNext("Hmmm...are you sure you have our designated coupon? Sorry but no haircut without it.")
        }
    } else {
        npc.sendOk("See you another time!")
    }
}
// Dye branch
else {
    let hair = Math.floor(plr.getHairStyle() / 10) * 10
    const colors = [hair + 0, hair + 1, hair + 2, hair + 3, hair + 4, hair + 5]
    hair = colors[Math.floor(Math.random() * colors.length)]

    const ok = npc.sendYesNo("If you use a regular coupon your hair will change RANDOMLY. Do you still want to use #b#t5151035##k and change it up?")
    if (ok) {
        if (plr.itemCount(couponDye) > 0) {
            plr.giveItem(couponDye, -1)
            plr.setHairStyle(hair)
            npc.sendNext("Hey, here's the mirror. What do you think of your new haircolor? I know it wasn't the smoothest of all, but didn't it come out pretty nice? Come back later when you need to change it up again!")
        } else {
            npc.sendNext("Hmmm...are you sure you have our designated coupon? Sorry but no dye your hair without it.")
        }
    } else {
        npc.sendOk("See you another time!")
    }
}