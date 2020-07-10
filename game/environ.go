package game

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/golog"
	"github.com/jinzhu/gorm"
	"github.com/yenkeia/yams/game/cm"
	"github.com/yenkeia/yams/game/orm"
	"github.com/yenkeia/yams/game/proto/client"
	"github.com/yenkeia/yams/game/proto/server"
)

const (
	// LOGIN 客户端连接到服务器，正在输入账号密码的状态
	LOGIN = iota
	// SELECT 选角色状态
	SELECT
	// GAME 进入游戏状态
	GAME
	// DISCONNECTED 应该是小退后的状态
	DISCONNECTED
)

var log = golog.New("yams.game")
var sessionPlayer = make(map[int64]*player)
var db *gorm.DB
var data *mirData
var conf *Config

// Environ 主游戏环境
type Environ struct {
	Peer cellnet.GenericPeer
	maps map[int]*mirMap // MapInfo.ID: Map
}

// NewEnviron 初始化
func NewEnviron(c *Config) *Environ {
	conf = c
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", conf.Mysql.Username, conf.Mysql.Password, conf.Mysql.Host, conf.Mysql.Port, conf.Mysql.DB))
	defer db.Close()
	if err != nil {
		panic(err)
	}
	data = newmirData()
	env := &Environ{}
	env.initMaps()
	return env
}

func (env *Environ) initMaps() {
	uppercaseNameRealNameMap := map[string]string{}
	files := cm.GetFiles(conf.Assets+"/Maps/", []string{".map"})
	for _, f := range files {
		uppercaseNameRealNameMap[strings.ToUpper(filepath.Base(f))] = f
	}
	// FIXME 开发只加载部分地图
	// allowarr := []int{1, 2, 3, 4, 6, 7, 8, 10, 11, 12, 13, 15, 16, 17, 18, 19, 20, 21, 22, 24, 26, 27, 28, 29, 30, 31, 32, 25, 144, 384}
	allowarr := []int{1}
	allow := map[int]bool{}
	for _, v := range allowarr {
		allow[v] = true
	}
	env.maps = map[int]*mirMap{}
	for _, mi := range data.mapInfos {
		if _, ok := allow[mi.ID]; !ok {
			continue
		}
		m := loadMap(uppercaseNameRealNameMap[strings.ToUpper(mi.Filename+".map")])
		m.info = mi
		m.info.Filename = strings.ToUpper(mi.Filename)
		// if err := m.InitAll(); err != nil {
		// 	panic(err)
		// }
		env.maps[mi.ID] = m
	}
}

// Update 更新游戏状态
func (env *Environ) Update() {
	// log.Debugln("Update")
}

// HandleEvent 处理客户端包
func (env *Environ) HandleEvent(e cellnet.Event) {
	s := e.Session()

	switch msg := e.Message().(type) {
	case *cellnet.SessionAccepted: // 有新的连接
		s.Send(&server.Connected{})
	case *cellnet.SessionClosed: // 有连接断开
		// sessionClosed(s, msg)
	case *client.ClientVersion:
		clientVersion(s, msg)
	case *client.KeepAlive:
		s.Send(&server.KeepAlive{Time: 0})
	case *client.NewAccount:
		newAccount(s, msg)
	case *client.ChangePassword:
		changePassword(s, msg)
	case *client.Login:
		login(s, msg, env)
	case *client.NewCharacter:
		newCharacter(s, msg, env)
	case *client.DeleteCharacter:
		deleteCharacter(s, msg)
	case *client.StartGame:
		startGame(s, msg)
	case *client.LogOut:
		logout(s, msg)
	default:
		if !checkGameStage(s, GAME) {
			return
		}
		p := sessionPlayer[s.ID()]
		handleEvent(p, e, s)
	}
}

func checkGameStage(s cellnet.Session, gameStage int) bool {
	player, ok := sessionPlayer[s.ID()]
	if !ok {
		return false
	}
	if player.gameStage != gameStage {
		return false
	}
	return true
}

func clientVersion(s cellnet.Session, msg *client.ClientVersion) {
	player := new(player)
	player.gameStage = LOGIN
	sessionPlayer[s.ID()] = player
	s.Send(&server.ClientVersion{Result: 1})
}

func newAccount(s cellnet.Session, msg *client.NewAccount) {
	if !checkGameStage(s, LOGIN) {
		return
	}
	res := uint8(0)
	a := new(orm.Account)
	db.Table("account").Where("username = ?", msg.UserName).Find(a)
	if a.ID == 0 {
		a.Username = msg.AccountID
		a.Password = msg.Password
		db.Table("account").Create(a)
		res = 8
	}
	s.Send(&server.NewAccount{Result: res})
}

func changePassword(s cellnet.Session, msg *client.ChangePassword) {
	if !checkGameStage(s, LOGIN) {
		return
	}
	res := uint8(5)
	a := new(orm.Account)
	db.Table("account").Where("username = ? AND password = ?", msg.AccountID, msg.CurrentPassword).Find(a)
	if a.ID != 0 {
		a.Password = msg.NewPassword
		db.Table("account").Model(a).Updates(orm.Account{Password: msg.NewPassword})
		res = 6
	}
	s.Send(&server.ChangePassword{Result: res})
}

