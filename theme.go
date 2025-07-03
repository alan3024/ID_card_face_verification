package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// chineseTheme 是一个自定义主题，它继承了 Fyne 的默认主题，
// 但覆盖了字体设置，以使用我们自己打包的中文字体。
type chineseTheme struct {
	fyne.Theme
}

// newChineseTheme 创建一个新的自定义中文主题实例。
func newChineseTheme() fyne.Theme {
	return &chineseTheme{Theme: theme.DefaultTheme()}
}

// Font 返回给定文本样式的字体资源。
// 在这里，我们强制它返回我们打包的中文字体（来自 bundled.go）。
// 注意：`resourceChineseTtf` 变量将在您运行 `fyne bundle` 命令后，在 `bundled.go` 文件中自动生成。
func (t *chineseTheme) Font(style fyne.TextStyle) fyne.Resource {
	return resourceChineseTtf
}

// 为了确保其他主题元素（如颜色、图标、尺寸）保持默认外观，
// 我们需要明确地从嵌入的默认主题中调用它们的方法。
func (t *chineseTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return t.Theme.Color(name, variant)
}

func (t *chineseTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.Theme.Icon(name)
}

func (t *chineseTheme) Size(name fyne.ThemeSizeName) float32 {
	return t.Theme.Size(name)
}
