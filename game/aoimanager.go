package game

import (
	"fmt"

	"github.com/yenkeia/yams/game/cm"
)

type aoiManager struct {
	minX  int              //区域左边界坐标
	maxX  int              //区域右边界坐标
	cntsX int              //x方向格子的数量
	minY  int              //区域上边界坐标
	maxY  int              //区域下边界坐标
	cntsY int              //y方向的格子数量
	grids map[int]*aoiGrid //当前区域中都有哪些格子，key=格子ID， value=格子对象
}

var gridLength = 10

func newAOIManager(width, height int) *aoiManager {
	minX := 0
	maxX := width
	cntsX := width / gridLength
	minY := 0
	maxY := height
	cntsY := height / gridLength
	return _newAOIManager(minX, maxX, cntsX, minY, maxY, cntsY)
}

// 初始化一个AOI区域
func _newAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *aoiManager {
	aoiMgr := &aoiManager{
		minX:  minX,
		maxX:  maxX,
		cntsX: cntsX,
		minY:  minY,
		maxY:  maxY,
		cntsY: cntsY,
		grids: make(map[int]*aoiGrid),
	}

	//给AOI初始化区域中所有的格子
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			//计算格子ID
			//格子编号：id = idy *nx + idx  (利用格子坐标得到格子编号)
			gid := y*cntsX + x

			//初始化一个格子放在AOI中的map里，key是当前格子的ID
			aoiMgr.grids[gid] = newAOIGrid(gid,
				aoiMgr.minX+x*aoiMgr.gridWidth(),
				aoiMgr.minX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.minY+y*aoiMgr.gridHeight(),
				aoiMgr.minY+(y+1)*aoiMgr.gridHeight())
		}
	}

	return aoiMgr
}

func (m *aoiManager) getGridByPoint(pos cm.Point) *aoiGrid {
	x := int(pos.X)
	xx := x / m.gridWidth()
	y := int(pos.Y)
	yy := y / m.gridHeight()
	gid := yy*m.cntsX + xx
	return m.grids[gid]
}

func (m *aoiManager) getSurroundGridsByPoint(pos cm.Point) (grids []*aoiGrid) {
	g := m.getGridByPoint(pos)
	return m.getSurroundGridsByGid(g.gID)
}

//根据格子的gID得到当前周边的九宫格信息
func (m *aoiManager) getSurroundGridsByGid(gID int) (grids []*aoiGrid) {
	//判断gID是否存在
	if _, ok := m.grids[gID]; !ok {
		return
	}

	//将当前gid添加到九宫格中
	grids = append(grids, m.grids[gID])

	//根据gid得到当前格子所在的X轴编号
	idx := gID % m.cntsX

	//判断当前idx左边是否还有格子
	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}
	//判断当前的idx右边是否还有格子
	if idx < m.cntsX-1 {
		grids = append(grids, m.grids[gID+1])
	}

	//将x轴当前的格子都取出，进行遍历，再分别得到每个格子的上下是否有格子

	//得到当前x轴的格子id集合
	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.gID)
	}

	//遍历x轴格子
	for _, v := range gidsX {
		//计算该格子处于第几列
		idy := v / m.cntsX

		//判断当前的idy上边是否还有格子
		if idy > 0 {
			grids = append(grids, m.grids[v-m.cntsX])
		}
		//判断当前的idy下边是否还有格子
		if idy < m.cntsY-1 {
			grids = append(grids, m.grids[v+m.cntsX])
		}
	}
	return
}

//得到每个格子在x轴方向的宽度
func (m *aoiManager) gridWidth() int {
	return (m.maxX - m.minX) / m.cntsX
}

//得到每个格子在x轴方向的长度
func (m *aoiManager) gridHeight() int {
	return (m.maxY - m.minY) / m.cntsY
}

//打印信息方法
func (m *aoiManager) String() string {
	s := fmt.Sprintf("AOIManagr:\nminX:%d, maxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d\n Grids in AOI Manager:\n",
		m.minX, m.maxX, m.cntsX, m.minY, m.maxY, m.cntsY)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}
