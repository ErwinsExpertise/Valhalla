npc.sendSelection("#L0#Can I leave please?#l\r\n#L1#Sorry to have bothered you.#l");

if (npc.selection() === 0) {
    plr.warp(680000000);
}
