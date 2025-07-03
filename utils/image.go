package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"

	// 注册解码器，以便 image.Decode 能处理多种格式
	_ "image/png"
	"os"

	"github.com/nfnt/resize"
)

// maxDimension 定义了图片在调整大小后的最大宽度或高度。
// 1024x1024 像素对于人脸识别来说足够清晰，且能显著减小文件大小。
const maxDimension uint = 1024

// ProcessAndEncodeImage 读取、压缩并编码图片。
// 它会打开一个图片文件，如果图片的尺寸超过 maxDimension，
// 就会按比例将其缩小，然后将结果编码为 JPEG 格式的 Base64 字符串。
func ProcessAndEncodeImage(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 解码图片以获取其数据和尺寸
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// 如果图片的宽度或高度大于我们的阈值，则进行压缩
	if uint(img.Bounds().Dx()) > maxDimension || uint(img.Bounds().Dy()) > maxDimension {
		img = resize.Thumbnail(maxDimension, maxDimension, img, resize.Lanczos3)
	}

	// 将处理后的图片编码为 JPEG 格式存入缓冲区
	// JPEG 格式很适合照片，并且提供了很好的压缩率
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 90}); err != nil {
		return "", err
	}

	// 返回缓冲区的 Base64 编码字符串
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// ImageFileToBase64 读取图片文件并将其转换为 base64 编码的字符串. (此函数将被新的 ProcessAndEncodeImage 替代)
func ImageFileToBase64(filePath string) (string, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}
