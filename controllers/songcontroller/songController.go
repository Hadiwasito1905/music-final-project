package songcontroller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"music-final-project/common"
	"music-final-project/configuration"
	"music-final-project/model"
	"net/http"
	"strconv"
)

func Index(c *gin.Context) {
	var songs []model.Song
	var albums []model.Album

	db := configuration.DB

	// Get query parameters for page number, page size, and artist ID
	pageNumber, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	albumID := c.Query("album_id")
	artistID := c.Query("artist_id")

	if artistID != "" {
		db.Where("artist_id = ?", artistID).Find(&albums)
	}
	if albumID != "" {
		db = db.Where("album_id = ?", albumID)
	}

	if pageNumber == 0 {
		pageNumber = 1
	}

	if limit == 0 {
		limit = 10
	}

	// Get the total number of songs
	var count int64
	db.Model(&model.Song{}).Count(&count)

	countInt := int(count)

	// Calculate the offset for the query
	offset := (pageNumber - 1) * limit

	if len(albums) > 0 {
		albumIDs := make([]uint, len(albums))
		for i, a := range albums {
			albumIDs[i] = uint(a.Id)
		}

		// Execute the paginated query with the album IDs as a filter
		db.Where("album_id IN (?)", albumIDs).Limit(limit).Offset(offset).Find(&songs)
	} else {
		// Execute the paginated query without any album ID filter
		db.Limit(limit).Offset(offset).Find(&songs)
	}

	// Build the URLs for the previous and next pages
	var prevURL, nextURL string
	if pageNumber > 1 {
		prevURL = common.BuildUrl(c, pageNumber-1, limit)
	}
	if (pageNumber * limit) < countInt {
		nextURL = common.BuildUrl(c, pageNumber+1, limit)
	}

	common.SendResponse(c, http.StatusOK, true, "Success get songs data", nil, gin.H{
		"songs": songs,
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
	var song model.Song

	id := c.Param("id")

	if err := configuration.DB.First(&song, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			common.SendResponse(c, http.StatusNotFound, false, "Failed to get song data", []string{"Data not found"}, nil)
			return
		default:
			common.SendResponse(c, http.StatusInternalServerError, false, "Failed to get song data", []string{err.Error()}, nil)
			return
		}
	}

	common.SendResponse(c, http.StatusOK, true, "Success get song data", nil, gin.H{
		"song": song,
	})
}

func Create(c *gin.Context) {
	var song model.Song

	if err := c.ShouldBindJSON(&song); err != nil {
		common.SendResponse(c, http.StatusBadRequest, false, "Failed to create song data", []string{err.Error()}, nil)
		return
	}

	createSong := configuration.DB.Create(&song)

	if createSong.Error != nil {
		common.SendResponse(c, http.StatusBadRequest, false, "Failed to create song data", []string{createSong.Error.Error()}, nil)
		return
	}

	common.SendResponse(c, http.StatusOK, true, "Success created song data", nil, gin.H{
		"song": song,
	})
}

func Update(c *gin.Context) {
	var song model.Song

	id := c.Param("id")

	if err := c.ShouldBindJSON(&song); err != nil {
		common.SendResponse(c, http.StatusBadRequest, false, "Failed to update song data", []string{err.Error()}, nil)
		return
	}

	result := configuration.DB.Model(&song).Where("id = ?", id).Updates(&song)

	if result.Error != nil {
		common.SendResponse(c, http.StatusBadRequest, false, "Failed to update song data", []string{result.Error.Error()}, nil)
		return
	}

	common.SendResponse(c, http.StatusOK, true, "Success updated song data", nil, nil)
}

func Delete(c *gin.Context) {
	var song model.Song

	id := c.Param("id")

	if configuration.DB.Delete(&song, id).RowsAffected == 0 {
		common.SendResponse(c, http.StatusNotFound, false, "Song not found", []string{"Data not found"}, nil)
		return
	}

	common.SendResponse(c, http.StatusOK, true, "Success deleted song data", nil, nil)
}
