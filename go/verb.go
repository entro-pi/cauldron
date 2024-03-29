package main

import (
	"github.com/SolarLune/dngn"
	"fmt"
	"strconv"
	"os"
	"bufio"
	"strings"
  "math/rand"
	term "github.com/nsf/termbox-go"
	"context"
	"time"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

func hash(value string) string {
  newVal := ""
  for i := 0;i < len(value);i++ {
    newVal += strconv.Itoa(int(value[i])*32+100)
  }
  return newVal
}
func lookupPlayerByHash(playerHash string) Player {

  userFile, err := os.Open("creds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@sncn-hifs4.mongodb.net/test?retryWrites=true&w=majority"))

  if err != nil {
    panic(err)
  }
  ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
  err = client.Connect(ctx)
  if err != nil {
    panic(err)
  }
  var player Player
  collection := client.Database("pfiles").Collection("Players")

  result  := collection.FindOne(context.Background(), bson.M{"playerhash": bson.M{"$eq":playerHash}})
  if err != nil {
    panic(err)
  }
  err = result.Decode(&player)

  if err != nil {
    fmt.Println("\033[38:2:150:0:150mPlayerfile requested was not found\033[0m")
    var noob Player
    noob.PlayerHash = "2"
		panic(err)
    return noob
  }
	player = decompInv(player)
	player = decompEq(player)
  return player
}

func LookupPlayer(name string, password string) Player {
  userFile, err := os.Open("creds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@sncn-hifs4.mongodb.net/test?retryWrites=true&w=majority"))

  if err != nil {
    panic(err)
  }
  ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
  err = client.Connect(ctx)
  if err != nil {
    panic(err)
  }
  var player Player
  collection := client.Database("pfiles").Collection("Players")

  result  := collection.FindOne(context.Background(), bson.M{"playerhash": bson.M{"$eq":hash(name+password)}})
  if err != nil {
    panic(err)
  }
  err = result.Decode(&player)

	if err != nil {
    fmt.Println("\033[38:2:150:0:150mPlayerfile requested was not found\033[0m")
    var noob Player
    noob.PlayerHash = "2"
		panic(err)
    return noob
  }
	player = decompInv(player)
	player = decompEq(player)
  return player

}
func getBroadcasts() []Broadcast {
	userFile, err := os.Open("creds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@sncn-hifs4.mongodb.net/test?retryWrites=true&w=majority"))

	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	findOptions := options.Find()
	findOptions.SetLimit(1000)
	collection := client.Database("broadcasts").Collection("general")

	result, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		panic(err)
	}
//	fmt.Println("\033[38:2:255:0:0m", result, "\033[0m")
	var container []Broadcast

	err = result.All(context.Background(), &container)
	if err != nil {
		panic(err)
	}
	//	fmt.Print("\033[38:2:0:0:200m",container, "\033[0m")
	return container
}

func updateChat() []Broadcast {
	userFile, err := os.Open("creds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@sncn-hifs4.mongodb.net/test?retryWrites=true&w=majority"))

	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	findOptions := options.Find()
	findOptions.SetLimit(1000)
	collection := client.Database("broadcasts").Collection("general")

	result, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		panic(err)
	}
//	fmt.Println("\033[38:2:255:0:0m", result, "\033[0m")
	var container []Broadcast

	err = result.All(context.Background(), &container)
	if err != nil {
		panic(err)
	}
	//	fmt.Print("\033[38:2:0:0:200m",container, "\033[0m")
	return container
}

func digDug(pos []int, play Player, digFrame [][]int, digNums string, digZone string, digNum int, populated []Space) (int, Space) {
	digVnumEnd := strings.Split(digNums, "-")[1]
	dg, digNum := initDigRoom(digFrame, digNums, digZone, play, digNum)
	play.CurrentRoom = dg
	for len(populated) <= digNum {
		populated = append(populated, dg)
	}
	populated[digNum] = dg
	dg.Vnum = digNum
	digFrame[pos[0]][pos[1]] = 8
	dg.ZonePos = dg.ZonePos[:0]
	dg.ZonePos = append(dg.ZonePos, pos[0])
	dg.ZonePos = append(dg.ZonePos, pos[1])
	fmt.Println("dug ", dg)
	drawDig(digFrame, dg.ZonePos)
	//updateRoom(play, populated)
	fmt.Println("Dug ", digNum, " rooms of ", digVnumEnd)
	return digNum, dg
}


