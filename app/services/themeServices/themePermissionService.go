package themeServices

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"wejh-go/app/models"
	"wejh-go/config/database"
)

func AddThemePermission(themeID int, reqStudentIDs []string) error {
	var studentIDs []string
	if len(reqStudentIDs) == 0 {
		var users []models.User
		err := database.DB.Select("student_id").Find(&users).Error
		if err != nil {
			return err
		}
		for _, user := range users {
			studentIDs = append(studentIDs, user.StudentID)
		}
	} else {
		studentIDs = reqStudentIDs
	}

	var permissions []models.ThemePermission
	err := database.DB.Where("student_id IN ?", studentIDs).Find(&permissions).Error
	if err != nil {
		return err
	}
	permissionMap := make(map[string]*models.ThemePermission)
	for i, permission := range permissions {
		permissionMap[permission.StudentID] = &permissions[i]
	}

	for _, studentID := range studentIDs {
		permission, exist := permissionMap[studentID]
		if !exist {
			themePermissionData := models.ThemePermissionData{
				ThemeIDs: []int{themeID},
			}
			newPermissionData, err := json.Marshal(themePermissionData)
			if err != nil {
				return err
			}
			newPermission := models.ThemePermission{
				StudentID:       studentID,
				CurrentThemeID:  themeID,
				ThemePermission: string(newPermissionData),
			}
			permission = &newPermission
			permissions = append(permissions, newPermission)
		}

		var themePermissionData models.ThemePermissionData
		err = json.Unmarshal([]byte(permission.ThemePermission), &themePermissionData)
		if err != nil {
			return err
		}
		if !containThemeID(themePermissionData.ThemeIDs, themeID) {
			themePermissionData.ThemeIDs = append(themePermissionData.ThemeIDs, themeID)
			newPermission, err := json.Marshal(themePermissionData)
			if err != nil {
				return err
			}
			permission.ThemePermission = string(newPermission)
		}
	}

	err = database.DB.Save(&permissions).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateCurrentTheme(id int, studentID string) error {
	result := database.DB.Model(models.ThemePermission{}).
		Where("student_id = ?", studentID).
		Update("current_theme_id", id)
	return result.Error
}

func DeleteThemePermission(studentID string) error {
	result := database.DB.Where("student_id = ?", studentID).Delete(&models.ThemePermission{})
	return result.Error
}

func GetThemePermissionByStudentID(studentID string) (models.ThemePermission, error) {
	var record models.ThemePermission
	result := database.DB.Model(models.ThemePermission{}).Where("student_id = ?", studentID).First(&record)
	return record, result.Error
}

func GetThemeNameByID(themePermission models.ThemePermission) ([]string, error) {
	var themePermissionData models.ThemePermissionData
	err := json.Unmarshal([]byte(themePermission.ThemePermission), &themePermissionData)
	if err != nil {
		return nil, err
	}
	var themes []models.Theme
	err = database.DB.Where("id IN ?", themePermissionData.ThemeIDs).Find(&themes).Error
	if err != nil {
		return nil, err
	}
	var themeNames []string
	for _, theme := range themes {
		themeNames = append(themeNames, theme.Name)
	}
	return themeNames, nil
}

func GetThemesByID(themePermission models.ThemePermission) ([]models.Theme, error) {
	var themePermissionData models.ThemePermissionData
	err := json.Unmarshal([]byte(themePermission.ThemePermission), &themePermissionData)
	if err != nil {
		return nil, err
	}
	var themes []models.Theme
	err = database.DB.Where("id IN ?", themePermissionData.ThemeIDs).Find(&themes).Error
	if err != nil {
		return nil, err
	}
	return themes, nil
}

func AddDefaultThemePermission(studentID string) error {
	var existingPermission models.ThemePermission
	err := database.DB.Where("student_id = ?", studentID).First(&existingPermission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			var themes []models.Theme
			err := database.DB.Where("type = ?", "all").Find(&themes).Error
			if err != nil {
				return err
			}
			if len(themes) == 0 {
				return errors.New("no themes found with type 'all'")
			}

			var themeIDs []int
			for _, theme := range themes {
				themeIDs = append(themeIDs, theme.ID)
			}
			var themePermissionData models.ThemePermissionData
			themePermissionData.ThemeIDs = themeIDs
			permission, err := json.Marshal(themePermissionData)
			if err != nil {
				return err
			}

			newPermission := models.ThemePermission{
				StudentID:       studentID,
				CurrentThemeID:  themeIDs[0],
				ThemePermission: string(permission),
			}

			result := database.DB.Create(&newPermission)
			return result.Error
		} else {
			return err
		}
	} else {
		return nil
	}
}

func containThemeID(themeIDs []int, id int) bool {
	for _, themeID := range themeIDs {
		if themeID == id {
			return true
		}
	}
	return false
}
