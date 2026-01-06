package http

import (
	"net/http"
	"path/filepath"
	"strconv"

	"icafe-registration/internal/config"
	"icafe-registration/internal/domain"
	"icafe-registration/pkg/response"

	"github.com/gin-gonic/gin"
)

// FileHandler represents the HTTP handler for files
type FileHandler struct {
	fileUsecase  domain.FileUsecase
	uploadConfig *config.UploadConfig
}

// NewFileHandler creates a new file handler
func NewFileHandler(router *gin.RouterGroup, engine *gin.Engine, uc domain.FileUsecase, uploadConfig *config.UploadConfig) {
	handler := &FileHandler{
		fileUsecase:  uc,
		uploadConfig: uploadConfig,
	}

	// File upload and management routes
	router.POST("/files/upload", handler.UploadFile)
	router.POST("/videos/upload", handler.UploadVideo)
	router.GET("/files", handler.GetAllFiles)
	router.GET("/videos", handler.GetAllVideos)
	router.GET("/files/:id", handler.GetFileByID)
	router.DELETE("/files/:id", handler.DeleteFile)

	// Static file serving for downloads and streaming
	filesPath := filepath.Join(uploadConfig.Path, "files")
	videosPath := filepath.Join(uploadConfig.Path, "videos")

	router.Static("/files/download", filesPath)
	router.Static("/videos/stream", videosPath)

	// Alternative: serve files with custom headers for proper download/streaming
	router.GET("/files/serve/:filename", handler.ServeFile)
	router.GET("/videos/serve/:filename", handler.ServeVideo)

	// Download by id
	router.GET("/files/download-by-id/:id", handler.DownloadFileByID)
}

// UploadFile godoc
// @Summary Upload a file
// @Description Upload a document file
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formance file true "File to upload"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /files/upload [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "No file provided", err.Error())
		return
	}

	uploadedFile, err := h.fileUsecase.Upload(c.Request.Context(), file, domain.FileTypeDocument)
	if err != nil {
		switch err {
		case domain.ErrFileTooLarge:
			response.BadRequest(c, "File too large", err.Error())
		case domain.ErrInvalidFileType:
			response.BadRequest(c, "Invalid file type", err.Error())
		default:
			response.InternalServerError(c, "Failed to upload file", err.Error())
		}
		return
	}

	response.Created(c, "File uploaded successfully", uploadedFile)
}

// UploadVideo godoc
// @Summary Upload a video
// @Description Upload a video file
// @Tags videos
// @Accept multipart/form-data
// @Produce json
// @Param file formance file true "Video file to upload"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /videos/upload [post]
func (h *FileHandler) UploadVideo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "No file provided", err.Error())
		return
	}

	uploadedFile, err := h.fileUsecase.Upload(c.Request.Context(), file, domain.FileTypeVideo)
	if err != nil {
		switch err {
		case domain.ErrFileTooLarge:
			response.BadRequest(c, "File too large", err.Error())
		case domain.ErrInvalidFileType:
			response.BadRequest(c, "Invalid file type", err.Error())
		default:
			response.InternalServerError(c, "Failed to upload video", err.Error())
		}
		return
	}

	response.Created(c, "Video uploaded successfully", uploadedFile)
}

// GetAllFiles godoc
// @Summary Get all files
// @Description Get all document files with pagination
// @Tags files
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /files [get]
func (h *FileHandler) GetAllFiles(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	files, total, err := h.fileUsecase.GetAll(c.Request.Context(), domain.FileTypeDocument, limit, offset)
	if err != nil {
		response.InternalServerError(c, "Failed to get files", err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Files retrieved successfully", files, &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// GetAllVideos godoc
// @Summary Get all videos
// @Description Get all video files with pagination
// @Tags videos
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /videos [get]
func (h *FileHandler) GetAllVideos(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	files, total, err := h.fileUsecase.GetAll(c.Request.Context(), domain.FileTypeVideo, limit, offset)
	if err != nil {
		response.InternalServerError(c, "Failed to get videos", err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Videos retrieved successfully", files, &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// GetFileByID godoc
// @Summary Get a file by ID
// @Description Get file information by ID
// @Tags files
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /files/{id} [get]
func (h *FileHandler) GetFileByID(c *gin.Context) {
	id := c.Param("id")

	file, err := h.fileUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrInvalidID:
			response.BadRequest(c, "Invalid ID format", err.Error())
		case domain.ErrNotFound:
			response.NotFound(c, "File not found")
		default:
			response.InternalServerError(c, "Failed to get file", err.Error())
		}
		return
	}

	response.OK(c, "File retrieved successfully", file)
}

// DeleteFile godoc
// @Summary Delete a file
// @Description Delete a file by ID
// @Tags files
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /files/{id} [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
	id := c.Param("id")

	err := h.fileUsecase.Delete(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrInvalidID:
			response.BadRequest(c, "Invalid ID format", err.Error())
		case domain.ErrNotFound:
			response.NotFound(c, "File not found")
		default:
			response.InternalServerError(c, "Failed to delete file", err.Error())
		}
		return
	}

	response.OK(c, "File deleted successfully", nil)
}

// ServeFile serves a file for download
func (h *FileHandler) ServeFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(h.uploadConfig.Path, "files", filename)

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")

	c.File(filePath)
}

// ServeVideo serves a video for streaming
func (h *FileHandler) ServeVideo(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(h.uploadConfig.Path, "videos", filename)

	// Set headers for video streaming
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Type", "video/mp4")

	c.File(filePath)
}

func (h *FileHandler) DownloadFileByID(c *gin.Context) {
	id := c.Param("id")

	// 1. Lấy thông tin file từ DB
	file, err := h.fileUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrInvalidID:
			response.BadRequest(c, "Invalid ID format", err.Error())
		case domain.ErrNotFound:
			response.NotFound(c, "File not found")
		default:
			response.InternalServerError(c, "Failed to get file", err.Error())
		}
		return
	}

	// 2. Build path
	filePath := filepath.Join(h.uploadConfig.Path, "files", file.FileName)

	// 3. Set header download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", `attachment; filename="`+file.OriginalName+`"`)
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Transfer-Encoding", "binary")

	// 4. Stream file
	c.File(filePath)
}
