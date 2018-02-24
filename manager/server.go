package manager

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	. "github.com/leetpy/venv/internal"
)

const (
	_errorKey = "error"
	_infoKey  = "message"
)

var manager *Manager

type Manager struct {
	cfg       *Config
	engine    *gin.Engine
	adapter   dbAdapter
	httpClent *http.Client
}

func GetManager(cfg *Config) *Manager {
	if manager != nil {
		return manager
	}

	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	s := &Manager{
		cfg: cfg,
	}
	adapter, err := makeDBAdapter(cfg.DB.DBType, cfg.DB.DBFile)
	if err != nil {
		//logger.Errorf("Error initializing DB adapter: %s", err.Error())
		return nil
	}
	s.setDBAdapter(adapter)
	s.engine = gin.New()
	s.engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{_infoKey: "pong"})
	})

	s.engine.GET("/server", s.listAllServer)
	s.engine.DELETE("/server", s.delServer)
	s.engine.PATCH("/server", s.updateServer)
	s.engine.POST("/server", s.addServer)

	s.engine.GET("/types", s.listAllType)
	s.engine.POST("/type", s.addType)

	manager = s
	return s
}

func (s *Manager) setDBAdapter(adapter dbAdapter) {
	s.adapter = adapter
}

func (s *Manager) Run() {
	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Addr, s.cfg.Server.Port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (s *Manager) returnErrJSON(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{
		_errorKey: err.Error(),
	})
}

func (s *Manager) listAllServer(c *gin.Context) {
	//servers, err := s.adapter.ListServer()
}

func (s *Manager) delServer(c *gin.Context) {

}

func (s *Manager) updateServer(c *gin.Context) {

}

func (s *Manager) addServer(c *gin.Context) {

}

func (s *Manager) listAllType(c *gin.Context) {
	nodeTypeList, err := s.adapter.ListTypes()
	if err != nil {
		err := fmt.Errorf("failed to list all Type status: %s",
			err.Error(),
		)
		c.Error(err)
		s.returnErrJSON(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, nodeTypeList)
}

func (s *Manager) addType(c *gin.Context) {
	var _type NodeType
	c.BindJSON(&_type)
	newType, err := s.adapter.CreateType(_type)
	if err != nil {
		err := fmt.Errorf("failed to register worker: %s",
			err.Error(),
		)
		c.Error(err)
		s.returnErrJSON(c, http.StatusInternalServerError, err)
		return
	}

	//logger.Noticef("Worker <%s> registered", _worker.ID)
	// create workerCmd channel for this worker
	c.JSON(http.StatusOK, newType)
}