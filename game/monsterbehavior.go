package game

import "time"

type behaviorTree struct {
	monster *monster
	root    behavior
}

func newBehaviorTree(m *monster) *behaviorTree {
	return &behaviorTree{
		monster: m,
		root:    newRootNode(m),
	}
}

func (bt *behaviorTree) update(now time.Time) {

}

type behavior interface {
	tick(time.Time) status
}

type node struct {
	children []behavior
}

func (n *node) addChild(b behavior) {
	n.children = append(n.children, b)
}

// 控制节点 - 选择
// 顺序执行所有的子节点，当一个子节点执行结果为 SUCCESS 的时候终止执行并返回 SUCCESS
// 选择节点可以被理解为一个或门（OR gate）
type selectNode struct {
	node
}

func (n *selectNode) tick(now time.Time) status {
	res := FAILED
	for _, child := range n.children {
		s := child.tick(now)
		if s == SUCCESS {
			return SUCCESS
		}
		if s == RUNNING {
			res = RUNNING
		}
	}
	return res
}

// 控制节点 - 序列
// 将其所有子节点依次执行，即当前执行的一个子节点返回成功后，再执行下一个子节点
// 顺序依次执行子节点，如果所有子节点都返回 SUCCESS，则向其父节点返回 SUCCESS
// 序列节点可以理解为与门（AND gate）
type sequenceNode struct {
	node
}

func (n *sequenceNode) tick(now time.Time) status {
	res := SUCCESS
	for _, child := range n.children {
		s := child.tick(now)
		if s != SUCCESS {
			return s
		}
	}
	return res
}

// 控制节点 - 并行
// 将其所有子节点都运行一遍，不管运行结果
type parallelNode struct {
	node
}

// 条件节点 执行返回 status
type conditionNode struct {
	node
	fn func() bool
}

func (n *conditionNode) tick(now time.Time) status {
	if n.fn() {
		return SUCCESS
	}
	return FAILED
}

// 行为节点 执行返回 status
type actionNode struct {
	node
	fn func() status
}

func (n *actionNode) tick(now time.Time) status {
	return n.fn()
}

func newRootNode(m *monster) behavior {
	switch m.info.AI {
	default:
		return defaultRoot()
	}
}

func defaultRoot() behavior {
	return nil
}
