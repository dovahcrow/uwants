package loginpool

import (
	//"fmt"
	"myfans"
	"uwants"
)

var Uwants map[string]*uwants.Uwants
var Myfans map[string]*myfans.Myfans

func init() {
	Uwants = make(map[string]*uwants.Uwants)
	Myfans = make(map[string]*myfans.Myfans)
}

func GetUwants(username, password string) *uwants.Uwants {
	//fmt.Printf("|%v|%v|\n", username, password)
	if cl, ok := Uwants[username]; ok {
		return cl
	} else {
		cl := uwants.New(username, password)
		Uwants[username] = cl
		cl.Login()
		return cl
	}

}
