package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/studio-senkou/lentera-cendekia-be/app/models"
	"github.com/studio-senkou/lentera-cendekia-be/database"
	"github.com/studio-senkou/lentera-cendekia-be/utils/storage"
)

type StaticAssetController struct {
	staticAssetRepository *models.StaticAssetRepository
}

func NewStaticAssetController() *StaticAssetController {
	db := database.GetDB()
	staticAssetRepository := models.NewStaticAssetRepository(db)
	
	return &StaticAssetController{
		staticAssetRepository: staticAssetRepository,
	}
}

func (sac *StaticAssetController) CreateStaticAsset(c *fiber.Ctx) error {

	asset, err := c.FormFile("asset")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to upload asset",
			"error":   err.Error(),
		})
	} else if asset == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "No asset provided",
		})
	}

	if !storage.IsValidImageExtension(asset.Filename) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid asset type",
		})
	}

	maxSize := int64(0.5 * 1024 * 1024)
	if asset.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Asset size exceeds the limit of 500KB",
		})
	}

	uploadedStoragePath, err := storage.UploadFileToStorage(asset, "static_assets", "STATIC", nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to upload asset",
			"error":   err.Error(),
		})
	}

	staticAsset := &models.StaticAsset{
		AssetName:      asset.Filename,
		AssetType:      "image", // Assuming all assets are images for now
		AssetURL:       uploadedStoragePath,
	}

	if err := sac.staticAssetRepository.Create(staticAsset); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to create static asset",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Asset uploaded successfully",
	})
}

func (sac *StaticAssetController) GetAllStaticAssets(c *fiber.Ctx) error {
	assets, err := sac.staticAssetRepository.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to retrieve static assets",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Static assets retrieved successfully",
		"data":    assets,
	})
}

func (sac *StaticAssetController) GetStaticAssetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid static asset ID",
			"error":   err.Error(),
		})
	}

	asset, err := sac.staticAssetRepository.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to retrieve static asset",
			"error":   err.Error(),
		})
	}

	if asset == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Static asset not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Static asset retrieved successfully",
		"data":    asset,
	})
}

func (sac *StaticAssetController) DeleteStaticAsset(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid static asset ID",
			"error":   err.Error(),
		})
	}

	if err := sac.staticAssetRepository.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to delete static asset",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Static asset deleted successfully",
	})
}