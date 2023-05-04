package main

import (
	"log"
	"net"
	"fmt"
	"time"
	"image/color"
	"os"
	"bufio"
	"strings"
	"os/exec"
	"math/rand"
	"strconv"
	"syscall"

	"github.com/go-toast/toast"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/canvas"
)

var IP string = "1.12.245.232"

func create_remo_port() string {
	rand.Seed(time.Now().UnixNano())

	nums := rand.Intn(9000) + 1000

	port := strconv.Itoa(nums)

	return port
}

func Connect_cloud(ip string, port string) bool {
	ip1 := ip + ":" + port
	//创建socket套接字并连接
	conn, err := net.Dial("tcp", ip1)
	if err != nil {
		fmt.Println("Error Connecttion", err)
		return false
	}
	defer conn.Close()
	return true
}

var REMOTE_PORT string
var KEY bool = false

func Change_Text(port string) string {
	filePath := "bin/frp_0.44.0_windows_amd64/frpc.ini"
	file, err:= os.OpenFile(filePath,os.O_RDWR,0644)
	if err != nil{
		log.Fatalf("Failed to open file:%s",err)
	}
	defer file.Close()

// 读取文件内容
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Failed to read file:", err)
	}

	//修改指定内容的行
	var newLines []string
	var nums string = create_remo_port()
	for _, line := range lines {
		if strings.HasPrefix(line, "local_port = ") {
			line = "local_port = "+port
		}
		if strings.HasPrefix(line, "remote_port = "){
			line = "remote_port = "+nums
			REMOTE_PORT = nums
		}
		newLines = append(newLines,line)
	}

	// 将修改后的内容写回文件
	file, err = os.Create("bin/frp_0.44.0_windows_amd64/frpc.ini")
	if err != nil {
		fmt.Println("Failed to create file:", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range newLines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()

	// 打开修改后的文件
	file1,err := os.OpenFile(filePath,os.O_RDWR,0644)
	if err != nil{
		log.Fatalf("Failed to open file:%s",err)
	}
	defer file1.Close()
	fmt.Println("file ok")
	fmt.Printf("File content after modification: %sand%s\n", port,nums)
	return nums
}

func Open_cmd(value bool,key bool){
	path := "bin/frp_0.44.0_windows_amd64/"
	cmd := exec.Command("cmd","/c","frpc.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd1 := exec.Command("taskkill", "/F", "/IM", "frpc.exe")

	cmd.Dir = path
	if value==true && key == true{
		cmd.Start()
	}
	if value==false && key == true{
		cmd1.Run()
	}

}

func Noti_Win_True(){
	ip := IP+":"+REMOTE_PORT
	notification := toast.Notification{
		AppID: "Microsoft.Windows.Shell.RunDialog",
		Title: "SplashNirvana2K",
		Message: "Your IP and Port is:"+ip,
		Actions: []toast.Action{
			{"protocol", "My GitHub", "https://github.com/Pumnn1ayLee/PMLfrp-v1.0"},
		},
	}
	err := notification.Push()
	if err != nil{
		log.Fatalln(err)
	}
}

func Noti_Win_False(){
	notification := toast.Notification{
		AppID: "Microsoft.Windows.Shell.RunDialog",
		Title: "SplashNirvana2K",
		Message: "Your IP and Port are exactly wrong!Please modify your IP and PORT! ",
		Actions: []toast.Action{
			{"protocol", "My GitHub", "https://github.com/Pumnn1ayLee/PMLfrp-v1.0"},
		},
	}
	err := notification.Push()
	if err != nil{
		log.Fatalln(err)
	}
}

func Init_Console(w fyne.Window){

	//连接按钮
	check := widget.NewCheck("Guangzhou",func(value bool){
		log.Println("Check set to",value)
		if value && KEY{
			Noti_Win_True()
			Open_cmd(true,KEY)
		}
		if value && KEY == false{
			Noti_Win_False()
		}
		if value == false && KEY == true{
			Open_cmd(false,KEY)
		}
		if value == false && KEY == false{
			Open_cmd(false,KEY)
		}
		})

	text1 := widget.NewLabel("The list of available Channels is as follows:")

	content1 := container.New(layout.NewBorderLayout(text1,nil,check,nil),layout.NewSpacer(),text1,check)

	entry1 := widget.NewEntry()
	entry1.SetPlaceHolder("Input ip")

	entry2 := widget.NewEntry()
	entry2.SetPlaceHolder("Input Port Default:22")

	entry3 := widget.NewEntry()
	entry3.SetPlaceHolder("Local Port")

	progress := widget.NewProgressBar()
    progress.Resize(fyne.NewSize(300,300))

	//成功连接画面
	suceText := canvas.NewText("Connection Successful Welcome to SplashNirvana2K,Please wait a few times!", color.Black)
	content_suce := container.New(layout.NewCenterLayout(), suceText)

	//失败连接画面
	faiText := canvas.NewText("Your Ip or Port maybe wrong,Please wait a few times!",color.Black)
	content_fai := container.New(layout.NewCenterLayout(), faiText)


	//窗格
	tabs := container.NewAppTabs(
        container.NewTabItem("Channel Lists", content1),
		container.NewTabItem("ReadMe Important!",widget.NewLabel("Please first open your MC and then open the LAN port,\n fill in the port and IP address to enter HOME, and then\n check the button in Channels!\n It's best not to repeatedly check!")),
        container.NewTabItem("About us", widget.NewLabel("Welcome to SplashNirvana2K and En_nuyeux\n\n\n Github:https://github.com/Pumnn1ayLee\n\nGithub:https://github.com/Ennuyeux233")),
    )

	button1 := widget.NewButton("connect", func() {
		log.Println("Input ip is:", entry1.Text)
		log.Println("Input port is:",entry2.Text)
		log.Println("Input local_port is:", entry3.Text)
		var connect bool = Connect_cloud(entry1.Text, entry2.Text)
		KEY = connect
		Change_Text(entry3.Text)

		/*progress2 := container.New(layout.NewCenterLayout(), progress)*/
		if connect {
			go func() {
				for i := 0.0; i < 1.0; i += 0.01 {
					time.Sleep(time.Millisecond * 100)
					progress.SetValue(i)
				}
			}()
			w.SetContent(progress)
			time.Sleep(13 * time.Second)
			w.SetContent(content_suce)
			time.Sleep(3 * time.Second)
			w.SetContent(tabs)
		}else{
			w.SetContent(progress)
			w.SetContent(content_fai)
			time.Sleep(3 * time.Second)
			w.SetContent(tabs)
		}
	})

	content := container.NewVBox(entry1, entry2, entry3, button1)

	tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), content))

	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)
	w.Resize(fyne.Size{Width: 500, Height: 300})
}

func main() {
    a := app.New()
    w := a.NewWindow("SplashNirvana2K")
	Init_Console(w)
    w.ShowAndRun()
}