package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
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
}

type exifData struct {
	Width  int `json:"ImageWidth"`
	Height int `json:"ImageHeight"`
}

func (svc *ServiceContext) getFeatures(c *gin.Context) {
	user := c.GetString("user")
	log.Printf("INFO: %s is requesting a list of features", user)
	type featureResp struct {
		ID   int64  `db:"id" json:"id"`
		Name string `db:"name" json:"name"`
	}
	var features []featureResp
	q := svc.DB.NewQuery("select id,name from features order by name asc")
	err := q.All(&features)
	if err != nil {
		log.Printf("ERROR: unable to get features: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, features)
}

func (svc *ServiceContext) getCollections(c *gin.Context) {
	user := c.GetString("user")
	log.Printf("INFO: %s is requesting a list of collections", user)
	type collResp struct {
		ID    int64  `db:"id" json:"id"`
		Title string `db:"title" json:"title"`
	}
	q := svc.DB.NewQuery("select id,title from collections order by title asc")
	var recs []collResp
	err := q.All(&recs)
	if err != nil {
		log.Printf("ERROR: unable to get collections: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, recs)
}

func (svc *ServiceContext) getCollectionDetails(c *gin.Context) {
	user := c.GetString("user")
	id := c.Param("id")
	log.Printf("INFO: %s is is requesting collection %s details", user, id)

	var rec collectionRec
	q := svc.DB.NewQuery("select * from collections where id={:cid}")
	q.Bind(dbx.Params{"cid": id})
	err := q.One(&rec)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("INFO: no collection context found for %s", id)
			c.String(http.StatusNotFound, "not found")
		} else {
			log.Printf("ERROR: contexed lookup for %s failed: %s", id, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	out := collectionfromDB(rec)
	out.getFeatures(svc.DB)
	out.getImages(svc.DB, svc.BaseImageURL)

	c.JSON(http.StatusOK, out)
}

func (svc *ServiceContext) deleteCollection(c *gin.Context) {
	user := c.GetString("user")
	id := c.Param("id")
	log.Printf("INFO: %s requests delete of collection %s", user, id)
	q := svc.DB.NewQuery("select * from collections where id={:cid}")
	q.Bind(dbx.Params{"cid": id})
	var rec collectionRec
	err := q.One(&rec)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("WARNING: collection %s not found", id)
			c.String(http.StatusNotFound, fmt.Sprintf("%s not found", id))
		} else {
			log.Printf("ERROR: unable to find collection %s: %s", id, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	err = svc.DB.Model(&rec).Delete()
	if err != nil {
		log.Printf("ERROR: unable to delete ollection %s: %s", id, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
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

	updateRec := collectionRec{ID: int64(req.ID), Title: req.Title, ItemLabel: req.ItemLabel, FilterName: req.Filter}
	if req.Description != "" {
		updateRec.Description.String = req.Description
		updateRec.Description.Valid = true
	}
	if req.StartDate != "" {
		updateRec.StartDate.String = req.StartDate
		updateRec.StartDate.Valid = true
	}
	if req.EndDate != "" {
		updateRec.EndDate.String = req.EndDate
		updateRec.EndDate.Valid = true
	}

	if req.ID == 0 {
		log.Printf("INFO: %s add collection %+v", user, req)
		addErr := svc.DB.Model(&updateRec).Insert()
		if addErr != nil {
			log.Printf("ERROR: %s add %v failed: %s", user, req, addErr.Error())
			c.String(http.StatusInternalServerError, addErr.Error())
			return
		}
	} else {
		log.Printf("INFO: %s update collection %+v", user, req)
		upErr := svc.DB.Model(&updateRec).Exclude("ID").Update()
		if upErr != nil {
			log.Printf("ERROR: %s update %v failed: %s", user, req, upErr.Error())
			c.String(http.StatusInternalServerError, upErr.Error())
			return
		}
		q := svc.DB.NewQuery("delete from collection_features where collection_id={:cid}")
		q.Bind(dbx.Params{"cid": updateRec.ID})
		_, upErr = q.Execute()
		if upErr != nil {
			log.Printf("ERROR: %s reset features failed: %s", user, upErr.Error())
			c.String(http.StatusInternalServerError, upErr.Error())
			return
		}
	}

	log.Printf("INFO: adding features to collection %d", updateRec.ID)
	qStr := "insert into collection_features (collection_id, feature_id) values "
	vals := make([]string, 0)
	for _, featureID := range req.FeatureIDs {
		vals = append(vals, fmt.Sprintf("(%d,%d)", updateRec.ID, featureID))
	}
	qStr += strings.Join(vals, ",")
	q := svc.DB.NewQuery(qStr)
	_, err = q.Execute()
	if err != nil {
		log.Printf("ERROR: %s add features failed: %s", user, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if req.ImageStatus != "no_change" {
		imgName := req.ImageFileName
		log.Printf("INFO: update image for collection %d to '%s'", updateRec.ID, imgName)
		iq := svc.DB.NewQuery("select * from images where collection_id={:cid} limit 1")
		iq.Bind(dbx.Params{"cid": updateRec.ID})
		var img imageRec
		err = iq.One(&img)
		if err == nil {
			// an image exists for this collection, see if the filename matches -- or the new filename is blank
			// if so, delete the current image
			delErr := svc.DB.Model(&img).Delete()
			if delErr != nil {
				log.Printf("ERROR: unable to delete image rec for %s: %s", img.Filename, delErr.Error())
			}
		}

		log.Printf("INFO: add logo %s to collection %d", imgName, updateRec.ID)
		newImage := imageRec{CollectionID: updateRec.ID, Filename: req.ImageFileName}
		if req.ImageTitle != "" {
			newImage.Title.String = req.ImageTitle
			newImage.Title.Valid = true
		}
		if req.ImageAlt != "" {
			newImage.AltText.String = req.ImageAlt
			newImage.AltText.Valid = true
		}

		if req.ImageStatus == "new" {
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
			log.Printf("INFO: lookup w/h of existing image %s", imgName)
			iq := svc.DB.NewQuery("select * from images where filename={:fn} order by collection_id asc limit 1")
			iq.Bind(dbx.Params{"fn": imgName})
			var img imageRec
			err := iq.One(&img)
			if err == nil {
				newImage.Width = img.Width
				newImage.Height = img.Height
			} else {
				log.Printf("ERROR: unable to get existing image %s: %s", imgName, err.Error())
			}
		}

		imgErr := svc.DB.Model(&newImage).Insert()
		if imgErr != nil {
			log.Printf("ERROR: unable to add image record for %s: %s", req.ImageFileName, imgErr.Error())
		}
	}

	// convert DB structs into JSON response
	out := collectionfromDB(updateRec)
	out.getFeatures(svc.DB)
	out.getImages(svc.DB, svc.BaseImageURL)
	c.JSON(http.StatusOK, out)
}

func (svc *ServiceContext) uploadImageToS3(filename string) error {
	log.Printf("INFO: upload logo %s to S3", filename)
	imgBytes, err := ioutil.ReadFile(fmt.Sprintf("/tmp/%s", filename))
	if err != nil {
		return err
	}
	upParams := s3manager.UploadInput{
		Bucket: aws.String(svc.S3ImageBucket),
		Key:    aws.String(filename),
		ACL:    aws.String("public-read"),
		Body:   bytes.NewReader(imgBytes),
	}

	start := time.Now()
	_, err = svc.S3Uploader.Upload(&upParams)
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
	b := aws.String(svc.S3ImageBucket)
	loInput := &s3.ListObjectsV2Input{Bucket: b}
	resp, err := svc.S3Service.ListObjectsV2(loInput)
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

	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("ERROR: unable to get upload file: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	filename := filepath.Base(file.Filename)
	dest := fmt.Sprintf("/tmp/%s", filename)
	if _, err := os.Stat(dest); err == nil {
		log.Printf("ERROR: File %s already exists", filename)
		c.String(http.StatusConflict, fmt.Sprintf("%s already exists", filename))
		return
	}
	log.Printf("INFO: receiving log file %s for collection %s", filename, id)
	if err := c.SaveUploadedFile(file, dest); err != nil {
		log.Printf("ERROR: unable to receive logo %s: %s", filename, err.Error())
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
