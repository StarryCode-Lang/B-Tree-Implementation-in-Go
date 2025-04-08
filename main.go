package main

import (
	"fmt"
	"sort"
)

// MaxKeys 定义每个节点能够存储的最大关键字数量（适用于叶节点和内部节点）
const MaxKeys = 3

// getMinKeys 计算并返回叶节点和内部节点所需的最小关键字数量
func getMinKeys() int {
	// 如果最大关键字数为偶数，则最小值为其一半
	if MaxKeys%2 == 0 {
		return MaxKeys / 2
	}
	// 如果最大关键字数为奇数，则最小值为其一半向上取整
	return (MaxKeys + 1) / 2
}

// Node 表示 B+ 树的节点
type Node struct {
	isLeaf   bool    // 是否为叶节点
	keys     []int   // 对于叶节点：存储键；对于内部节点：每个关键词为对应子节点的最大键
	parent   *Node   // 指向父节点
	values   []int   // 仅叶节点有效：保存对应的值
	next     *Node   // 仅叶节点有效：链表指针
	children []*Node // 仅内部节点有效：指向子节点
}

// NewNode 创建一个新节点
func NewNode(isLeaf bool) *Node {
	return &Node{
		isLeaf:   isLeaf,
		keys:     make([]int, 0, MaxKeys),
		parent:   nil,
		values:   make([]int, 0, MaxKeys),
		next:     nil,
		children: make([]*Node, 0, MaxKeys+1),
	}
}

// BPlusTree 表示 B+ 树
type BPlusTree struct {
	root *Node
}

// NewBPlusTree 创建一个新的 B+ 树
func NewBPlusTree() *BPlusTree {
	return &BPlusTree{
		root: NewNode(true),
	}
}

// 更新内部节点的关键词：每个关键词等于对应子节点的最大键
func (bpt *BPlusTree) updateInternalKeys(node *Node) {
	if node == nil || node.isLeaf {
		return
	}
	node.keys = []int{}
	for _, child := range node.children {
		// 每个子节点至少有一个键
		node.keys = append(node.keys, child.keys[len(child.keys)-1])
	}
}

// 若孩子结点的最大键发生变化，则向上更新父节点中的对应关键词
func (bpt *BPlusTree) updateParent(child *Node) {
	if child.parent == nil {
		return
	}
	parent := child.parent
	for i, c := range parent.children {
		if c == child {
			parent.keys[i] = child.keys[len(child.keys)-1]
			break
		}
	}
	bpt.updateParent(parent)
}

// 从根开始查找应存放 key 的叶节点
func (bpt *BPlusTree) findLeaf(node *Node, key int) *Node {
	if node.isLeaf {
		return node
	}
	for i, k := range node.keys {
		if key <= k {
			return bpt.findLeaf(node.children[i], key)
		}
	}
	return bpt.findLeaf(node.children[len(node.children)-1], key)
}

// 叶节点分裂：当叶节点中键数超过 MaxKeys 时
func (bpt *BPlusTree) splitLeaf(leaf *Node) {
	newLeaf := NewNode(true)
	newLeaf.parent = leaf.parent
	total := len(leaf.keys)
	mid := total / 2 // 前 mid 个保留，后半部分移至 newLeaf

	// 分裂时同步分裂 keys 与 values
	newLeaf.keys = append(newLeaf.keys, leaf.keys[mid:]...)
	newLeaf.values = append(newLeaf.values, leaf.values[mid:]...)
	leaf.keys = leaf.keys[:mid]
	leaf.values = leaf.values[:mid]

	// 调整链表指针
	newLeaf.next = leaf.next
	leaf.next = newLeaf

	if leaf.parent == nil {
		// 当前叶为根，则构造新根（内部节点）
		newRoot := NewNode(false)
		newRoot.children = append(newRoot.children, leaf)
		newRoot.children = append(newRoot.children, newLeaf)
		newRoot.keys = append(newRoot.keys, leaf.keys[len(leaf.keys)-1])
		newRoot.keys = append(newRoot.keys, newLeaf.keys[len(newLeaf.keys)-1])
		leaf.parent = newRoot
		newLeaf.parent = newRoot
		bpt.root = newRoot
	} else {
		parent := leaf.parent
		// 在父节点中找到 leaf 的位置，并在其后插入 newLeaf
		pos := 0
		for pos < len(parent.children) && parent.children[pos] != leaf {
			pos++
		}
		// 插入子节点到 children 切片
		parent.children = append(parent.children, nil)
		copy(parent.children[pos+2:], parent.children[pos+1:])
		parent.children[pos+1] = newLeaf
		// 插入关键字到 keys 切片
		parent.keys = append(parent.keys, 0)
		copy(parent.keys[pos+2:], parent.keys[pos+1:])
		parent.keys[pos+1] = newLeaf.keys[len(newLeaf.keys)-1]
		parent.keys[pos] = leaf.keys[len(leaf.keys)-1]
		newLeaf.parent = parent
		if len(parent.children) > MaxKeys {
			bpt.splitInternal(parent)
		} else {
			bpt.updateParent(newLeaf)
		}
	}
}

