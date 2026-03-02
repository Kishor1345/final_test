// Package  controllerpdf for pdf generation.
//
// --- Creator's Info ---
//
// Creator: Sridharan
//
// Created On:28/11/2025
//
//modify by : kishorekumar
//
//modify on :02/02/2026
package controllerspdf
import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	credentials "Hrmodule/dbconfig"

	// Custom Modules
	"Hrmodule/auth"
	"Hrmodule/utils"

	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	_ "github.com/lib/pq"
	"github.com/pkg/sftp"
	"github.com/skip2/go-qrcode"
	"golang.org/x/crypto/ssh"
)

// Template struct
type Template struct {
	OrderNo       string
	Order_type    string
	OrderDate     string
	ToColumn      string
	Subject       string
	Reference     string
	BodyHTML      string
	SignatureHTML string
	CCTo          string
	FooterHTML    string
	ProcessCode   string
	EmployeeID    string
	TaskID        string
}

// ApproverDetail struct
type ApproverDetail struct {
	UserDisplay string
	UserRole    string
	Remarks     string
	UpdatedOn   string
}

// PDFRequest struct
type PDFRequest struct {
	Data string `json:"Data"`
}

// DecryptedPDFData struct - actual payload after decryption
type DecryptedPDFData struct {
	Token        string `json:"token"`
	OrderNo      string `json:"order_no"`
	ProcessID    int    `json:"process_id"`
	TaskID       string `json:"task_id"`
	Status       string `json:"status"`
	TemplateType string `json:"templatetype"`
}

// APIResponse struct for encrypted response
type APIResponse struct {
	Status  int         `json:"Status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data"`
}

