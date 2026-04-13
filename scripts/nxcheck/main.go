package main

import (
	"fmt"
	"path/filepath"

	"github.com/Hucaru/Valhalla/nx"
)

func main() {
	nx.LoadFile(filepath.Join("..", "v48", "wz", "nx"))

	item, itemErr := nx.GetItem(2000000)
	field, mapErr := nx.GetMap(100000000)
	mob, mobErr := nx.GetMob(100100)
	quest, questErr := nx.GetQuest(1000)

	fmt.Printf("item(2000000): name=%q err=%v\n", item.Name, itemErr)
	fmt.Printf("map(100000000): town=%t name=%q err=%v\n", field.Town, field.MapName, mapErr)
	fmt.Printf("mob(100100): hp=%d err=%v\n", mob.MaxHP, mobErr)
	fmt.Printf("quest(1000): name=%q err=%v\n", quest.Name, questErr)
	fmt.Printf("commodities=%d packages=%d reactors=%d maps=%d\n", len(nx.GetCommodities()), len(nx.GetPackages()), len(nx.GetReactorInfoList()), len(nx.GetMaps()))
}
