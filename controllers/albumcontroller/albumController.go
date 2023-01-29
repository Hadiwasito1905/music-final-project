package albumcontroller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"music-final-project/common"
	"music-final-project/configuration"
	"music-final-project/model"
	"net/http"
	"strconv"
)

type Album struct {
	Id        int64  `json:"id"`
	Artist_Id int64  `json:"artist_id"`
	Title     string `json:"title"`
	Price     int64  `json:"price"`
}

func Index(c *gin.Context) {
	var albums []Album

	db := configuration.DB

	// Get query parameters for page number and page size
	pageNumber, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if artistID := c.Query("artist_id"); artistID != "" {
		db = db.Where("artist_id = ?", artistID)
	}

	if pageNumber == 0 {
		pageNumber = 1
	}

	if limit == 0 {
		limit = 10
	}

	// Get the total number of albums
	var count int64
	db.Model(&model.Album{}).Count(&count)

	countInt := int(count)

	// Calculate the offset for the query
	offset := (pageNumber - 1) * limit

	// Execute the paginated query
	db.Limit(limit).Offset(offset).Find(&albums)

	// Build the URLs for the previous and next pages
	var prevURL, nextURL string
	if pageNumber > 1 {
		prevURL = common.BuildUrl(c, pageNumber-1, limit)
	}
	if (pageNumber * limit) < countInt {
		nextURL = common.BuildUrl(c, pageNumber+1, limit)
	}

	common.SendResponse(c, http.StatusOK, true, "Success get albums data", nil, gin.H{
		"albums": albums,
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
	var album model.Album

	id := c.Param("id")

	if err := configuration.DB.Preload("Songs").First(&album, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			common.SendResponse(c, http.StatusNotFound, false, "Failed to get album data", []string{"Data not found"}, nil)
			return
		default:
			common.SendResponse(c, http.StatusInternalServerError, false, "Failed to get album data", []string{err.Error()}, nil)
			return
		}
	}

	common.SendResponse(c, http.StatusOK, true, "Success get album data", nil, gin.H{
		"album": album,
	})
}

func Create(c *gin.Context) {
	var album model.Album

	if err := c.ShouldBindJSON(&album); err != nil {
		common.SendResponse(c, http.StatusBadRequest, false, "Failed to create album data", []string{err.Error()}, nil)
		return
	}

	createAlbum := configuration.DB.Create(&album)

	if createAlbum.Error != nil {
		common.SendResponse(c, http.StatusBadRequest, false, "Failed to create album data", []string{createAlbum.Error.Error()}, nil)
		return
	}

	common.SendResponse(c, http.StatusOK, true, "Success created album data", nil, gin.H{
		"album": album,
	})
}

func Update(c *gin.Context) {
	var album model.Album

	id := c.Param("id")

	if err := c.ShouldBindJSON(&album); err != nil {
		common.SendResponse(c, http.StatusBadRequest, false, "Failed to update album data", []string{err.Error()}, nil)
		return
	}

	result := configuration.DB.Model(&album).Where("id = ?", id).Updates(&album)
	if result.Error != nil {
		common.SendResponse(c, http.StatusBadRequest, false, "Failed to update album data", []string{result.Error.Error()}, nil)
		return
	}

	common.SendResponse(c, http.StatusOK, true, "Success updated album data", nil, nil)
}

func Delete(c *gin.Context) {
	id := c.Param("id")

	configuration.DB.Where("album_id = ?", id).Delete(&model.Song{})
	result := configuration.DB.Where("id = ?", id).Delete(&model.Album{})

	if result.RowsAffected == 0 {
		common.SendResponse(c, http.StatusNotFound, false, "Album not found", []string{"Data not found"}, nil)
		return
	}

	if result.Error != nil {
		common.SendResponse(c, http.StatusBadRequest, false, "Failed to delete album data", []string{result.Error.Error()}, nil)
		return
	}

	common.SendResponse(c, http.StatusOK, true, "Success deleted album data", nil, nil)
}
