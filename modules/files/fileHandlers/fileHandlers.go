package fileHandlers

import (
	"fmt"
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/entities"
	"go_learn_project_rest_api/modules/files"
	"go_learn_project_rest_api/modules/files/fileUsecases"
	"go_learn_project_rest_api/pkgs/utils"
	"math"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type IFileHandlers interface {
	UploadFiles(fiber.Ctx) error
	DeleteFile(fiber.Ctx) error
}

type fileHandlerErrCode string

const (
	uploadErr fileHandlerErrCode = "files-001"
	deleteErr fileHandlerErrCode = "files-002"
)

type fileHandlers struct {
	cfg          config.IConfig
	fileUsecases fileUsecases.IFileUsecases
}

func FileHandlers(cfg config.IConfig, fileUsecases fileUsecases.IFileUsecases) IFileHandlers {
	return &fileHandlers{
		cfg:          cfg,
		fileUsecases: fileUsecases,
	}
}

func (h *fileHandlers) UploadFiles(c fiber.Ctx) error {
	req := make([]*files.FileReq, 0)

	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(uploadErr),
			err.Error(),
		).Res()
	}

	filesReq := form.File["files"]
	destination := c.FormValue("destination")

	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
	}

	for _, file := range filesReq {
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(uploadErr),
				"extension is not acceptable",
			).Res()
		}

		if file.Size > int64(h.cfg.App().FileLimit()) {
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(uploadErr),
				fmt.Sprintf("file size must less than %d MiB", int(math.Ceil(float64(h.cfg.App().FileLimit())/math.Pow(1024, 2)))),
			).Res()
		}

		fileName := utils.RandFileName(ext)
		req = append(req, &files.FileReq{
			File:        file,
			Destination: destination,
			FileName:    fileName,
			Extension:   ext,
		})
	}
	res, err := h.fileUsecases.UploadToGCP(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(uploadErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).SuccessResponse(fiber.StatusCreated, res).Res()
}

func (h *fileHandlers) DeleteFile(c fiber.Ctx) error {
	req := make([]*files.DeleteFileReq, 0)
	if err := c.Bind().Body(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(deleteErr),
			err.Error(),
		).Res()
	}

	if err := h.fileUsecases.DeleteFileOnGCP(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(deleteErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).SuccessResponse(fiber.StatusOK, nil).Res()
}
