package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-face-id-validator/api"
	"go-face-id-validator/utils"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// 请将 '你自己的AppCode' 替换为你的 AppCode
const defaultAppCode = "xxxxxxxxxxxx"

type appState struct {
	selectedFilePath string
	apiClient        api.Client
	window           fyne.Window
}

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(newChineseTheme()) // 应用自定义的中文主题
	myWindow := myApp.NewWindow("人脸身份证比对 (Go/Fyne)")
	myWindow.Resize(fyne.NewSize(600, 650)) // 增加窗口高度以容纳新控件

	// 初始化 API 客户端
	// 这里可以轻松换成其他实现了 api.Client 接口的客户端
	apiClient := api.NewAliyunClient(defaultAppCode)

	state := &appState{
		apiClient: apiClient,
		window:    myWindow,
	}

	// --- UI 组件 ---
	appCodeEntry := widget.NewPasswordEntry() // 使用密码样式输入框隐藏AppCode
	appCodeEntry.SetPlaceHolder("在此输入您的AppCode")
	// 尝试从配置文件加载AppCode
	if conf, err := LoadConfig(); err == nil {
		appCodeEntry.SetText(conf.AppCode)
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("请输入姓名")

	idCardEntry := widget.NewEntry()
	idCardEntry.SetPlaceHolder("请输入身份证号码")

	filePathLabel := widget.NewLabel("未选择图片")
	filePathLabel.Wrapping = fyne.TextWrapWord

	// 使用一个可滚动的标签来显示结果，以获得更好的字体渲染效果
	resultLabel := widget.NewLabel("识别结果将显示在这里...")
	resultLabel.Wrapping = fyne.TextWrapWord
	resultScroll := container.NewScroll(resultLabel)

	// 图片选择按钮
	selectBtn := widget.NewButton("选择图片", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if reader == nil {
				return
			}
			state.selectedFilePath = reader.URI().Path()
			filePathLabel.SetText(state.selectedFilePath)
			defer reader.Close()
		}, myWindow)
		// 设置文件过滤器，移除.bmp格式
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fd.Show()
	})

	// 提交按钮
	submitBtn := widget.NewButton("开始比对", func() {
		// 在新的 goroutine 中运行验证，避免 UI 阻塞
		// 从UI控件获取所有需要的数据
		go state.validateIdentity(appCodeEntry.Text, nameEntry.Text, idCardEntry.Text, resultLabel)
	})

	// --- 布局 ---
	form := widget.NewForm(
		widget.NewFormItem("AppCode", appCodeEntry),
		widget.NewFormItem("姓名", nameEntry),
		widget.NewFormItem("身份证号", idCardEntry),
	)

	// 使用 VBox 来更好地控制文件路径标签的换行
	fileSelectBox := container.NewVBox(
		container.NewHBox(widget.NewLabel("人脸图片:"), selectBtn),
		filePathLabel,
	)

	// 按钮的视觉美化
	submitBtnContainer := container.NewCenter(
		container.NewMax(
			widget.NewButton("", func() {}), // Background button
			canvas.NewRectangle(color.NRGBA{R: 0x4C, G: 0xAF, B: 0x50, A: 0xFF}),
			submitBtn,
		),
	)

	// 结果区域
	resultContainer := container.NewBorder(
		widget.NewLabel("识别结果:"),
		nil, nil, nil,
		resultScroll, // 使用可滚动容器
	)

	// 整体布局
	content := container.New(
		layout.NewBorderLayout(form, submitBtnContainer, fileSelectBox, nil),
		form,
		fileSelectBox,
		submitBtnContainer,
		resultContainer, // 将结果区域放在底部
	)

	containerWithPadding := container.New(layout.NewPaddedLayout(), content)

	myWindow.SetContent(containerWithPadding)
	myWindow.ShowAndRun()
}

