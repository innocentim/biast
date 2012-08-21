package main

import (
	"bytes"
	"net/http"
	"sync"
)

var indexCache bytes.Buffer
var indexCacheMutex sync.RWMutex

func updateIndex() {
	// TODO pager
	indexList := artMgr.atomGetAllArticles()
	qsortForArticleList(indexList, 0, len(indexList)-1)
	indexCacheMutex.Lock()
	indexCache.Reset()
	tmpl.ExecuteTemplate(&indexCache, "index", map[string]interface{}{
		"config":   config,
		"articles": indexList,
	})
	indexCacheMutex.Unlock()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexCacheMutex.RLock()
	w.Write(indexCache.Bytes())
	indexCacheMutex.RUnlock()
}

func initPageIndex() {
	http.HandleFunc(config["RootUrl"], indexHandler)
	updateIndex()
}

func qsortForArticleList(a []*Article, l, r int) {
	if l > r {
		return
	}
	i := l
	j := (r-l)/2 + l
	a[i], a[j] = a[j], a[i]
	j = l
	for i = l + 1; i <= r; i++ {
		if a[i].Id < a[l].Id {
			j++
			a[j], a[i] = a[i], a[j]
		}
	}
	a[j], a[l] = a[l], a[j]
	qsortForArticleList(a, l, j-1)
	qsortForArticleList(a, j+1, r)
}
