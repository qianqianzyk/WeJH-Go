package themeServices

import (
	"encoding/json"
	"fmt"
	"strconv"
	"wejh-go/app/config"
	"wejh-go/app/models"
	"wejh-go/config/database"
)

func CheckThemeExist(ids ...int) error {
	var themes []models.Theme
	result := database.DB.Model(&models.Theme{}).Where("id IN ?", ids).Find(&themes)
	if result.Error != nil {
		return result.Error
	}
	if len(themes) != len(ids) {
		return fmt.Errorf("some theme IDs do not exist")
	}
	return nil
}

func CreateTheme(themeName, themeType string, hasDarkMode bool, themeConfigData models.ThemeConfigData) error {
	themeConfig, err := json.Marshal(themeConfigData)
	if err != nil {
		return err
	}
	record := models.Theme{
		Name:        themeName,
		Type:        themeType,
		IsDarkMode:  hasDarkMode,
		ThemeConfig: string(themeConfig),
	}
	result := database.DB.Create(&record)
	return result.Error
}

func UpdateTheme(themeID int, themeName string, hasDarkMode bool, themeConfigData models.ThemeConfigData) error {
	themeConfig, err := json.Marshal(themeConfigData)
	if err != nil {
		return err
	}
	record := models.Theme{
		Name:        themeName,
		IsDarkMode:  hasDarkMode,
		ThemeConfig: string(themeConfig),
	}
	result := database.DB.Model(models.Theme{}).Where(&models.Theme{ID: themeID}).Updates(&record)
	return result.Error
}

func GetThemeByID(id int) (map[string]interface{}, string, bool, error) {
	var record models.Theme
	result := database.DB.Model(models.Theme{}).Where(&models.Theme{ID: id}).First(&record)

	var themeConfig models.ThemeConfigData
	if err := json.Unmarshal([]byte(record.ThemeConfig), &themeConfig); err != nil {
		return nil, "", false, err
	}

	parsedTheme := map[string]interface{}{
		"name":         record.Name,
		"id":           record.ID,
		"theme_config": themeConfig,
		"type":         record.Type,
		"is_dark_mode": record.IsDarkMode,
	}
	return parsedTheme, record.Type, record.IsDarkMode, result.Error
}

func GetThemes() ([]map[string]interface{}, error) {
	var themes []models.Theme
	result := database.DB.Model(models.Theme{}).Find(&themes)
	if result.Error != nil {
		return nil, result.Error
	}

	var parsedThemes []map[string]interface{}

	for _, theme := range themes {
		var themeConfig models.ThemeConfigData
		if err := json.Unmarshal([]byte(theme.ThemeConfig), &themeConfig); err != nil {
			return nil, err
		}

		parsedTheme := map[string]interface{}{
			"name":         theme.Name,
			"theme_id":     theme.ID,
			"theme_config": themeConfig,
			"is_dark_mode": theme.IsDarkMode,
		}
		parsedThemes = append(parsedThemes, parsedTheme)
	}

	return parsedThemes, nil
}

func DeleteTheme(id int, themeType string, isDarkMode bool) error {
	tx := database.DB.Begin()
	if err := tx.Delete(&models.Theme{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	var theme models.Theme
	if err := tx.Where("type = ?", "all").First(&theme).Error; err != nil {
		tx.Rollback()
		return err
	}

	var defaultThemeID int
	defaultThemeIDStr := config.GetDefaultThemeKey()
	if defaultThemeIDStr != "" {
		defaultThemeID, _ = strconv.Atoi(defaultThemeIDStr)
		if id == defaultThemeID {
			defaultThemeID = theme.ID
			err := config.SetDefaultThemeKey(strconv.Itoa(defaultThemeID))
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	} else {
		defaultThemeID = theme.ID
	}

	updateField := "current_theme_id"
	if isDarkMode {
		updateField = "current_theme_dark_id"
	}
	if err := tx.Model(&models.ThemePermission{}).
		Where(updateField+" = ?", id).
		Update(updateField, defaultThemeID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if themeType == "all" {
		tx.Commit()
		return nil
	}

	var permissions []models.ThemePermission
	result := tx.Model(models.ThemePermission{}).Find(&permissions)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	updatedPermissionMap := make(map[string]models.ThemePermissionData)
	for _, permission := range permissions {
		var themePermissionData models.ThemePermissionData
		err := json.Unmarshal([]byte(permission.ThemePermission), &themePermissionData)
		if err != nil {
			tx.Rollback()
			return err
		}
		updatedThemeIDs := removeThemeID(themePermissionData.ThemeIDs, id)
		if len(updatedThemeIDs) != len(themePermissionData.ThemeIDs) {
			themePermissionData.ThemeIDs = updatedThemeIDs
			if len(updatedThemeIDs) == 0 {
				themePermissionData.ThemeIDs = []int{}
			}
			updatedPermissionMap[permission.StudentID] = themePermissionData
		}
	}
	for studentID, data := range updatedPermissionMap {
		newPermission, err := json.Marshal(data)
		if err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&models.ThemePermission{}).
			Where("student_id = ?", studentID).
			Update("theme_permission", string(newPermission)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func removeThemeID(themeIDs []int, id int) []int {
	var updatedThemeIDs []int
	for _, themeID := range themeIDs {
		if themeID != id {
			updatedThemeIDs = append(updatedThemeIDs, themeID)
		}
	}
	return updatedThemeIDs
}
