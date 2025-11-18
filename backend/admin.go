package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type updateCollectionRequest struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	ItemLabel     string `json:"itemLabel"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	Filter        string `json:"filter"`
	FeatureIDs    []int  `json:"features"`
	ImageFileName string `json:"imageFile"`
	ImageTitle    string `json:"imageTitle"`
	ImageAlt      string `json:"imageAlt"`
	ImageURL      string `json:"imageURL"`
	ImageStatus   string `json:"imageStatus"`
	Active        bool   `json:"active"`
}

type exifData struct {
	Width  int `json:"ImageWidth"`
	Height int `json:"ImageHeight"`
}

func (svc *ServiceContext) getFeatures(c *gin.Context) {
	user := c.GetString("user")
	log.Printf("INFO: %s is requesting a list of features", user)
	var features []feature
	dbResp := svc.GDB.Order("name asc").Find(&features)
	if dbResp.Error != nil {
		log.Printf("ERROR: unable to get features: %s", dbResp.Error.Error())
		c.String(http.StatusInternalServerError, dbResp.Error.Error())
		return
	}

	c.JSON(http.StatusOK, features)
}

func (svc *ServiceContext) getCollectionDetails(c *gin.Context) {
	user := c.GetString("user")
	id, _ := strconv.Atoi(c.Param("id"))
	log.Printf("INFO: %s is is requesting collection %d details", user, id)

	var rec collection
	dbResp := svc.GDB.Preload("Image").Preload("Features").First(&rec, id)
	if dbResp.Error != nil {
		if errors.Is(dbResp.Error, gorm.ErrRecordNotFound) {
			log.Printf("INFO: no collection context found for %d", id)
			c.String(http.StatusNotFound, "not found")
		} else {
			log.Printf("ERROR: contexed lookup for %d failed: %s", id, dbResp.Error.Error())
			c.String(http.StatusInternalServerError, dbResp.Error.Error())
		}
		return
	}

	if rec.Image != nil {
		rec.Image.URL = fmt.Sprintf("%s/%s", svc.BaseImageURL, rec.Image.Filename)
	}

	c.JSON(http.StatusOK, rec)
}

func (svc *ServiceContext) deleteCollection(c *gin.Context) {
	user := c.GetString("user")
	id, _ := strconv.Atoi(c.Param("id"))
	log.Printf("INFO: %s requests delete of collection %d", user, id)
	dbResp := svc.GDB.Delete(&collection{}, id)
	if dbResp.Error != nil {
		if errors.Is(dbResp.Error, gorm.ErrRecordNotFound) {
			log.Printf("WARNING: collection %d not found", id)
			c.String(http.StatusNotFound, fmt.Sprintf("%d not found", id))
		} else {
			log.Printf("ERROR: unable to find collection %d: %s", id, dbResp.Error.Error())
			c.String(http.StatusInternalServerError, dbResp.Error.Error())
		}
		return
	}

	c.String(http.StatusOK, "deleted")
}

func (svc *ServiceContext) addOrUpdateCollection(c *gin.Context) {
	user := c.GetString("user")
	var req updateCollectionRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("ERROR: %s update request with invalid payload: %v", user, err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("INFO: User %s add/update collection %+v", user, req)

	updateRec := collection{ID: int64(req.ID), Active: req.Active, Title: req.Title, ItemLabel: req.ItemLabel,
		FilterName: req.Filter, Description: req.Description, StartDate: req.StartDate, EndDate: req.EndDate}
	if req.ID == 0 {
		log.Printf("INFO: %s add collection %+v", user, req)
		addResp := svc.GDB.Create(&updateRec)
		if addResp.Error != nil {
			log.Printf("ERROR: %s add %v failed: %s", user, req, addResp.Error.Error())
			c.String(http.StatusInternalServerError, addResp.Error.Error())
			return
		}
	} else {
		log.Printf("INFO: %s update collection %+v", user, req)
		upResp := svc.GDB.Omit("ID").Updates(&updateRec)
		if upResp.Error != nil {
			log.Printf("ERROR: %s update %v failed: %s", user, req, upResp.Error.Error())
			c.String(http.StatusInternalServerError, upResp.Error.Error())
			return
		}
		svc.GDB.Model(&updateRec).Association("Features").Clear()
	}

	log.Printf("INFO: adding features to collection %d", updateRec.ID)
	qStr := "insert into collection_features (collection_id, feature_id) values "
	vals := make([]string, 0)
	for _, featureID := range req.FeatureIDs {
		vals = append(vals, fmt.Sprintf("(%d,%d)", updateRec.ID, featureID))
	}
	qStr += strings.Join(vals, ",")
	dbResp := svc.GDB.Exec(qStr)
	if dbResp.Error != nil {
		log.Printf("ERROR: %s add features failed: %s", user, dbResp.Error.Error())
		c.String(http.StatusInternalServerError, dbResp.Error.Error())
		return
	}

	if req.ImageStatus != "no_change" {
		imgName := req.ImageFileName
		log.Printf("INFO: update image for collection %d to '%s'", updateRec.ID, imgName)
		var exitImg image
		dbResp := svc.GDB.Where("collection_id=?", updateRec.ID).First(&exitImg)
		if dbResp.Error == nil {
			// only one image per collection
			delResp := svc.GDB.Delete(&exitImg)
			if delResp.Error != nil {
				log.Printf("ERROR: unable to delete image rec for %s: %s", exitImg.Filename, delResp.Error.Error())
			}
		}

		newImage := image{CollectionID: updateRec.ID, Filename: req.ImageFileName, Title: req.ImageTitle, AltText: req.ImageAlt}
		if req.ImageStatus == "new" {
			log.Printf("INFO: add logo %s to collection %d; lookup sizing", imgName, updateRec.ID)
			tmpImagePath := fmt.Sprintf("/tmp/%s", req.ImageFileName)
			cmdArray := []string{"-json", "-ImageWidth", "-ImageHeight", tmpImagePath}
			stdout, err := exec.Command("exiftool", cmdArray...).Output()
			if err != nil {
				log.Printf("WARNINIG: unable to get %s metadata: %s", req.ImageFileName, err.Error())
			} else {
				var parsed []exifData
				json.Unmarshal(stdout, &parsed)
				newImage.Width = parsed[0].Width
				newImage.Height = parsed[0].Height
			}
			s3Err := svc.uploadImageToS3(req.ImageFileName)
			if s3Err != nil {
				log.Printf("ERROR: unable to upload %s to S3: %s", req.ImageFileName, s3Err.Error())
			}

			log.Printf("INFO: cleaning up temp image file %s", tmpImagePath)
			os.Remove(tmpImagePath)
		} else {
			log.Printf("INFO: reuse existing logo %s in collection %d", req.ImageFileName, updateRec.ID)
			var origImg image
			imgResp := svc.GDB.Where("filename=?", req.ImageFileName).First(&origImg)
			if imgResp.Error != nil {
				log.Printf("ERROR: unable to get existing image %s: %s", req.ImageFileName, imgResp.Error.Error())
			} else {
				newImage.Width = origImg.Width
				newImage.Height = origImg.Height
			}
		}

		imgResp := svc.GDB.Omit("URL").Create(&newImage)
		if imgResp.Error != nil {
			log.Printf("ERROR: unable to add image record for %s: %s", req.ImageFileName, imgResp.Error.Error())
		}
	}

	var out collection
	svc.GDB.Preload("Image").Preload("Features").First(&out, updateRec.ID)
	out.Image.URL = fmt.Sprintf("%s/%s", svc.BaseImageURL, out.Image.Filename)
	c.JSON(http.StatusOK, out)
}

func (svc *ServiceContext) uploadImageToS3(filename string) error {
	log.Printf("INFO: upload logo %s to S3", filename)
	logoPath := fmt.Sprintf("/tmp/%s", filename)
	logoFile, err := os.Open(logoPath)
	if err != nil {
		return err
	}
	start := time.Now()
	_, err = svc.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(svc.S3ImageBucket),
		Key:    aws.String(filename),
		Body:   logoFile,
	})

	if err != nil {
		return err
	}

	duration := time.Since(start)
	log.Printf("INFO: upload of %s complete in %0.2f seconds", filename, duration.Seconds())
	return nil
}

func (svc *ServiceContext) getLogos(c *gin.Context) {
	user := c.GetString("user")
	log.Printf("INFO: %s is requesting a list of logos", user)
	resp, err := svc.S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(svc.S3ImageBucket),
	})
	if err != nil {
		log.Printf("ERROR: list %s object failed: %s", svc.S3ImageBucket, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]string, 0)
	for _, item := range resp.Contents {
		out = append(out, fmt.Sprintf("%s/%s", svc.BaseImageURL, *item.Key))
	}
	c.JSON(http.StatusOK, out)
}

func (svc *ServiceContext) uploadLogo(c *gin.Context) {
	user := c.GetString("user")
	id := c.Param("id")
	log.Printf("INFO: %s is uploading a new logo for collection %s", user, id)

	formFile, err := c.FormFile("file")
	if err != nil {
		log.Printf("ERROR: unable to get upload file: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	filename := filepath.Base(formFile.Filename)
	dest := fmt.Sprintf("/tmp/%s", filename)
	if _, err := os.Stat(dest); err == nil {
		log.Printf("ERROR: File %s already exists", filename)
		c.String(http.StatusConflict, fmt.Sprintf("%s already exists", filename))
		return
	}
	log.Printf("INFO: receiving log file %s for collection %s", filename, id)
	frmFile, err := formFile.Open()
	if err != nil {
		log.Printf("ERROR: unable to open uploaded file %s: %s", formFile.Filename, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer frmFile.Close()
	out, err := os.Create(dest)
	if err != nil {
		log.Printf("ERROR: unable to create temp file %s: %s", dest, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer out.Close()
	_, err = io.Copy(out, frmFile)
	if err != nil {
		log.Printf("ERROR: unable to save %s: %s", formFile.Filename, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("INFO: done receiving %s", filename)
	c.String(http.StatusOK, "Submitted")
}

func (svc *ServiceContext) deletePendingLogo(c *gin.Context) {
	user := c.GetString("user")
	id := c.Param("id")
	filename := c.Param("fn")
	log.Printf("INFO: %s is deleting pending logo %s from collection %s", user, filename, id)
	dest := fmt.Sprintf("/tmp/%s", filename)
	if _, err := os.Stat(dest); err == nil {
		os.Remove(dest)
		return
	}
	c.String(http.StatusOK, "ok")
}
