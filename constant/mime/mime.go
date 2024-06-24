package mime

// Type represents Internet Media Types as defined by IANA.
// Reference: https://www.iana.org/assignments/media-types/media-types.xhtml
type Type string

// noinspection GoUnusedConst
const (
	AACAudio                 Type = "audio/aac"
	AVIFImage                Type = "image/avif"
	AVIVideo                 Type = "video/x-msvideo"
	AbiWordDocument          Type = "application/x-abiword"
	AmazonKindleEBook        Type = "application/vnd.amazon.ebook"
	AppleInstallerPackage    Type = "application/vnd.apple.installer+xml"
	ArchiveDocument          Type = "application/x-freearc"
	BZip2Archive             Type = "application/x-bzip2"
	BZipArchive              Type = "application/x-bzip"
	BinaryData               Type = "application/octet-stream"
	BitmapImage              Type = "image/bmp"
	BourneShellScript        Type = "application/x-sh"
	CDAudio                  Type = "application/x-cdf"
	CSS                      Type = "text/css"
	CSV                      Type = "text/csv"
	CShellScript             Type = "application/x-csh"
	EPUB                     Type = "application/epub+zip"
	GIF                      Type = "image/gif"
	GZipCompressedArchive    Type = "application/gzip"
	HTML                     Type = "text/html"
	ICalendar                Type = "text/calendar"
	IconFormat               Type = "image/vnd.microsoft.icon"
	JPEG                     Type = "image/jpeg"
	JSON                     Type = "application/json"
	JSONLD                   Type = "application/ld+json"
	JavaArchive              Type = "application/java-archive"
	JavaScript               Type = "text/javascript"
	JavaScriptModule         Type = "text/javascript"
	MIDI                     Type = "audio/midi"
	MP3Audio                 Type = "audio/mpeg"
	MP4Video                 Type = "video/mp4"
	MPEGTransportStream      Type = "video/mp2t"
	MPEGVideo                Type = "video/mpeg"
	MSEmbeddedOpenTypeFonts  Type = "application/vnd.ms-fontobject"
	MSExcel                  Type = "application/vnd.ms-excel"
	MSExcelOpenXML           Type = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	MSPowerPoint             Type = "application/vnd.ms-powerpoint"
	MSPowerPointOpenXML      Type = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	MSVisio                  Type = "application/vnd.visio"
	MSWord                   Type = "application/msword"
	MSWordOpenXML            Type = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	OGG                      Type = "application/ogg"
	OGGAudio                 Type = "audio/ogg"
	OGGVideo                 Type = "video/ogg"
	OpenDocumentPresentation Type = "application/vnd.oasis.opendocument.presentation"
	OpenDocumentSpreadsheet  Type = "application/vnd.oasis.opendocument.spreadsheet"
	OpenDocumentText         Type = "application/vnd.oasis.opendocument.text"
	OpenTypeFont             Type = "font/otf"
	OpusAudio                Type = "audio/opus"
	PDF                      Type = "application/pdf"
	PHP                      Type = "application/x-httpd-php"
	PNG                      Type = "image/png"
	RARArchive               Type = "application/vnd.rar"
	RichTextFormat           Type = "application/rtf"
	SVG                      Type = "image/svg+xml"
	SevenZipArchive          Type = "application/x-7z-compressed"
	TARArchive               Type = "application/x-tar"
	TIFF                     Type = "image/tiff"
	Text                     Type = "text/plain"
	ThreeG2AudioVideo        Type = "video/3gpp2"
	ThreeGPAudioVideo        Type = "video/3gpp"
	TrueTypeFont             Type = "font/ttf"
	WAVAudio                 Type = "audio/wav"
	WEBMAudio                Type = "audio/webm"
	WEBMVideo                Type = "video/webm"
	WEBPImage                Type = "image/webp"
	WOFF                     Type = "font/woff"
	WOFF2                    Type = "font/woff2"
	XHTML                    Type = "application/xhtml+xml"
	XML                      Type = "application/xml"
	XUL                      Type = "application/vnd.mozilla.xul+xml"
	ZIPArchive               Type = "application/zip"
)

// String returns the string representation of the Internet Media Type.
func (t Type) String() string {
	return string(t)
}

// Parse parses the string into a Type.
func Parse(value string) Type {
	return Type(value)
}
