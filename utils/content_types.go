package utils

var ExtensionToContentType = map[string]string{
	".txt":  "text/plain",
	".csv":  "text/csv",
	".html": "text/html",
	".htm":  "text/html",
	".json": "application/json",
	".xml":  "application/xml",

	// Images
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".svg":  "image/svg+xml",
	".webp": "image/webp",

	// Documents
	".pdf":  "application/pdf",
	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xls":  "application/vnd.ms-excel",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".ppt":  "application/vnd.ms-powerpoint",
	".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",

	// Archives
	".zip": "application/zip",
	".rar": "application/vnd.rar",
	".7z":  "application/x-7z-compressed",
	".tar": "application/x-tar",
	".gz":  "application/gzip",

	// Audio
	".mp3": "audio/mpeg",
	".wav": "audio/wav",
	".ogg": "audio/ogg",

	// Video
	".mp4": "video/mp4",
	".avi": "video/x-msvideo",
	".mov": "video/quicktime",
	".mkv": "video/x-matroska",
}

func IsValidExtension(ext string) bool {
	switch ext {
	case ".pdf", ".jpg", ".jpeg", ".png":
		return true
	}
	return false
}