// 内部节点分裂：当内部节点的子节点数超过 MaxKeys 时
func (bpt *BPlusTree) splitInternal(node *Node) {
	newNode := NewNode(false)
	newNode.parent = node.parent
	totalChildren := len(node.children)
	mid := totalChildren / 2 // 左侧保留 mid 个子节点，右侧移至 newNode
	newNode.children = append(newNode.children, node.children[mid:]...)
	for _, child := range newNode.children {
		child.parent = newNode
	}
	node.children = node.children[:mid]
	bpt.updateInternalKeys(node)
	bpt.updateInternalKeys(newNode)

	if node.parent == nil {
		newRoot := NewNode(false)
		newRoot.children = append(newRoot.children, node)
		newRoot.children = append(newRoot.children, newNode)
		newRoot.keys = append(newRoot.keys, node.keys[len(node.keys)-1])
		newRoot.keys = append(newRoot.keys, newNode.keys[len(newNode.keys)-1])
		node.parent = newRoot
		newNode.parent = newRoot
		bpt.root = newRoot
	} else {
		parent := node.parent
		pos := 0
		for pos < len(parent.children) && parent.children[pos] != node {
			pos++
		}
		// 插入子节点到 children 切片
		parent.children = append(parent.children, nil)
		copy(parent.children[pos+2:], parent.children[pos+1:])
		parent.children[pos+1] = newNode
		// 插入关键字到 keys 切片
		parent.keys = append(parent.keys, 0)
		copy(parent.keys[pos+2:], parent.keys[pos+1:])
		parent.keys[pos+1] = newNode.keys[len(newNode.keys)-1]
		parent.keys[pos] = node.keys[len(node.keys)-1]
		newNode.parent = parent
		if len(parent.children) > MaxKeys {
			bpt.splitInternal(parent)
		} else {
			bpt.updateParent(newNode)
		}
	}
}

