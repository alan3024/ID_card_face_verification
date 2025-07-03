# Go 人脸身份验证应用
阿里云市场购买API  https://marketnext.console.aliyun.com/
【精准识别】人像比对-人证比对-人脸身份证比对-人脸身份证比对-人脸三要素对比-人脸识别-人脸识别-人脸实名认证-人像比对
这是一个使用 Go 语言和 Fyne UI 库编写的桌面应用程序，用于通过姓名、身份证号和人脸照片来验证身份。

## 主要功能
-   **图形用户界面**：简洁的 UI，用于输入 AppCode、姓名和身份证号。
-   **配置自动保存**：在程序目录下自动创建 `config.json` 文件，用于保存您输入的 AppCode，下次启动时会自动加载，无需重复输入。
-   **图片选择与压缩**：支持从本地选择图片，并能自动压缩过大的图片，以避免上传时因文件过大导致的 `413` 错误。
-   **清晰的结果展示**：使用专门的标签和滚动条来显示 API 返回的结果，字体渲染清晰。
-   **中文字体支持**：通过打包"思源黑体"字体，完美解决了中文乱码问题。
-   **隐藏控制台**：在 Windows 上编译出的 `.exe` 文件运行时，不会弹出黑色的控制台窗口。
-   **API 接口预留**：代码中预留了 `api.Client` 接口，方便开发者替换为自己的后端服务。
## 使用说明
### 方式一：直接运行可执行文件 (推荐)
1.  将编译好的 `.exe` 文件放在任意文件夹中。
2.  双击运行 `go-face-id-validator.exe` 文件。
3.  **首次运行**：
    -   在 `AppCode` 输入框中填入您的阿里云 AppCode。
    -   此后，程序会在同目录下自动创建一个 `config.json` 文件来保存您的 AppCode。
4.  **后续运行**：
    -   程序会自动从 `config.json` 读取 AppCode 并填入，您无需再次输入。
5.  在界面中填入姓名和身份证号。
6.  点击"选择图片"按钮选择一张人脸照片。
7.  点击"开始比对"进行验证。
### 方式二：从源码编译 (高级用户)
如果您想自行修改代码或编译，请遵循以下步骤。
#### 1. 环境准备
-   安装 [Go](https://golang.org/doc/install) (版本 1.18 或更高)。
-   安装 Fyne 的[系统依赖](https://developer.fyne.io/started/#prerequisites)。对于 Windows，通常需要安装 TDM-GCC 或 MSYS2 环境中的 mingw-w64。
#### 2. 准备中文字体
为解决中文乱码问题，您需要手动准备一个字体文件。
1.  下载一个免费开源的中文字体。推荐使用[**思源黑体 (Source Han Sans SC)**](https://github.com/adobe-fonts/source-han-sans/tree/release)。
    -   在该页面找到 **Region-specific Subset OTFs** -> **China (中国)** 并下载 ZIP 文件。
2.  解压下载的 `SourceHanSansSC.zip` 文件。
3.  从解压后的文件中，找到 `SourceHanSansSC-Regular.otf`。
4.  将 `SourceHanSansSC-Regular.otf` 复制到本项目的根目录下，并**将其重命名为 `chinese.ttf`**。
#### 3. 编译步骤
1.  在您的电脑上打开一个终端 (PowerShell 或 CMD)。
2.  进入本项目的根目录：
    ```sh
    cd /path/to/your/go-face-id-validator
    ```
3.  安装 Fyne 命令行工具：
    ```sh
    go install fyne.io/tools/cmd/fyne@latest
    ```
4.  使用 Fyne 工具将字体文件打包成 Go 代码：
    ```sh
    fyne bundle -o bundled.go chinese.ttf
    ```
    > 这步执行成功后，项目下会生成一个 `bundled.go` 文件。
5.  下载所有项目依赖：
    ```sh
    go mod tidy
    ```
6.  编译生成最终的可执行文件：
    ```sh
    # 该命令会生成一个运行时不带黑色控制台窗口的 .exe 文件
    go build -ldflags="-H windowsgui"
    ```
7.  编译成功后，您就会在项目目录下找到 `go-face-id-validator.exe` 文件。

## API 客户端定制

该应用程序被设计为可扩展的。如果你想使用自己的 API 服务，可以遵循以下步骤：

1.  **实现 `api.Client` 接口**：在 `api/client.go` 文件中定义了 `Client` 接口。您可以创建自己的结构体和方法来实现这个接口。
2.  **在 `main.go` 中使用你的客户端**：在 `main.go` 文件中，找到创建 API 客户端实例的地方 (`api.NewAliyunClient(...)`)，并将其替换为您自己的客户端实例。 
