package factory

import (
	"github.com/komadiina/spelltext/utils/singleton/logging"
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
	Logger    *logging.Logger
	factories map[string]PageFactory
	closers   map[string]OnClose
	refresh   map[string]Refresher
	cache     map[string]tview.Primitive
	stack     []string
}

// creates a new pagemanager instance
func NewPageManager(logger *logging.Logger, app *tview.Application) *PageManager {
	return &PageManager{
		App:       app,
		Pages:     tview.NewPages(),
		Logger:    logger,
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

func (pm *PageManager) Push(pageName string, keepCached bool) {
	if len(pm.stack) > 0 {
		cur := pm.stack[len(pm.stack)-1]
		if closer, ok := pm.closers[cur]; ok && closer != nil {
			closer()
		}
		pm.Pages.HidePage(cur)
	}

	if _, exists := pm.cache[pageName]; !exists {
		pm.showFresh(pageName, keepCached)
	} else {
		if _, hasRef := pm.refresh[pageName]; !hasRef {
			pm.showFresh(pageName, keepCached)
		} else {
			if p := pm.cache[pageName]; p != nil {
				if r, ok := pm.refresh[pageName]; ok && r != nil {
					r(p)
				}
			}
		}
	}

	pm.Pages.ShowPage(pageName)
	pm.stack = append(pm.stack, pageName)
}

// navigate backward
func (pm *PageManager) Pop() int {
	pm.Logger.Infof("navigating backward from: %v", pm.stack[len(pm.stack)-1])

	if len(pm.stack) <= 2 {
		return -1
	}

	top := pm.stack[len(pm.stack)-1]
	if closer, ok := pm.closers[top]; ok && closer != nil {
		closer()
	}
	pm.stack = pm.stack[:len(pm.stack)-1]
	pm.Pages.HidePage(top)

	prev := pm.stack[len(pm.stack)-1]

	if _, ok := pm.cache[prev]; ok {
		if _, hasRef := pm.refresh[prev]; !hasRef {
			pm.showFresh(prev, true)
		} else {
			if r, ok := pm.refresh[prev]; ok && r != nil {
				r(pm.cache[prev])
			}
		}
	} else {
		pm.showFresh(prev, true)
	}

	pm.Pages.ShowPage(prev)
	return 0
}

// check if page exists (is factory for pageKey registered)
func (pm *PageManager) HasPage(pageKey string) bool {
	_, ok := pm.factories[pageKey]
	return ok
}