// 删除后对节点进行借补或合并，保证节点达到最少关键字数要求
func (bpt *BPlusTree) rebalance(node *Node) {
	// 若 node 为根节点，特殊处理
	if node == bpt.root {
		// 若根为内部节点且只有一个子节点，则下降为新根
		if !node.isLeaf && len(node.children) == 1 {
			newRoot := node.children[0]
			newRoot.parent = nil
			bpt.root = newRoot
			// 在 Go 中，内存由垃圾回收器管理，不需要显式删除
		}
		return
	}
	minRequired := getMinKeys() // 对于叶节点与内部节点均采用同一标准（非根节点最少关键字数）
	if len(node.keys) >= minRequired {
		return // 已满足最小要求
	}

	parent := node.parent
	// 在父节点中找到 node 的位置
	index := 0
	for index < len(parent.children) && parent.children[index] != node {
		index++
	}
	var leftSibling *Node
	var rightSibling *Node
	if index-1 >= 0 {
		leftSibling = parent.children[index-1]
	}
	if index+1 < len(parent.children) {
		rightSibling = parent.children[index+1]
	}

	if node.isLeaf {
		// 叶节点：先尝试从左侧兄弟借补
		if leftSibling != nil && len(leftSibling.keys) > minRequired {
			// 从左侧兄弟借最后一个键值对
			borrowedKey := leftSibling.keys[len(leftSibling.keys)-1]
			borrowedValue := leftSibling.values[len(leftSibling.values)-1]
			leftSibling.keys = leftSibling.keys[:len(leftSibling.keys)-1]
			leftSibling.values = leftSibling.values[:len(leftSibling.values)-1]
			node.keys = append([]int{borrowedKey}, node.keys...)
			node.values = append([]int{borrowedValue}, node.values...)
			bpt.updateInternalKeys(parent)
			return
		} else if rightSibling != nil && len(rightSibling.keys) > minRequired {
			// 从右侧兄弟借第一个键值对
			borrowedKey := rightSibling.keys[0]
			borrowedValue := rightSibling.values[0]
			rightSibling.keys = rightSibling.keys[1:]
			rightSibling.values = rightSibling.values[1:]
			node.keys = append(node.keys, borrowedKey)
			node.values = append(node.values, borrowedValue)
			bpt.updateInternalKeys(parent)
			return
		} else {
			// 无法借补，则合并节点（优先与左侧合并）
			if leftSibling != nil {
				// 将当前节点的内容合并到左侧兄弟
				leftSibling.keys = append(leftSibling.keys, node.keys...)
				leftSibling.values = append(leftSibling.values, node.values...)
				leftSibling.next = node.next
				// 在父节点中删除当前节点对应的指针和关键字
				parent.children = append(parent.children[:index], parent.children[index+1:]...)
				parent.keys = append(parent.keys[:index], parent.keys[index+1:]...)
				// 在 Go 中，不需要显式删除节点，垃圾回收器会处理
				bpt.updateInternalKeys(parent)
				bpt.rebalance(parent)
			} else if rightSibling != nil {
				// 将右侧兄弟合并到当前节点
				node.keys = append(node.keys, rightSibling.keys...)
				node.values = append(node.values, rightSibling.values...)
				node.next = rightSibling.next
				parent.children = append(parent.children[:index+1], parent.children[index+2:]...)
				parent.keys = append(parent.keys[:index+1], parent.keys[index+2:]...)
				// 在 Go 中，不需要显式删除节点，垃圾回收器会处理
				bpt.updateInternalKeys(parent)
				bpt.rebalance(parent)
			}
		}
	} else {
		// 内部节点：处理方式与叶节点类似，不过借补或合并时调整的是子节点指针
		if leftSibling != nil && len(leftSibling.keys) > minRequired {
			// 从左侧兄弟借出其最后一个子节点
			borrowedChild := leftSibling.children[len(leftSibling.children)-1]
			leftSibling.children = leftSibling.children[:len(leftSibling.children)-1]
			leftSibling.keys = leftSibling.keys[:len(leftSibling.keys)-1]
			node.children = append([]*Node{borrowedChild}, node.children...)
			borrowedChild.parent = node
			bpt.updateInternalKeys(leftSibling)
			bpt.updateInternalKeys(node)
			bpt.updateInternalKeys(parent)
			return
		} else if rightSibling != nil && len(rightSibling.keys) > minRequired {
			// 从右侧兄弟借出其第一个子节点
			borrowedChild := rightSibling.children[0]
			rightSibling.children = rightSibling.children[1:]
			rightSibling.keys = rightSibling.keys[1:]
			node.children = append(node.children, borrowedChild)
			borrowedChild.parent = node
			bpt.updateInternalKeys(rightSibling)
			bpt.updateInternalKeys(node)
			bpt.updateInternalKeys(parent)
			return
		} else {
			// 合并内部节点（优先与左侧合并）
			if leftSibling != nil {
				// 将当前节点的所有子节点合并到左侧兄弟
				for _, child := range node.children {
					leftSibling.children = append(leftSibling.children, child)
					child.parent = leftSibling
				}
				bpt.updateInternalKeys(leftSibling)
				parent.children = append(parent.children[:index], parent.children[index+1:]...)
				parent.keys = append(parent.keys[:index], parent.keys[index+1:]...)
				// 在 Go 中，不需要显式删除节点，垃圾回收器会处理
				bpt.updateInternalKeys(parent)
				bpt.rebalance(parent)
			} else if rightSibling != nil {
				// 将右侧兄弟的所有子节点合并到当前节点
				for _, child := range rightSibling.children {
					node.children = append(node.children, child)
					child.parent = node
				}
				bpt.updateInternalKeys(node)
				parent.children = append(parent.children[:index+1], parent.children[index+2:]...)
				parent.keys = append(parent.keys[:index+1], parent.keys[index+2:]...)
				// 在 Go 中，不需要显式删除节点，垃圾回收器会处理
				bpt.updateInternalKeys(parent)
				bpt.rebalance(parent)
			}
		}
	}
}

// Insert 插入操作：在叶节点中插入 key 与 value，并在必要时分裂
func (bpt *BPlusTree) Insert(key, value int) {
	leaf := bpt.findLeaf(bpt.root, key)
	pos := sort.SearchInts(leaf.keys, key)

	// Insert key and value
	leaf.keys = append(leaf.keys, 0)
	copy(leaf.keys[pos+1:], leaf.keys[pos:])
	leaf.keys[pos] = key
	leaf.values = append(leaf.values, 0)
	copy(leaf.values[pos+1:], leaf.values[pos:])
	leaf.values[pos] = value

	// Update parent only if the new key is the maximum and differs from the old maximum
	if pos == len(leaf.keys)-1 && (len(leaf.keys) == 1 || key > leaf.keys[len(leaf.keys)-2]) {
		bpt.updateParent(leaf)
	}

	if len(leaf.keys) > MaxKeys {
		bpt.splitLeaf(leaf)
	}
}