// AssetIndex struct for database insertion
type AssetIndex struct {
	ID             int       `json:"id"`
	TaskID         string    `json:"task_id"`
	ProcessName    string    `json:"process_name"`
	OriginalName   string    `json:"original_name"`
	StoredName     string    `json:"stored_name"`
	FilePath       string    `json:"file_path"`
	FileSize       int64     `json:"file_size"`
	MimeType       string    `json:"mime_type"`
	UploadedBy     string    `json:"uploaded_by"`
	RoleRights     string    `json:"role_rights"`
	Rights         string    `json:"rights"`
	ChecksumSHA256 string    `json:"checksum_sha256"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Status         string    `json:"status"`
}

// ConnectSFTP establishes an SFTP connection using credentials from .env
func ConnectSFTP() (*sftp.Client, *ssh.Client, error) {
	// 1. Try to get from .env
	host := os.Getenv("SFTP_HOST")
	port := os.Getenv("SFTP_PORT")
	user := os.Getenv("SFTP_USERNAME")
	pass := os.Getenv("SFTP_PASSWORD")

	// 2. Fallback if .env failed to load (Fixes "dial tcp :0" error)
	if host == "" {
		host = "wfvault.iitm.ac.in"
	}
	if port == "" {
		port = "22"
	}
	if user == "" {
		user = "NAS_Admin"
	}
	if pass == "" {
		pass = "@dm1n#N@&(102424)"
	}

	address := fmt.Sprintf("%s:%s", host, port)
	log.Printf("🔌 Connecting to SFTP Server: %s (User: %s)", address, user)

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// SSH connection
	sshConn, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, nil, fmt.Errorf("SSH dial failed to %s: %v", address, err)
	}

	// Create SFTP
	client, err := sftp.NewClient(sshConn)
	if err != nil {
		sshConn.Close()
		return nil, nil, fmt.Errorf("SFTP client creation failed: %v", err)
	}

	return client, sshConn, nil
}

// sftpMkdirAll mimics os.MkdirAll for SFTP
func sftpMkdirAll(client *sftp.Client, path string) error {
	// Normalize path separators
	path = strings.ReplaceAll(path, "\\", "/")
	parts := strings.Split(path, "/")

	currentPath := ""
	for _, part := range parts {
		if part == "" {
			currentPath += "/"
			continue
		}
		currentPath = filepath.Join(currentPath, part)
		currentPath = strings.ReplaceAll(currentPath, "\\", "/") // Ensure forward slashes

		// Try to get stat, if fails, create dir
		_, err := client.Stat(currentPath)
		if err != nil {
			if os.IsNotExist(err) {
				if err := client.Mkdir(currentPath); err != nil {
					return fmt.Errorf("failed to create directory %s: %v", currentPath, err)
				}
			} else {
				// Some other error
				return err
			}
		}
	}
	return nil
}

// saveToSFTP uploads bytes to a specific path on SFTP
func saveToSFTP(client *sftp.Client, fullPath string, data []byte) error {
	// Ensure path uses forward slashes
	fullPath = strings.ReplaceAll(fullPath, "\\", "/")

	f, err := client.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file on sftp %s: %v", fullPath, err)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("failed to write data to sftp file: %v", err)
	}
	return nil
}

// Fetch template from DB using stored procedure
func getTemplateFromDB(orderNo string, processID int, taskID string) (*Template, error) {
	//Database Connection
    db := credentials.GetDB()

	var err error
	var query string
	var rows *sql.Rows

	if orderNo != "" && orderNo != "null" && (taskID == "" || taskID == "null") {
		log.Printf("📋 Fetching template using order_no: %s, process_id: %d", orderNo, processID)
		query = `SELECT order_no, order_type, order_date, to_column, subject, reference, 
				body_html, signature_html, cc_to, footer_html, 
				process_code, employee_id, task_id
				FROM meivan.pdftemplaterecordsnew($1, $2, NULL)
				LIMIT 1`
		rows, err = db.Query(query, orderNo, processID)
	} else if (orderNo == "" || orderNo == "null") && taskID != "" && taskID != "null" {
		log.Printf("📋 Fetching template using process_id: %d, task_id: %s", processID, taskID)
		query = `SELECT order_no, order_type, order_date, to_column, subject, reference, 
				body_html, signature_html, cc_to, footer_html, 
				process_code, employee_id, task_id
				FROM meivan.pdftemplaterecordsnew(NULL, $1, $2)
				LIMIT 1`
		rows, err = db.Query(query, processID, taskID)
	} else {
		return nil, fmt.Errorf("invalid parameter combination: provide either (order_no + process_id) or (process_id + task_id)")
	}

	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	t := &Template{}
	if rows.Next() {
		var orderNo, orderType, orderDate, toColumn, subject, reference sql.NullString
		var bodyHTML, signatureHTML, ccTo, footerHTML, processCode, employeeID, taskID sql.NullString

		err = rows.Scan(
			&orderNo, &orderType, &orderDate, &toColumn, &subject, &reference,
			&bodyHTML, &signatureHTML, &ccTo, &footerHTML,
			&processCode, &employeeID, &taskID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		t.OrderNo = nullStringToString(orderNo)
		t.Order_type = nullStringToString(orderType)
		t.OrderDate = nullStringToString(orderDate)
		t.ToColumn = nullStringToString(toColumn)
		t.Subject = nullStringToString(subject)
		t.Reference = nullStringToString(reference)
		t.BodyHTML = nullStringToString(bodyHTML)
		t.SignatureHTML = nullStringToString(signatureHTML)
		t.CCTo = nullStringToString(ccTo)
		t.FooterHTML = nullStringToString(footerHTML)
		t.ProcessCode = nullStringToString(processCode)
		t.EmployeeID = nullStringToString(employeeID)
		t.TaskID = nullStringToString(taskID)
	} else {
		return nil, fmt.Errorf("no template found")
	}

	if t.OrderDate != "" {
		var parsedTime time.Time
		var parseErr error

		if strings.Contains(t.OrderDate, "T") {
			parsedTime, parseErr = time.Parse(time.RFC3339, t.OrderDate)
		}

		if parseErr != nil || parsedTime.IsZero() {
			parsedTime, parseErr = time.Parse("2006-01-02", t.OrderDate)
		}

		if parseErr != nil || parsedTime.IsZero() {
			parsedTime, parseErr = time.Parse("02-01-2006", t.OrderDate)
		}

		if parseErr == nil && !parsedTime.IsZero() {
			t.OrderDate = parsedTime.Format("02-01-2006")
		}
	}

	t.BodyHTML = html.UnescapeString(t.BodyHTML)
	t.BodyHTML = strings.ReplaceAll(t.BodyHTML, `\n`, " ")
	t.SignatureHTML = cleanHTML(t.SignatureHTML)
	t.ToColumn = cleanHTML(t.ToColumn)
	t.CCTo = cleanHTML(t.CCTo)
	t.FooterHTML = cleanHTML(t.FooterHTML)

	log.Printf("✅ Template fetched successfully - Order No: %s, Task ID: %s, Order Date: %s", t.OrderNo, t.TaskID, t.OrderDate)
	return t, nil
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func getApproverDetails(orderNo string, processID int, taskID string) ([]ApproverDetail, error) {
	//Database Connection
    db := credentials.GetDB()

	var err error
	var query string
	var rows *sql.Rows

	if (orderNo == "" || orderNo == "null") && taskID != "" && taskID != "null" {
		log.Printf("📋 Fetching comments using process_id: %d, task_id: %s", processID, taskID)
		query = `SELECT user_display, user_role, remarks, updated_on
				FROM meivan.getcomments($1, $2)
				ORDER BY updated_on DESC`
		rows, err = db.Query(query, processID, taskID)
	} else if orderNo != "" && orderNo != "null" {
		log.Printf("📋 Fetching task_id for order_no: %s", orderNo)
		var fetchedTaskID string
		taskQuery := `SELECT task_id FROM meivan.pdftemplaterecordsnew($1, $2, NULL) LIMIT 1`
		err = db.QueryRow(taskQuery, orderNo, processID).Scan(&fetchedTaskID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch task_id for order_no: %w", err)
		}

		log.Printf("📋 Fetching comments using process_id: %d, task_id: %s", processID, fetchedTaskID)
		query = `SELECT user_display, user_role, remarks, updated_on
				FROM meivan.getcomments($1, $2)
				ORDER BY updated_on DESC`
		rows, err = db.Query(query, processID, fetchedTaskID)
	} else {
		return nil, fmt.Errorf("invalid parameter combination: task_id is required")
	}

	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var approvers []ApproverDetail
	for rows.Next() {
		var approver ApproverDetail
		var userDisplay, userRole, remarks, updatedOn sql.NullString

		err := rows.Scan(&userDisplay, &userRole, &remarks, &updatedOn)
		if err != nil {
			log.Printf("⚠️ Failed to scan approver row: %v", err)
			continue
		}

		if userDisplay.Valid && strings.TrimSpace(userDisplay.String) != "" {
			approver.UserDisplay = userDisplay.String
			approver.UserRole = nullStringToString(userRole)
			approver.Remarks = nullStringToString(remarks)

			updatedOnStr := nullStringToString(updatedOn)
			if updatedOnStr != "" {
				var parsedTime time.Time
				var parseErr error

				if strings.Contains(updatedOnStr, "T") {
					parsedTime, parseErr = time.Parse(time.RFC3339, updatedOnStr)
				}

				if parseErr != nil || parsedTime.IsZero() {
					parsedTime, parseErr = time.Parse("2006-01-02", updatedOnStr)
				}

				if parseErr != nil || parsedTime.IsZero() {
					parsedTime, parseErr = time.Parse("02-01-2006", updatedOnStr)
				}

				if parseErr == nil && !parsedTime.IsZero() {
					approver.UpdatedOn = parsedTime.Format("02-01-2006")
				} else {
					approver.UpdatedOn = updatedOnStr
				}
			}

			approvers = append(approvers, approver)
			log.Printf("✅ Added approver: UserDisplay=%s, Role=%s, Remarks=%s, UpdatedOn=%s",
				approver.UserDisplay, approver.UserRole, approver.Remarks, approver.UpdatedOn)
		}
	}

	log.Printf("✅ Found %d approver(s)", len(approvers))
	return approvers, nil
}

func cleanHTML(s string) string {
	s = strings.ReplaceAll(s, `\n`, " ")
	re := regexp.MustCompile(`>\s+<`)
	s = re.ReplaceAllString(s, "><")
	return strings.TrimSpace(s)
}

func generateQRCode(text string) (string, error) {
	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		return "", err
	}

	qr.DisableBorder = false
	pngBytes, err := qr.PNG(256)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(pngBytes), nil
}

func fetchSignatureFromAPI(role string) ([]byte, error) {
	encodedRole := strings.ReplaceAll(role, " ", "%20")
	url := fmt.Sprintf("https://wftest1.iitm.ac.in:7007/download-signature?role=%s", encodedRole)

	log.Printf("🔄 Fetching signature .dat file from API for role: '%s'", role)
	log.Printf("📍 API URL: %s", url)
	// 🚨 Disable SSL verification (TEST ENV ONLY)
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // ❗ Skips SSL validation
		},
	}

	client := &http.Client{
		Transport: customTransport,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET request failed: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("📥 API Response Status: %d %s", resp.StatusCode, resp.Status)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ API Error Response: %s", string(data))
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(data))
	}

	log.Printf("✅ Retrieved signature .dat file (%d bytes)", len(data))
	return data, nil
}

func decryptSignatureDat(datBytes []byte, outputPath string) error {
	imageFormat := "unknown"
	if len(datBytes) >= 8 {
		pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
		if bytes.Equal(datBytes[:8], pngHeader) {
			imageFormat = "png"
		} else if len(datBytes) >= 12 && string(datBytes[4:12]) == "ftypavif" {
			imageFormat = "avif"
		} else if len(datBytes) >= 2 && datBytes[0] == 0xFF && datBytes[1] == 0xD8 {
			imageFormat = "jpeg"
		}
	}

	log.Printf("✅ Detected image format from .dat: %s", imageFormat)

	if imageFormat == "avif" {
		tempAvifPath := outputPath + ".temp.avif"
		err := ioutil.WriteFile(tempAvifPath, datBytes, 0644)
		if err != nil {
			return fmt.Errorf("failed to save temp AVIF: %w", err)
		}
		defer os.Remove(tempAvifPath)

		err = convertAVIFtoPNG(tempAvifPath, outputPath)
		if err != nil {
			return fmt.Errorf("failed to convert AVIF to PNG: %w", err)
		}
	} else {
		err := ioutil.WriteFile(outputPath, datBytes, 0644)
		if err != nil {
			return fmt.Errorf("failed to save image: %w", err)
		}
	}

	return nil
}

func convertAVIFtoPNG(avifPath, pngPath string) error {
	cmd := exec.Command("ffmpeg", "-i", avifPath, "-y", pngPath)
	err := cmd.Run()
	if err == nil {
		log.Printf("✅ Converted AVIF to PNG using ffmpeg")
		return nil
	}
	log.Printf("⚠️ ffmpeg conversion failed: %v", err)

	cmd = exec.Command("convert", avifPath, pngPath)
	err = cmd.Run()
	if err == nil {
		log.Printf("✅ Converted AVIF to PNG using ImageMagick")
		return nil
	}
	log.Printf("⚠️ ImageMagick conversion failed: %v", err)

	return fmt.Errorf("no suitable converter found (tried ffmpeg and ImageMagick)")
}

func extractRoleName(signatureHTML string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	plainText := re.ReplaceAllString(signatureHTML, "")

	lines := strings.Split(plainText, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			return line
		}
	}

	return strings.TrimSpace(plainText)
}

func sanitizeFilename(filename string) string {
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, ":", "_")
	filename = strings.ReplaceAll(filename, "*", "_")
	filename = strings.ReplaceAll(filename, "?", "_")
	filename = strings.ReplaceAll(filename, "\"", "_")
	filename = strings.ReplaceAll(filename, "<", "_")
	filename = strings.ReplaceAll(filename, ">", "_")
	filename = strings.ReplaceAll(filename, "|", "_")
	filename = strings.TrimSpace(filename)
	return filename
}

func generateApproversTableHTML(approvers []ApproverDetail) string {
	if len(approvers) == 0 {
		return ""
	}

	var tableHTML strings.Builder
	tableHTML.WriteString(`<div style="margin-top: 20px;">
		<table style="border-collapse: collapse; width: 100%; border: 1px solid #000;">
			<thead>
				<tr style="background-color: #f0f0f0;">
					<th colspan="4" style="border: 1px solid #000; padding: 8px; text-align: center; font-weight: bold;">Approvers Details</th>
				</tr>
				<tr style="background-color: #f0f0f0;">
					<th style="border: 1px solid #000; padding: 8px; text-align: center;">User</th>
					<th style="border: 1px solid #000; padding: 8px; text-align: center;">Role</th>
					<th style="border: 1px solid #000; padding: 8px; text-align: center;">Remarks</th>
					<th style="border: 1px solid #000; padding: 8px; text-align: center;">Updated On</th>
				</tr>
			</thead>
			<tbody>`)

	for _, approver := range approvers {
		tableHTML.WriteString(fmt.Sprintf(`
				<tr>
					<td style="border: 1px solid #000; padding: 8px; text-align: left;">%s</td>
					<td style="border: 1px solid #000; padding: 8px; text-align: left;">%s</td>
					<td style="border: 1px solid #000; padding: 8px; text-align: left;">%s</td>
					<td style="border: 1px solid #000; padding: 8px; text-align: center;">%s</td>
				</tr>`,
			approver.UserDisplay, approver.UserRole, approver.Remarks, approver.UpdatedOn))
	}

	tableHTML.WriteString(`
			</tbody>
		</table>
	</div>`)

	return tableHTML.String()
}

func insertAssetIndex(asset *AssetIndex) error {
	//Database Connection
    db := credentials.GetDB()
 
	query := `
		INSERT INTO meivan.asset_index(
			task_id, process_name, original_name, stored_name, file_path, 
			file_size, mime_type, uploaded_by, role_rights, rights, 
			checksum_sha256, created_at, updated_at, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at
	`
	var err error
	err = db.QueryRow(
		query,
		asset.TaskID,
		asset.ProcessName,
		asset.OriginalName,
		asset.StoredName,
		asset.FilePath,
		asset.FileSize,
		asset.MimeType,
		asset.UploadedBy,
		asset.RoleRights,
		asset.Rights,
		asset.ChecksumSHA256,
		time.Now(),
		time.Now(),
		asset.Status,
	).Scan(&asset.ID, &asset.CreatedAt, &asset.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert asset index: %w", err)
	}

	log.Printf("✅ Asset index inserted with ID: %d", asset.ID)
	return nil
}

// processPDFMetadataAndUploadJSON handles database insertion and uploads JSON to SFTP
func processPDFMetadataAndUploadJSON(sftpClient *sftp.Client, pdfBytes []byte, sftpPdfPath string, taskID, processName, originalName, uploadedBy string) error {
	fileSize := int64(len(pdfBytes))

	// Calculate checksum in memory
	hash := sha256.New()
	hash.Write(pdfBytes)
	checksum := hex.EncodeToString(hash.Sum(nil))

	asset := &AssetIndex{
		TaskID:         taskID,
		ProcessName:    processName,
		OriginalName:   originalName,
		StoredName:     originalName,
		FilePath:       sftpPdfPath, // Storing the SFTP path in DB
		FileSize:       fileSize,
		MimeType:       "application/pdf",
		UploadedBy:     uploadedBy,
		RoleRights:     "admin",
		Rights:         "read,write",
		ChecksumSHA256: checksum,
		Status:         "active",
	}

	err := insertAssetIndex(asset)
	if err != nil {
		return fmt.Errorf("failed to insert asset index: %w", err)
	}

	// Create JSON
	metaJSON, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal meta JSON: %w", err)
	}

	// Define JSON Path on SFTP (Same dir as PDF)
	dir := filepath.Dir(sftpPdfPath)
	baseName := filepath.Base(sftpPdfPath)
	jsonFileName := strings.TrimSuffix(baseName, filepath.Ext(baseName)) + "_meta.json"

	// Ensure separators are correct for SFTP (forward slash)
	jsonFilePath := filepath.Join(dir, jsonFileName)
	jsonFilePath = strings.ReplaceAll(jsonFilePath, "\\", "/")

	// Upload JSON to SFTP
	err = saveToSFTP(sftpClient, jsonFilePath, metaJSON)
	if err != nil {
		return fmt.Errorf("failed to save meta JSON to SFTP: %w", err)
	}

	log.Printf("✅ Meta JSON saved to SFTP: %s", jsonFilePath)
	return nil
}

// wrapFooterHTML wraps plain text footer with styled HTML if it doesn't already have HTML tags
func wrapFooterHTML(footer string) string {
	footer = strings.TrimSpace(footer)

	// Check if footer is empty
	if footer == "" {
		return ""
	}

	// Check if footer already contains HTML tags
	hasHTMLTags := strings.Contains(footer, "<div") ||
		strings.Contains(footer, "<p") ||
		strings.Contains(footer, "<table") ||
		strings.Contains(footer, "<ul") ||
		strings.Contains(footer, "<ol")

	// If it already has HTML, return as is
	if hasHTMLTags {
		return footer
	}

	// Otherwise, wrap plain text with styled HTML
	// Split by newlines to preserve paragraph structure
	lines := strings.Split(footer, "\n")

	var wrappedContent strings.Builder
	wrappedContent.WriteString(`<div style="border: 1px solid black; padding: 10px; max-width: 1250px;">`)

	var currentParagraph strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			// Empty line - close current paragraph if any
			if currentParagraph.Len() > 0 {
				wrappedContent.WriteString("<p>" + currentParagraph.String() + "</p>")
				currentParagraph.Reset()
			}
		} else {
			// Check if line looks like a list item (starts with a), b), c), etc. or 1., 2., etc.)
			isListItem := regexp.MustCompile(`^[a-z]\)|^\d+\.`).MatchString(line)

			if isListItem {
				// Close current paragraph if any
				if currentParagraph.Len() > 0 {
					wrappedContent.WriteString("<p>" + currentParagraph.String() + "</p>")
					currentParagraph.Reset()
				}

				// Make list items bold and italic
				wrappedContent.WriteString("<p><b><i>" + html.EscapeString(line) + "</i></b><br></p>")
			} else {
				// Add to current paragraph
				if currentParagraph.Len() > 0 {
					currentParagraph.WriteString(" ")
				}
				currentParagraph.WriteString(html.EscapeString(line))
			}
		}
	}

	// Close any remaining paragraph
	if currentParagraph.Len() > 0 {
		wrappedContent.WriteString("<p>" + currentParagraph.String() + "</p>")
	}

	wrappedContent.WriteString("</div>")

	return wrappedContent.String()
}

/*	correct working code for office order and amendment */
//smart spacing generateHTML  -> add the processID  next to []ApproverDetail,
func generateHTML(t *Template, logoBase64, qrBase64, signatureBase64 string, isDraft, includeSignature, includeApprovers bool, approvers []ApproverDetail, processID int) string {
	ccToHTML := ""
	if strings.TrimSpace(t.CCTo) != "" {
		recipients := strings.Split(t.CCTo, ",")
		var listItems strings.Builder
		for i, recipient := range recipients {
			recipient = strings.TrimSpace(recipient)
			if recipient != "" {
				cleanRecipient := strings.TrimSuffix(recipient, ".")
				if i > 0 {
					listItems.WriteString("<br>")
				}
				listItems.WriteString(cleanRecipient)
			}
		}
		ccToHTML = listItems.String()
	}

	formattedReference := strings.ReplaceAll(t.Reference, "\n", "<br>")

	isAmendmentOrCancellation := strings.EqualFold(strings.TrimSpace(t.Order_type), "Amendment") ||
		strings.EqualFold(strings.TrimSpace(t.Order_type), "Cancellation")

	signatureHTML := ""
	if includeSignature && signatureBase64 != "" {
		signatureHTML = fmt.Sprintf(`<div style="text-align: right;"><img src="file://%s" style="max-width: 150px; height: auto; display: inline-block; margin-bottom: 5px;" alt="Signature"><br>%s</div>`, signatureBase64, t.SignatureHTML)
	} else if includeSignature {
		signatureHTML = fmt.Sprintf(`<div style="text-align: right;">%s</div>`, t.SignatureHTML)
	} else if isDraft {
		signatureHTML = fmt.Sprintf(`<div style="text-align: right;">%s</div>`, t.SignatureHTML)
	}

	ccToAfterSignature := ""
	if isAmendmentOrCancellation && ccToHTML != "" {
		ccToAfterSignature = fmt.Sprintf(`
		<div class="copyto" style="margin-top: 20px;">
			<strong>Copy to:</strong><br>
			%s
		</div>`, ccToHTML)
	}

	watermarkStyle := ""
	watermarkHTML := ""
	if isDraft {
		watermarkStyle = `
		.watermark {
			position: fixed;
			font-size: 80pt;
			color: rgba(200, 200, 200, 0.15);
			font-weight: bold;
			transform: rotate(45deg);
			z-index: 9999;
			pointer-events: none;
			user-select: none;
		}`

		var wmBuilder strings.Builder
		positions := []struct{ top, left string }{
			{"10%", "10%"}, {"30%", "30%"}, {"50%", "50%"}, {"70%", "70%"}, {"90%", "20%"},
			{"20%", "80%"}, {"60%", "10%"}, {"80%", "40%"}, {"40%", "70%"},
		}
		for _, pos := range positions {
			wmBuilder.WriteString(fmt.Sprintf(`<div class="watermark" style="top:%s; left:%s;">DRAFT</div>`, pos.top, pos.left))
		}
		watermarkHTML = wmBuilder.String()
	}

	approversTableHTML := ""
	if includeApprovers {
		approversTableHTML = generateApproversTableHTML(approvers)
	}

	if isAmendmentOrCancellation {
		html := fmt.Sprintf(`<!DOCTYPE html>

<html>
<head>
<meta charset="UTF-8">
<title>Office Order</title>
<style>

	@page {
		size: A4;
		margin: 20mm;
	}

	@page {
		@bottom-center {
			content: "Page " counter(page) " of " counter(pages);
			font-size: 10pt;
			font-family: Arial, sans-serif;
			color: #000;
		}
	}

	body {
		font-family: Arial, sans-serif;
		font-size: 12pt;
		margin: 0;
		padding: 0;
		line-height: 1.5;
		position: relative;
		text-align: justify;
	}

	.header-container {
		display: table;
		width: 100%%;
		margin-bottom: 20px;
		padding-bottom: 10px;
	}

	.header-left {
		display: table-cell;
		width: 85%%;
		text-align: center;
		vertical-align: middle;
	}

	.header-right {
		display: table-cell;
		width: 15%%;
		text-align: center;
		vertical-align: top;
	}

	.header-logo {
		max-width: 1000px;
		height: auto;
	}

	.qr-code {
		width: 100px;
		height: 100px;
		margin-top: 5px;
	}

	.order-info {
		display: table;
		width: 100%%;
		margin-top: 15px;
		margin-bottom: 15px;
	}

	.order-left {
		display: table-cell;
		width: 50%%;
		vertical-align: top;
	}

	.order-right {
		display: table-cell;
		width: 50%%;
		text-align: right;
		vertical-align: top;
		padding-right: 10px;
	}

	table {
		border-collapse: collapse;
		width: 100%%;
		margin-top: 10px;
		margin-bottom: 10px;
	}

	th, td {
		border: 1px solid #000;
		padding: 8px;
		text-align: center;
		font-size: 12pt;
	}

	th {
		background-color: #f0f0f0;
		font-weight: bold;
	}

	.copyto {
		font-size: 12pt;
		margin-top: 20px;
	}

	.footer {
		font-size: 12pt;
		margin-top: 20px;
		padding: 10px;
		max-width: 100%%;
		text-align: justify;
	}

	.footer-section {
		font-size: 12pt;
		border: 1px solid #000;
		margin-top: 20px;
		padding: 10px;
		max-width: 100%%;
	}

	ul {
		margin: 10px 0;
		padding-left: 25px;
	}

	li {
		margin: 6px 0;
		line-height: 1.4;
	}

	p {
		margin: 8px 0;
		line-height: 1.5;
	}

	strong {
		font-weight: bold;
	}

	.signature-section {
		text-align: right;
		margin-top: 30px;
		margin-bottom: 20px;
	}

	.signature-section img {
		max-width: 200px;
		height: auto;
		display: inline-block;
		margin-bottom: 5px;
	}

	.recipient-section {
		margin-bottom: 3px;
		line-height: 1.6;
	}

%s
</style>
</head>
<body>
%s

<div class="single-page">

	<div class="header-container">
		<div class="header-left">
			<img src="data:image/png;base64,%s" class="header-logo">
		</div>
		<div class="header-right">
		</div>
	</div>

	<div class="order-info">
		<div class="order-left">
			<strong>%s</strong><br>
			%s
		</div>
		<br>

		<div class="order-right">
			<strong>Date:</strong> %s<br>
			<img src="data:image/png;base64,%s" class="qr-code">
		</div>
	</div>

	<div style="margin-bottom: 20px;">
		<div style="margin-bottom: 5px;"><strong>Sub:</strong> %s</div>

	<table style="width: 45%%; border-collapse: collapse; border: none;margin-left:-3px">
			<tbody>
				<tr>
					<td style="padding:0px;border:none;vertical-align: top;">
						<strong>Ref:</strong>
					</td>
					<td style="vertical-align: top; padding: 0;text-align: left;border:none;">
						%s
					</td>
				</tr>
			</tbody>
		</table>
	</div>

	<div>%s</div>

	<div class="signature-section"><strong>%s</strong></div>

	%s

	<div class="footer">
		%s
	</div>

	%s

</div>

</body>
</html>`,

			watermarkStyle,
			watermarkHTML,
			logoBase64,
			t.OrderNo,
			t.ToColumn,
			t.OrderDate,
			qrBase64,
			t.Subject,
			formattedReference,
			t.BodyHTML,
			signatureHTML,
			ccToAfterSignature,
			wrapFooterHTML(t.FooterHTML),
			approversTableHTML)

		return html
	}
	//
	// Add this code right after the isAmendmentOrCancellation check and before the "Logic for Standard Office Order" section

	// Check if this should use smart spacing based on process_id
	// Smart spacing applies when process_id is NOT 1 or 2
	isSmartSpacing := processID != 1 && processID != 2

	if isSmartSpacing {
		// Build Subject section conditionally
		subjectSection := ""
		if strings.TrimSpace(t.Subject) != "" {
			subjectSection = fmt.Sprintf(`<div style="margin-bottom: 5px;"><strong>Sub:</strong> %s</div>`, t.Subject)
		}

		// Build Reference section conditionally
		referenceSection := ""
		if strings.TrimSpace(t.Reference) != "" {
			referenceSection = fmt.Sprintf(`
	<table style="width: 45%%; border-collapse: collapse; border: none;margin-left:-3px">
		<tbody>
			<tr>
				<td style="padding:0px;border:none;vertical-align: top;">
					<strong>Ref:</strong>
				</td>
				<td style="vertical-align: top; padding: 0;text-align: left;border:none;">
					%s
				</td>
			</tr>
		</tbody>
	</table>`, formattedReference)
		}

		// Combine subject and reference sections
		subRefSection := ""
		if subjectSection != "" || referenceSection != "" {
			subRefSection = fmt.Sprintf(`<div style="margin-bottom: 20px;">
		%s
		%s
	</div>`, subjectSection, referenceSection)
		}

		// Build Copy to section
		copyToSection := ""
		if ccToHTML != "" {
			copyToSection = fmt.Sprintf(`
	<div class="copyto" style="margin-top: 20px;">
		<strong>Copy to:</strong><br>
		%s
	</div>`, ccToHTML)
		}

		// Build Footer section - adjust margin based on whether copy to exists
		footerSection := ""
		if strings.TrimSpace(t.FooterHTML) != "" {
			footerMargin := "20px"
			if ccToHTML == "" {
				footerMargin = "0px"
			}
			footerSection = fmt.Sprintf(`
	<div class="footer" style="margin-top: %s; padding: 10px; max-width: 100%%; text-align: justify; font-size: 12pt;">
		%s
	</div>`, footerMargin, wrapFooterHTML(t.FooterHTML))
		}

		// Build Approvers section - adjust margin based on what's above it
		approversSection := ""
		if includeApprovers && len(approvers) > 0 {
			approversMargin := "20px"
			if ccToHTML == "" && strings.TrimSpace(t.FooterHTML) == "" {
				approversMargin = "0px"
			}
			approversSection = fmt.Sprintf(`
	<div style="margin-top: %s;">
		%s
	</div>`, approversMargin, approversTableHTML)
		}

		html := fmt.Sprintf(`<!DOCTYPE html>

<html>
<head>
<meta charset="UTF-8">
<title>Office Order</title>
<style>

	@page {
		size: A4;
		margin: 20mm;
	}

	@page {
		@bottom-center {
			content: "Page " counter(page) " of " counter(pages);
			font-size: 10pt;
			font-family: Arial, sans-serif;
			color: #000;
		}
	}

	body {
		font-family: Arial, sans-serif;
		font-size: 12pt;
		margin: 0;
		padding: 0;
		line-height: 1.5;
		position: relative;
		text-align: justify;
	}

	.header-container {
		display: table;
		width: 100%%;
		margin-bottom: 20px;
		padding-bottom: 10px;
	}

	.header-left {
		display: table-cell;
		width: 85%%;
		text-align: center;
		vertical-align: middle;
	}

	.header-right {
		display: table-cell;
		width: 15%%;
		text-align: center;
		vertical-align: top;
	}

	.header-logo {
		max-width: 1000px;
		height: auto;
	}

	.qr-code {
		width: 100px;
		height: 100px;
		margin-top: 5px;
	}

	.order-info {
		display: table;
		width: 100%%;
		margin-top: 15px;
		margin-bottom: 15px;
	}

	.order-left {
		display: table-cell;
		width: 50%%;
		vertical-align: top;
	}

	.order-right {
		display: table-cell;
		width: 50%%;
		text-align: right;
		vertical-align: top;
		padding-right: 10px;
	}

	table {
		border-collapse: collapse;
		width: 100%%;
		margin-top: 10px;
		margin-bottom: 10px;
	}

	th, td {
		border: 1px solid #000;
		padding: 8px;
		text-align: center;
		font-size: 12pt;
	}

	th {
		background-color: #f0f0f0;
		font-weight: bold;
	}

	.copyto {
		font-size: 12pt;
	}

	.footer {
		font-size: 12pt;
	}

	.footer-section {
		font-size: 12pt;
		border: 1px solid #000;
		margin-top: 20px;
		padding: 10px;
		max-width: 100%%;
	}

	ul {
		margin: 10px 0;
		padding-left: 25px;
	}

	li {
		margin: 6px 0;
		line-height: 1.4;
	}

	p {
		margin: 8px 0;
		line-height: 1.5;
	}

	strong {
		font-weight: bold;
	}

	.signature-section {
		text-align: right;
		margin-top: 30px;
		margin-bottom: 20px;
	}

	.signature-section img {
		max-width: 200px;
		height: auto;
		display: inline-block;
		margin-bottom: 5px;
	}

	.recipient-section {
		margin-bottom: 3px;
		line-height: 1.6;
	}

%s
</style>
</head>
<body>
%s

<div class="smart-spacing-page">

	<div class="header-container">
		<div class="header-left">
			<img src="data:image/png;base64,%s" class="header-logo">
		</div>
		<div class="header-right">
		</div>
	</div>

	<div class="order-info">
		<div class="order-left">
			<strong>%s</strong><br>
			%s
		</div>
		<br>

		<div class="order-right">
			<strong>Date:</strong> %s<br>
			<img src="data:image/png;base64,%s" class="qr-code">
		</div>
	</div>

	%s

	<div>%s</div>

	<div class="signature-section"><strong>%s</strong></div>

	%s

	%s

	%s

</div>

</body>
</html>`,
			watermarkStyle,
			watermarkHTML,
			logoBase64,
			t.OrderNo,
			t.ToColumn,
			t.OrderDate,
			qrBase64,
			subRefSection,
			t.BodyHTML,
			signatureHTML,
			copyToSection,
			footerSection,
			approversSection)

		return html
	}
	//smart spacing generateHTML  end

	// Logic for Standard Office Order (Default Template)
	html := fmt.Sprintf(`<!DOCTYPE html>

<html>
<head>
<meta charset="UTF-8">
<title>Office Order</title>
<style>

	@page {
		size: A4;
		margin: 20mm;
	}

	@page {
		@bottom-center {
			content: "Page " counter(page) " of " counter(pages);
			font-size: 10pt;
			font-family: Arial, sans-serif;
			color: #000;
		}
	}

	body {
		font-family: Arial, sans-serif;
		font-size: 12pt;
		margin: 0;
		padding: 0;
		line-height: 1.5;
		position: relative;
		text-align: justify;
	}

	.header-container {
		display: table;
		width: 100%%;
		margin-bottom: 20px;
		padding-bottom: 10px;
	}

	.header-left {
		display: table-cell;
		width: 85%%;
		text-align: center;
		vertical-align: middle;
	}

	.header-right {
		display: table-cell;
		width: 15%%;
		text-align: center;
		vertical-align: top;
	}

	.header-logo {
		max-width: 1000px;
		height: auto;
	}

	.qr-code {
		width: 100px;
		height: 100px;
		margin-top: 5px;
	}

	.order-info {
		display: table;
		width: 100%%;
		margin-top: 15px;
		margin-bottom: 15px;
	}

	.order-left {
		display: table-cell;
		width: 50%%;
		vertical-align: top;
	}

	.order-right {
		display: table-cell;
		width: 50%%;
		text-align: right;
		vertical-align: top;
		padding-right: 10px;
	}

	table {
		border-collapse: collapse;
		width: 100%%;
		margin-top: 10px;
		margin-bottom: 10px;
	}

	th, td {
		border: 1px solid #000;
		padding: 8px;
		text-align: center;
		font-size: 12pt;
	}

	th {
		background-color: #f0f0f0;
		font-weight: bold;
	}

	.copyto {
		font-size: 12pt;
		margin-top: 20px;
	}

	.footer {
		font-size: 12pt;
		margin-top: 20px;
		padding: 10px;
		max-width: 100%%;
		text-align: justify;
	}

	.footer-section {
		font-size: 12pt;
		border: 1px solid #000;
		margin-top: 20px;
		padding: 10px;
		max-width: 100%%;
	}

	ul {
		margin: 10px 0;
		padding-left: 25px;
	}

	li {
		margin: 6px 0;
		line-height: 1.4;
	}

	p {
		margin: 8px 0;
		line-height: 1.5;
	}

	strong {
		font-weight: bold;
	}

	.signature-section {
		text-align: right;
		margin-top: 30px;
		margin-bottom: 20px;
	}

	.signature-section img {
		max-width: 200px;
		height: auto;
		display: inline-block;
		margin-bottom: 5px;
	}

	.recipient-section {
		margin-bottom: 3px;
		line-height: 1.6;
	}

	.page-break {
		page-break-after: always;
		break-after: page;
		height: 0;
		margin: 0;
		padding: 0;
	}

%s
</style>
</head>
<body>
%s

<div class="page-one">

	<div class="header-container">
		<div class="header-left">
			<img src="data:image/png;base64,%s" class="header-logo">
		</div>
		<div class="header-right">
		</div>
	</div>

	<div class ="order-info">
		<div class="order-left">
			<strong>%s</strong><br>
			%s
		</div>
		<br>

		<div class="order-right">
			<strong>Date:</strong> %s<br>
			<img src="data:image/png;base64,%s" class="qr-code">
		</div>
	</div>

	<div style="margin-bottom: 20px;">
		<div style="margin-bottom: 5px;"><strong>Sub:</strong> %s</div>
		<div><strong>Ref:</strong> %s</div>
	</div>

	<div>%s</div>

	<div class="signature-section"><strong>%s</strong></div>

</div>

<div class="page-break"></div>

<div class="page-two">

	<div class="copyto">
		<strong>Copy to:</strong><br>
		%s
	</div>

	<div class="footer">
		%s
	</div>

	%s

</div>

</body>
</html>`,

		watermarkStyle,
		watermarkHTML,
		logoBase64,
		t.OrderNo,
		t.ToColumn,
		t.OrderDate,
		qrBase64,
		t.Subject,
		formattedReference,
		t.BodyHTML,
		signatureHTML,
		ccToHTML,
		wrapFooterHTML(t.FooterHTML),
		approversTableHTML)

	return html
}


func generatePDF(html string) (*pdf.PDFGenerator, error) {
	pdfg, err := pdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	page := pdf.NewPageReader(bytes.NewReader([]byte(html)))
	page.EnableLocalFileAccess.Set(true)

	page.FooterCenter.Set("Page [page] of [toPage]")
	page.FooterFontSize.Set(9)
	page.FooterSpacing.Set(5)

	page.DisableSmartShrinking.Set(false)
	page.PrintMediaType.Set(true)
	page.Zoom.Set(1.0)
	page.NoStopSlowScripts.Set(false)
	page.LoadErrorHandling.Set("ignore")

	pdfg.AddPage(page)

	pdfg.Dpi.Set(300)
	pdfg.MarginTop.Set(10)
	pdfg.MarginBottom.Set(20)
	pdfg.MarginLeft.Set(15)
	pdfg.MarginRight.Set(15)
	pdfg.Grayscale.Set(false)

	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	return pdfg, nil
}

// PDF generation handler - Returns JSON paths for completed, streams PDF for draft
func GeneratePDFHandler(w http.ResponseWriter, r *http.Request) {
	// 1️⃣ Validate Request Method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	// 2️⃣ Read and Parse Request Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req PDFRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// 3️⃣ Split and decrypt (Using utils) - INPUT IS STILL ENCRYPTED
	parts := strings.Split(req.Data, "||")
	if len(parts) != 2 {
		http.Error(w, "Invalid Data format", http.StatusBadRequest)
		return
	}
	pid := parts[0]
	encryptedPart := parts[1]

	key, err := utils.GetDecryptKey(pid)
	if err != nil {
		http.Error(w, "Decryption key fetch failed", http.StatusUnauthorized)
		return
	}

	decryptedJSON, err := utils.DecryptAES(encryptedPart, key)
	if err != nil {
		http.Error(w, "Decryption failed", http.StatusUnauthorized)
		return
	}

	// Unmarshal to map first for Token extraction
	var decryptedMap map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedJSON), &decryptedMap); err != nil {
		http.Error(w, "Invalid decrypted data", http.StatusBadRequest)
		return
	}

	// Extract Token and set in Header
	token, ok := decryptedMap["token"].(string)
	if !ok || token == "" {
		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}
	r.Header.Set("token", token)

	// 4️⃣ Authentication check
	if !auth.HandleRequestfor_apiname_ipaddress_token(w, r) {
		return
	}

	// 5️⃣ Log Request and Process Logic
	auth.LogRequestInfo(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth.IsValidIDFromRequest(r); err != nil {
			http.Error(w, "Invalid TOKEN", http.StatusBadRequest)
			return
		}

		// --- START BUSINESS LOGIC ---

		// Re-unmarshal into specific struct for PDF logic
		var decryptedData DecryptedPDFData
		if err := json.Unmarshal([]byte(decryptedJSON), &decryptedData); err != nil {
			http.Error(w, "Invalid data structure", http.StatusBadRequest)
			return
		}

		orderNo := decryptedData.OrderNo
		processID := decryptedData.ProcessID
		taskID := decryptedData.TaskID
		status := strings.ToLower(strings.TrimSpace(decryptedData.Status))

		// Validate required parameters
		if processID == 0 {
			http.Error(w, "process_id is required", http.StatusBadRequest)
			return
		}

		hasOrderNo := orderNo != "" && orderNo != "null"
		hasTaskID := taskID != "" && taskID != "null"

		if !hasOrderNo && !hasTaskID {
			http.Error(w, "Either (order_no + process_id) or (process_id + task_id) is required", http.StatusBadRequest)
			return
		}

		if hasOrderNo && hasTaskID {
			http.Error(w, "Provide either (order_no + process_id) or (process_id + task_id), not both", http.StatusBadRequest)
			return
		}

		if status != "draft" && status != "completed" {
			http.Error(w, "status must be either 'Draft' or 'Completed'", http.StatusBadRequest)
			return
		}

		log.Printf("📥 Received request - OrderNo: %s, ProcessID: %d, TaskID: %s, Status: %s", orderNo, processID, taskID, status)

		// Fetch template
		t, err := getTemplateFromDB(orderNo, processID, taskID)
		if err != nil {
			http.Error(w, "Template not found: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Fetch approvers
		approvers, err := getApproverDetails(orderNo, processID, taskID)
		if err != nil {
			log.Printf("⚠️ Warning: Failed to fetch approver details: %v", err)
		}

		// Read Logo
		logoBytes, err := ioutil.ReadFile("/var/www/html/go_projects/HRMODULE/hrnew/7000port/newpdfpoc/iitm_logo.png")
		if err != nil {
			http.Error(w, "Logo file not found: "+err.Error(), http.StatusInternalServerError)
			return
		}
		logoBase64 := base64.StdEncoding.EncodeToString(logoBytes)

		// Generate QR
		qrBase64, err := generateQRCode(t.OrderNo)
		if err != nil {
			http.Error(w, "QR code generation error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Fetch Signature
		roleName := extractRoleName(t.SignatureHTML)
		var signatureFilePath string
		signatureDatBytes, err := fetchSignatureFromAPI(roleName)
		if err == nil {
			// Signatures still saved locally temporarily for inclusion in PDF generation (wkhtmltopdf needs local path or http)
			signaturesDir := "/var/www/html/go_projects/HRMODULE/hrnew/9000/PDF/signatures"
			os.MkdirAll(signaturesDir, 0755)
			sanitizedRole := sanitizeFilename(roleName)
			signatureDatPath := filepath.Join(signaturesDir, fmt.Sprintf("signature_%s.dat", sanitizedRole))
			ioutil.WriteFile(signatureDatPath, signatureDatBytes, 0644)

			tempSigDir := "/tmp/pdf_signatures"
			os.MkdirAll(tempSigDir, 0755)
			signatureFilePath = filepath.Join(tempSigDir, fmt.Sprintf("signature_%s_temp.png", sanitizedRole))
			decryptSignatureDat(signatureDatBytes, signatureFilePath)
		}

		sanitizedOrderNo := sanitizeFilename(t.OrderNo)
		sanitizedTaskID := sanitizeFilename(t.TaskID)
		sanitizedEmployeeID := sanitizeFilename(t.EmployeeID)
		sanitizedProcessCode := sanitizeFilename(t.ProcessCode)
		currentYear := time.Now().Format("2006")

		// Variables to hold response data
		var finalPDFBytes []byte
		var finalFilename string
		var savedPaths map[string]string

		// --- Logic for Draft vs Completed ---
		if status == "draft" {
			draftHTML := generateHTML(t, logoBase64, qrBase64, "", true, false, false, nil, processID) //smart spacing processID
			draftPDF, err := generatePDF(draftHTML)
			if err != nil {
				http.Error(w, "Failed to generate draft PDF: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Set data for streaming (get bytes, don't save)
			finalPDFBytes = draftPDF.Buffer().Bytes()
			finalFilename = fmt.Sprintf("draft_%s.pdf", sanitizedOrderNo)

		} else if status == "completed" {
			// 📡 Connect to SFTP
			sftpClient, sshClient, err := ConnectSFTP()
			if err != nil {
				http.Error(w, "SFTP Connection Error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer sshClient.Close()
			defer sftpClient.Close()

			// Define Remote Paths
			// Path 1: /HR_Test/meivan/off_order/{process_code}/{year}/{taskid}/
			officeOrderDir := fmt.Sprintf("/HR_Test/meivan/off_order/%s/%s/%s", sanitizedProcessCode, currentYear, sanitizedTaskID)

			// Path 2: /HR_Test/employee/{employeeid}/off_order/{process_code}/{taskid}/
			employeeDir := fmt.Sprintf("/HR_Test/employee/%s/off_order/%s", sanitizedEmployeeID, sanitizedProcessCode)

			// Create Remote Directories
			if err := sftpMkdirAll(sftpClient, officeOrderDir); err != nil {
				http.Error(w, "Failed to create office dir on SFTP: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if err := sftpMkdirAll(sftpClient, employeeDir); err != nil {
				http.Error(w, "Failed to create employee dir on SFTP: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// --- Generate & Upload Office Copy ---
			officeHTML := generateHTML(t, logoBase64, qrBase64, signatureFilePath, false, true, true, approvers, processID) //smart spacing processID
			officePDF, err := generatePDF(officeHTML)
			if err != nil {
				http.Error(w, "Failed to generate office copy: "+err.Error(), http.StatusInternalServerError)
				return
			}
			officeCopyBytes := officePDF.Buffer().Bytes()
			officePath := fmt.Sprintf("%s/officecopy_%s.pdf", officeOrderDir, sanitizedOrderNo)

			err = saveToSFTP(sftpClient, officePath, officeCopyBytes)
			if err != nil {
				http.Error(w, "Failed to save office copy to SFTP: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// --- Generate & Upload User Copy ---
			userHTML := generateHTML(t, logoBase64, qrBase64, signatureFilePath, false, true, false, nil, processID) //smart spacing processID
			userPDF, err := generatePDF(userHTML)
			if err != nil {
				http.Error(w, "Failed to generate user copy: "+err.Error(), http.StatusInternalServerError)
				return
			}
			userCopyBytes := userPDF.Buffer().Bytes()

			// Save User Copy to Office Dir
			userPathOffice := fmt.Sprintf("%s/usercopy_%s.pdf", officeOrderDir, sanitizedOrderNo)
			err = saveToSFTP(sftpClient, userPathOffice, userCopyBytes)
			if err != nil {
				http.Error(w, "Failed to save user copy (office dir) to SFTP: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Save User Copy to Employee Dir
			userPathEmployee := fmt.Sprintf("%s/usercopy_%s.pdf", employeeDir, sanitizedOrderNo)
			err = saveToSFTP(sftpClient, userPathEmployee, userCopyBytes)
			if err != nil {
				http.Error(w, "Failed to save user copy (employee dir) to SFTP: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// --- Handle Metadata (Checksum, DB, JSON to SFTP) ---
			// Stored name in DB will be the user copy name
			storedName := fmt.Sprintf("usercopy_%s.pdf", sanitizedOrderNo)

			// Process metadata and upload JSON to SFTP (using the Employee directory path for the file)
			err = processPDFMetadataAndUploadJSON(sftpClient, userCopyBytes, userPathEmployee, t.TaskID, t.ProcessCode, storedName, t.EmployeeID)
			if err != nil {
				log.Printf("⚠️ Warning: Failed to save metadata: %v", err)
			}

			// Store the saved paths for response
			savedPaths = map[string]string{
				"office_copy_path":        officePath,
				"user_copy_path_office":   userPathOffice,
				"employee_user_copy_path": userPathEmployee,
				"order_no":                t.OrderNo,
				"task_id":                 t.TaskID,
			}

			log.Printf("✅ PDFs saved to SFTP successfully - Office: %s, Emp: %s", officePath, userPathEmployee)
		}

		// 6️⃣ Handle Response based on status
		if status == "completed" {
			// Return JSON with saved paths
			responseData := map[string]interface{}{
				"status":  "success",
				"message": "PDF files generated and saved successfully to SFTP",
				"data":    savedPaths,
			}

			// Log the response
			auth.SaveResponseLog(
				r,
				responseData,
				http.StatusOK,
				"application/json",
				0,
				string(body),
			)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(responseData)

		} else {
			// For draft, stream the PDF
			auth.SaveResponseLog(
				r,
				map[string]string{"message": "PDF File Served", "filename": finalFilename},
				http.StatusOK,
				"application/pdf",
				len(finalPDFBytes),
				string(body),
			)

			w.Header().Set("Content-Type", "application/pdf")
			w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", finalFilename))
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(finalPDFBytes)))
			w.WriteHeader(http.StatusOK)
			w.Write(finalPDFBytes)
		}

	})).ServeHTTP(w, r)
}



