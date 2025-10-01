package factory

import (
	"fmt"
	"os"

	"github.com/rivo/tview"
)

// PageFactory returns a tview.Primitive, encapsulates page creation logic
type PageFactory func() tview.Primitive

// Refresher updates an existing primitive when shown
type Refresher func(p tview.Primitive)

// gets called when the page[pageKey] is closed (out of view)
type OnClose func()

// PageManager holds factories, optional (!!!) cache, and a navigation stack
type PageManager struct {
	App       *tview.Application
	Pages     *tview.Pages
	factories map[string]PageFactory
	closers   map[string]OnClose
	refresh   map[string]Refresher
	cache     map[string]tview.Primitive
	stack     []string
}

// creates a new pagemanager instance
func NewPageManager(app *tview.Application) *PageManager {
	return &PageManager{
		App:       app,
		Pages:     tview.NewPages(),
		factories: make(map[string]PageFactory),
		closers:   make(map[string]OnClose),
		refresh:   make(map[string]Refresher),
		cache:     make(map[string]tview.Primitive),
		stack:     []string{},
	}
}

// RegisterFactory registers how to construct a page and an optional refresher
// if refresher is present, it will be called when page is shown
func (pm *PageManager) RegisterFactory(name string, factory PageFactory, refresher Refresher, onClose OnClose) {
	pm.factories[name] = factory
	pm.closers[name] = onClose
	if refresher != nil {
		pm.refresh[name] = refresher
	}
}

// showFresh creates (or recreates) primitive and adds it to pages
func (pm *PageManager) showFresh(name string, keepCached bool) {
	if pm.Pages.HasPage(name) {
		pm.Pages.RemovePage(name)
		delete(pm.cache, name)
	}

	factory, ok := pm.factories[name]
	if !ok {
		return
	}

	p := factory()
	pm.Pages.AddPage(name, p, true, false)

	if keepCached {
		pm.cache[name] = p
	}
}

// navigates to pageName. if cached primitive exists, it will be refreshed (if refresher present)
// elif no cached primitive, creates it via corresponding factory
// keepCached controls whether to keep instance
func (pm *PageManager) Push(pageName string, keepCached bool) {
	if len(pm.stack) > 0 {
		cur := pm.stack[len(pm.stack)-1]
		pm.closers[cur]()
		pm.Pages.HidePage(cur)
	}

	if _, exists := pm.cache[pageName]; !exists {
		pm.showFresh(pageName, keepCached)
	}

	if p, exists := pm.cache[pageName]; exists {
		if r, ok := pm.refresh[pageName]; ok && r != nil {
			r(p)
		}
	}
	pm.Pages.ShowPage(pageName)
	pm.stack = append(pm.stack, pageName)
}

// navigate backwards once
func (pm *PageManager) Pop() {
	if len(pm.stack) <= 1 {
		os.Exit(1)
	}

	top := pm.stack[len(pm.stack)-1]
	fmt.Println(top)
	pm.closers[top]()
	pm.stack = pm.stack[:len(pm.stack)-1]
	pm.Pages.HidePage(top)
	prev := pm.stack[len(pm.stack)-1]

	if p, ok := pm.cache[prev]; ok {
		if refr, ok2 := pm.refresh[prev]; ok2 && refr != nil {
			refr(p)
		}
	}

	pm.Pages.ShowPage(prev)
}

// check if page exists (is factory for pageKey registered)
func (pm *PageManager) HasPage(pageKey string) bool {
	_, ok := pm.factories[pageKey]
	return ok
}