func AssembleComposeCel(chatMess Chat, row int) (string, int) {
	var cel string
	inWord := chatMess.Message
	wor := ""
	word := ""
	words := ""
	if len(inWord) > 68 {
		return "DONE COMPOSTING", 0
	}
	if len(inWord) > 28 && len(inWord) > 54 {
		wor += inWord[:28]
		word += inWord[28:54]
		words += inWord[54:]
		for i := len(words); i <= 28; i++ {
			words += " "
		}
	}
	if len(inWord) > 28 && len(inWord) < 54 {
		wor += inWord[:28]
		word += inWord[28:]
		for i := len(word); i <= 28; i++ {
			word += " "
		}
		words = "                            "

	}
	if len(inWord) <= 28 {
		wor = "                            "
		word += ""
		word += inWord
		for i := len(word); i <= 28; i++ {
			word += " "
		}
		words = "                            "
	}
	timeString := strings.Split(chatMess.Time.String(), " ")
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";180H\033[48;2;10;255;20m \033[48;2;10;10;20m", wor, "\033[48;2;10;255;20m \033[0m")
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";180H\033[48;2;10;255;20m \033[48;2;10;10;20m", word, "\033[48;2;10;255;20m \033[0m"+timeString[1])
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";180H\033[48;2;10;255;20m \033[48;2;10;10;20m", words, "\033[48;2;10;255;20m \033[0m"+timeString[0])
	row++
	namePlate := "                            "[len(chatMess.User.Name):]
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";180H\033[48;2;10;255;20m\033[38:2:50:0:50m@"+chatMess.User.Name+namePlate+"\033[48;2;10;255;20m \033[0m")

	return cel, row
	//	fmt.Println(cel)
}
func sendBroadcast(bcast Broadcast) Broadcast {
  userFile, err := os.Open("creds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@sncn-hifs4.mongodb.net/test?retryWrites=true&w=majority"))

  if err != nil {
    panic(err)
  }
  ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
  err = client.Connect(ctx)
  if err != nil {
    panic(err)
  }
  update := bson.M{"event":bcast.Event,"ref":bcast.Ref,"payload":bson.M{
		"channel":bcast.Payload.Channel,"id":bcast.Payload.ID,
		 "message":bcast.Payload.Message,"game":bcast.Payload.Game,
		 "name":bcast.Payload.Name,"bigmessage":bcast.Payload.BigMessage,
		 "transaction":bcast.Payload.Transaction,"cpu":bcast.Payload.CPU,
		 "fights":bcast.Payload.Fights,
		 "plaincoreboard":bcast.Payload.PlainCoreBoard,"coreboard":bcast.Payload.CoreBoard}}
  collection := client.Database("broadcasts").Collection("general")
  _, err = collection.InsertOne(context.Background(), update)
  if err != nil {
    panic(err)
  }
  fmt.Println("Upserted the broadcast")
  return bcast

}
func onlineHash(value string) string {
  newVal := ""
  for i := 0;i < len(value);i++ {
    newVal += strconv.Itoa(int(value[i])*24+240)
  }
  return newVal
}
func updateBroadcast(bcast Broadcast) Broadcast {
  userFile, err := os.Open("creds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@sncn-hifs4.mongodb.net/test?retryWrites=true&w=majority"))

  if err != nil {
    panic(err)
  }
  ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
  err = client.Connect(ctx)
  if err != nil {
    panic(err)
  }
	filter := bson.M{"ref":bson.M{"$eq":bcast.Ref}}
  update := bson.M{"$set":bson.M{"event":bcast.Event,"ref":bcast.Ref,
		"payload":bson.M{"channel":bcast.Payload.Channel,"id":bcast.Payload.ID,
			 "message":bcast.Payload.Message,"game":bcast.Payload.Game,"name":bcast.Payload.Name,
			 "bigmessage":bcast.Payload.BigMessage,"transaction":bcast.Payload.Transaction,
			 "cpu":bcast.Payload.CPU,"plaincoreboard":bcast.Payload.PlainCoreBoard,
			 "coreboard":bcast.Payload.CoreBoard,"fights":bcast.Payload.Fights}}}
  collection := client.Database("broadcasts").Collection("general")
  _, err = collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
  if err != nil {
    panic(err)
  }
  fmt.Println("Upserted the broadcast")
  return bcast

}