func (bpt *BPlusTree) Remove(key int) error {
	leaf := bpt.findLeaf(bpt.root, key)
	pos := sort.SearchInts(leaf.keys, key)
	if pos >= len(leaf.keys) || leaf.keys[pos] != key {
		return fmt.Errorf("删除失败：未找到 key = %d", key)
	}

	leaf.keys = append(leaf.keys[:pos], leaf.keys[pos+1:]...)
	leaf.values = append(leaf.values[:pos], leaf.values[pos+1:]...)
	if pos == len(leaf.keys) {
		bpt.updateParent(leaf)
	}
	if leaf != bpt.root && len(leaf.keys) < getMinKeys() {
		bpt.rebalance(leaf)
	}
	return nil
}

func (bpt *BPlusTree) Modify(key, newValue int) error {
	leaf := bpt.findLeaf(bpt.root, key)
	pos := sort.SearchInts(leaf.keys, key)
	if pos >= len(leaf.keys) || leaf.keys[pos] != key {
		return fmt.Errorf("修改失败：未找到 key = %d", key)
	}
	leaf.values[pos] = newValue
	return nil
}

// Search 查找操作：返回 key 对应的 value；若不存在返回 -1
func (bpt *BPlusTree) Search(key int) int {
	leaf := bpt.findLeaf(bpt.root, key)

	// 查找键位置
	for i, k := range leaf.keys {
		if k == key {
			return leaf.values[i]
		}
	}

	return -1
}

// PrintTree 打印整棵树（层次遍历，用于调试）
func (bpt *BPlusTree) PrintTree() {
	current := []*Node{bpt.root}
	for len(current) > 0 {
		var next []*Node
		for _, node := range current {
			nodeType := "Leaf"
			if !node.isLeaf {
				nodeType = "Internal"
			}
			parentKey := -1
			if node.parent != nil && len(node.parent.keys) > 0 {
				parentKey = node.parent.keys[0]
			}
			fmt.Printf("[%s, %d: ", nodeType, parentKey)
			for _, k := range node.keys {
				fmt.Printf("%d ", k)
			}
			fmt.Print("]")
			if !node.isLeaf {
				fmt.Print("  ")
				next = append(next, node.children...)
			}
		}
		fmt.Println()
		current = next
	}
}

// PrintLeafValues 新增函数：一次性输出所有叶节点对应的值
func (bpt *BPlusTree) PrintLeafValues() {
	// 从根节点一路向下找到最左侧叶节点
	node := bpt.root
	for !node.isLeaf {
		node = node.children[0]
	}
	fmt.Print("所有叶节点对应的值：")
	for node != nil {
		for _, value := range node.values {
			fmt.Printf("%d ", value)
		}
		node = node.next
	}
	fmt.Println()
}

func main() {
	tree := NewBPlusTree()

	// 插入测试数据，覆盖多种插入情况
	fmt.Println("插入数据：")
	insertions := []struct{ key, value int }{
		{4, 40}, {7, 70}, {1, 10}, {9, 90}, {2, 20}, {5, 50}, {8, 80}, {3, 30}, {6, 60},
	}
	for _, pair := range insertions {
		fmt.Printf("插入 (%d, %d)\n", pair.key, pair.value)
		tree.Insert(pair.key, pair.value)
	}
	fmt.Println("初始B+树结构：")
	tree.PrintTree()
	tree.PrintLeafValues()

	// 删除测试，覆盖多种删除情况
	fmt.Println("\n删除 keys: 9, 1, 5, 2")
	tree.Remove(9) // 删除最大键，更新父节点
	fmt.Println("删除 9 后的结构：")
	tree.PrintTree()
	tree.PrintLeafValues()

	tree.Remove(1) // 删除不需调整的键
	fmt.Println("删除 1 后的结构：")
	tree.PrintTree()
	tree.PrintLeafValues()

	tree.Remove(5) // 触发借补
	fmt.Println("删除 5 后的结构：")
	tree.PrintTree()
	tree.PrintLeafValues()

	tree.Remove(2) // 触发合并
	fmt.Println("删除 2 后的结构：")
	tree.PrintTree()
	tree.PrintLeafValues()

	// 查找测试
	fmt.Printf("\n查找 key 7 对应的 value: %d\n", tree.Search(7))
}