func login(s cellnet.Session, msg *client.Login, env *Environ) {
	if !checkGameStage(s, LOGIN) {
		return
	}
	a := new(orm.Account)
	db.Table("account").Where("username = ? AND password = ?", msg.AccountID, msg.Password).Find(a)
	if a.ID == 0 {
		s.Send(&server.Login{Result: uint8(4)})
		return
	}
	res := new(server.LoginSuccess)
	res.Characters = nil // TODO 查询角色
	s.Send(res)
}

func newCharacter(s cellnet.Session, msg *client.NewCharacter, env *Environ) {
	if !checkGameStage(s, SELECT) {
		return
	}
	player := sessionPlayer[s.ID()]

	ac := make([]orm.AccountCharacter, 3)
	db.Table("account_character").Where("account_id = ?", player.accountID).Limit(3).Find(&ac)
	if len(ac) >= 3 {
		s.Send(&server.NewCharacter{Result: uint8(4)})
		return
	}

	// c := new(orm.Character)
	// TODO 判断角色名字是否重复

	res := new(server.NewCharacterSuccess)
	res.CharInfo = server.SelectInfo{
		Index:      1, // TODO 顺序 uint32(c.ID)
		Name:       msg.Name,
		Level:      1,          // uint16
		Class:      msg.Class,  // cm.MirClass
		Gender:     msg.Gender, // cm.MirGender
		LastAccess: 0,          // int6
	}
	s.Send(res)
}

func deleteCharacter(s cellnet.Session, msg *client.DeleteCharacter) {
	if !checkGameStage(s, SELECT) {
		return
	}
	c := new(orm.Character)
	db.Table("character").Where("id = ?", msg.CharacterIndex).Find(c)
	if c.ID == 0 {
		res := new(server.DeleteCharacter)
		res.Result = 4
		s.Send(res)
		return
	}
	db.Table("character").Delete(c)
	db.Table("account_character").Where("character_id = ?", c.ID).Delete(orm.Character{})
	res := new(server.DeleteCharacterSuccess)
	res.CharacterIndex = msg.CharacterIndex
	s.Send(res)
}

func startGame(s cellnet.Session, msg *client.StartGame) {
	if !checkGameStage(s, SELECT) {
		return
	}
	// TODO
	p := sessionPlayer[s.ID()]
	p.receiveChat("[欢迎进入游戏，如有任何建议、疑问欢迎交流。联系QQ群：32309474]", cm.ChatTypeHint)
	// p.enqueueItemInfos()
	// p.refreshStats()
	// p.enqueueQuestInfo()
	p.enqueue(&server.MapInformation{
		FileName:     p.currentMap.info.Filename,
		Title:        p.currentMap.info.Title,
		MiniMap:      uint16(p.currentMap.info.MiniMap),
		BigMap:       uint16(p.currentMap.info.BigMap),
		Lights:       cm.LightSetting(p.currentMap.info.Light),
		Lightning:    true,
		MapDarkLight: 0,
		Music:        uint16(p.currentMap.info.Music),
	})
	p.enqueue(&server.UserInformation{})
	// p.enqueue(&server.TimeOfDay{Lights: env.Lights})
	// p.enqueue(&server.ChangeAMode{Mode: p.AMode})
	// p.enqueue(&server.ChangePMode{Mode: p.PMode})
	// p.enqueue(&server.SwitchGroup{AllowGroup: p.AllowGroup})
	// p.enqueueAreaObjects(nil, p.GetCell())
	// p.enqueue(ServerMessage{}.NPCResponse([]string{}))
	// p.broadcast(ServerMessage{}.ObjectPlayer(p))
}

func logout(s cellnet.Session, msg *client.LogOut) {

}

func handleEvent(p *player, e cellnet.Event, s cellnet.Session) {
	switch msg := e.Message().(type) {
	case *client.Turn:
		p.turn(msg)
	case *client.Walk:
		p.walk(msg)
	case *client.Run:
		p.run(msg)
	case *client.Chat:
		p.chat(msg)
	case *client.MoveItem:
		p.moveItem(msg)
	case *client.StoreItem:
		p.storeItem(msg)
	case *client.DepositRefineItem:
		p.depositRefineItem(msg)
	case *client.RetrieveRefineItem:
		p.retrieveRefineItem(msg)
	case *client.RefineCancel:
		p.refineCancel(msg)
	case *client.RefineItem:
		p.refineItem(msg)
	case *client.CheckRefine:
		p.checkRefine(msg)
	case *client.ReplaceWedRing:
		p.replaceWedRing(msg)
	case *client.DepositTradeItem:
		p.depositTradeItem(msg)
	case *client.RetrieveTradeItem:
		p.retrieveTradeItem(msg)
	case *client.TakeBackItem:
		p.takeBackItem(msg)
	case *client.MergeItem:
		p.mergeItem(msg)
	case *client.EquipItem:
		p.equipItem(msg)
	case *client.RemoveItem:
		p.removeItem(msg)
	case *client.RemoveSlotItem:
		p.removeSlotItem(msg)
	case *client.SplitItem:
		p.splitItem(msg)
	case *client.UseItem:
		p.useItem(msg)
	case *client.DropItem:
		p.dropItem(msg)
	case *client.DropGold:
		p.dropGold(msg)
	case *client.PickUp:
		p.pickUp(msg)
	case *client.Inspect:
		p.inspect(msg)
	// case *client.ChangeAMode:
	// 	p.ChangeAMode(msg)
	default:
		log.Debugln("default:", msg)
		//MessageQueue.Enqueue(string.Format("Invalid packet received. Index : {0}", p.Index));
	}
}