// validateIdentity 从UI收集数据，调用API，并显示结果.
func (s *appState) validateIdentity(appCode, name, idCardNo string, resultLabel *widget.Label) {
	// 更新UI显示加载状态
	resultLabel.SetText("正在识别中，请稍候...")

	// 尝试保存AppCode
	if appCode != "" {
		// 不论API调用是否成功，都保存用户输入的AppCode供下次使用
		// 忽略错误，因为这只是一个便利功能
		_ = SaveConfig(Config{AppCode: appCode})
	}

	if name == "" || idCardNo == "" || s.selectedFilePath == "" {
		dialog.ShowInformation("输入不完整", "请输入姓名、身份证号并选择图片。", s.window)
		resultLabel.SetText("")
		return
	}

	if s.apiClient == nil {
		dialog.ShowError(fmt.Errorf("API客户端未初始化"), s.window)
		resultLabel.SetText("")
		return
	}

	// 检查AppCode是否有效
	if appCode == "" || appCode == "你自己的AppCode" {
		dialog.ShowInformation("配置错误", "请在UI界面输入有效的AppCode。", s.window)
		resultLabel.SetText("配置错误")
		return
	}
	// 在验证前，使用从UI获取的AppCode设置客户端
	s.apiClient.SetAppCode(appCode)

	// 将图片转换为Base64，此函数现在包含压缩逻辑
	facePhotoBase64, err := utils.ProcessAndEncodeImage(s.selectedFilePath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("处理图片失败: %w", err), s.window)
		resultLabel.SetText(fmt.Sprintf("图片处理失败: %v", err))
		return
	}

	// 调用API
	result, err := s.apiClient.Validate(name, idCardNo, facePhotoBase64)
	if err != nil {
		// 对特定错误进行更友好的提示
		if strings.Contains(err.Error(), "400") || strings.Contains(err.Error(), "姓名或身份证号有误") {
			dialog.ShowInformation("输入错误", "输入信息有误，请仔细检查您的姓名和身份证号码是否正确、一致。", s.window)
		} else {
			dialog.ShowError(fmt.Errorf("API 请求失败: %w", err), s.window)
		}
		resultLabel.SetText(fmt.Sprintf("请求出错: %v", err))
		return
	}

	// 格式化并显示结果
	resultLabel.SetText(formatResult(result))
}

// formatResult 将 API 结果格式化为用户友好的字符串.
func formatResult(res *api.ValidationResult) string {
	var builder strings.Builder

	builder.WriteString("--- 比对结论 ---\n")
	if res.Success {
		conclusion := "未知"
		switch res.ResultCode {
		case 1:
			conclusion = "✅ 认证成功：同一人"
		case 2:
			conclusion = "❌ 认证失败：不同人"
		case 3:
			conclusion = "❓ 无法确认"
		default:
			conclusion = fmt.Sprintf("未知状态码: %d", res.ResultCode)
		}
		builder.WriteString(fmt.Sprintf("判定结果: %s\n", conclusion))

		if res.Message != "" {
			builder.WriteString(fmt.Sprintf("系统信息: %s\n", res.Message))
		}
		if res.Score > 0 {
			builder.WriteString(fmt.Sprintf("人脸相似度: %.2f%%\n", res.Score*100))
		}

		builder.WriteString("\n--- 身份信息 ---\n")
		builder.WriteString(fmt.Sprintf("性别: %s\n", formatValue(res.Sex)))
		builder.WriteString(fmt.Sprintf("生日: %s\n", formatValue(res.Birthday)))
		builder.WriteString(fmt.Sprintf("地址: %s\n", formatValue(res.Address)))

	} else {
		builder.WriteString(fmt.Sprintf("❌ API 调用失败: %s\n", res.Message))
	}

	builder.WriteString("\n--- 原始API响应 ---\n")
	// 尝试美化JSON输出
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(res.RawResponse), "", "    "); err == nil {
		builder.Write(prettyJSON.Bytes())
	} else {
		builder.WriteString(res.RawResponse)
	}

	return builder.String()
}

// formatValue 检查值是否为空，如果为空则返回 "未提供".
func formatValue(v string) string {
	if v == "" {
		return "未提供"
	}
	return v
}