func AssembleBroadside(broadside Broadcast, row int, col int) (string) {
	var cel string
	colString := strconv.Itoa(col)
	inWord := broadside.Payload.Message
	if len(inWord) < 15 {
		inWord = " "+inWord
		for len(inWord) < 90 {
			inWord += " "+inWord+" "
		}
	}else if len(inWord) < 30 {
		inWord = " "+inWord
		for len(inWord) < 90 {
			inWord = " "+inWord+" "
		}
	}else {
		inWord = " "+inWord
		for len(inWord) < 90 {
			inWord = inWord + " "
		}
	}
	wor := inWord[0:30]
	word := inWord[29:59]
	words := inWord[59:89]

	numString := strconv.Itoa(broadside.Payload.ID)

	row++
	if broadside.Payload.Selected {
		cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;200;25;150m ", wor, "\033[48;2;200;25;150m \033[0m")
	}else {
		cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;20;255;50m \033[48;2;10;10;20m", wor, "\033[48;2;20;255;50m \033[0m")
	}

	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;20;255;50m \033[48;2;10;10;20m", word, "\033[48;2;20;255;50m \033[0m")
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;20;255;50m \033[48;2;10;10;20m", words, "\033[48;2;20;255;50m \033[0m")
	row++
	if broadside.Payload.Game == "" {
		if broadside.Payload.Selected {
			broadside.Payload.Game = "SELECTED"
		} else {
			broadside.Payload.Game = "snowcrash"
		}
	}

	namePlate := "                            "[len(broadside.Payload.Name+numString):]
	if broadside.Payload.Selected {
		cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;200;25;150m @"+broadside.Payload.Name+"@"+numString+namePlate+"\033[48;2;200;25;50m \033[0m")
	}else {
		cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;20;255;50m@"+broadside.Payload.Name+"@"+numString+namePlate+"\033[48;2;20;255;50m \033[0m")

	}

	broadRow := 0
	if broadside.Payload.Selected && len(strings.Split(broadside.Payload.BigMessage, "\n")) > 1 {
		bigSplit := strings.Split(broadside.Payload.BigMessage, "\n")
		for i := 0;i < len(bigSplit);i++ {
			cel += fmt.Sprint("\033["+strconv.Itoa(24+broadRow)+";53H\033[0m"+bigSplit[broadRow]+"\033[0m")
			broadRow++
		}
		if !broadside.Payload.Transaction.Sold && broadside.Payload.Transaction.Item.Name != "" {
			cel += fmt.Sprintf("\033[44;53H\033[38:2:150:50:150m{{FOR SALE %v}} \033[38:2:75:75:0m||'BUY' for %v credbits||\033[0m", broadside.Payload.Transaction.Item.Name, broadside.Payload.Transaction.Price)
		}
	}



	return cel
	//	fmt.Println(cel)
}
func AssembleBM(broadside Broadcast, row int, col int) (string) {
	broadside.Payload.Selected = true
	var cel string
	colString := strconv.Itoa(col)
	inWord := broadside.Payload.Message
	if len(inWord) < 15 {
		inWord = " "+inWord
		for len(inWord) < 90 {
			inWord += " "+inWord+" "
		}
	}else if len(inWord) < 30 {
		inWord = " "+inWord
		for len(inWord) < 90 {
			inWord = " "+inWord+" "
		}
	}else {
		inWord = " "+inWord
		for len(inWord) < 90 {
			inWord = inWord + " "
		}
	}
	wor := inWord[0:30]
	word := inWord[29:59]
	words := inWord[59:89]

	numString := strconv.Itoa(broadside.Payload.ID)

	row++
	if broadside.Payload.Selected {
		cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;200;25;150m ", wor, "\033[48;2;200;25;150m \033[0m")
	}else {
		cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;20;255;50m \033[48;2;10;10;20m", wor, "\033[48;2;20;255;50m \033[0m")
	}

	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;20;255;50m \033[48;2;10;10;20m", word, "\033[48;2;20;255;50m \033[0m")
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;20;255;50m \033[48;2;10;10;20m", words, "\033[48;2;20;255;50m \033[0m")
	row++
	if broadside.Payload.Game == "" {
		if broadside.Payload.Selected {
			broadside.Payload.Game = "SELECTED"
		} else {
			broadside.Payload.Game = "snowcrash"
		}
	}

	namePlate := "                            "[len(broadside.Payload.Name+numString):]
	if broadside.Payload.Selected {
		cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;200;25;150m @"+broadside.Payload.Name+"@"+numString+namePlate+"\033[48;2;200;25;50m \033[0m")
	}else {
		cel += fmt.Sprint("\033["+strconv.Itoa(row)+";"+colString+"H\033[48;2;20;255;50m@"+broadside.Payload.Name+"@"+numString+namePlate+"\033[48;2;20;255;50m \033[0m")

	}

	broadRow := 0
	if broadside.Payload.Selected && len(strings.Split(broadside.Payload.BigMessage, "\n")) > 1 {
		bigSplit := strings.Split(broadside.Payload.BigMessage, "\n")
		for i := 0;i < len(bigSplit);i++ {
			cel += fmt.Sprint("\033["+strconv.Itoa(24+broadRow)+";53H\033[0m"+bigSplit[broadRow]+"\033[0m")
			broadRow++
		}
		if !broadside.Payload.Transaction.Sold && broadside.Payload.Transaction.Item.Name != "" {
			cel += fmt.Sprintf("\033[44;53H\033[38:2:150:50:150m{{FOR SALE %v}} \033[38:2:75:75:0m||'BUY' for %v credbits||\033[0m", broadside.Payload.Transaction.Item.Name, broadside.Payload.Transaction.Price)
		}
	}



	return cel
	//	fmt.Println(cel)
}

func AssembleDescCel(room Space, row int) (string) {
	var cel string
	inWord := strings.Split(room.Desc, "\n")
	for len(strings.Split(room.Desc, "\n")) < 9 {
		room.Desc += "\n"
		inWord = strings.Split(room.Desc, "\n")
	}
	for len(inWord[0]) < 100 {
		inWord[0] += " "
	}
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row+20)+";51H\033[48;2;10;255;20m \033[48;2;10;10;20m", inWord[0], "\033[48;2;10;255;20m \033[0m")
	for len(inWord[1]) < 100 {
		inWord[1] += " "
	}
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row+20)+";51H\033[48;2;10;255;20m \033[48;2;10;10;20m", inWord[1], "\033[48;2;10;255;20m \033[0m")
	for len(inWord[2]) < 100 {
		inWord[2] += " "
	}
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row+20)+";51H\033[48;2;10;255;20m \033[48;2;10;10;20m", inWord[2], "\033[48;2;10;255;20m \033[0m")
	for len(inWord[3]) < 100 {
		inWord[3] += " "
	}
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row+20)+";51H\033[48;2;10;255;20m \033[48;2;10;10;20m", inWord[3], "\033[48;2;10;255;20m \033[0m")
	for len(inWord[4]) < 100 {
		inWord[4] += " "
	}
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row+20)+";51H\033[48;2;10;255;20m \033[48;2;10;10;20m", inWord[4], "\033[48;2;10;255;20m \033[0m")
	for len(inWord[5]) < 100 {
		inWord[5] += " "
	}
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row+20)+";51H\033[48;2;10;255;20m \033[48;2;10;10;20m", inWord[5], "\033[48;2;10;255;20m \033[0m")
	for len(inWord[6]) < 100 {
		inWord[6] += " "
	}
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row+20)+";51H\033[48;2;10;255;20m \033[48;2;10;10;20m", inWord[6], "\033[48;2;10;255;20m \033[0m")
	for len(inWord[7]) < 100 {
		inWord[7] += " "
	}
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row+20)+";51H\033[48;2;10;255;20m\033[38:2:50:0:50m@"+room.Zone+"#"+strconv.Itoa(room.Vnum), inWord[7][len(room.Zone)+len(strconv.Itoa(room.Vnum)):], "\033[48;2;10;255;20m \033[0m")

	return cel
}
func countKeys() {
  keys := "abcdefghijklmnopqrstuvwxyz0123456789"
  fmt.Println("\033[38:2:150:0:150m",len(keys),"in :",keys)

  keys = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
  fmt.Println("\033[38:2:175:0:150m",len(keys),"in :",keys)

  keys = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
  fmt.Println("\033[38:2:185:0:150m",len(keys),"in :",keys)
}

func goTo(dest int, play Player, populated []Space) (Player, []Space) {

	for i := 0;i < len(populated);i++ {
		if dest == populated[i].Vnum {
			play.CurrentRoom = populated[i]
			fmt.Print(populated[i].Vnum, populated[i].Vnums, populated[i].Zone)
			showDesc(play.CurrentRoom)
			DescribePlayer(play)
			fmt.Printf("\033[0;0H\033[38:2:0:255:0mPASS\033[0m")
			savePfile(play)
			break
		}else {
			fmt.Printf("\033[0;0H\033[38:2:255:0:0mERROR\033[0m")
		}
	}
	return play, populated
}


func mergeMaps(source [][]int, dest [][]int) ([][]int) {
  for i := 0;i < len(source);i++ {
    for c := 0;c < len(source[i]);c++ {
      if source[i][c] == 1 {
        dest[i][c] = 1
      }
    }
  }
  return dest
}
func improvedTargeting(play Player, target string) (Player) {

	if strings.Contains(target, "|") {
		tarX, err := strconv.Atoi(strings.Split(target, "|")[0])
		if err != nil {
			panic(err)
		}
		tarY, err := strconv.Atoi(strings.Split(target, "|")[1])
		if err != nil {
			panic(err)
		}

		play.TarX = tarX
		play.TarY = tarY
	}else {
		play.OldX, play.OldY = play.TarX, play.TarY

		switch target {
		case "8":
			play.TarY -= 1
		case "2":
			play.TarY += 1
		case "4":
			play.TarX -= 1
		case "6":
			play.TarX += 1
		}

	}
	targ := ""
//	fmt.Print(play.CPU)
	splitCPU := strings.Split(play.CPU, "\n")
	CPU:
	for i := 0;i < len(splitCPU);i++ {
		for r := 0;r < len(splitCPU[i]);r++ {
			if play.TarX == r && play.TarY == i {
				if string(splitCPU[i][r]) == "%" {
					play.TarX, play.TarY = play.OldX, play.OldY
					targ = fmt.Sprint("\033["+strconv.Itoa(play.TarY+20)+";"+strconv.Itoa(play.TarX+54)+"H\033[48:2:175:0:150m"+string(splitCPU[play.TarY][play.TarX])+"\033[0m")
					break CPU
				}else {
//					fmt.Print("\033["+strconv.Itoa(i+20)+";"+strconv.Itoa(r+54)+"H\033[48:2:175:0:150m"+string(splitCPU[play.TarY][play.TarX])+"\033[0m")
					play.TargetLong = string(splitCPU[play.TarY][play.TarX])

				}

				targ = fmt.Sprint("\033["+strconv.Itoa(i+20)+";"+strconv.Itoa(r+54)+"H\033[48:2:175:0:150m"+string(splitCPU[play.TarY][play.TarX])+"\033[0m")

			}else {
				fmt.Print("\033["+strconv.Itoa(i+20)+";"+strconv.Itoa(r+54)+"H"+string(splitCPU[i][r]))
			}
		}
	}

	play.Target = targ
	return play
}

func setCoreBoard(play Player, broad Broadcast) Player {
	play.CPU = broad.Payload.CPU
	play.PlainCoreBoard = broad.Payload.PlainCoreBoard
	play.Fights = broad.Payload.Fights
	play.CoreBoard = broad.Payload.CoreBoard
	return play
}

func genCoreBoard(sizeX int, sizeY int, broad Broadcast) (string, Broadcast) {
	//Create a room map
	Room := dngn.NewRoom(sizeX, sizeY)
	splits := rand.Intn(75)
	Room.GenerateBSP('%', 'D', splits)
//	_, err = collection.InsertOne(context.Background(), bson.M{"room":Room})
//	if err != nil {
//		panic(err)
//	}
  newValue := ""
  outVal := ""
	broad.Payload.Fights = InitFight()
//	fmt.Println("Generating and populating map")
	for i := 0;i < len(Room.Data);i++ {

	//				fmt.Println(populated[0].Room.Data[populated[0].Room.Width-1][i])
			value := string(Room.Data[i])
//      newValue = ""
			for s := 0;s < len(value);s++ {
				if string(value[s]) == " " {
					ChanceTreasure := "T"
					if rand.Intn(100) > 98 {
							newValue += ChanceTreasure
							tiara := InitObject()
							tiara.X = s
							tiara.Y = i
							broad.Payload.Fights.Treasure = append(broad.Payload.Fights.Treasure, tiara)
							continue
					}
					if rand.Intn(100) > 95 {

						TL := ""
						roll := rand.Intn(100)
						if roll <= 30 {
							TL = "A Rabid Ferret"
							newValue += "F"
							ferret := InitMob()
							ferret.Name = TL
							ferret.X = s
							ferret.Y = i
							ferret.AC = 3
							ferret.Char = "F"
							ferret.Corpse = strings.ToLower(TL)
							broad.Payload.Fights.Oppose = append(broad.Payload.Fights.Oppose, ferret)

						}
						if roll <= 60 && roll > 30 {
							TL = "A Wild Boar"
							newValue += "B"
							boar := InitMob()
							boar.Name = TL
							boar.X = s
							boar.Y = i
							boar.AC = 2
							boar.Rezz += rand.Intn(15)
							boar.MaxRezz = boar.Rezz
							boar.Char = "B"
							boar.Corpse = strings.ToLower(TL)
							broad.Payload.Fights.Oppose = append(broad.Payload.Fights.Oppose, boar)

						}
						if roll > 60 {
							TL = "A Razor Beast"
							newValue += "R"
							razor := InitMob()
							razor.Name = TL
							razor.X = s
							razor.Y = i
							razor.AC = 10
							razor.Rezz += rand.Intn(25)
							razor.MaxRezz = razor.Rezz
							razor.Char = "R"
							razor.Corpse = strings.ToLower(TL)
							broad.Payload.Fights.Oppose = append(broad.Payload.Fights.Oppose, razor)
						}
						continue
					}else {
						newValue += string(value[s])
					}
				}else {
					newValue += string(value[s])
				}

			}
      newValue += "\n"
    }
		broad.Payload.CPU = newValue + "\n"
    broad.Payload.PlainCoreBoard = newValue
    broad.Payload.CoreBoard = newValue
    newValue = strings.ReplaceAll(newValue, "T", "\033[48;2;200;150;0mT\033[0m")

    broad.Payload.CoreBoard = newValue
    newValue = strings.ReplaceAll(newValue, "M", "\033[48;2;200;50;50mM\033[0m")

    broad.Payload.CoreBoard = newValue
		newValue = strings.ReplaceAll(newValue, "%", "\033[38;2;0;150;150m%\033[0m")

    broad.Payload.CoreBoard = newValue
		newValue = strings.ReplaceAll(newValue, "D", "\033[48;2;200;150;150mD\033[0m")

    broad.Payload.CoreBoard = newValue
		newValue = strings.ReplaceAll(newValue, " ", "\033[48:2:0:200:150m \033[0m")
		broad.Payload.CoreBoard = newValue
    outVal += newValue + "\n"
	return outVal, broad
}

func resetCraft() string {
	val := ""
	val += "\033[0;53H\033[48:2:200:120:0m                                                                              \033[0m"
	val += "\033[32;53H\033[48:2:200:120:0m                                                                              \033[0m"
	for i := 2;i < 32;i++ {

		if i == 2{
			val += "\033["+strconv.Itoa(i)+";53H\033[48:2:200:120:0m \033[0m                                Name                                        \033[48:2:200:120:0m \033[0m"
		}else if i == 5 {
			val += "\033["+strconv.Itoa(i)+";53H\033[48:2:200:120:0m \033[0m                             Description                                    \033[48:2:200:120:0m \033[0m"
			}else if i == 4 || i > 11 && i <= 15 {
			if i == 4 {
				val += fmt.Sprint("\033[38:2:225:0:225m\033["+strconv.Itoa(i)+";53H\033[48:2:200:120:0m$                                                                             \033[0m")

			}else if i > 11 && i < 15 {
				val += fmt.Sprint("\033[38:2:225:0:225m\033["+strconv.Itoa(i)+";53H\033[48:2:200:120:0m&                                                                             \033[0m")

			}else if i == 15 {
				val += "\033["+strconv.Itoa(i)+";53H\033[48:2:200:120:0m \033[0m                               Stats                                        \033[48:2:200:120:0m \033[0m"

				}else {
				val += fmt.Sprint("\033[38:2:225:0:225m\033["+strconv.Itoa(i)+";53H\033[48:2:200:120:0m                                                                              \033[0m")

			}
		}else {
			val += "\033["+strconv.Itoa(i)+";53H\033[48:2:200:120:0m \033[0m                                                                            \033[48:2:200:120:0m \033[0m"
			val += "\033["+strconv.Itoa(i+1)+";53H\"exit\" to end"
		}
	}
	return val
}
func craftObject() Object {
	var obj Object
	namePos := ""
	name := ""
	longName := ""
	named := false
	longNamePos := ""
	longNamed := false

         err := term.Init()
         if err != nil {
                 panic(err)
         }

         defer term.Close()
				 val := resetCraft()
				 fmt.Print(val)
         for {
                 switch ev := term.PollEvent(); ev.Type {
                 case term.EventKey:
                         switch ev.Key {
                         case term.KeyEsc:
													 			named = false
																name = ""
																longNamed = false
																longName = ""
//                                 break keyPressListenerLoop
                         case term.KeyBackspace:
													 clearDirty()
													 val := resetCraft()
													 val += fmt.Sprint(namePos, name)
													 val += fmt.Sprint(longNamePos, longName)
													 fmt.Print(val)
														if len(name) <= 0 {
															name = " "
														}
														if len(longName) <= 0 {
															longName = " "
														}
													if !named {
														 fmt.Println("BACKSPACE")
														 name = name[:len(name)-1]
														 if len(name) <= 0 {
															 name = " "
														 }
														 fmt.Printf("\033[3;80H%v                                                            ", name)
													 }else if !longNamed {
														 longName = longName[:len(longName)-1]
														 if len(longName) <= 0 {
															longName = " "
														 }
														 fmt.Printf("\033[6;70H%v                                                            ", longName)
													 }
												 case term.KeyBackspace2:
																 clearDirty()
																 val := resetCraft()
																 val += fmt.Sprint(namePos, name)
 	 															 val += fmt.Sprint(longNamePos, longName)
 																 fmt.Print(val)
													 				if len(name) <= 0 {
																		name = " "
																	}
																	if len(longName) <= 0 {
																		longName = " "
																	}
													 			if !named {
																	 fmt.Println("BACKSPACE")
																	 name = name[:len(name)-1]
																	 if len(name) <= 0 {
																		 name = " "
																	 }
																	 fmt.Printf("\033[3;80H%v                                                            ", name)
																 }else if !longNamed {
																	 longName = longName[:len(longName)-1]
																	 if len(longName) <= 0 {
																	 	longName = " "
																	 }
																	 fmt.Printf("\033[7;70H%v                                                            ", longName)
																 }
//																 fmt.Printf("\033[3;80H%v                                                  doot      ", name)

	//															 fmt.Printf("\033[6;70H%v                                                        ", longName)

												 case term.KeySpace:
													 if !named {
														 name += " "
													 }else if !longNamed {
														 longName += " "
													 }
                         case term.KeyEnter:
                                 val := resetCraft()

	 															 val += fmt.Sprint(namePos, name)
	 															 val += fmt.Sprint(longNamePos, longName)
																 fmt.Print(val)
																 fmt.Println("Value Accepted.")
																 if !named {
																	 obj.Name = name
																	 named = true
																 }else if !longNamed && len(longName) > 5 {
																	 obj.LongName = longName
 															 		longNamed = true
																 }
																 if named && longNamed {
																	 fmt.Println("Naming complete, are you happy with these changes?")
																	 fmt.Print("@ to save and exit, escape to discard changes.")
																	 key := ""
																	 _, err := fmt.Scan(&key)
																	 if err != nil {
																		 panic(err)
																	 }
																	 if key == "@" {
																		 return obj
																	 }
																 }
                         default:
                                 // we only want to read a single character or one key pressed event
                                 val := resetCraft()
																 fmt.Print(val)

																fmt.Print(namePos, name)
																fmt.Print(longNamePos, longName)
															 	if !named {
																	name += string(ev.Ch)
															 		namePos = fmt.Sprint("\033[3;80H")
																	fmt.Print(namePos, name)
															 	}
																if !longNamed && named {
																	fmt.Print(namePos, name)
																	longName += string(ev.Ch)
															 		longNamePos = fmt.Sprint("\033[7;70H")
															 		fmt.Print(longNamePos, longName)
															 	}

                         }
                 case term.EventError:
                         panic(ev.Err)
                 }
         }

return obj

}

func craftMobile() Mobile {
	var mob Mobile
	namePos := ""
	name := ""
	longName := ""
	named := false
	longNamePos := ""
	longNamed := false

         err := term.Init()
         if err != nil {
                 panic(err)
         }

         defer term.Close()
				 val := resetCraft()
				 fmt.Print(val)
         for {
                 switch ev := term.PollEvent(); ev.Type {
                 case term.EventKey:
                         switch ev.Key {
                         case term.KeyEsc:
													 			named = false
																name = ""
																longNamed = false
																longName = ""
//                                 break keyPressListenerLoop
                         case term.KeyBackspace:
													 clearDirty()
													 val := resetCraft()
													 val += fmt.Sprint(namePos, name)
													 val += fmt.Sprint(longNamePos, longName)
													 fmt.Print(val)
														if len(name) <= 0 {
															name = " "
														}
														if len(longName) <= 0 {
															longName = " "
														}
													if !named {
														 fmt.Println("BACKSPACE")
														 name = name[:len(name)-1]
														 if len(name) <= 0 {
															 name = " "
														 }
														 fmt.Printf("\033[3;80H%v                                                            ", name)
													 }else if !longNamed {
														 longName = longName[:len(longName)-1]
														 if len(longName) <= 0 {
															longName = " "
														 }
														 fmt.Printf("\033[6;70H%v                                                            ", longName)
													 }
												 case term.KeyBackspace2:
																 clearDirty()
																 val := resetCraft()
																 val += fmt.Sprint(namePos, name)
 	 															 val += fmt.Sprint(longNamePos, longName)
 																 fmt.Print(val)
													 				if len(name) <= 0 {
																		name = " "
																	}
																	if len(longName) <= 0 {
																		longName = " "
																	}
													 			if !named {
																	 fmt.Println("BACKSPACE")
																	 name = name[:len(name)-1]
																	 if len(name) <= 0 {
																		 name = " "
																	 }
																	 fmt.Printf("\033[3;80H%v                                                            ", name)
																 }else if !longNamed {
																	 longName = longName[:len(longName)-1]
																	 if len(longName) <= 0 {
																	 	longName = " "
																	 }
																	 fmt.Printf("\033[7;70H%v                                                            ", longName)
																 }
//																 fmt.Printf("\033[3;80H%v                                                  doot      ", name)

	//															 fmt.Printf("\033[6;70H%v                                                        ", longName)

												 case term.KeySpace:
													 if !named {
														 name += " "
													 }else if !longNamed {
														 longName += " "
													 }
                         case term.KeyEnter:
                                 val := resetCraft()

	 															 val += fmt.Sprint(namePos, name)
	 															 val += fmt.Sprint(longNamePos, longName)
																 fmt.Print(val)
																 fmt.Println("Value Accepted.")
																 if !named {
																	 mob.Name = name
																	 named = true
																 }else if !longNamed && len(longName) > 5 {
																	 mob.LongName = longName
 															 		longNamed = true
																 }
																 if named && longNamed {
																	 fmt.Println("Naming complete, are you happy with these changes?")
																	 fmt.Print("@ to save and exit, escape to discard changes.")
																	 key := ""
																	 _, err := fmt.Scan(&key)
																	 if err != nil {
																		 panic(err)
																	 }
																	 if key == "@" {
																		 return mob
																	 }
																 }
                         default:
                                 // we only want to read a single character or one key pressed event
                                 val := resetCraft()
																 fmt.Print(val)

																fmt.Print(namePos, name)
																fmt.Print(longNamePos, longName)
															 	if !named {
																	name += string(ev.Ch)
															 		namePos = fmt.Sprint("\033[3;80H")
																	fmt.Print(namePos, name)
															 	}
																if !longNamed && named {
																	fmt.Print(namePos, name)
																	longName += string(ev.Ch)
															 		longNamePos = fmt.Sprint("\033[7;70H")
															 		fmt.Print(longNamePos, longName)
															 	}

                         }
                 case term.EventError:
                         panic(ev.Err)
                 }
         }

return mob

}

//TODO make this modular
func createChat(message string, play Player) {
	userFile, err := os.Open("creds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@sncn-hifs4.mongodb.net/test?retryWrites=true&w=majority"))

	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	//process the strings
	if len(message) >= 180 {
		message = message[:180]
	}

	collection := client.Database("chat").Collection("log")
	_, err = collection.InsertOne(context.Background(), bson.M{"name":play.Name,
						"message":message, "time":time.Now(), "user":play })
	if err != nil {
		panic(err)
	}
}

//TODO make this modular
func createMobiles(name string) {
	userFile, err := os.Open("creds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@sncn-hifs4.mongodb.net/test?retryWrites=true&w=majority"))

	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	collection := client.Database("mobiles").Collection("lvl1")
	_, err = collection.InsertOne(context.Background(), bson.M{"name":name,
						"str": 1, "int": 1, "dex": 1, "wis": 1, "con":1, "cha":1, "challengedice":1 })
}

func addPfile(play Player) {
	userFile, err := os.Open("creds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@sncn-hifs4.mongodb.net/test?retryWrites=true&w=majority"))

	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	collection := client.Database("pfiles").Collection("Players")
	_, err = collection.InsertOne(context.Background(), bson.M{"playerhash":play.PlayerHash,
		"name":play.Name,"title":play.Title,"inventory":play.Inventory, "equipped":play.Equipped,
		"coreboard": play.CoreBoard, "str": play.Str, "int": play.Int, "dex": play.Dex,
		 "wis": play.Wis, "con":play.Con, "cha":play.Cha })
		if err != nil {
			panic(err)
		}
}
func savePfile(play Player) {
	userFile, err := os.Open("creds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@sncn-hifs4.mongodb.net/test?retryWrites=true&w=majority"))

	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	filter := bson.M{"playerhash":bson.M{"$eq":play.PlayerHash}}
	collection := client.Database("pfiles").Collection("Players")
	_, err = collection.UpdateOne(context.Background(),filter, bson.M{"$set":bson.M{
		"playerhash":play.PlayerHash,"name":play.Name,
		"title":play.Title,"inventory":play.Inventory,"level":play.Level,
		 "equipped":play.Equipped,"coreboard": play.CoreBoard,
		 "bankaccount":bson.M{"owner":play.Name,"amount":play.BankAccount.Amount},
		  "str": play.Str, "int": play.Int, "dex": play.Dex, "wis": play.Wis,
			"battling":play.Battling,"battlingmob":play.BattlingMob,"rezz":play.Rezz,
			"tech":play.Tech,
			 "con":play.Con, "cha":play.Cha, "classes": play.Classes }}, options.Update().SetUpsert(true))
	if err != nil {
		panic(err)
	}
	savePinv(play)
	savePeq(play)
}
func savePinv(play Player) {
	play = composeInv(play)
	userFile, err := os.Open("creds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@sncn-hifs4.mongodb.net/test?retryWrites=true&w=majority"))

	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	filter := bson.M{"playerhash":bson.M{"$eq": play.PlayerHash}}
	update := bson.M{"$set":bson.M{"itembank":bson.M{"slotone":play.ItemBank.SlotOne,
		"slotoneamount":play.ItemBank.SlotOneAmount,"slottwo":play.ItemBank.SlotTwo,
		"slottwoamount":play.ItemBank.SlotTwoAmount,"slotthree":play.ItemBank.SlotThree,
		"slotthreeamount":play.ItemBank.SlotThreeAmount,"slotfour":play.ItemBank.SlotFour,"slotfouramount":play.ItemBank.SlotFourAmount,
		"slotfiveamount":play.ItemBank.SlotFiveAmount,"slotsix":play.ItemBank.SlotSix,
		"slotsixamount":play.ItemBank.SlotSixAmount,"slotseven":play.ItemBank.SlotSeven,
		"slotsevenamount":play.ItemBank.SlotSevenAmount,"sloteight":play.ItemBank.SlotEight,
		"sloteightamount":play.ItemBank.SlotEightAmount,"slotnine":play.ItemBank.SlotNine,
		"slotnineamount":play.ItemBank.SlotNineAmount,"slotten":play.ItemBank.SlotTen,
		"slottenamount":play.ItemBank.SlotTenAmount}}}
	collection := client.Database("pfiles").Collection("Players")
	_, err = collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		panic(err)
	}
}

func savePeq(play Player) {
	play = composeEq(play)
	userFile, err := os.Open("creds")
  if err != nil {
    panic(err)
  }
  defer userFile.Close()
  scanner := bufio.NewScanner(userFile)
  scanner.Scan()
  user := scanner.Text()
  scanner.Scan()
  pass := scanner.Text()
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://"+user+":"+pass+"@sncn-hifs4.mongodb.net/test?retryWrites=true&w=majority"))

	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	filter := bson.M{"playerhash":bson.M{"$eq": play.PlayerHash}}
	update := bson.M{"$set":bson.M{"eqbank":bson.M{"slotone":play.EqBank.SlotOne,
		"slotoneamount":play.EqBank.SlotOneAmount,"slottwo":play.EqBank.SlotTwo,
		"slottwoamount":play.EqBank.SlotTwoAmount,"slotthree":play.EqBank.SlotThree,
		"slotthreeamount":play.EqBank.SlotThreeAmount,"slotfour":play.EqBank.SlotFour,"slotfouramount":play.EqBank.SlotFourAmount,
		"slotfiveamount":play.EqBank.SlotFiveAmount,"slotsix":play.EqBank.SlotSix,
		"slotsixamount":play.EqBank.SlotSixAmount,"slotseven":play.EqBank.SlotSeven,
		"slotsevenamount":play.EqBank.SlotSevenAmount,"sloteight":play.EqBank.SlotEight,
		"sloteightamount":play.EqBank.SlotEightAmount,"slotnine":play.EqBank.SlotNine,
		"slotnineamount":play.EqBank.SlotNineAmount,"slotten":play.EqBank.SlotTen,
		"slottenamount":play.EqBank.SlotTenAmount}}}
	collection := client.Database("pfiles").Collection("Players")
	_, err = collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		panic(err)
	}
					}
