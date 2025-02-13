package servers

import (
	"go_learn_project_rest_api/modules/files/fileHandlers"
	"go_learn_project_rest_api/modules/files/fileUsecases"
)

type IFilesModule interface {
	Init()
	Usecase() fileUsecases.IFileUsecases
	Handler() fileHandlers.IFileHandlers
}

type filesModule struct {
	*moduleFactory
	usecase fileUsecases.IFileUsecases
	handler fileHandlers.IFileHandlers
}

func (m *moduleFactory) FilesModule() IFilesModule {
	usecase := fileUsecases.FileUsecases(m.server.cfg)
	handlers := fileHandlers.FileHandlers(m.server.cfg, usecase)
	return &filesModule{
		moduleFactory: m,
		usecase:       usecase,
		handler:       handlers,
	}
}

func (f *filesModule) Init() {
	router := f.router.Group("/files")
	router.Post("/upload", f.handler.UploadFiles, f.mid.JwtAuth(), f.mid.Authorize(2))
	router.Post("/delete", f.handler.DeleteFile, f.mid.JwtAuth(), f.mid.Authorize(2))
}

func (f *filesModule) Usecase() fileUsecases.IFileUsecases {
	return f.usecase
}

func (f *filesModule) Handler() fileHandlers.IFileHandlers {
	return f.handler
}
