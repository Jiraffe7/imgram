package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/Jiraffe7/imgram/user"
	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
)

const (
	CaptionLimitBytes = 1000
	FileLimitBytes    = 100 << 20 // 100 MB
)

type Response struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
}

func respondError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.Encode(Response{Error: err.Error()})
}

type Post struct {
	UserID   uint64 `db:"user_id"`
	Caption  string
	Filepath string
}

// Create a post with at most one image and a caption
func PostImage(w http.ResponseWriter, r *http.Request) {
	user, ok := user.FromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, errors.New("PostImage: not authenticated"))
		return
	}

	reader, err := r.MultipartReader()
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	const (
		FormNameCaption = "caption"
		FormNameFile    = "file"
	)

	var post Post
	post.UserID = user.ID

	for part, err := reader.NextPart(); err == nil; part, err = reader.NextPart() {
		log.Printf("PostImage: part: formname=%s, filename=%s\n", part.FormName(), part.FileName())

		switch part.FormName() {
		case FormNameCaption:
			post.Caption, err = readCaption(part)
			if err != nil {
				log.Printf("PostImage: error reading caption: %v\n", err)
				respondError(w, http.StatusBadRequest, err)
				return
			}
		case FormNameFile:
			post.Filepath, err = readFile(part, user.ID)
			if err != nil {
				log.Printf("PostImage: error reading file: %v\n", err)
				respondError(w, http.StatusBadRequest, err)
				return
			}
		}
	}

	_, err = app.db.NamedExec("insert into posts (user_id, caption, filepath) values (:user_id, :caption, :filepath)", post)
	if err != nil {
		log.Printf("PostImage: error persisting to DB: %v\n", err)
	}

	w.WriteHeader(http.StatusOK)
}

// readCaption reads the caption from the part up to a limit of CaptionLimitBytes.
// Truncates the caption if it exceeds the limit.
func readCaption(part *multipart.Part) (string, error) {
	r := io.LimitReader(part, CaptionLimitBytes)

	var caption bytes.Buffer
	_, err := io.Copy(&caption, r)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return caption.String(), nil
}

// readFile reads the file from the part up to a limit of FileLimitBytes.
func readFile(part *multipart.Part, userID uint64) (string, error) {
	var (
		dataDir  = app.dataDir
		userDir  = strconv.FormatUint(userID, 10)
		filename = part.FileName()
		ext      = path.Ext(filename)
		filepath = path.Join(dataDir, userDir, filename)
		dir      = path.Dir(filepath)
	)

	var imgDecoder func(io.Reader) (image.Image, error)
	switch ext {
	case ".jpg":
		imgDecoder = jpeg.Decode
	case ".png":
		imgDecoder = png.Decode
	case ".bmp":
		imgDecoder = bmp.Decode
	default:
		return "", errors.New("invalid file format: " + ext)
	}

	err := os.MkdirAll(dir, 0766)
	if err != nil {
		return "", err
	}

	f, err := os.Create(filepath)
	if err != nil {
		return "", err
	}

	r := io.LimitReader(part, FileLimitBytes)
	img, err := imgDecoder(r)
	if err != nil {
		return "", err
	}

	// resize image
	scaled := image.NewRGBA(image.Rect(0, 0, 600, 600))
	draw.NearestNeighbor.Scale(scaled, scaled.Rect, img, img.Bounds(), draw.Over, nil)

	w := bufio.NewWriter(f)
	defer w.Flush()

	err = jpeg.Encode(w, scaled, nil)
	if err != nil {
		return "", err
	}

	return filepath, nil
}
