package helpers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func PrintRoutes(router chi.Router) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		fmt.Printf("%s -> %s %s\n", handlerName, method, route)
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Fatalf("Logging err: %s\n", err.Error())
	}
}

func GetIDFromParams(r *http.Request) (int, error) {
	idParams := chi.URLParam(r, "id")
	id, error := strconv.Atoi(idParams)

	return id, error
}

func PopulateStudent(r *http.Request, target interface{}) (map[string]*multipart.FileHeader, error) {

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, err
	}

	// Get the value of the target
	value := reflect.ValueOf(target).Elem()
	targetType := value.Type()

	// Populate form values into the target struct
	for i := 0; i < value.NumField(); i++ {
		field := targetType.Field(i)
		formValue := r.FormValue(strings.ToLower(field.Name))
		if formValue != "" {
			value.Field(i).SetString(formValue)
		}
	}

	// Return parsed files for further processing
	files := map[string]*multipart.FileHeader{}
	for key, headers := range r.MultipartForm.File {
		if len(headers) > 0 {
			files[key] = headers[0]
		}
	}

	return files, nil
}

func UploadImageToSupabase(file multipart.File, handler *multipart.FileHeader, bucketName string) (string, error) {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		return "", fmt.Errorf("failed to load environment variables: %w", err)
	}

	// Read the file into a buffer
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	// Generate a unique file name
	formattedFileName := fmt.Sprintf("student_%d%s", time.Now().Unix(), filepath.Ext(handler.Filename))

	// Initialize Supabase storage client
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseAPIKey := os.Getenv("SUPABASE_API_KEY")

	// create upload url for supabase storage
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, bucketName, formattedFileName)

	// create a request
	req, err := http.NewRequest(http.MethodPost, uploadURL, strings.NewReader(string(fileBytes)))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", supabaseAPIKey))
	req.Header.Set("Content-Type", handler.Header.Get("Content-Type"))
	req.Header.Set("Cache-Control", "3600") // 1 hour

	//send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}

	defer resp.Body.Close()

	// check the response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Construct the public URL
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, bucketName, formattedFileName)
	return publicURL, nil
}
