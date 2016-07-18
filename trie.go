package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//Trie 树
type Trie struct {
	Root *TrieNode
}

//TrieNode 子节点
type TrieNode struct {
	Children map[rune]*TrieNode
	IsWord   bool
	End      bool
}

//NewTrie 新建Trie数
func NewTrie() Trie {
	var r Trie
	r.Root = NewTrieNode()
	return r
}

//NewTrieNode 新建TrieNode节点
func NewTrieNode() *TrieNode {
	n := new(TrieNode)
	n.Children = make(map[rune]*TrieNode)
	return n
}

//Inster 插入单词
func (t *Trie) Inster(word string) {
	if len(word) < 1 {
		return
	}

	node := t.Root
	key := []rune(word)
	klen := len(key)
	for i := 0; i < klen; i++ {
		if _, exists := node.Children[key[i]]; !exists {
			node.Children[key[i]] = NewTrieNode()
		}
		node = node.Children[key[i]]
		//当前为完整词
		if i == klen-1 {
			node.IsWord = true
		}
	}

	node.End = true
}

//Find 查找是否包含某个关键字
func (t *Trie) Find(word string) bool {
	if len(word) < 1 {
		return false
	}

	node := t.Root
	key := []rune(word)
	klen := len(key)

	for i := 0; i < klen; i++ {
		if _, exists := node.Children[key[i]]; exists {
			node = node.Children[key[i]]

			for k := i + 1; k < klen; k++ {

				if _, exists := node.Children[key[k]]; !exists {
					break
				}

				node = node.Children[key[k]]

				if node.IsWord == true {
					return true
				}
			}
		}
	}

	return false
}

var trie Trie

func init() {
	fmt.Println("开始初始化关键词词库...")
	trie = NewTrie()
	trie.Inster("中华")
	trie.Inster("中国")
	trie.Inster("中国人")
	fmt.Println("初始化完成")
}

func main() {
	fmt.Println("开始启动HttpServer...")
	http.HandleFunc("/trie", trieHandler)
	http.ListenAndServe(":8080", nil)
}

//Result 输出结构体
type Result struct {
	Code int
	Flag bool
}

func trieHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	r.ParseForm()
	word := r.Form.Get("word")
	flag := trie.Find(word)

	if flag {
		fmt.Printf("\"%s\" 是关键词\n", word)
	} else {
		fmt.Printf("\"%s\" 不是关键词\n", word)
	}

	out := &Result{1, flag}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	w.Write(b)
}
