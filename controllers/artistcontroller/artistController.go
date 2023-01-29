package artistcontroller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"music-final-project/common"
	"music-final-project/configuration"
	"music-final-project/model"
	"net/http"
	"strconv"
)

type Artist struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func Index(c *gin.Context) {
	var artists []Artist

	// Get query parameters for page number and page size
	pageNumber, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if pageNumber == 0 {
		pageNumber = 1
	}

	if limit == 0 {
		limit = 10
	}

	// Get the total number of artist
	var count int64
	configuration.DB.Model(&model.Artist{}).Count(&count)

	countInt := int(count)

	// Calculate the offset for the query
	offset := (pageNumber - 1) * limit

	// Execute the paginated query
	configuration.DB.Select("id", "name").Limit(limit).Offset(offset).Find(&artists)

	// Build the URLs for the previous and next pages
	var prevURL, nextURL string
	if pageNumber > 1 {
		prevURL = common.BuildUrl(c, pageNumber-1, limit)
	}
	if (pageNumber * limit) < countInt {
		nextURL = common.BuildUrl(c, pageNumber+1, limit)
	}

	common.SendResponse(c, http.StatusOK, true, "Success get artists data", nil, gin.H{
		"artists": artists,
		"pagination": gin.H{
			"page":          pageNumber,
			"limit":         limit,
			"total":         count,
			"prev_page_url": prevURL,
			"next_page_url": nextURL,
		},
	})
}

func Show(c *gin.Context) {
	var artist model.Artist

	id := c.Param("id")

	if err := configuration.DB.Preload("Albums.Songs").First(&artist, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			common.SendResponse(c, http.StatusNotFound, false, "Failed to get artist data", []string{"Data not found"}, nil)
			return
		default:
			common.SendResponse(c, http.StatusInternalServerError, false, "Failed to get artist data", []string{err.Error()}, nil)
			return
		}
	}

	common.SendResponse(c, http.StatusOK, true, "Success get artist data", nil, gin.H{
		"artist": artist,
	})
}

func Create(c *gin.Context) {
	var artist model.Artist

	if err := c.ShouldBindJSON(&artist); err != nil {
		common.SendResponse(c, http.StatusBadRequest, false, "Failed to create artist data", []string{err.Error()}, nil)
		return
	}

	createArtist := configuration.DB.Create(&artist)

	if createArtist.Error != nil {
		common.SendResponse(c, http.StatusBadRequest, false, "Failed to create artist data", []string{createArtist.Error.Error()}, nil)
		return
	}

	common.SendResponse(c, http.StatusOK, true, "Success created artist data", nil, gin.H{
		"artist": artist,
	})
}

func Update(c *gin.Context) {
	var artist model.Artist

	id := c.Param("id")

	if err := c.ShouldBindJSON(&artist); err != nil {
		common.SendResponse(c, http.StatusBadRequest, false, "Failed to update artist data", []string{err.Error()}, nil)
		return
	}

	result := configuration.DB.Model(&artist).Where("id = ?", id).Updates(&artist)

	if result.Error != nil {
		common.SendResponse(c, http.StatusBadRequest, false, "Failed to update artist data", []string{result.Error.Error()}, nil)
		return
	}

	common.SendResponse(c, http.StatusOK, true, "Success updated artist data", nil, nil)
}

func Delete(c *gin.Context) {
	var artist model.Artist
	id := c.Param("id")

	var albums []model.Album
	if err := configuration.DB.Where("artist_id = ?", id).Find(&albums).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if configuration.DB.Where("id = ?", id).Delete(&artist).RowsAffected == 0 {
				common.SendResponse(c, http.StatusNotFound, false, "Artist not found", []string{"Data not found"}, nil)
				return
			}
			common.SendResponse(c, http.StatusOK, true, "Success deleted artist data", nil, nil)
			return
		} else {
			common.SendResponse(c, http.StatusBadRequest, false, "Failed to find albums", []string{err.Error()}, nil)
			return
		}
	}

	for _, album := range albums {
		if err := configuration.DB.Where("album_id = ?", album.Id).Delete(&model.Song{}).Error; err != nil {
			common.SendResponse(c, http.StatusBadRequest, false, "Failed to delete songs", []string{err.Error()}, nil)
			return
		}
	}

	if err := configuration.DB.Where("artist_id = ?", id).Delete(&model.Album{}).Error; err != nil {
		common.SendResponse(c, http.StatusBadRequest, false, "Failed to delete albums", []string{err.Error()}, nil)
		return
	}

	if configuration.DB.Where("id = ?", id).Delete(&artist).RowsAffected == 0 {
		common.SendResponse(c, http.StatusNotFound, false, "Artist not found", []string{"Data not found"}, nil)
		return
	}

	common.SendResponse(c, http.StatusOK, true, "Success deleted artist data", nil, nil)
}
